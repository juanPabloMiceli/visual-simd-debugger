import React, { Component } from 'react';
import TextInput from './TextInput';
import TextOutput from './TextOutput';
import axios from 'axios'
let local_host = "http://localhost:8080"

axios.defaults.headers.common = {
	"Content-Type": "application/json"
  }

class Cell extends Component {

	constructor(props){
		super(props)
		this.state={
			code: ''
		}

	}

	updateCodeFromChild = (_code) => {
		this.setState({code: _code})
		this.props.parentUpdateCode(this.props.id, _code)
	}

    onSubmitHandler = (e) => {
		e.preventDefault()
		let data =  {CodeText: this.state.code}
		axios.post(local_host + "/codeSave", JSON.stringify(data))
		.then(response => {
			console.log("Salio bien")
			console.log(response)
		})
		.catch(error => {
			console.log(error)
			console.log("Salio mal")
		})
	}


    render() {
        return (
            <div id={'Cell'}>
                <TextInput parentUpdateCode={this.updateCodeFromChild}/>
                <TextOutput/>
				<div><h3>{this.state.code}</h3></div>
            </div>
        );
    }
}

export default Cell;
