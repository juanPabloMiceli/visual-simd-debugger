import React from 'react';
import Cell from './Cell'
import { useContextData, useContextSubmit, useContextNewCell, useContextDeleteCell, useContextCopyToClipBoard } from "./Context" 
import TextOutput from './TextOutput';

function SimdVisualizer(){

    const submitCode = useContextSubmit() 
    const newCell = useContextNewCell()
    const deleteCell = useContextDeleteCell()
    const copyToClipBoard = useContextCopyToClipBoard()
    const VisualizerData = useContextData()

    console.log(VisualizerData)

    let CellComponents = VisualizerData.CellsData.map(cell => {
        return (<div>
                <Cell
                    key={cell.id}
                    id={cell.id}
                    />
                </div>)
    })


    return (
        <div>
            <button className="btn btn-success" id="submitButton" onClick={submitCode}>Submit Code</button>
            {CellComponents}
            <button className="btn btn-newCell" id="newCellButton" onClick={newCell}>New Cell</button>
            <button className="btn btn-DelCell" id="delCellButton" onClick={deleteCell}>Delete Cell</button>
            <button className="btn btn-DelCell" id="copyToClipButton" onClick={copyToClipBoard}>Copy code to clipboard</button>
            <TextOutput/>
        </div>
    );
}

export default SimdVisualizer;
