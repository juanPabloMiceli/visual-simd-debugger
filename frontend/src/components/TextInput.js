import React, { useEffect, useRef } from 'react'
import autosize from 'autosize'
import {useContextData, useContextUpdateData} from './Context'

export default function TextInput(props){

	const textArea = useRef(null)

	const VisualizerData = useContextData()
	const updateData = useContextUpdateData()

	useEffect(() => {
		textArea.current.focus()
		autosize(textArea.current)
	}, []);

	
	function updateCode(e){
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
			id="code"
			onChange={updateCode}
			placeholder={"Code goes here"}
			/>		
		</div>
	)
}


