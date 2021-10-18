import React from 'react';
import Cell from './Cell'
import { useContextData, useContextSubmit, useContextCopyToClipBoard, useContextCleanCode, useContextToggleHelp } from "./Context" 
import TextOutput from './TextOutput';
import DataCell from './DataCell'

import GitHubCorners from '@uiw/react-github-corners';
import HelpPopUp from './HelpPopUp';
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"


function SimdVisualizer(){

    const submitCode = useContextSubmit() 
    const copyToClipBoard = useContextCopyToClipBoard()
    const VisualizerData = useContextData()
    const cleanCode = useContextCleanCode()
    const toggleHelp = useContextToggleHelp()


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
                    href="https://github.com/juanPabloMiceli/visual-simd-debugger/"
                    target="_blank"
                />
                <div className="topButtonsContainer">
                    <button className="btn-clean-code" onClick={cleanCode}>Clean Code</button>
                    <button className="btn-help" onClick={toggleHelp}><FontAwesomeIcon icon="info" fixedWidth/></button>
                    {/* <button className="btn-open-code" onClick={cleanCode}>Open Code</button>
                    <button className="btn-new-code" onClick={cleanCode}>New Code</button> */}
                </div>
                <DataCell cellNumber={0} id={"dataCell"}/>
                {CellComponents}
                <div className="copyContainer">
                    <button className="btn btn-Copy" id="copyToClipButton" onClick={copyToClipBoard}>Copy code to clipboard</button>
                    <button className="btn-submit" id="submitButton" onClick={submitCode}>Run Code</button>
                </div>
                <TextOutput/>
                <HelpPopUp/>
            </div>
        
    );
}



export default SimdVisualizer;
