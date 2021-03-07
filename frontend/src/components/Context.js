import axios from 'axios'
import autosize from 'autosize'
import React, {useContext, useState} from 'react'

const ContextData  = React.createContext()
const ContextUpdateData = React.createContext()
const ContextXMMRegisters = React.createContext()
const ContextSubmit = React.createContext()
const ContextNewCell = React.createContext()
const ContextNewCellDown = React.createContext()
const ContextDeleteCell = React.createContext()
const ContextCopyToClipBoard = React.createContext()
const local_host = "http://localhost:8080"

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
    
    const [CellsData, setCellsData] = useState([initCell(0)]);
    const [TotalCells, setTotalCells] = useState(1);
    const [ConsoleOutput, setConsoleOutput] = useState("");

    function initCell(_id){
        return {
            id: _id,
            code: "",
            output: []
        }
    }

    function submitCode(e){
        e.preventDefault()
        setCellsData(cleanOutput(CellsData))
        axios.post(local_host + "/codeSave", JSON.stringify({
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

export default Provider;