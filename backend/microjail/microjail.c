#include <stddef.h>
#include <stdlib.h>
#include <stdio.h>   
#include <unistd.h>  
#include "seccomp-bpf.h"
#include <sys/syscall.h>
#include <sys/mman.h>
#include <sys/random.h>

       

static int install_syscall_filter(char* msg, long unsigned int *randomNumbers){

    struct sock_filter filter[] = {
        VALIDATE_ARCHITECTURE,
        EXAMINE_SYSCALL,
        ALLOW_SYSCALL(exit),

        ///Execve
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, __NR_execve, 0, 4), //if execve sys, don't jump
        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(0)),    //Load arg0 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)msg, 1, 0), //If arg0 = msg => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(1)),    //Load arg1 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)NULL, 1, 0), //If arg0 = msg => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(2)),    //Load arg2 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)NULL, 1, 0), //If arg0 = msg => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(3)),    //Load arg3 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)randomNumbers[0], 1, 0), //If arg3 = random1 => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(4)),    //Load arg4 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)randomNumbers[1], 1, 0), //If arg4 = random2 => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, syscall_arg(5)),    //Load arg5 to W
        BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, (long unsigned int)randomNumbers[2], 1, 0), //If arg5 = random3 => jump (do not kill)
        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),            //Else => kill

        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_ALLOW),
        KILL_PROCESS

    };

    struct sock_fprog prog = {
		.len = (unsigned short)(sizeof(filter)/sizeof(filter[0])),
		.filter = filter,
	};

    if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0)) {
		perror("prctl(NO_NEW_PRIVS)");
		if (errno == EINVAL){
		    fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
            fflush(stdout);
        }
	    return 1;
	}
	if (prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &prog)) {
		perror("prctl(SECCOMP)");
		if (errno == EINVAL){
		    fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
            fflush(stdout);

        }
	    return 1;
	}
    
	return 0;
}

int main(int argc, char *argv[]){
    

    long pageSize = sysconf(_SC_PAGE_SIZE);

    
    char *execPath = mmap(NULL, pageSize, PROT_WRITE | PROT_READ, MAP_PRIVATE | MAP_ANONYMOUS, 0, 0);

    if (execPath == (void *)-1){
        printf("Mapping error. pointer: %p\n", execPath);
        fflush(stdout);
        return 1;
    }

    char *dest = strcpy(execPath,argv[1]);
    if(dest != execPath){
        printf("Wanted dest %p. Got %p.", execPath, dest);
        fflush(stdout);
        return 1;
    }

    // printf("argv: %p\n", &argv[1]);
    // printf("exec path: %p\n", &execPath[0]);
    // sleep(20);

    long unsigned int randomNumbers[3*8];
    size_t bufLen = sizeof randomNumbers; 
    ssize_t err = getrandom(randomNumbers, bufLen, 0);
    if(err < bufLen){
        printf("Wanted %ld numbers. Got %ld numbers", bufLen, err);
        fflush(stdout);
        return 1;
    }  
    
    if(install_syscall_filter(execPath, randomNumbers)){
        return 1;
    }

    syscall(__NR_execve, execPath, NULL, NULL, randomNumbers[0],randomNumbers[1],randomNumbers[2]);

    return 0;
}
