import React from 'react';
import Big from 'big.js'

function numberToString(number, bits, base, symbol){
    let digits = bits
    if(base === 16){
        digits = bits/4
    }
    let stringRes = symbol

    let bigNumber = Big(number) 

    if(bigNumber.lt(0)){
        let exponente = Big(2)
        exponente = exponente.pow(bits)
        bigNumber = bigNumber.plus(exponente)

    }

    let rawNumber = ""
    let counter = 1
    
    while(bigNumber.gte(base)){
        counter++
        let modulo = bigNumber.mod(base).toString()
        rawNumber += parseInt(modulo, 10).toString(base).toUpperCase()
        bigNumber = bigNumber.div(base).round()
    }

    rawNumber += parseInt(bigNumber.toString(), 10).toString(base).toUpperCase()

    while(counter < digits){
        stringRes += "0"
        counter++
    }

    rawNumber = rawNumber.split("").reverse().join("")
    stringRes += rawNumber

    stringRes = stringRes.split("").reverse().join("")

    counter = 4

    while(counter < stringRes.length-2){
        stringRes = [stringRes.slice(0, counter), " ", stringRes.slice(counter)].join("")
        counter += 5
    }

    stringRes = stringRes.split("").reverse().join("")


    return stringRes
}


export default function XMM(props){
    Big.RM = 0
    let len = props.data.length
    let bits = 128/len
    let xmmID = "XMMVal-"+len
    let XMMReg    

    if(!props.printFormat || props.printFormat === "/d" || props.printFormat === "/u"){
        XMMReg = props.data.map(value =>
            <div id={xmmID}>{value}</div>
        )
    }else if(props.printFormat === "/x"){
        XMMReg = props.data.map(value =>{
            return (<div id={xmmID}>{numberToString(value, bits, 16, "0x")}</div>)
        })
    }else if(props.printFormat === "/t"){

        XMMReg = props.data.map(value =>{
            return (<div id={xmmID}>{numberToString(value, bits, 2, "0b")}</div>)
        })
    }
    


    return (
        <div>
            <div id="XMMName">
                {props.name}
            </div>
            {XMMReg}
        </div>
    );
}
