import React, {useEffect} from 'react';
import TextInput from './TextInput';
import {useContextXMMRegisters, useContextData} from './Context'
import XMMS from './XMMS'

export default function Cell(props){

    const VisualizerData = useContextData()

    // let outputData

    // useEffect(() => {
    //     console.log("XMMS: " + props.id)
    //     console.log(VisualizerData.CellsData[props.id].output)
    //     outputData = <XMMS data={VisualizerData.CellsData[props.id].output}/>
    // }, [VisualizerData.CellsData[props.id].output]);
    
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


// class Cell extends Component {

// 	// updateCodeFromChild = (_code) => {
// 	// 	this.props.parentUpdateCode(this.props.id, _code)
// 	// }


//     render() {
//         return (
//             <div id={'Cell'}>
//                 hols
//                 <TextInput parentUpdateCode={this.updateCodeFromChild}/>
//                 {/* <TextOutput text={this.props.output}/> */}
//             </div>
//         );
//     }
// }

// export default Cell;
