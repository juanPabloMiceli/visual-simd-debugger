export function submitCode(event){
    return event.key === "Enter" && event.ctrlKey
}

export function newCellAbove(event){
    return event.key === "ArrowUp" && event.ctrlKey
}

export function newCellBelow(event){
    return event.key === "ArrowDown" && event.ctrlKey
}

export function deleteCell(event){
    return event.key.toLowerCase() === "d" && event.ctrlKey && event.altKey
}


export function moveUp(event){
    return (event.key === "ArrowUp" && event.altKey) || event.keyCode === 33
}

export function moveDown(event){
    return (event.key === "ArrowDown" && event.altKey) || event.keyCode === 34
}