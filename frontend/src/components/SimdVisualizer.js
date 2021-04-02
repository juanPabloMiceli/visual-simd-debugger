import React from 'react';
import Cell from './Cell'
import { useContextData, useContextSubmit, useContextCopyToClipBoard, useContextCleanCode } from "./Context" 
import TextOutput from './TextOutput';
import DataCell from './DataCell'

import GitHubCorners from '@uiw/react-github-corners';


function SimdVisualizer(){

    const submitCode = useContextSubmit() 
    const copyToClipBoard = useContextCopyToClipBoard()
    const VisualizerData = useContextData()
    const cleanCode = useContextCleanCode()


    console.log(VisualizerData)

    let CellComponents = VisualizerData.CellsData.map(cell => {
        if(cell.id === 0){
            return(null)
        }
        return (<Cell
                    key={cell.id}
                    id={"cell"+cell.id.toString()}
                    cellNumber={cell.id}/>)
    })


    return (
        
            <div>
                <GitHubCorners
                    position="right"
                    href="https://gitlab.com/juampi_miceli/visual-simd-debugger"
                    target="_blank"
                />
                <div className="cleanCodeContainer">
                    <button className="btn-clean-code" onClick={cleanCode}>Clean Code</button>
                </div>
                <DataCell cellNumber={0} id={"dataCell"}/>
                {CellComponents}
                <div className="copyContainer">
                    <button className="btn btn-Copy" id="copyToClipButton" onClick={copyToClipBoard}>Copy code to clipboard</button>
                    <button className="btn-submit" id="submitButton" onClick={submitCode}>Run Code</button>
                </div>
                <TextOutput/>
            </div>
        
    );
}



export default SimdVisualizer;
