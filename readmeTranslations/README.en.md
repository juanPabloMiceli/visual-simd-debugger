# Visual SIMD Debugger

Supported languages:

[🇪🇸](/../../README.md)[🇬🇧](/../../readmes/README.en.md)[🇵🇹](/readmes/README.pt.md)[🇫🇷](/../readmes/README.fr.md)[🇮🇳](/../readmes/README.hi.md)

## Introduction:

This is a development environment specially designed for working with SIMD instructions.

It has a format similar to that of jupyter notebook, separated into cells where the code that is being developed will be placed in them, for example the core of a cycle in an image filter, and between them it can be shown as In a friendly way, the value of the XMM records, in the desired format, of which you want to keep track.

## How to use:

### Data section:

The first cell always corresponds to the data cell and is signaled by the label`section .data`, the data is defined exactly as we would in assembler.

It is very important to note that only this cell can be used as a data cell, and furthermore, this must be the only function of it. That is, you cannot add code in it.

This cell cannot be deleted.

### Text section:

The beginning of this section is marked by the tag`section .text`, that is, all the text that is in the cells below this label will be considered code.

Deleting all text cells is not allowed, so there will always be at least one cell in the text section.

### Code in C:

For security reasons this environment does not allow any calls to libc functions. For today you will only have to program in assembler.

### Syscalls:

This environment does not accept any syscall. If you try to run any of these the program will terminate.

### Add / Delete cells:

All cells allow you to add cells, both above and below them. To achieve this you must press the button`+ Code`found between cells. It is invisible until the mouse is located at that height.
To delete a cell, it is as simple as pressing the button with the shape of a garbage can that each cell has at the top right.

### Print records:

XMM records accept the following formats for printing:

-   16 8-bit integer registers`.v16_int8`
-   8 16-bit integer registers`.v8_int16`
-   4 32-bit integer registers`.v4_int32`
-   2 64-bit integer registers`.v2_int64`
-   4 single precision floating number registers (32 bits)`.v4_float`
-   2 double precision floating number registers (64 bit)`.v2_double`

In turn, integer formats can be printed in the following number bases:

-   Signed Base 10`/d`
-   Base 10 unsigned`/u`
-   Base 16`/x`
-   Base 2 in complement A2`/t`

To request that a record be printed, the same format used in GDB must be used:

`;print<NB> xmm<PR><BF>`

Where`NB`is the format of the number base and is optional, by default this is`/d`. But this default value is updated when using the statement`;print`.

`PR`is the record to print.

`BF`is the bit format and is also optional and its default value is`.v16_int8`. It is updated when using the instruction`;print`. You can set this format for all records at the same time, omitting the value of the record to be printed.

`;print`is analogous to`;p`

Some examples of this:

`;p xmm.v4_int16`: I set the default value of the bits for all the registers to be printed in 16-bit integers.

`;p/x xmm1.v2_int64`: I print the xmm1 register as 64 bit integers in hexadecimal

`;print xmm2.v2_double`: I print the xmm2 register as 64 bit floating numbers.

Records are printed in 2 cases.

1) If the value of the register was modified in the execution of the cell.
2) If the user requests that the record be printed.

If the printing is done by the first case, the default values ​​that the registry has at that moment will be used.

A way to "comment" on`print`is by writing any character in the middle. For instance:

`;=p/x xmm0.v2_int64`

### Copy code:

El botón "copy to clipboard" hace exactamente lo que pensás. Copia todo el código que haya en las celdas al porta papeles para que lo puedas pegar cómodamente en el assembler. Un detalle de esto es que si hay alguna celda que no querés copiar, podés agregar un:
`;nope`in that cell and it will not appear in the copy.

### Clean code:

As clearing the code from all cells at once is kind of awkward, you can press the button`Clean Code`to perform this action.**It cannot be undone**, for which the use of this button must be confirmed by means of a pop-up.

### Hide records:

To hide records that you do not want to print there are 2 alternatives: the first is to use the command`hide`which is explained in the command section, this command lets us hide records individually. If what you want to do is hide all the records coming from a specific cell, you can press the button with the shape of an eye that is found in the upper right of the cell.

### Keyboard Shortcuts:

-   `Ctrl+Enter`: Execute code.
-   `Ctrl+ArrowDown`: Insert a cell below the current cell. If you are not standing on any cell, it will be inserted at the end.
-   `Ctrl+ArrowUp`: Insert a cell above the current cell. If you are not standing on any cell, it will be inserted at the beginning.
-   `Ctrl+Alt+D`: Delete the current cell.**This action can not be undone.**
-   `Alt+ArrowDown`: Moves the cursor one cell below the current cell.
-   `Alt+ArrowUp`: Moves the cursor one cell above the current cell.
