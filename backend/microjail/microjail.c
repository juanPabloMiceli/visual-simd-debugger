
#include <stddef.h>
#include <stdlib.h>
#include <stdio.h>   /* printf */
#include <unistd.h>  /* dup2: just for test */
// #include <seccomp.h> /* libseccomp */
#include "seccomp-bpf.h"

static int install_syscall_filter(void){

    

    struct sock_filter filter[] = {
        VALIDATE_ARCHITECTURE,
        EXAMINE_SYSCALL,
        ALLOW_SYSCALL(exit),
        KILL_PROCESS,

    };

    struct sock_fprog prog = {
		.len = (unsigned short)(sizeof(filter)/sizeof(filter[0])),
		.filter = filter,
	};

    if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0)) {
		perror("prctl(NO_NEW_PRIVS)");
		if (errno == EINVAL){
		    fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
        }
	    return 1;
	}
	if (prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &prog)) {
		perror("prctl(SECCOMP)");
		if (errno == EINVAL){
		    fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
        }
	    return 1;
	}
	return 0;
}

int main(int argc, char *argv[]){
    if(install_syscall_filter()){
        return 1;
    }

    char** exec_argv = NULL;
    char** envp = NULL;



    execve(argv[1], exec_argv, envp);

    return 1;
}