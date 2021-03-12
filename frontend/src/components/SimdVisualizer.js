import React from 'react';
import Cell from './Cell'
import { useContextData, useContextSubmit, useContextCopyToClipBoard, useContextCleanCode } from "./Context" 
import TextOutput from './TextOutput';
import DataCell from './DataCell'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"
import { faGithub } from '@fortawesome/free-brands-svg-icons';




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

    function openGit(e){
        e.preventDefault()
        window.open('https://gitlab.com/juampi_miceli/visual-simd-debugger', '_blank')
    }


    return (
        
            <div>
                <div className="submitContainer">
                    <button className="btn-submit" id="submitButton" onClick={submitCode}>Run Code</button>
                    <button className="btn-git" onClick={openGit}><FontAwesomeIcon icon={faGithub}/></button>
                    <button className="btn-clean-code" onClick={cleanCode}>Clean Code</button>
                </div>
                <DataCell cellNumber={0} id={"dataCell"}/>
                {CellComponents}
                <div className="copyContainer">
                    <button className="btn btn-DelCell" id="copyToClipButton" onClick={copyToClipBoard}>Copy code to clipboard</button>
                </div>
                <TextOutput/>
            </div>
        
    );
}



export default SimdVisualizer;
