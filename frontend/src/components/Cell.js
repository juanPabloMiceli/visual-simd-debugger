import React, {} from 'react';
import TextInput from './TextInput';
// import TextInputV2 from './TextInputV2';
import { useContextData, useContextDeleteCell, useContextNewCell} from './Context'
import XMMS from './XMMS'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"


export default function Cell(props){

    const VisualizerData = useContextData()
    const newCell = useContextNewCell()
    const deleteCell = useContextDeleteCell()


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
        if(dataLength > 0){
            return(
                <XMMS data={VisualizerData.CellsData[props.cellNumber].output}/>
            )
        }else{
            <div></div>
        }
    }

    return(
        <div id={'Cell'}>
            {isFirstCell(props.cellNumber)}
            <div className="delCell">
                <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.cellNumber)}><FontAwesomeIcon icon="trash-alt"/></button>
            </div>
            <TextInput id={props.cellNumber} />
            {hasXMMData(VisualizerData.CellsData[props.cellNumber].output.length)}
            {isLastCell(props.cellNumber, VisualizerData.TotalCells)}           
        </div>
    )
}
