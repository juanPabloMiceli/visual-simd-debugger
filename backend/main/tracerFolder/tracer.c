#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <sys/user.h>   
#include <sys/reg.h>   
#include <stdio.h>
#include <sys/syscall.h>   /* For SYS_write etc */

typedef struct s_XMM_Reg {
	uint32_t TopQuarter;
	uint32_t MiddleTopQuarter;
	uint32_t MiddleBottomQuarter;
	uint32_t BottomQuarter;
} __attribute__((packed)) XMM_Reg;

void newXMMReg(XMM_Reg *reg, uint32_t TopQuarter, uint32_t MiddleTopQuarter, uint32_t MiddleBottomQuarter, uint32_t BottomQuarter);
void printXMMReg(XMM_Reg *reg);

typedef struct s_XMM_Regs {
	XMM_Reg XMM0;
	XMM_Reg XMM1;
	XMM_Reg XMM2;
	XMM_Reg XMM3;
	XMM_Reg XMM4;
	XMM_Reg XMM5;
	XMM_Reg XMM6;
	XMM_Reg XMM7;
	XMM_Reg XMM8;
	XMM_Reg XMM9;
	XMM_Reg XMM10;
	XMM_Reg XMM11;
	XMM_Reg XMM12;
	XMM_Reg XMM13;
	XMM_Reg XMM14;
	XMM_Reg XMM15;
} __attribute__((packed)) XMM_Regs;


void newXMMRegs(XMM_Regs *xmm_regs, uint32_t *p_regs);
void printXMMRegs(XMM_Regs *xmm_regs);

void newXMMReg(XMM_Reg *reg, uint32_t TopQuarter, uint32_t MiddleTopQuarter, uint32_t MiddleBottomQuarter, uint32_t BottomQuarter){
	reg->TopQuarter = TopQuarter;
	reg->MiddleTopQuarter = MiddleTopQuarter;
	reg->MiddleBottomQuarter = MiddleBottomQuarter;
	reg->BottomQuarter = BottomQuarter;
}

void printXMMReg(XMM_Reg *reg){
	printf("0x%08x%08x%08x%08x", reg->TopQuarter, reg->MiddleTopQuarter, reg->MiddleBottomQuarter, reg->BottomQuarter);
}

void newXMMRegs(XMM_Regs *xmm_regs, uint32_t *p_regs){
	newXMMReg(&xmm_regs->XMM0, p_regs[0], p_regs[1], p_regs[2], p_regs[3]);
	newXMMReg(&xmm_regs->XMM1, p_regs[4], p_regs[5], p_regs[6], p_regs[7]);
	newXMMReg(&xmm_regs->XMM2, p_regs[8], p_regs[9], p_regs[10], p_regs[11]);
	newXMMReg(&xmm_regs->XMM3, p_regs[12], p_regs[13], p_regs[14], p_regs[15]);
	newXMMReg(&xmm_regs->XMM4, p_regs[16], p_regs[17], p_regs[18], p_regs[19]);
	newXMMReg(&xmm_regs->XMM5, p_regs[20], p_regs[21], p_regs[22], p_regs[23]);
	newXMMReg(&xmm_regs->XMM6, p_regs[24], p_regs[25], p_regs[26], p_regs[27]);
	newXMMReg(&xmm_regs->XMM7, p_regs[28], p_regs[29], p_regs[30], p_regs[31]);
	newXMMReg(&xmm_regs->XMM8, p_regs[32], p_regs[33], p_regs[34], p_regs[35]);
	newXMMReg(&xmm_regs->XMM9, p_regs[36], p_regs[37], p_regs[38], p_regs[39]);
	newXMMReg(&xmm_regs->XMM10, p_regs[40], p_regs[41], p_regs[42], p_regs[43]);
	newXMMReg(&xmm_regs->XMM11, p_regs[44], p_regs[45], p_regs[46], p_regs[47]);
	newXMMReg(&xmm_regs->XMM12, p_regs[48], p_regs[49], p_regs[50], p_regs[51]);
	newXMMReg(&xmm_regs->XMM13, p_regs[52], p_regs[53], p_regs[54], p_regs[55]);
	newXMMReg(&xmm_regs->XMM14, p_regs[56], p_regs[57], p_regs[58], p_regs[59]);
	newXMMReg(&xmm_regs->XMM15, p_regs[60], p_regs[61], p_regs[62], p_regs[63]);
}

void printXMMRegs(XMM_Regs *xmm_regs){
	printf("XMM0:  ");
	printXMMReg(&xmm_regs->XMM0);
	printf("\nXMM1:  ");
	printXMMReg(&xmm_regs->XMM1);
	printf("\nXMM2:  ");
	printXMMReg(&xmm_regs->XMM2);
	printf("\nXMM3:  ");
	printXMMReg(&xmm_regs->XMM3);
	printf("\nXMM4:  ");
	printXMMReg(&xmm_regs->XMM4);
	printf("\nXMM5:  ");
	printXMMReg(&xmm_regs->XMM5);
	printf("\nXMM6:  ");
	printXMMReg(&xmm_regs->XMM6);
	printf("\nXMM7:  ");
	printXMMReg(&xmm_regs->XMM7);
	printf("\nXMM8:  ");
	printXMMReg(&xmm_regs->XMM8);
	printf("\nXMM9:  ");
	printXMMReg(&xmm_regs->XMM9);
	printf("\nXMM10: ");
	printXMMReg(&xmm_regs->XMM10);
	printf("\nXMM11: ");
	printXMMReg(&xmm_regs->XMM11);
	printf("\nXMM12: ");
	printXMMReg(&xmm_regs->XMM12);
	printf("\nXMM13: ");
	printXMMReg(&xmm_regs->XMM13);
	printf("\nXMM14: ");
	printXMMReg(&xmm_regs->XMM14);
	printf("\nXMM15: ");
	printXMMReg(&xmm_regs->XMM15);
	printf("\n");
}

void main(int argc, char *argv[]){
    struct user_fpregs_struct fpregs;
    
    char *p;
    int pid;
    long conv = strtol(argv[1], &p, 10);
    pid = conv;
    ptrace(PTRACE_ATTACH, pid, NULL, NULL);
    sleep(1);
    ptrace(PTRACE_GETFPREGS, pid, NULL, &fpregs);
    printf("\n\n%d\n\n", pid);
    XMM_Regs XMM;
    newXMMRegs(&XMM, fpregs.xmm_space);
    printXMMRegs(&XMM);
}