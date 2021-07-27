section .data
  message: db "Primer celda", 10, 0
  emessage equ $ - message

global _start

section .text
  _start:
;Print
mov rax, 1
mov rdi, 1
mov rsi, message
mov rdx, emessage
syscall
mov rax, 60
mov rdi, 0
syscall







; mov rax, 39
; syscall
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+6], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+5], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+4], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+3], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+2], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message+1], dl
; mov rdx, 0
; mov rbx, 10
; div rbx
; add rdx, 48
; mov [message], dl

; movdqu xmm0, [bytes]
; mov rax, 1
; mov rdi, 1
; mov rsi, message
; mov rdx, emessage
; syscall
