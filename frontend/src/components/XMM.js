import React from 'react';


export default function XMM(props){

    let len = props.data.length
    let xmmID = "XMMVal-"+len
    
    let XMMReg = props.data.map(value =>
        <div id={xmmID}>{value}</div>
    )

    return (
        <div>
            <div id="XMMName">
                {props.name}
            </div>
            {XMMReg}
        </div>
    );
}
