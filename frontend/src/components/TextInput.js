import React from 'react'
import Editor from 'react-simple-code-editor'
import {useContextData, useContextUpdateData} from './Context'
import "highlight.js/lib/languages/x86asm"

export default function TextInput(props){

	const VisualizerData = useContextData()
	const updateData = useContextUpdateData()
	
	const hljs = require('highlight.js')
	hljs.registerLanguage('x86asm', require('highlight.js/lib/languages/x86asm'))

	return(
		<Editor
			className={"code"}
			textareaClassName={"code"}
			textareaId={"code"+props.id}
			tabSize={4}
			value={VisualizerData.CellsData[props.id].code}
			onValueChange={(e) => updateData(props.id, e)}
			highlight={(e) => hljs.highlight('x86asm', e, true).value}
			padding={9}
			/>
	)
}
