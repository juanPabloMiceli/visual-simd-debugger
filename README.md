# Visual SIMD Debugger

## Introduccion:

Este es un entorno de desarrollo especialmente pensado para trabajar con instrucciones SIMD.

El mismo cuenta con un formato similar al de jupyter notebook, separándose en celdas donde en las mismas se pondrá el código que se está desarrollando, por ejemplo el core de un ciclo en un filtro de una imagen, y entre las mismas se puede mostrar de manera amigable el valor de los registros XMM, en el formato deseado, de los cuales se quiera seguir el rastro.

## Modo de uso:

### Sección de datos:

Para agregar una sección de datos a nuestro código debemos iniciar la primer celda con: `;data`. Luego de eso podemos definir los datos exactamente como lo haríamos en assembler.

Es muy importante notar que solo esta celda puede ser utilizada como celda de datos, y además, esta debe ser la única función de la misma. Es decir, no se puede agregar código en la misma.

### Sección de texto:

Para esta sección alcanza con empezar a escribir código. La sección de texto empezará en la primer celda si no existe sección de datos, y en la segunda en caso contrario.

### Agregar/Borrar celdas:

Todas las celdas permiten agregar celdas, tanto arriba como debajo de ellas. Para conseguirlo se debe presionar el botón `+ Code` que se encuentra entre las celdas. El mismo es invisible hasta que se ubica el mouse a esa altura.
Para eliminar una celda es tan sencillo como presionar el botón con forma de tacho de basura que cada celda tiene arriba a la derecha.

### Imprimir registros:

Los registros XMM aceptan los siguientes formatos para ser impresos:

* 16 registros de números enteros de 8 bits `.v16_int8`
* 8 registros de números enteros de 16 bits `.v8_int16`
* 4 registros de números enteros de 32 bits `.v4_int32`
* 2 registros de números enteros de 64 bits `.v2_int64`
* 4 registros de números flotantes de precisión simple (32 bits) `.v4_float`
* 2 registros de números flotantes de precisión doble (64 bits) `.v2_double`

A su vez, los formatos de números enteros se pueden imprimir en las siguientes bases numericas:

* Base 10 con signo `/d`
* Base 10 sin signo `/u`
* Base 16 `/x`
* Base 2 en complemento A2 `/t`

Para solicitar que se imprime un registro se debe utilizar el mismo formato que se utiliza en GDB:

`;print<formato base numérica> xmm<registro a imprimir><formato de bits>`

Donde el formato de la base numérica es opcional, por defecto esta es `/d`. Pero este valor por defecto se actualiza al usar la instrucción `;print`.

El formato de los bits también es opcional y su valor por defecto es `.v16_int8`. El mismo se actualiza al usar la instrucción `;print`. Se puede setear este formato para todos los registros a la vez omitiendo el valor del registro a imprimir.

`;print` es análogo a `;p`

Algunos ejemplos de esto:

`;p xmm.v4_int16`: Seteo el valor por defecto de los bits para todos los registros a imprimir enteros de 16 bits.

`;p/x xmm1.v2_int64`: Imprimo el registro xmm1 como enteros de 64 bits en hexadecimal

`;print xmm2.v2_double`: Imprimo el registro xmm2 como números flotantes de 64 bits.

Los registros se imprimen en 2 casos.

1) Si el valor del registro fue modificado en la ejecución de la celda.
2) Si el usuario pide que se imprima el registro.

Si la impresión se realiza por el primer caso, se usarán los valores por defecto que tenga el registro en ese momento.

Una forma de "comentar" los `print` es escribiendo un caracter cualquiera en el medio. Por ejemplo:

`;=p/x xmm0.v2_int64`



### Copiar código:

El botón "copy to clipboard" hace exactamente lo que pensás. Copia todo el código que haya en las celdas al porta papeles para que lo puedas pegar cómodamente en el assembler. Un detalle de esto es que si hay alguna celda que no querés copiar, podés agregar un:
`;nope` en esa celda y esta no va a aparecer en la copia.


### Shortcuts de teclado:

* `Ctrl+Enter`: Ejecutar código
* `Ctrl+ArrowDown`: Inserta una celda debajo de la celda actual. Si no se está parado sobre ninguna celda, la misma se insertará al final.
* `Ctrl+ArrowUp`: Inserta una celda arriba de la celda actual. Si no se está parado sobre ninguna celda, la misma se insertará al principio.
* `Ctrl+Alt+D`: Elimina la celda actual. **Esta acción no se puede deshacer.**