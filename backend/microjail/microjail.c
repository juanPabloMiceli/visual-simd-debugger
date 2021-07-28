#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <sys/random.h>
#include <sys/syscall.h>
#include <unistd.h>

#include "seccomp-bpf.h"

/**
 * Activa los filtros BPF para el proceso dejando activos unicamente exit (sin
 * restricciones) y execve (con restricciones en sus 6 parametros) msg: Posicion
 * del string donde se encuentra el path al archivo. Toda llamada a execve debe
 * tener el path en esta posicion. randomNumbers: Puntero a 3 numeros de 64 bits
 * generados de forma aleatoria. Toda llamada a execve debe tener estos 3
 * numeros en los parametros no utilizados. Dificultando enormemente el uso de
 * la misma fuera de este programa. Si la ejecucion es exitosa retorna 0. Caso
 * contrario retorna 1.
 */
static int install_syscall_filter(char *msg, u_int64_t *randomNumbers)
{
  struct sock_filter filter[] = {
      VALIDATE_ARCHITECTURE,

      // Load syscall number
      /*0*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_nr), // Load syscall number to W

      // Exit
      /*1*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_exit, 0, 1), // if syscall is exit, don't jump. Else, jump to statement 3.
      /*2*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_ALLOW),         // Allow syscall.(exit is always allowed)

      /// Execve
      /*3*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_execve, 0, 19),   // if syscall is execve, don't jump. Else, jump to statement 23. (there could be another syscall to inspect)
      /*4*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(0)),        // Load arg0 to W
      /*5*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)msg, 1, 0), // If arg0 == &msg => jump to statement 7. (first check passed)
      /*6*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),               // Kill process. Invalid arg0

      /*7*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(1)),         // Load arg1 to W
      /*8*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)NULL, 1, 0), // If arg1 == NULL => jump to statement 10. (second check passed)
      /*9*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),                // Kill process. Invalid arg1

      /*10*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(2)),         // Load arg2 to W
      /*11*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)NULL, 1, 0), // If arg2 == msg => jump to statement 13. (third check passed)
      /*12*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),                // Kill process. Invalid arg2

      /*13*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(3)),                     // Load arg3 to W
      /*14*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)randomNumbers[0], 1, 0), // If arg3 == random1 => jump to statement 16. (fourth check passed)
      /*15*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),                            // Kill process. Invalid arg3

      /*16*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(4)),                     // Load arg4 to W
      /*17*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)randomNumbers[1], 1, 0), // If arg4 = random2 => jump to statement 19. (fifth check passed)
      /*18*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),                            // Kill process. Invalid arg4

      /*19*/ BPF_STMT(BPF_LD + BPF_W + BPF_ABS, syscall_arg(5)),                     // Load arg5 to W
      /*20*/ BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, (u_int64_t)randomNumbers[2], 1, 0), // If arg5 = random3 => jump to statement 22. (all checks passed. Syscall is allowed)
      /*21*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),                            // Kill process. Invalid arg5

      /*22*/ BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_ALLOW), // Allow syscall
      /*23*/ KILL_PROCESS                                  // Kill process

  };

  struct sock_fprog prog = {
      .len = (unsigned short)(sizeof(filter) / sizeof(filter[0])),
      .filter = filter,
  };

  // Se setea el flag "no_new_privs" que indica que desde este punto, ni el
  // proceso ni ninguno de sus hijos podran obtener privilegios nuevos. Esto
  // perdura luego del exec.
  if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0))
  {
    perror("prctl(NO_NEW_PRIVS)");
    if (errno == EINVAL)
    {
      fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
      fflush(stdout);
    }
    return 1;
  }

  // Se activa el filtro creado.
  if (prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &prog))
  {
    perror("prctl(SECCOMP)");
    if (errno == EINVAL)
    {
      fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
      fflush(stdout);
    }
    return 1;
  }

  return 0;
}

int main(int argc, char *argv[])
{
  long pageSize = sysconf(_SC_PAGE_SIZE);

  char *execPath = mmap(/*address*/ NULL, pageSize, PROT_WRITE | PROT_READ, MAP_PRIVATE | MAP_ANONYMOUS, /*fd*/ -1, /*offset*/ 0);

  if (execPath == (void *)-1)
  {
    fprintf(stderr, "Mapping error. pointer: %p\n", execPath);
    exit(EXIT_FAILURE);
  }

  if (strlen(argv[1]) > pageSize / 2)
  {
    fprintf(stderr, "Path name is too long");
    exit(EXIT_FAILURE);
  }

  char *dest = strcpy(execPath, argv[1]);

  u_int64_t randomNumbers[3];
  size_t bufLen = sizeof randomNumbers;
  ssize_t err = getrandom(randomNumbers, bufLen, 0);
  if (err < bufLen)
  {
    fprintf(stderr, "Wanted %ld numbers. Got %ld numbers", bufLen, err);
    exit(EXIT_FAILURE);
  }

  if (install_syscall_filter(execPath, randomNumbers))
  {
    exit(EXIT_FAILURE);
  }

  syscall(__NR_execve, execPath, NULL, NULL, randomNumbers[0], randomNumbers[1], randomNumbers[2]);

  return 0;
}
