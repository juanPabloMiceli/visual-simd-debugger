package cellshandler

import (
	"strings"
)

//CellData ...
type CellData struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Output string `json:"output"`
}

//CellsData ...
type CellsData struct {
	Data        []CellData `json:"CellsData"`
	HasDataCell bool
}

//NewCellsData creates a new CellsData
func NewCellsData() CellsData {
	return CellsData{Data: make([]CellData, 0), HasDataCell: false}
}

//HandleCellsData edit cells code content such that the cells to source code convertion is direct
func (obj *CellsData) HandleCellsData() {
	obj.checkIfDataCellExists()
	obj.addDataSection()
	obj.addTextSection()
	obj.removeUserBreakpoints()
	obj.addCellsBreakpoints()
	obj.addExitSyscall()
}

//CellsData2SourceCode converts Cells Data to source code
func (obj *CellsData) CellsData2SourceCode() string {
	sourceCode := ""

	for _, cellData := range obj.Data {
		sourceCode += cellData.Code
	}

	return sourceCode
}

func (obj *CellsData) addExitSyscall() {

	var exitText string
	exitText += "mov rax, 60\n"
	exitText += "mov rdi, 0\n"
	exitText += "syscall\n"

	var lastIndex = len(obj.Data) - 1
	obj.Data[lastIndex].Code += exitText
}

func notDataCell(index int, hasDataCell bool) bool {
	return (index == 0 && !hasDataCell) || index != 0
}

func (obj *CellsData) removeUserBreakpoints() {
	for index, cellData := range obj.Data {
		if strings.Contains(cellData.Code, "int 3") {
			obj.Data[index].Code = strings.ReplaceAll(cellData.Code, "int 3", "")
		}
	}
}

func (obj *CellsData) addCellsBreakpoints() {
	for index := range obj.Data {
		if notDataCell(index, obj.HasDataCell) {
			obj.Data[index].Code += "\nint 3\n"
		}
	}
}

func containsAny(original string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(original, sub) {
			return true
		}
	}
	return false
}

func (obj *CellsData) addTextSection() {
	var startText string
	startText += "global _start\n"
	startText += "section .text\n"
	startText += "_start:\n"

	if obj.HasDataCell {
		obj.Data[1].Code = startText + obj.Data[1].Code
	} else {
		obj.Data[0].Code = startText + obj.Data[0].Code
	}
}

func (obj *CellsData) addDataSection() {
	if obj.HasDataCell {
		obj.Data[0].Code = strings.Replace(obj.Data[0].Code, ";data", "section .data", 1)
	}
}

func (obj *CellsData) checkIfDataCellExists() {
	if containsAny(obj.Data[0].Code, ";data", "section .data") {
		obj.HasDataCell = true
	}
}
