# Visual SIMD Debugger

## Introdução:

Este é um ambiente de desenvolvimento especialmente projetado para trabalhar com instruções SIMD.

Possui um formato semelhante ao do notebook jupyter, separado em células onde o código que está sendo desenvolvido será colocado nelas, por exemplo o núcleo de um ciclo em um filtro de imagem, e entre eles pode ser mostrado como Em um amigável forma, o valor dos registros XMM, no formato desejado, dos quais você deseja acompanhar.

## Modo de uso:

### Seção de dados:

A primeira célula sempre corresponde à célula de dados e é sinalizada pelo rótulo`section .data`, os dados são definidos exatamente como faríamos no assembler.

É muito importante notar que apenas esta célula pode ser usada como célula de dados e, além disso, esta deve ser a única função dela. Ou seja, você não pode adicionar código nele.

Esta célula não pode ser excluída.

### Seção de texto:

O início desta seção é marcado pela tag`section .text`, ou seja, todo o texto que estiver nas células abaixo deste rótulo será considerado código.

A exclusão de todas as células de texto não é permitida, portanto, sempre haverá pelo menos uma célula na seção de texto.

### Código em C:

Por razões de segurança, este ambiente não permite chamadas para funções libc. Por hoje você só terá que programar em assembler.

### Syscalls:

Este ambiente não aceita syscall. Se você tentar executar qualquer um deles, o programa será encerrado.

### Adicionar / excluir células:

Todas as células permitem que você adicione células, tanto acima quanto abaixo delas. Para conseguir isso, você deve pressionar o botão`+ Code`encontrado entre as células. Ele fica invisível até que o mouse esteja localizado nessa altura.
Para deletar uma célula basta apertar o botão em forma de lata de lixo que cada célula possui no canto superior direito.

### Imprimir registros:

Os registros XMM aceitam os seguintes formatos de impressão:

-   16 registradores inteiros de 8 bits`.v16_int8`
-   8 registradores inteiros de 16 bits`.v8_int16`
-   4 registradores inteiros de 32 bits`.v4_int32`
-   2 registradores inteiros de 64 bits`.v2_int64`
-   4 registradores de número flutuante de precisão única (32 bits)`.v4_float`
-   2 registradores de número flutuante de precisão dupla (64 bits)`.v2_double`

Por sua vez, os formatos inteiros podem ser impressos nas seguintes bases numéricas:

-   Base Assinada 10`/d`
-   Base 10 sem sinal`/u`
-   Base 16`/x`
-   Base 2 em complemento A2`/t`

Para solicitar a impressão de um registro, deve-se utilizar o mesmo formato utilizado no GDB:

`;print<formato base numérica> xmm<registro a imprimir><formato de bits>`

Onde o formato da base numérica é opcional, por padrão é`/d`. Mas este valor padrão é atualizado ao usar a instrução`;print`.

O formato de bits também é opcional e seu valor padrão é`.v16_int8`. É atualizado ao usar a instrução`;print`. Você pode definir este formato para todos os registros ao mesmo tempo, omitindo o valor do registro a ser impresso.

`;print`é análogo a`;p`

Alguns exemplos disso:

`;p xmm.v4_int16`: Eu defino o valor padrão dos bits para todos os registros a serem impressos em inteiros de 16 bits.

`;p/x xmm1.v2_int64`: Eu imprimo o registro xmm1 como inteiros de 64 bits em hexadecimal

`;print xmm2.v2_double`: Eu imprimo o registro xmm2 como números flutuantes de 64 bits.

Os registros são impressos em 2 casos.

1) Se o valor do registro foi alterado na execução da célula.
2) Se o usuário solicitar a impressão do registro.

Se a impressão for feita pelo primeiro caso, serão utilizados os valores padrão que o registro possui naquele momento.

Uma maneira de "comentar" sobre`print`é escrevendo qualquer caractere no meio. Por exemplo:

`;=p/x xmm0.v2_int64`

### Copiar código:

O botão "copiar para a área de transferência" faz exatamente o que você pensa. Copie todo o código que está nas células para a área de transferência para que você possa colá-lo facilmente no montador. Um detalhe disso é que se houver alguma célula que você não deseja copiar, você pode adicionar:`;nope`nessa célula e não aparecerá na cópia.

### Limpiar código:

Como limpar o código de todas as células de uma vez é meio estranho, você pode pressionar o botão`Clean Code`para realizar esta ação.**Não pode ser desfeito**, para o qual o uso deste botão deve ser confirmado por meio de um pop-up.

### Esconder registros:

Para ocultar os registros que você não deseja imprimir, existem 2 alternativas: a primeira é usar o comando`hide`que é explicado na seção de comandos, este comando nos permite ocultar os registros individualmente. Se o que você deseja fazer é ocultar todos os registros provenientes de uma célula específica, você pode pressionar o botão com a forma de um olho que se encontra no canto superior direito da célula.

### Shortcuts de teclado:

-   `Ctrl+Enter`: Execute o código.
-   `Ctrl+ArrowDown`: Insira uma célula abaixo da célula atual. Se você não estiver em nenhuma célula, ela será inserida no final.
-   `Ctrl+ArrowUp`: Insira uma célula acima da célula atual. Se você não estiver em nenhuma célula, ela será inserida no início.
-   `Ctrl+Alt+D`: Exclua a célula atual.**Essa ação não pode ser desfeita.**
-   `Alt+ArrowDown`: Move o cursor uma célula abaixo da célula atual.
-   `Alt+ArrowUp`: Move o cursor uma célula acima da célula atual.
