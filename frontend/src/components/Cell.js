import React, {} from 'react';
import TextInput from './TextInput';
import { useContextData, useContextDeleteCell, useContextNewCell} from './Context'
import XMMS from './XMMS'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"


export default function Cell(props){

    const VisualizerData = useContextData()
    const newCell = useContextNewCell()
    const deleteCell = useContextDeleteCell()


    function isFirstCell(id){
        if(id ===  0){
            return(
                <div className="newCellContainer">
                    <span><button className="btn-newCell" onClick={e => newCell(e, props.id)}>+ Code</button></span>
                </div>
            )
        }else{
            <div></div>
        }
    }

    function isLastCell(id, totalCells) {
        if(id === totalCells-1){
            return(
                <div className="newCellContainer">
                    <span><button className="btn-lastNewCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
                </div>   
            )
        }else{
            return(
                <div className="newCellContainer">
                    <span><button className="btn-newCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
                </div> 
            )
            
        }
    }

    function hasXMMData(dataLength) {
        if(dataLength > 0){
            return(
                <XMMS data={VisualizerData.CellsData[props.id].output}/>
            )
        }else{
            <div></div>
        }
    }

    return(
        <div id={'Cell'}>
            {isFirstCell(props.id)}
            <div className="delCell">
                <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
            </div>
            <TextInput id={props.id}/>
            {hasXMMData(VisualizerData.CellsData[props.id].output.length)}
            {isLastCell(props.id, VisualizerData.TotalCells)}           
        </div>
    )
    
    // if(VisualizerData.CellsData[props.id].output.length > 0){
    //     if(props.id === VisualizerData.TotalCells-1){
    //         return(
    //             <div id={'Cell'}>
    //                 <div className="delCell">
    //                     <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //                 </div>
    //                 <TextInput id={props.id}/>
    //                 <XMMS data={VisualizerData.CellsData[props.id].output}/>
    //                 <div className="newCellContainer">
    //                     <span><button className="btn-lastNewCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //                 </div>                
    //             </div>
    //         ) 
    //     }
    //     if(props.id === 0){
    //         return(
    //             <div id={'Cell'}>
    //                 <div className="newCellContainer">
    //                     <span><button className="btn-newCell" onClick={e => newCell(e, props.id)}>+ Code</button></span>
    //                 </div>
    //                 <div className="delCell">
    //                     <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //                 </div>
    //                 <TextInput id={props.id}/>
    //                 <XMMS data={VisualizerData.CellsData[props.id].output}/>
    //                 <div className="newCellContainer">
    //                         <span><button className="btn-newCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //                 </div>                
    //             </div>
    //         ) 
    //     }
    //     return(
    //         <div id={'Cell'}>
    //             <div className="delCell">
    //                 <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //             </div>
    //             <TextInput id={props.id}/>
    //             <XMMS data={VisualizerData.CellsData[props.id].output}/>
    //             <div className="newCellContainer">
    //                 <span><button className="btn-newCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //             </div>                
    //         </div>
    //     ) 
    // }
    // if(props.id === 0){
    //     return(
    //         <div id={'Cell'}>
    //             <div className="newCellContainer">
    //                 <span><button className="btn-newCell" onClick={e => newCell(e, props.id)}>+ Code</button></span>
    //             </div>
    //             <div className="delCell">
    //                 <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //             </div>
    //             <TextInput id={props.id}/>
    //             <div className="newCellContainer">
    //                 <span><button className="btn-newCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //             </div>                
    //         </div>
    //     ) 
    // }
    // if(props.id === VisualizerData.TotalCells-1){
    //     return(
    //         <div id={'Cell'}>
    //             <div className="delCell">
    //                 <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //             </div>
    //             <TextInput id={props.id}/>
    //             <div className="newCellContainer">
    //                 <span><button className="btn-lastNewCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //             </div>                
    //         </div>
    //     ) 
    // }
    // return(
    //     <div id={'Cell'}>
    //         <div className="delCell">
    //             <button className="btn btn-DelCell" id="delCellButton" onClick={e => deleteCell(e, props.id)}><FontAwesomeIcon icon="trash-alt"/></button>
    //         </div>
    //         <TextInput id={props.id}/>
    //         <div className="newCellContainer">
    //                 <span><button className="btn-newCell" onClick={e => newCell(e, props.id+1)}>+ Code</button></span>
    //         </div>                
    //     </div>
    // ) 
}
