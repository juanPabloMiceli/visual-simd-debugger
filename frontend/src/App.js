import React from 'react';
import "./components/icons";
import './App.css';
import SimdVisualizer from './components/SimdVisualizer';
import Provider from './components/Context'



export default function App(){
  return(
    <div id="App">
      <Provider>
        <SimdVisualizer/>
      </Provider>
    </div>
    
    
  )
}