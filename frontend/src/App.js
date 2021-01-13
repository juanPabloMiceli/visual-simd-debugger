import React, { Component } from 'react';
import './App.css';
import axios from 'axios'
import Cell from "./components/Cell"
let local_host = "http://localhost:8080"



class App extends Component{

  constructor(){
    super()

    this.state = {
      CellsData: [{id:0, code: '', output: ''}],
      totalCells: 1
    }

    
  }


	onSubmitHandler = (e) => {
    e.preventDefault()
		axios.post(local_host + "/codeSave", JSON.stringify({
      CellsData: this.state.CellsData
    }))
		.then(response => {
      let copy = this.state.CellsData
      copy[this.state.totalCells-1].output = response.data.ConsoleOut.replaceAll(String.fromCharCode(0), '')
      this.setState({CellsData: copy})      
		})
		.catch(error => {
			console.log(error)
		})
	}

  onNewCellHandler = (e) => {
    e.preventDefault()
    let len = this.state.CellsData.length
    let joined = this.state.CellsData.concat({id: len, code:'', output:''})
    let newLen = this.state.totalCells + 1
    this.setState({CellsData: joined, len: newLen})
  }

  onDelCellHandler = (e) => {
    e.preventDefault()
    let len = this.state.CellsData.length
    if(len > 1){
      let deleted = this.state.CellsData
      let newLen = this.state.totalCells - 1
      deleted.pop()
      this.setState({CellsData: deleted, len: newLen})
    }
  }

  updateCodeFromChild = (id, _code) => {
    let copy = this.state.CellsData
    copy[id].code = _code
    this.setState({CellsData: copy})
	}

  render() {
    
    let CellsComponents = this.state.CellsData.map(cell => 
      <div>
        <Cell 
          key={cell.id}
          id={cell.id}
          parentUpdateCode={this.updateCodeFromChild}
          output={cell.output}
        />
      </div>)
    return (
      <div className="App">
        <button className="btn btn-success" id="submitButton" onClick={this.onSubmitHandler}>Submit Code</button>
        {CellsComponents}
				<button className="btn btn-newCell" id="newCellButton" onClick={this.onNewCellHandler}>New Cell</button>
				<button className="btn btn-DelCell" id="DelCellButton" onClick={this.onDelCellHandler}>Delete Cell</button>
      </div>
      
    );
  }
}

export default App;
