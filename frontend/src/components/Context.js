import axios from 'axios'
import React, {useContext, useEffect, useState, useCallback} from 'react'
import * as HotKeys from './HotKeys'

const ContextData  = React.createContext()
const ContextUpdateData = React.createContext()
const ContextXMMRegisters = React.createContext()
const ContextSubmit = React.createContext()
const ContextNewCell = React.createContext()
const ContextNewCellDown = React.createContext()
const ContextDeleteCell = React.createContext()
const ContextCopyToClipBoard = React.createContext()
const ContextToggleCellRegs = React.createContext()
const ContextToggleHelp = React.createContext()
const ContextCleanCode = React.createContext()



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

export function useContextToggleCellRegs(){
    return useContext(ContextToggleCellRegs)
}

export function useContextToggleHelp(){
    return useContext(ContextToggleHelp)
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
    
    const [CellsData, setCellsData] = useState([initCell(0)]);
    const [TotalCells, setTotalCells] = useState(1);
    const [ConsoleOutput, setConsoleOutput] = useState("");
    const [HelpActive, setHelpActive] = useState(false)

    function initCell(_id){
        return {
            id: _id,
            code: "",
            output: [],
            isVisible: true
        }
    }

    const initExample = useCallback(() => {
        let cellsData = []

        for(let i = 0; i < 8; i++){
            cellsData.push(initCell(i))
        }

        cellsData[0].code = ";Este programa se encarga de pasar los 2 primeros pixeles en memoria a escala de grises.\n"
        cellsData[0].code += ";Esto se hace mediante la formula: gris = (rojo + 2 * verde + azul)/4.\n"
        cellsData[0].code += ";Una vez se obtienen los 2 valores se les hace un broadcast de tal modo que "
        cellsData[0].code += "la parte baja del registro este formada por el primer pixel gris repetido 2 veces y "
        cellsData[0].code += "la parte alta contenga el segundo pixel gris.\n\n"
        cellsData[0].code += "pixeles: db 173, 68, 144, 255, 16, 54, 231, 255, 29, 178, 50, 255, 79, 211, 203, 255\n"
        cellsData[0].code += "mascara: db 0, 0, 0, 2, 0, 0, 0, 2, 1, 1, 1, 2, 1, 1, 1, 2"
    
        cellsData[1].code = ";Cargo la mascara en xmm0 y los 2 primeros pixeles en xmm1 e imprimo los valores.\n\n"
        cellsData[1].code += "movdqu xmm0, [mascara]\n"
        cellsData[1].code += "pmovzxbw xmm1, [pixeles]\n\n"
        cellsData[1].code += ";p/u xmm0.v16_int8\n"
        cellsData[1].code += ";p/u xmm1.v8_int16"

        cellsData[2].code = ";Reemplazo el valor alfa de cada pixel con el valor verde del mismo. Hacemos esto para ambos pixeles.\n"
        cellsData[2].code += ";Esto nos permite conseguir el valor de gris haciendo sumas horizontales.\n\n"
        cellsData[2].code += "pshufhw xmm1, xmm1, 0b01100100\n"
        cellsData[2].code += "pshuflw xmm1, xmm1, 0b01100100"
    
        cellsData[3].code = ";Hacemos la primer suma horizontal.\n"
        cellsData[3].code += ";Ahora tenemos en los 2 registros de la derecha el los valores de R + 2G + B de cada pixel.\n\n"
        cellsData[3].code += "phaddw xmm1, xmm1\n"
        cellsData[3].code += "phaddw xmm1, xmm1"

        cellsData[4].code = ";Shifteamos a derecha para dividir los resultados por 4 y asi obtener el valor gris del pixel.\n\n"
        cellsData[4].code += "psrlw xmm1, 2"

        cellsData[5].code = ";Volvemos a almacenar los pixeles conseguidos como enteros de 8 bits. En este punto vuelvo a imprimir en 8 bits.\n\n"
        cellsData[5].code += "packuswb xmm1, xmm1\n"
        cellsData[5].code += ";p xmm1.v16_int8"

        cellsData[6].code = ";Inserto un 255 en el valor de alfa para restaurar el valor original.\n\n"
        cellsData[6].code += "xor rax, rax\n"
        cellsData[6].code += "dec rax\n"
        cellsData[6].code += "pinsrb xmm1, al, 2"

        cellsData[7].code = ";Uso la mascara en xmm0 para hacer un broadcast de los valores conseguidos y del alfa "
        cellsData[7].code += "como decia el enunciado.\n\n"
        cellsData[7].code += "pshufb xmm1, xmm0"

        return cellsData
      }, [])

    useEffect(() =>{
        console.log(localStorage.getItem('visualizerData'))
        if(localStorage.getItem('cellsData') && localStorage.getItem('totalCells') && localStorage.getItem('consoleOutput')){
            setCellsData(JSON.parse(localStorage.getItem('cellsData')))
            setTotalCells(JSON.parse(localStorage.getItem('totalCells')))
            setConsoleOutput(JSON.parse(localStorage.getItem('consoleOutput')))
        }else{
            let auxData = initExample()
            setCellsData(auxData)
            localStorage.setItem('cellsData', JSON.stringify(auxData))
            setTotalCells(auxData.length)
            localStorage.setItem('totalCells', JSON.stringify(auxData.length))
            setConsoleOutput("")
            localStorage.setItem('consoleOutput', JSON.stringify(""))

        }
        
    },[initExample])
    
    


    function getRequestObj(cellsData){
        let res = []
        cellsData.forEach(cellData => {
            res.push({
                id: cellData.id,
                code: cellData.code
            })
        });

        return res
    }


    function submitCode(){
        let url = new URL(window.location.href)
        url.port = "80"

        let requestObj = getRequestObj(CellsData)
        

        setCellsData(cleanOutput(CellsData))
        localStorage.setItem('cellsData', JSON.stringify(cleanOutput(CellsData)))
        axios.post(url + "codeSave", JSON.stringify({
            CellsData: requestObj
        }))
        .then(response =>{
            console.log(response)
            setConsoleOutput(response.data.ConsoleOut) 
            localStorage.setItem('consoleOutput', JSON.stringify(response.data.ConsoleOut))
            if(response.data.CellRegs){
                updateXMMData(response.data.CellRegs) 
            }
        }) 
        .catch(error => {
            setConsoleOutput(error.message) 
            localStorage.setItem('consoleOutput', JSON.stringify(error.message))
        })
    }

    function fixIndexing(cells){
        let copy = JSON.parse(JSON.stringify(cells))
        for(let i=0; i<cells.length; i++){
            copy[i].id = i
        }

        return copy
    }
    function newCell(e, buttonNumber){
        e.preventDefault()
        if(buttonNumber === 0) return;
        let copy = JSON.parse(JSON.stringify(CellsData))//This is the only way of doing a depth copy
        // setCellsData([initCell(0)])
        
        copy.splice(buttonNumber, 0, initCell(TotalCells))

        copy = fixIndexing(copy)
        setCellsData(copy)
        localStorage.setItem('cellsData', JSON.stringify(copy))

        setTotalCells(TotalCells + 1)
        localStorage.setItem('totalCells', JSON.stringify(TotalCells+1))

        let ca = document.getElementById("code"+buttonNumber.toString())
        if(ca){
            ca.focus()
        }
    }

    function focusTextElement(buttonNumber){
        let ca = document.getElementById("code"+buttonNumber.toString())
        while(!ca && buttonNumber > 0){
            buttonNumber--
            ca = document.getElementById("code"+buttonNumber.toString())
        }
        if(ca){
            ca.focus()
        }
    }

    function deleteCell(e, buttonNumber){
        e.preventDefault()
        
        if(TotalCells > 2 && buttonNumber > 0){
            let copy = JSON.parse(JSON.stringify(CellsData))//This is the only way of doing a depth copy
            copy.splice(buttonNumber, 1)
            copy = fixIndexing(copy)
            setCellsData(copy)
            localStorage.setItem('cellsData', JSON.stringify(copy))
            setTotalCells(TotalCells - 1)
            localStorage.setItem('totalCells', JSON.stringify(TotalCells-1))
        }
        focusTextElement(buttonNumber)
    }

    function toggleHelp(e){
        e.preventDefault()
        setHelpActive(!HelpActive)
    }

    function toggleCellRegs(e, buttonNumber){
        e.preventDefault()
        let copy = JSON.parse(JSON.stringify(CellsData))//This is the only way of doing a depth copy
        copy[buttonNumber].isVisible = !copy[buttonNumber].isVisible
        setCellsData(copy)
        localStorage.setItem('cellsData', JSON.stringify(copy))
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
            if(CellsData[i].isVisible && CellsData[i].code !== ""){
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
            localStorage.setItem('cellsData', JSON.stringify(copy))
            setTotalCells(2)
            localStorage.setItem('totalCells', JSON.stringify(2))
        }
        
    }
    function updateXMMData(data){
        let copy = JSON.parse(JSON.stringify(CellsData))
        for(let i = 0; i < copy.length; i++){
            copy[i].output = data[i]
        }

        setCellsData(copy)
        localStorage.setItem('cellsData', JSON.stringify(copy))

    }

    function updateCodeData(cellId, newCode){
        let copy = JSON.parse(JSON.stringify(CellsData))
        copy[cellId].code = newCode
        setCellsData(copy)
        localStorage.setItem('cellsData', JSON.stringify(copy))
    }

    function getXMMRegisters(cellId){
        console.log("GetXMMRegisters")
        console.log(CellsData[cellId].output)
        return CellsData[cellId].output
    }

    let VisualizerData = {
        CellsData: CellsData,
        TotalCells: TotalCells,
        ConsoleOutput: ConsoleOutput,
        HelpActive: HelpActive
    }

    document.onkeydown = checkKey

    function checkKey(event){
        if(event.repeat) return  
        if(event.ctrlKey && event.altKey && event.key.toLowerCase() === "h"){
            toggleHelp(event)
        }
        if(HelpActive) return
        let number;
        if(event.target.id){
            number =  parseInt(event.target.id.replace( /^\D+/g, ''));
        }
        if(HotKeys.submitCode(event)){//SubmitCode
            submitCode(event)
        }
        if(event.target.tagName === "TEXTAREA"){
            if(HotKeys.newCellBelow(event)){//NewCellDown
                newCell(event, number+1)
            }
            if(HotKeys.newCellAbove(event)){//NewCellUp
                newCell(event, number)
            }
            if(HotKeys.deleteCell(event)){//DeleteCurrentCell
                if(number !== "0"){//Won't delete data cell
                    deleteCell(event, number)
                }
            }
            if(HotKeys.moveUp(event)){//SelectCellAbove
                focusTextElement(number-1)
            }
            if(HotKeys.moveDown(event)){//SelectCellBelow
                focusTextElement(number+1)
            }
        }
        else{
            if(HotKeys.newCellBelow(event)){//NewCellAtEnd
                newCell(event, VisualizerData.TotalCells)
            }
            if(HotKeys.newCellAbove(event)){//NewCellAtBeginning
                newCell(event, 1)
            }
            if(HotKeys.moveUp(event)){//SelectLastCell
                focusTextElement(VisualizerData.TotalCells-1)
            }
            if(HotKeys.moveDown(event)){//SelectFirstCell
                focusTextElement(0)
            }
        }
    }

    //Saves current code state before window is closed
    // window.onbeforeunload = function(){
    //     localStorage.setItem('visualizerData', JSON.stringify(VisualizerData))
    // }

   

    return(
        <ContextData.Provider value={VisualizerData}>
            <ContextSubmit.Provider value={submitCode}>
                <ContextCleanCode.Provider value={cleanCode}>
                    <ContextNewCell.Provider value={newCell}>
                            <ContextDeleteCell.Provider value={deleteCell}>
                                <ContextCopyToClipBoard.Provider value={copyToClipBoard}>
                                    <ContextToggleCellRegs.Provider value={toggleCellRegs}>
                                        <ContextToggleHelp.Provider value={toggleHelp}>
                                            <ContextUpdateData.Provider value={updateCodeData}>
                                                <ContextXMMRegisters.Provider value={getXMMRegisters}>
                                                    {children}
                                                </ContextXMMRegisters.Provider>
                                            </ContextUpdateData.Provider>
                                        </ContextToggleHelp.Provider>
                                    </ContextToggleCellRegs.Provider>
                                </ContextCopyToClipBoard.Provider>
                            </ContextDeleteCell.Provider>
                    </ContextNewCell.Provider>
                </ContextCleanCode.Provider>
            </ContextSubmit.Provider>
        </ContextData.Provider>
    )
}

export default Provider;


