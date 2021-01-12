import React, { Component } from 'react';
import './App.css';
import axios from 'axios'
import Cell from "./components/Cell"
let local_host = "http://localhost:8080"



class App extends Component{

  constructor(){
    super()

    this.state = {
      CellsData: [{id:0, code: ''}]
    }

    
  }


	onSubmitHandler = (e) => {
		e.preventDefault()
		axios.post(local_host + "/codeSave", JSON.stringify(this.state))
		.then(response => {
			console.log(response)
		})
		.catch(error => {
			console.log(error)
		})
	}

  onNewCellHandler = (e) => {
    e.preventDefault()
    let len = this.state.CellsData.length
    let joined = this.state.CellsData.concat({id: len, code:''})
    this.setState({CellsData: joined})
  }

  onDelCellHandler = (e) => {
    e.preventDefault()
    let len = this.state.CellsData.length
    if(len > 1){
      let deleted = this.state.CellsData
      deleted.pop()
      this.setState({CellsData: deleted})
    }
  }

  updateCodeFromChild = (id, _code) => {
    let copy = this.state.CellsData
    copy[id] = {id: id,
                code: _code}
    this.setState({CellsData: copy})
    console.log(this.state.CellsData)
	}

  render() {

    let CellsComponents = this.state.CellsData.map(cell => 
      <div>
        <Cell 
          key={cell.id}
          id={cell.id}
          parentUpdateCode={this.updateCodeFromChild}
        />
      </div>)
    return (
      <div className="App">
        <button className="btn btn-success" id="submitButton" onClick={this.onSubmitHandler}>Submit Code</button>
        {CellsComponents}
				<button className="btn btn-newCell" id="newCellButton" onClick={this.onNewCellHandler}>New Cell</button>
				<button className="btn btn-DelCell" id="DelCellButton" onClick={this.onDelCellHandler}>Delete Cell</button>
        {this.state.CellsData.map(data => 
          <h3>
            {data.code}
          </h3>)}
      </div>
      
    );
  }
}

export default App;
