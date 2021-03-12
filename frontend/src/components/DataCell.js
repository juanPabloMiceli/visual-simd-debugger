import React, {} from 'react';
import TextInput from './TextInput';


export default function Cell(props){



    return(
        <div id={'Cell'}>
            <p id={"section-text"}><span id={"Section"}>Section</span> .data</p>
            <TextInput id={props.cellNumber} />
        </div>
    )
}
