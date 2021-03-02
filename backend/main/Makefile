make output: output.o
	ld -o output output.o

output.o: output.asm 
	nasm -f elf64 -g -F DWARF output.asm -o output.o

clean:
	rm output
	rm output.o