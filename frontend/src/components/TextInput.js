import React, { Component } from 'react'
import axios from 'axios'
import autosize from 'autosize'
import TextOutput from './TextOutput'

let local_host = "http://localhost:8080"

axios.defaults.headers.common = {
	"Content-Type": "application/json"
  }

export default class TextInput extends Component {

	componentDidMount(){
		this.textarea.focus()
		autosize(this.textarea)
	}

	constructor(){
		super()
		this.state={
			code: ""
		}
	}
	

	OnCodeChangeHandler = (e) => {
		this.setState({
			code: e.target.value
		})
		this.props.parentUpdateCode(e.target.value)
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
			<div className="inputContainer">
				<textarea
				style={style}
				ref={c => this.textarea = c}
				rows={1}
				id="code"
				value={this.state.code}
				onChange={this.OnCodeChangeHandler}
				placeholder={"Code goes here"}
				/>		
			</div>
		)

	}
}


