import React from 'react';
import Cell from './Cell'
import { useContextData, useContextSubmit, useContextCopyToClipBoard } from "./Context" 
import TextOutput from './TextOutput';
import FileInput from './FileInput'


function SimdVisualizer(){

    const submitCode = useContextSubmit() 
    const copyToClipBoard = useContextCopyToClipBoard()
    const VisualizerData = useContextData()

    console.log(VisualizerData)

    let CellComponents = VisualizerData.CellsData.map(cell => {
        return (<Cell
                    key={cell.id}
                    id={cell.id}/>)
    })


    return (
        <div>
            <div className="submitContainer">
                <button className="btn btn-success" id="submitButton" onClick={submitCode}>Submit Code</button>
            </div>
            {CellComponents}
            <div className="copyContainer">
                <button className="btn btn-DelCell" id="copyToClipButton" onClick={copyToClipBoard}>Copy code to clipboard</button>
            </div>
            <TextOutput/>
            <FileInput/>
        </div>
    );
}

export default SimdVisualizer;
