import React, { Component } from 'react'
import axios from 'axios'

let local_host = "http://localhost:8080"

axios.defaults.headers.common = {
	"Content-Type": "application/json"
  }

export default class TextInput extends Component {

	state={
		code: ""
	}

	onSubmitHandler = (e) => {
		e.preventDefault()
		let data =  {CodeText: this.state.code}
		// alert("Submit code")	
		axios.post(local_host + "/codeSave", JSON.stringify(data))
		.then(response => {
			console.log("Malas malas malas")
			console.log(response)
		})
		.catch(error => {
			console.log(error)
			console.log("Buenas buenas buenas")
		})
	}

	OnCodeChangeHandler = (e) => {
		this.setState({
			code: e.target.value
		})
	}

	render() {
		return (
			<div className="container">
				<div className="row">
					<div className="col-12 mt-5">
						<p className="lead d-block my-0" id="writeBelow"> Write below</p>
						<textarea type="text" id="code" value={this.state.code} onChange={this.OnCodeChangeHandler}>
						</textarea>
					</div>
				</div>
				<button className="btn btn-success" id="submitButton" onClick={this.onSubmitHandler}>Submit Code</button>
			</div>
		)

	}
}