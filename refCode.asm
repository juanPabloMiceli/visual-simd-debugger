section .data
  sumar_alfas:  db 0, 0, 0, 15, 0, 0, 0, 10, 0, 0, 0, 5, 0, 0, 0, 255
  msg: DB 'Hola mundo 0', 10, 0
  largo EQU $ - msg

  global _start

section .text
  _start:
    mov rax, -1
    movdqu xmm0, [sumar_alfas]
  	xor esi, esi

  ciclo:
    mov rax, 4
    mov rbx, 1
    mov rcx, msg
    mov rdx, largo
    int 0x80

    inc byte [msg+largo-3]
    inc esi
    cmp esi, 10
    jnz ciclo

    mov rdi, 0
    mov rax, 60
    syscall 
