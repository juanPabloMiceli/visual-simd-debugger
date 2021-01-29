import axios from 'axios'
import React, {useContext, useState} from 'react'

const ContextData  = React.createContext()
const ContextUpdateData = React.createContext()
const ContextXMMRegisters = React.createContext()
const ContextSubmit = React.createContext()
const ContextNewCell = React.createContext()
const ContextDeleteCell = React.createContext()
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

export function useContextDeleteCell(){
    return useContext(ContextDeleteCell)
}


axios.defaults.headers.common = {
	"Content-Type": "application/json"
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
        console.log(CellsData)
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

    function newCell(e){
        e.preventDefault()
        
        setCellsData(CellsData.concat(initCell(TotalCells)))
        setTotalCells(TotalCells + 1)
    }

    function deleteCell(e){
        e.preventDefault()

        if(TotalCells > 1){
            let copy = JSON.parse(JSON.stringify(CellsData))
            console.log(copy)
            copy.pop()
            setCellsData(copy)
            setTotalCells(TotalCells - 1)
        }
    }

    function updateXMMData(data){
        let copy = JSON.parse(JSON.stringify(CellsData))
        for(let i = 0; i < copy.length; i++){
            copy[i].output = data[i]
        }

        setCellsData(copy)

    }

    function updateCodeData(cellId, newCode){
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
                        <ContextUpdateData.Provider value={updateCodeData}>
                            <ContextXMMRegisters.Provider value={getXMMRegisters}>
                                {children}
                            </ContextXMMRegisters.Provider>
                        </ContextUpdateData.Provider>
                    </ContextDeleteCell.Provider>
                </ContextNewCell.Provider>
            </ContextSubmit.Provider>
        </ContextData.Provider>
    )
}

export default Provider;