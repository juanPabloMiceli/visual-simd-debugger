import React, {} from 'react';
import TextInput from './TextInput';
import { useContextData} from './Context'
import XMMS from './XMMS'

export default function Cell(props){

    const VisualizerData = useContextData()
    
    if(VisualizerData.CellsData[props.id].output.length > 0){
        return(
            <div id={'Cell'}>
                <TextInput id={props.id}/>
                <XMMS data={VisualizerData.CellsData[props.id].output}/>
            </div>
        ) 
    }
    return(
        <div id={'Cell'}>
            <TextInput id={props.id}/>
        </div>
    ) 
}
