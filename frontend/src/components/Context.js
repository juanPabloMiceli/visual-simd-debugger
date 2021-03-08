import axios from 'axios'
import React, {useContext, useEffect, useState, useRef} from 'react'
import { HotKeys } from 'react-hotkeys'
import {useKeyboardShortcut} from './useKeyboardShortcut'

const ContextData  = React.createContext()
const ContextUpdateData = React.createContext()
const ContextXMMRegisters = React.createContext()
const ContextSubmit = React.createContext()
const ContextNewCell = React.createContext()
const ContextNewCellDown = React.createContext()
const ContextDeleteCell = React.createContext()
const ContextCopyToClipBoard = React.createContext()
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
    const [TotalCells, setTotalCells] = useState(3);
    const [ConsoleOutput, setConsoleOutput] = useState("");

    // useKey(['Ctrl', 'Enter'], () => {
    //     alert("apretaste")
    // })
    // useKeyboardShortcut(["Control", "Enter"], () => {
    //     submitCode()
    // })
    // useKeyboardShortcut(["Control", "ArrowDown"], (e) =>{
    //     console.log(e)
    // })

    function initCell(_id){
        return {
            id: _id,
            code: "",
            output: []
        }
    }

    function initExample() {

        let cellsData = [initCell(0), initCell(1), initCell(2)]
    
        cellsData[0].code = ";data\n"
        cellsData[0].code += "words: dw 1, 2, 3, 4, 5, 6, 7, 8"
    
        cellsData[1].code = "movdqu xmm0,[words]\n"
        cellsData[1].code += ";p/u xmm0.v8_int16\n"
        cellsData[1].code += ";p/x xmm0.v8_int16"
    
        cellsData[2].code = "psllq xmm0, 4"
    
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
            console.log(error)
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
        if(TotalCells > 1){
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
        console.log(event.target.tagName)
        if(event.key === "Enter" && event.ctrlKey){
            submitCode(event)
        }
        if(event.target.tagName == "TEXTAREA"){
            if(event.key === "ArrowDown" && event.ctrlKey){
                newCell(event, parseInt(event.target.id)+1)
            }
            if(event.key === "ArrowUp" && event.ctrlKey){
                newCell(event, parseInt(event.target.id))
            }
        }
        else{
            if(event.key === "ArrowDown" && event.ctrlKey){
                newCell(event, VisualizerData.TotalCells)
            }
            if(event.key === "ArrowUp" && event.ctrlKey){
                newCell(event, 0)
            }
        }
    }
    // window.addEventListener("onkeypress", (event) => {
        
    // })

    return(
        <ContextData.Provider value={VisualizerData}>
            <ContextSubmit.Provider value={submitCode}>
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
            </ContextSubmit.Provider>
        </ContextData.Provider>
    )
}

// function useKey(keys, callback){
//     const callbackRef = useRef(callback)

//     useEffect(() => {
//         callbackRef.current = callback
//     })

//     useEffect(() => {
//         function handle(event){
//             if(event.code === key){
//                 callbackRef.current(event)
//             }
//         }
//         document.addEventListener("keypress", handle)
//         return() => document.removeEventListener("keypress", handle)
//     }, [key])

    
// }

export default Provider;