import React, {useEffect} from 'react';
import XMM from './XMM'


export default function XMMS(props){
    let XMMRegisters = props.data.map(XMMRegister => {
        return (
            <XMM name={XMMRegister.XmmID} data={XMMRegister.XmmValues}/>
        )
    })

    return (
        <div>
            {XMMRegisters}
        </div>
    );
}