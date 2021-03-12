import axios from 'axios'
import React, {useContext, useState} from 'react'

const ContextData  = React.createContext()
const ContextUpdateData = React.createContext()
const ContextXMMRegisters = React.createContext()
const ContextSubmit = React.createContext()
const ContextNewCell = React.createContext()
const ContextNewCellDown = React.createContext()
const ContextDeleteCell = React.createContext()
const ContextCopyToClipBoard = React.createContext()
const ContextCleanCode = React.createContext()
// const local_host = "http://localhost:8080"



export function useContextData(){
    return useContext(ContextData)
}

export function useContextUpdateData(){
    return useContext(ContextUpdateData)
}

export function useContextXMMRegisters(){
    return useContext(ContextXMMRegisters)
}

export function useContextSubmit(){
    return useContext(ContextSubmit)
}

export function useContextNewCell(){
    return useContext(ContextNewCell)
}

export function useContextNewCellDown(){
    return useContext(ContextNewCellDown)
}

export function useContextDeleteCell(){
    return useContext(ContextDeleteCell)
}

export function useContextCopyToClipBoard(){
    return useContext(ContextCopyToClipBoard)
}

export function useContextCleanCode(){
    return useContext(ContextCleanCode)
}


axios.defaults.headers.common = {
	"Content-Type": "application/json"
  }

function cleanOutput(cellsData){
    for(let i = 0; i < cellsData.length; i++){
        cellsData[i].output = []
    }

    return cellsData
}



function Provider({children}){
    
    const [CellsData, setCellsData] = useState(initExample());
    const [TotalCells, setTotalCells] = useState(10);
    const [ConsoleOutput, setConsoleOutput] = useState("");

    function initCell(_id){
        return {
            id: _id,
            code: "",
            output: []
        }
    }

    function initExample() {

        let cellsData = []

        for(let i = 0; i < 10; i++){
            cellsData.push(initCell(i))
        }

        cellsData[0].code = ";Este programa se encarga de pasar los 2 primeros pixeles en memoria a escala de grises.\n"
        cellsData[0].code += ";Esto se hace mediante la formula: gris = (rojo + 2 * verde + azul)/4.\n"
        cellsData[0].code += ";Una vez se obtienen los 2 valores se les hace un broadcast de tal modo que "
        cellsData[0].code += "la parte baja del registro este formada por el primer pixel gris repetido 2 veces y "
        cellsData[0].code += "la parte alta contenga el segundo pixel gris.\n\n"
        cellsData[0].code += "pixeles: db 173, 68, 144, 255, 16, 54, 231, 255, 29, 178, 50, 255, 79, 211, 203, 255\n"
        cellsData[0].code += "mascara: db 0, 0, 0, 2, 0, 0, 0, 2, 1, 1, 1, 2, 1, 1, 1, 2"
    
        cellsData[1].code = ";Cargo los pixeles en los registros xmm y los imprimo como enteros de 8 bits sin signo.\n"
        cellsData[1].code += "movdqu xmm0, [mascara]\n"
        cellsData[1].code += "movdqu xmm1, [pixeles]\n"
        cellsData[1].code += ";p/u xmm0.v16_int8\n"
        cellsData[1].code += ";p/u xmm1.v16_int8"
    
        cellsData[2].code = ";Me quedo solo con la parte baja de xmm1 en forma de enteros de 16 bits e imprimo el registro de esa forma.\n"
        cellsData[2].code += "pmovzxbw xmm1, xmm1\n"
        cellsData[2].code += ";p xmm1.v8_int16"

        cellsData[3].code = ";Reemplazo el valor alfa de cada pixel con el valor verde del mismo. Hacemos esto para ambos pixeles.\n"
        cellsData[3].code += ";Esto nos permite conseguir el valor de gris haciendo sumas horizontales.\n"
        cellsData[3].code += "pshufhw xmm1, xmm1, 0b01100100\n"
        cellsData[3].code += "pshuflw xmm1, xmm1, 0b01100100"
    
        cellsData[4].code = ";Hacemos la primer suma horizontal.\n"
        cellsData[4].code += "phaddw xmm1, xmm1"

        cellsData[5].code = ";Hacemos la segunda suma horizontal.\n"
        cellsData[5].code += ";Ahora tenemos en los 2 registros de la derecha el los valores de R + 2G + B de cada pixel.\n"
        cellsData[5].code += "phaddw xmm1, xmm1"

        cellsData[6].code = ";Shifteamos a derecha para dividir los resultados por 4 y asi obtener el valor gris del pixel.\n"
        cellsData[6].code += "psrlw xmm1, 2"

        cellsData[7].code = ";Volvemos a almacenar los pixeles conseguidos como enteros de 8 bits. En este puento vuelvo a imprimir en 8 bits.\n"
        cellsData[7].code += "packuswb xmm1, xmm1\n"
        cellsData[7].code += ";p xmm1.v16_int8"

        cellsData[8].code = ";Inserto un 255 en el valor de alfa para restaurar el valor original.\n"
        cellsData[8].code += "xor rax, rax\n"
        cellsData[8].code += "dec rax\n"
        cellsData[8].code += "pinsrb xmm1, al, 2"

        cellsData[9].code = ";Uso la mascara en xmm0 para hacer un broadcast de los valores conseguidos y del alfa "
        cellsData[9].code += "como decia el enunciado.\n"
        cellsData[9].code += "pshufb xmm1, xmm0"

        


        return cellsData
        
    }

    function submitCode(){
        let url = new URL(window.location.href)
        url.port = "8080"

        setCellsData(cleanOutput(CellsData))
        axios.post(url + "codeSave", JSON.stringify({
            CellsData: CellsData
        }))
        .then(response =>{
            setConsoleOutput(response.data.ConsoleOut) 
            updateXMMData(response.data.CellRegs) 
        }) 
        .catch(error => {
            setConsoleOutput(error) 
        })
    }



    function fixIndexing(cells){

        for(let i=0; i<cells.length; i++){
            cells[i] = {id: i, code: cells[i].code, output: cells[i].output}
        }

        return cells
    }
    function newCell(e, buttonNumber){
        e.preventDefault()
        let copy = JSON.parse(JSON.stringify(CellsData))//This is the only way of doing a depth copy
        setCellsData([initCell(0)])
        
        copy.splice(buttonNumber, 0, initCell(TotalCells))

        copy = fixIndexing(copy)
        setCellsData(copy)

        setTotalCells(TotalCells + 1)
        document.getElementById(buttonNumber.toString()).focus()
    }


    function deleteCell(e, buttonNumber){
        e.preventDefault()
        if(TotalCells > 2){
            let copy = JSON.parse(JSON.stringify(CellsData))//This is the only way of doing a depth copy
            copy.splice(buttonNumber, 1)
            copy = fixIndexing(copy)
            setCellsData(copy)
            setTotalCells(TotalCells - 1)
        }
    }

    function wantToPrint(cell){

        if(cell.code.toLowerCase().includes(';nope')){
            return false;
        }
        return true;
    }

    function copyStringToClipboard (str) {
        var el = document.createElement('textarea');
        el.value = str;
        el.setAttribute('readonly', '');
        el.style = {position: 'absolute', left: '-9999px'};
        document.body.appendChild(el);
        el.select();
        document.execCommand('copy');
        document.body.removeChild(el);
     }

    function copyToClipBoard(e){
        e.preventDefault()
        let resText = ""

        for(let i=0; i<TotalCells; i++){
            if(wantToPrint(CellsData[i]) && CellsData[i].code !== ""){
                resText += CellsData[i].code + "\n\n"
                console.log(CellsData[i].code)
            }
        }

        if(resText === ""){
            resText += " "
        }

        copyStringToClipboard(resText)
        
    }

    function cleanCode(e){
        if(window.confirm("Are you sure you want to delete ALL code?\nThere is no turning back.")){
            e.preventDefault()
            let copy = [initCell(0), initCell(1)]
            setCellsData(copy)
        }
        
    }
    function updateXMMData(data){
        let copy = JSON.parse(JSON.stringify(CellsData))
        for(let i = 0; i < copy.length; i++){
            copy[i].output = data[i]
        }

        setCellsData(copy)

    }

    function updateCodeData(cellId, newCode, textArea){
        let copy = JSON.parse(JSON.stringify(CellsData))
        copy[cellId].code = newCode
        setCellsData(copy)
    }

    function getXMMRegisters(cellId){
        console.log("GetXMMRegisters")
        console.log(CellsData[cellId].output)
        return CellsData[cellId].output
    }

    let VisualizerData = {
        CellsData: CellsData,
        TotalCells: TotalCells,
        ConsoleOutput: ConsoleOutput
    }

    document.onkeydown = checkKey

    function checkKey(event){
        if(event.repeat){
            return  
        } 
        if(event.key === "Enter" && event.ctrlKey){
            submitCode(event)
        }
        if(event.target.tagName === "TEXTAREA"){
            if(event.key === "ArrowDown" && event.ctrlKey){
                newCell(event, parseInt(event.target.id)+1)
            }
            if(event.key === "ArrowUp" && event.ctrlKey){
                newCell(event, parseInt(event.target.id))
            }
            if(event.key.toLowerCase() === "d" && event.ctrlKey && event.altKey){
                if(parseInt(event.target.id) !== 0){//Won't delete data cell
                    deleteCell(event, parseInt(event.target.id))
                }
            }
        }
        else{
            if(event.key === "ArrowDown" && event.ctrlKey){
                newCell(event, VisualizerData.TotalCells)
            }
            if(event.key === "ArrowUp" && event.ctrlKey){
                newCell(event, 1)
            }
        }
    }

   

    return(
        <ContextData.Provider value={VisualizerData}>
            <ContextSubmit.Provider value={submitCode}>
                <ContextCleanCode.Provider value={cleanCode}>
                    <ContextNewCell.Provider value={newCell}>
                            <ContextDeleteCell.Provider value={deleteCell}>
                                <ContextCopyToClipBoard.Provider value={copyToClipBoard}>
                                    <ContextUpdateData.Provider value={updateCodeData}>
                                        <ContextXMMRegisters.Provider value={getXMMRegisters}>
                                            {children}
                                        </ContextXMMRegisters.Provider>
                                    </ContextUpdateData.Provider>
                                </ContextCopyToClipBoard.Provider>
                            </ContextDeleteCell.Provider>
                    </ContextNewCell.Provider>
                </ContextCleanCode.Provider>
            </ContextSubmit.Provider>
        </ContextData.Provider>
    )
}

export default Provider;


