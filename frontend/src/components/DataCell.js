import React, {} from 'react';
import TextInput from './TextInput';
// import TextInputV2 from './TextInputV2';


export default function Cell(props){



    return(
        <div id={'Cell'}>
            <p id={"section-text"}><span id={"Section"}>Section</span> .data</p>
            {/* <TextInputV2 id={props.cellNumber} /> */}
            <TextInput id={props.cellNumber} />
        </div>
    )
}
