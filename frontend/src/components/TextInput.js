import React, { useEffect, useRef } from 'react'
import {useContextData, useContextUpdateData} from './Context'

export default function TextInput(props){


	const textArea = useRef(null)

	const VisualizerData = useContextData()
	const updateData = useContextUpdateData()
	let currentCode = VisualizerData.CellsData[props.id].code




	useEffect(() => {
		setInputHeight(textArea.current, '38px')
	}, []);

	useEffect(() => {
		setInputHeight(textArea.current, '38px')
	}, [currentCode]);

	
	function setInputHeight(element, defaultHeight){
		if(element){
			const target = element.target ? element.target : element
			target.style.height = defaultHeight
			let height = parseInt(target.scrollHeight)+2
			target.style.height = height.toString()+'px'
		}
	}
	

	function updateCode(e){
		setInputHeight(e, '38px')
		updateData(props.id, e.target.value)
	}

	const style = {
		minHeight: '38px',
		resize: 'none',
		padding: '9px',
		boxSizing: 'border-box',
		fontSize: '15px'
	}

	return (
		<div className="inputContainer">
			<textarea
			style={style}
			ref={textArea}
			rows={1}
			id={props.id}
			className={"code"}
			onChange={updateCode}
			value={VisualizerData.CellsData[props.id].code}
			placeholder={"Code goes here"}
			/>		
		</div>
	)
}





// import React, { useEffect, useRef } from 'react'
// import autosize from 'autosize'
// import {useContextData, useContextUpdateData, useState} from './Context'

// export default function TextInput(props){


// 	const textArea = useRef(null)

// 	const VisualizerData = useContextData()
// 	const updateData = useContextUpdateData()
// 	// const [code, setcode] = useState(VisualizerData.CellsData[props.id].code);



// 	useEffect(() => {
		
// 		textArea.current.focus()		
// 		autosize(textArea.current)
// 	}, []);

// 	useEffect(() => {
		
// 		textArea.current.focus()	
// 		autosize(textArea.current)
// 	}, [VisualizerData]);
	
// 	function updateCode(e){
// 		console.log("Changing")

// 		updateData(props.id, e.target.value, textArea)
// 	}

// 	const style = {
// 		minHeight: '38px',
// 		resize: 'none',
// 		padding: '9px',
// 		boxSizing: 'border-box',
// 		fontSize: '15px'
// 	}

// 	return (
// 		<div className="inputContainer">
// 			<textarea
// 			style={style}
// 			ref={textArea}
// 			rows={1}
// 			id="code"
// 			onChange={updateCode}
// 			value={VisualizerData.CellsData[props.id].code}
// 			placeholder={"Code goes here"}
// 			/>		
// 		</div>
// 	)
// }


