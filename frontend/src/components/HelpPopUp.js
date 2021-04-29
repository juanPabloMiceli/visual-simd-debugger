import React from 'react';
import ReactDOM from 'react-dom'
import { useContextData, useContextToggleHelp } from "./Context" 
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"


export default function HelpPopUp(){
    const VisualizerData = useContextData()
    const toggleHelp = useContextToggleHelp()


    if(!VisualizerData.HelpActive) return null

    return ReactDOM.createPortal(
        <>
            <div className='overlay' onClick={toggleHelp}></div>
            <div className='help-popup'>
            <button className="btn-close-help" id="closeHelpButton" onClick={toggleHelp}><FontAwesomeIcon icon="window-close" fixedWidth/></button>
            {explanation}
            {shortcutsTable}
            {commands}
            </div>
        </>, document.getElementById('portal')
    );
}

let explanation = (<div className='commandsContainer' align='justify'><h2>Introducción:</h2>
<p>Este es un entorno de desarrollo especialmente pensado para trabajar con instrucciones SIMD.</p>
<p>El mismo cuenta con un formato similar al de jupyter notebook, separándose en celdas donde en las mismas se pondrá el código que se está desarrollando, por ejemplo el core de un ciclo en un filtro de una imagen, y entre las mismas se puede mostrar de manera amigable el valor de los registros XMM, en el formato deseado, de los cuales se quiera seguir el rastro.</p>
<h2>Modo de uso:</h2>
<h3>Sección de datos:</h3>
<p>La primer celda siempre corresponde a la celda de datos y está señalizada por la etiqueta <kbd>section .data</kbd>, los datos se definen exactamente como lo haríamos en assembler.</p>
<p>Es muy importante notar que solo esta celda puede ser utilizada como celda de datos, y además, esta debe ser la única función de la misma. Es decir, no se puede agregar código en la misma.</p>
<p>Esta celda no se puede eliminar.</p>
<h3>Sección de texto:</h3>
<p>El inicio de esta sección está marcado por la etiqueta <kbd>section .text</kbd>, es decir que todo el texto que esté en las celdas debajo de esta etiqueta será considerado código.</p>
<p>No se permite eliminar todas las celdas de texto, por lo cual siempre habrá como mínimo una celda en la sección de texto.</p>
<h3>Código en C:</h3>
<p>Por motivos de seguridad este entorno no permite ningún llamado a funciones de la <kbd>libc</kbd>. Por hoy vas a tener que programar únicamente en assembler.</p>
<h3>Syscalls:</h3>
<p>Este entorno no acepta ninguna syscall. Si tratas de ejecutar cualquiera de estas el programa va a terminar.</p>
<h3>Agregar/Borrar celdas:</h3>
<p>Todas las celdas permiten agregar celdas, tanto arriba como debajo de ellas. Para conseguirlo se debe presionar el botón <kbd>+ Code</kbd> que se encuentra entre las celdas. El mismo es invisible hasta que se ubica el mouse a esa altura.
Para eliminar una celda es tan sencillo como presionar el botón con forma de tacho de basura que cada celda tiene arriba a la derecha.</p>
<h3>Copiar código:</h3>
<p>El botón <kbd>copy to clipboard</kbd> hace exactamente lo que pensás. Copia todo el código que haya en las celdas al porta papeles para que lo puedas pegar cómodamente en el assembler. Si no querés que alguna celda se copie al assembler podés apretar el botón con forma de ojo arriba a la derecha de la celda.</p>
<h3>Limpiar código:</h3>
<p>Como limpiar el código de todas las celdas a la vez es medio incómodo, se puede apretar el botón <kbd>Clean Code</kbd> para realizar esta acción. <strong>La misma no se puede deshacer</strong>, por lo cual se debe confirmar el uso de este botón mediante un pop-up.</p>
<h3>Esconder registros:</h3>
<p>Para esconder registros que no se quieran imprimir hay 2 alternativas: la primera es usar el comando <kbd>hide</kbd> que se explica en la sección de comandos, este comando nos deja esconder registros individualmente. Si lo que se quiere hacer es esconder todos los registros provenientes de una celda específica se puede presionar el botón con forma de ojo que se encuentra arriba a la derecha en la celda.</p>
<h3>Guardado automático:</h3>
<p>Está activado el guardado automático, es decir que todo cambio que se realice en la página se guardará en caso de cerrar y volver a abrir el navegador. <span style={{fontWeight:'bold'}}>El guardado automático funciona para una sola pestaña a la vez.</span> Si se está trabajando con 2 pestañas a la vez y se reinicia una pestaña, los cambios conservados van a ser los de la última pestaña modificada. <span style={{fontWeight:'bold'}}>Se recomienda fuertemente trabajar con una pestaña a la vez.</span> Si se quiere recuperar el ejemplo inicial se debe limpiar el <kbd>localStorage</kbd> del navegador. Esto debería ser modificado en un futuro.</p>
</div>)

let commands = (<div className='commandsContainer'>
    <p id='commandsTitle'>Comandos</p>
    <div style={{minHeight: '100px'}}>
        <p id='printTitle'>Imprimir registros</p>
        <p id='examplePrint'>;p<span style={{color: '#A00'}}>/u</span> <span style={{color: '#0A0'}}>xmm0</span><span style={{color: '#00A'}}>.v8_int16</span></p>
    </div>
    <div>
        <p style={{paddingLeft: '22px'}}>Si se quiere imprimir un registro en determinada celda, alcanza con incluir esta instrucción en esa celda especificando la base, el registro y el formato deseados de impresión. Los parámetros se especifican de la siguiente manera:</p>

        <ul id='listForPrint'>
            <li style={{listStyleType: 'square', color: '#A00'}}><span style={{fontWeight: 'bold'}}>Base: </span>Elegimos la base en la cual imprimir el registro. Este parámetro es opcional y es ignorado si se imprimen números flotantes. Bases posibles:</li>
            <ul style={{color: '#A00'}}>
                <li style={{paddingTop: '8px'}}>/d: Base 10 con signo. Esta es la base por defecto.</li>
                <li style={{paddingTop: '8px'}}>/u: Base 10 sin signo.</li>
                <li style={{paddingTop: '8px'}}>/t: Base 2 en complemento A2.</li>
                <li style={{paddingTop: '8px', paddingBottom: '8px'}}>/x: Base 16.</li>
            </ul>
            <li style={{listStyleType: 'square', color: '#0A0', paddingBottom: '8px'}}><span style={{fontWeight: 'bold'}}>Registro: </span>Elegimos el registro a imprimir. Los registros validos van desde xmm0 hasta xmm15 inclusive. Este parámetro es obligatorio.</li>
            <li style={{listStyleType: 'square', color: '#00A', paddingBottom: '8px'}}><span style={{fontWeight: 'bold'}}>Formato: </span>Elegimos el formato en el cual queremos subdividir el registro. Este parámetro es obligatorio. Los formatos posibles son:</li>
            <ul style={{color: '#00A'}}>
                <li style={{paddingTop: '8px'}}>.v16_int8: El XMM se subdivide en 16 registros enteros de 8 bits cada uno.</li>
                <li style={{paddingTop: '8px'}}>.v8_int16: El XMM se subdivide en 8 registros enteros de 16 bits cada uno.</li>
                <li style={{paddingTop: '8px'}}>.v4_int32: El XMM se subdivide en 4 registros enteros de 32 bits cada uno.</li>
                <li style={{paddingTop: '8px'}}>.v2_int64: El XMM se subdivide en 3 registros enteros de 64 bits cada uno.</li>
                <li style={{paddingTop: '8px'}}>.v4_float: El XMM se subdivide en 4 registros de punto flotante de 32 bits cada uno.</li>
                <li style={{paddingTop: '8px', paddingBottom: '8px'}}>v2_double:  El XMM se subdivide en 2 registros de punto flotante de 64 bits cada uno.</li>
            </ul>
        </ul>
    </div>
    <div style={{minHeight: '100px'}}>
        <p id='printTitle'>Esconder registros</p>
        <p id='examplePrint'>;hide <span style={{color: '#A00'}}>xmm0</span></p>
    </div>
    <div>
        <p style={{paddingLeft: '22px'}}>En caso que no se quiera mostrar algún registro modificado se puede agregar esta instrucción para quitar su impresión.</p>
        <ul id='listForPrint'>
            <li style={{color: '#A00'}}><span style={{fontWeight: 'bold'}}>Registro: </span>Elegimos el registro a imprimir. Los registros validos van desde xmm0 hasta xmm15 inclusive. Este parámetro es obligatorio.</li>
        </ul>
    </div>
    
</div>)

let shortcutsTable = (<table id="shortcuts">
<thead>
    <tr>
        <th colSpan='2'>Shortcuts de teclado</th>
    </tr>
</thead>
<tbody>
    <tr>
        <td className='definition'>Abrir/Cerrar ayuda</td>
        <td className='hotkey'><span><kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>H</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Correr código</td>
        <td className='hotkey'><span><kbd>Ctrl</kbd>+<kbd>Enter</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Insertar celda debajo</td>
        <td className='hotkey'><span><kbd>Ctrl</kbd>+<kbd>ArrowDown</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Insertar celda encima</td>
        <td className='hotkey'><span><kbd>Ctrl</kbd>+<kbd>ArrowUp</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Borrar celda actual</td>
        <td className='hotkey'><span><kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>D</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Moverse a la celda de abajo</td>
        <td className='hotkey'><span><kbd>Alt</kbd>+<kbd>ArrowDown</kbd></span></td>
    </tr>
    <tr>
        <td className='definition'>Moverse a la celda de arriba</td>
        <td className='hotkey'><span><kbd>Alt</kbd>+<kbd>ArrowUp</kbd></span></td>
    </tr>
</tbody>
</table>)