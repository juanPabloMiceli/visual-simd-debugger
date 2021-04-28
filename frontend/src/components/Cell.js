import React, {} from 'react';
import TextInput from './TextInput';
import { useContextData, useContextDeleteCell, useContextNewCell, useContextToggleCellRegs} from './Context'
import XMMS from './XMMS'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"


export default function Cell(props){

    const VisualizerData = useContextData()
    const newCell = useContextNewCell()
    const deleteCell = useContextDeleteCell()
    const toggleCellRegs = useContextToggleCellRegs()


    function isFirstCell(id){
        if(id ===  1){
            return(
            <div className="section-new-cell-container">
                <div className="newCellContainer">
                    <span><button className="btn-newCell" onClick={e => newCell(e, props.cellNumber)}>+ Code</button></span>
                </div>
                <p className={"section-text"}><span id={"Section"}>Section</span> .text</p>
            </div>
            )
        }else{
            <div></div>
        }
    }

    function isLastCell(id, totalCells) {
        if(id === totalCells-1){
            return(
                <div>
                    <div className="section-last-new-cell-container">
                        <span><button className="newCellContainer" onClick={e => newCell(e, props.cellNumber+1)}>+ Code</button></span>
                    </div>   
                </div>
            )
        }else{
            return(
                <div className="section-new-cell-container">
                        <span><button className="newCellContainer" onClick={e => newCell(e, props.cellNumber+1)}>+ Code</button></span>
                </div>
            )
            
        }
    }

    function hasXMMData(dataLength) {
        if(dataLength > 0 && VisualizerData.CellsData[props.cellNumber].isVisible){
            return(
                <XMMS data={VisualizerData.CellsData[props.cellNumber].output}/>
            )
        }else{
            <div></div>
        }
    }

    function getEyeIcon(iconState){
        if(iconState) return(<FontAwesomeIcon icon="eye" fixedWidth/>)
        return(<FontAwesomeIcon icon="eye-slash" fixedWidth/>)
    }

    return(
        <div id={'Cell'}>
            {isFirstCell(props.cellNumber)}
            <div className="buttons">
                <div className="showCellRegs">
                    <button className="btn btn-DelCell" id="showCellRegsButton" onClick={e => toggleCellRegs(e, props.cellNumber)}>{getEyeIcon(VisualizerData.CellsData[props.cellNumber].isVisible)}</button>
                </div>
                <div className="delCell">
                    <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.cellNumber)}><FontAwesomeIcon icon="trash-alt" fixedWidth/></button>
                </div>
            </div>
            <TextInput id={props.cellNumber} />
            {hasXMMData(VisualizerData.CellsData[props.cellNumber].output.length)}
            {isLastCell(props.cellNumber, VisualizerData.TotalCells)}           
        </div>
    )
}
