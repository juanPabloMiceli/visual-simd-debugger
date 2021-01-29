import React, { useEffect, useRef } from 'react';
import autosize from 'autosize'
import { useContextData } from "./Context" 


function TextOutput(){

	const VisualizerData = useContextData()
	const textArea = useRef(null)

	useEffect(
		() => {
			autosize(textArea.current)
		}, []
	);

	useEffect(
		() => {
			autosize.update(textArea.current)
		}, [VisualizerData.ConsoleOutput]
	);

	function updateOutput(){
		autosize.update(textArea.current)
	}

	const style = {
		minHeight: '38px',
		resize: 'none',
		padding: '9px',
		boxSizing: 'border-box',
		fontSize: '15px'
	}

	return (
		<div className="outputContainer">
			<textarea
			style={style}
			ref={textArea}
			rows={1} 
			id="output_text" 
			value={VisualizerData.ConsoleOutput}
			onChange={updateOutput}
			placeholder={"Console output goes here"}
			/>
		</div>
	)	
}

export default TextOutput;
