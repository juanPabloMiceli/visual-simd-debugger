section .data
  timeval:
    tv_sec  dd 0
    tv_usec dd 0

  bmessage  db "Sleep", 10, 0
  bmessagel equ $ - bmessage

  emessage  db "Continue", 10, 0
  emessagel equ $ - emessage
  filename db "archivo.txt", 10, 0
  p: db 12, 32,1 , 12, 321, 123,2
  global _start

section .text
  _start:
;print "sleep"
  mov rax, 4
  mov rbx, 1
  mov rcx, bmessage
  mov rdx, bmessagel
  int 0x80

;sleep for 5 seconds and 0 nanoseconds
  mov dword [tv_sec], 1
  mov dword [tv_usec], 0
  mov rax, 162
  mov rbx, timeval
  mov rcx, 0
  int 0x80

  loop:
;print "continue"
  mov rax, 4
  mov rbx, 1
  mov rcx, emessage
  mov rdx, emessagel
  int 0x80
;jmp loop

  mov rax,85;syscall number for create()
  mov rdi,filename;argv[1], the file name
  mov esi,00644q;rw,r,r
  syscall;call the kernel 

  movdqu xmm0, [p]

  mov rax, 60
mov rdi, 0
syscall
