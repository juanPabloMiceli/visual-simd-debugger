package models

//ResponseObj is the object send to the client as a JSON.
//This contains the console error and the info of every register to print.
type ResponseObj struct {
	ConsoleOut string
	CellRegs   []CellRegisters //Could be a slice of any of int or float types
}
