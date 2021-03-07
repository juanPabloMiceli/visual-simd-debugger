all: install output hello

install: microjail
	sudo cp microjail /usr/bin/

microjail: microjail.c
	gcc $^ -g  -o microjail -lseccomp

output: output.o
	ld -o output output.o

output.o: output.asm 
	nasm -f elf64 -g -F DWARF output.asm -o output.o

hello: hello.o
	ld -o hello hello.o

hello.o: hello.asm
	nasm -f elf64 -g -F DWARF hello.asm -o hello.o


clean:
	rm microjail
	rm output
	rm output.o
	rm hello.o
	rm hello