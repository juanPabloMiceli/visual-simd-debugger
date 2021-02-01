package cellshandler

import (
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/xmmhandler"
)

//XMMData ...
type XMMData struct {
	XmmID       string
	XmmValues   interface{}
	PrintFormat string
}

//CellData ...
type CellData struct {
	ID     int       `json:"id"`
	Code   string    `json:"code"`
	Output []XMMData `json:"output"`
}

//CellsData ...
type CellsData struct {
	Data                  []CellData `json:"CellsData"`
	HasDataCell           bool
	Requests              []XmmRequests
	DefaultDataFormat     []string
	DefaultPrintingFormat []string
}

//XmmRequest ...
type XmmRequest struct {
	XmmNumber   int
	XmmID       string
	DataFormat  string
	PrintFormat string
}

//XmmRequests ...
type XmmRequests []XmmRequest

//NewCellsData creates a new CellsData
func NewCellsData() CellsData {

	defaultDataFormat := make([]string, xmmhandler.XMMREGISTERS)
	defaultPrintingFormat := make([]string, xmmhandler.XMMREGISTERS)

	for i := range defaultDataFormat {
		defaultDataFormat[i] = xmmhandler.INT8STRING
		defaultPrintingFormat[i] = xmmhandler.SIGNEDFORMAT
	}

	return CellsData{
		Data:                  make([]CellData, 0),
		HasDataCell:           false,
		Requests:              make([]XmmRequests, 0),
		DefaultDataFormat:     defaultDataFormat,
		DefaultPrintingFormat: defaultPrintingFormat,
	}
}

//CellsData2SourceCode converts Cells Data to source code
func (obj *CellsData) CellsData2SourceCode() string {
	sourceCode := ""

	for _, cellData := range obj.Data {
		sourceCode += cellData.Code
	}

	return sourceCode
}

//HandleCellsData edit cells code content such that the cells to source code convertion is direct
func (obj *CellsData) HandleCellsData() {
	obj.toLowerCase()
	obj.fixCommentInstructions()
	obj.handleAllXmmRequests()
	obj.checkIfDataCellExists()
	obj.addDataSection()
	obj.addTextSection()
	obj.removeUserBreakpoints()
	obj.addCellsBreakpoints()
	obj.addExitSyscall()
}

func (obj *CellsData) toLowerCase() {
	for index := range obj.Data {
		obj.Data[index].Code = strings.ToLower(obj.Data[index].Code)
	}
}

func (obj *CellsData) fixCommentInstructions() {
	r := regexp.MustCompile("(( |\\t)+)?;(( |\\t)+)?")
	for index := range obj.Data {
		obj.Data[index].Code = r.ReplaceAllString(obj.Data[index].Code, ";")
	}
}

//GetGroupValues returns a map with each value found in the regexp match
func GetGroupValues(r *regexp.Regexp, match []string) map[string]string {

	values := make(map[string]string)

	for i, name := range r.SubexpNames() {
		values[name] = match[i]
	}

	return values
}

//XmmID2Number receives a string with the format "xmm<number>" and returns the number as an int
func XmmID2Number(xmmID string) int {
	runes := []rune(xmmID)
	xmmSafeString := string(runes[3:])
	xmmNumber, _ := strconv.Atoi(xmmSafeString)

	return xmmNumber
}

func (obj *CellsData) handleAllXmmRequests() {
	r := regexp.MustCompile(";(print|p)(?P<printFormat>\\/(d|x|t|u))?(( |\\t)+)?(?P<xmmID>xmm([0-9]|1[0-5]))\\.(?P<dataFormat>v16_int8|v8_int16|v4_int32|v2_int64|v4_float|v2_double)")
	for cellIndex := range obj.Data {
		obj.Requests = append(obj.Requests, make(XmmRequests, 0))
		obj.handleCellXmmRequests(r, cellIndex)
	}

}

func (obj *CellsData) handleCellXmmRequests(r *regexp.Regexp, cellIndex int) {
	matches := r.FindAllStringSubmatch(obj.Data[cellIndex].Code, -1)
	for _, match := range matches {
		obj.handleXmmRequest(r, match, cellIndex)
	}
}

func (obj *CellsData) handleXmmRequest(r *regexp.Regexp, match []string, cellIndex int) {

	var xmmRequest XmmRequest

	values := GetGroupValues(r, match)

	xmmRequest.XmmID = strings.ToUpper(values["xmmID"])
	xmmRequest.XmmNumber = XmmID2Number(xmmRequest.XmmID)
	xmmRequest.DataFormat = values["dataFormat"]
	xmmRequest.PrintFormat = values["printFormat"]
	obj.DefaultDataFormat[xmmRequest.XmmNumber] = xmmRequest.DataFormat

	obj.Requests[cellIndex] = append(obj.Requests[cellIndex], xmmRequest)
}

func (obj *CellsData) checkIfDataCellExists() {
	if containsAny(obj.Data[0].Code, ";data", "section .data") {
		obj.HasDataCell = true
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

func (obj *CellsData) addDataSection() {
	if obj.HasDataCell {
		obj.Data[0].Code = strings.Replace(obj.Data[0].Code, ";data", "section .data", 1)
	}
}

func (obj *CellsData) addTextSection() {
	var startText string
	startText += "\nglobal _start\n"
	startText += "section .text\n"
	startText += "_start:\n"

	if obj.HasDataCell {
		obj.Data[1].Code = startText + obj.Data[1].Code
	} else {
		obj.Data[0].Code = startText + obj.Data[0].Code
	}
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

func notDataCell(index int, hasDataCell bool) bool {
	return (index == 0 && !hasDataCell) || index != 0
}

func (obj *CellsData) addExitSyscall() {

	var exitText string
	exitText += "mov rax, 60\n"
	exitText += "mov rdi, 0\n"
	exitText += "syscall\n"

	var lastIndex = len(obj.Data) - 1
	obj.Data[lastIndex].Code += exitText
}
