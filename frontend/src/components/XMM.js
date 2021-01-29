import React from 'react';


export default function XMM(props){
    
    let XMMReg = props.data.map(value =>
        <div id="XMMVal">{value}</div>
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
