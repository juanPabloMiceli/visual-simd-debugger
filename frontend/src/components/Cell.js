import React, { Component } from 'react';
import TextInput from './TextInput';
import TextOutput from './TextOutput';
import axios from 'axios'
let local_host = "http://localhost:8080"

axios.defaults.headers.common = {
	"Content-Type": "application/json"
  }

class Cell extends Component {

	updateCodeFromChild = (_code) => {
		this.props.parentUpdateCode(this.props.id, _code)
	}


    render() {
        return (
            <div id={'Cell'}>
                <TextInput parentUpdateCode={this.updateCodeFromChild}/>
                <TextOutput text={this.props.output}/>
            </div>
        );
    }
}

export default Cell;
