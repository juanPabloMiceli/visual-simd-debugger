import React, { Component } from 'react';
import autosize from 'autosize'

class TextOutput extends Component {

	componentDidMount(){
		this.textarea.focus()
		autosize(this.textarea)
	}

	componentDidUpdate(){
		autosize.update(this.textarea)
	}

	constructor(props){
		super(props)
		this.state={
			id: 0
		}
	}
	

	render() {
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
				ref={c => this.textarea = c}
				rows={1} 
				id="output_text" 
				value={this.props.text}
				placeholder={"Results go here"}
				/>
			</div>
		)

	}
}

export default TextOutput;
