package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/zchee/kube-timeleap/pkg/ptrace"
	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/cellshandler"
	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/xmmhandler"
	"golang.org/x/sys/unix"
)

//XMMData ...
type XMMData struct {
	XmmID       string
	XmmValues   interface{}
	PrintFormat string
}

//CellRegisters contains the different XMMData in a cell.
type CellRegisters []XMMData

//Contains returns true if CellRegisters contains XMMData input
func (cellRegs *CellRegisters) Contains(newXmmData *XMMData) bool {

	for _, xmmData := range *cellRegs {
		if xmmData.XmmID == newXmmData.XmmID {
			return true
		}
	}
	return false
}

//ResponseObj ...
type ResponseObj struct {
	ConsoleOut string
	CellRegs   []CellRegisters //Could be a slice of any of int or float types
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func response(w *http.ResponseWriter, obj interface{}) {

	responseJSON, err := json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(responseJSON)
}

func getRequestedRegisters(cellsData *cellshandler.CellsData, xmmHandler *xmmhandler.XMMHandler, cellIndex int) CellRegisters {
	cellRegisters := CellRegisters{}

	for _, request := range cellsData.Requests[cellIndex] {
		fmt.Println("Request: ", request.PrintFormat)
		if request.PrintFormat == "" {
			request.PrintFormat = cellsData.DefaultPrintingFormat[request.XmmNumber]
		}
		fmt.Println("Final request: ", request.PrintFormat)
		xmmData := XMMData{XmmID: request.XmmID, XmmValues: xmmHandler.GetXMMData(request.XmmNumber, request.DataFormat, request.PrintFormat), PrintFormat: request.PrintFormat}
		cellRegisters = append(cellRegisters, xmmData)
	}

	return cellRegisters
}

func getChangedRegisters(oldXmmHandler *xmmhandler.XMMHandler, newXmmHandler *xmmhandler.XMMHandler, cellsData *cellshandler.CellsData, cellIndex int) CellRegisters {
	cellRegisters := CellRegisters{}

	if cellIndex > 0 {
		for index := range oldXmmHandler.Xmm {
			oldXmm := oldXmmHandler.Xmm[index]
			newXmm := newXmmHandler.Xmm[index]
			if !oldXmm.Equals(newXmm) {
				xmmString := "XMM" + strconv.Itoa(index)
				xmmData := XMMData{XmmID: xmmString, XmmValues: newXmmHandler.GetXMMData(index, cellsData.DefaultDataFormat[index], cellsData.DefaultPrintingFormat[index]), PrintFormat: cellsData.DefaultPrintingFormat[index]}
				cellRegisters = append(cellRegisters, xmmData)
			}
		}
	}

	return cellRegisters
}

func getXMMRegs(pid int) xmmhandler.XMMHandler {
	var unixRegs unix.PtraceRegs
	ptrace.GetFPRegs(pid, &unixRegs)

	fmt.Printf("\nAddress: %p\n", &unixRegs)

	fpPointer := (*ptrace.FPRegs)(unsafe.Pointer(&unixRegs))
	xmmSlice := fpPointer.XMMSpace[:]
	return xmmhandler.NewXMMHandler(&xmmSlice)
}

func joinWithPriority(cellRegs1 *CellRegisters, cellRegs2 *CellRegisters) CellRegisters {

	resCellRegisters := *cellRegs1

	for _, xmmData := range *cellRegs2 {
		if !resCellRegisters.Contains(&xmmData) {
			resCellRegisters = append(resCellRegisters, xmmData)
		}
	}

	return resCellRegisters
}

func setDefaultDataFormat(cellsData *cellshandler.CellsData, newDataFormat string) {
	for i := range cellsData.DefaultDataFormat {
		cellsData.DefaultDataFormat[i] = newDataFormat
	}
}

func setDefaultPrintFormat(cellsData *cellshandler.CellsData, newPrintFormat string) {
	for i := range cellsData.DefaultPrintingFormat {
		cellsData.DefaultPrintingFormat[i] = newPrintFormat
	}
}

func updatePrintFormat(cellsData *cellshandler.CellsData, cellIndex int) {

	r := regexp.MustCompile(";(print|p)(?P<printFormat>\\/(d|x|t|u))?(( |\\t)+)?(?P<xmmID>xmm([0-9]|1[0-5])?)\\.(?P<dataFormat>v16_int8|v8_int16|v4_int32|v2_int64|v4_float|v2_double)")
	matches := r.FindAllStringSubmatch(cellsData.Data[cellIndex].Code, -1)

	if len(matches) > 0 {

		for _, match := range matches {
			values := cellshandler.GetGroupValues(r, match)
			if values["xmmID"] != "xmm" {
				//This only changes one register
				fmt.Println("Quiero imprimir: ", values["xmmID"])
				xmmNumber := cellshandler.XmmID2Number(values["xmmID"])
				if !(values["printFormat"] == "") {
					fmt.Println("El valor estaba vacio, nuevo valor: ", values["printFormat"])
					cellsData.DefaultPrintingFormat[xmmNumber] = values["printFormat"]
				}
				cellsData.DefaultDataFormat[xmmNumber] = values["dataFormat"]

			} else {
				//I want to change all defaultValues

				setDefaultDataFormat(cellsData, values["dataFormat"])
				setDefaultPrintFormat(cellsData, values["printFormat"])

			}
		}

	}

}

func cellsLoop(cellsData *cellshandler.CellsData, pid int) ResponseObj {

	res := ResponseObj{CellRegs: make([]CellRegisters, 0)}

	cellIndex := 0

	if cellsData.HasDataCell {
		res.CellRegs = append(res.CellRegs, CellRegisters{})
		cellIndex++
	}

	oldXmmHandler := getXMMRegs(pid)

	ptrace.Cont(pid, 0)

	var ws syscall.WaitStatus

	syscall.Wait4(pid, &ws, syscall.WALL, nil)

	for !ws.Exited() {
		newXmmHandler := getXMMRegs(pid)
		updatePrintFormat(cellsData, cellIndex)
		requestedCellRegisters := getRequestedRegisters(cellsData, &newXmmHandler, cellIndex)
		changedCellRegisters := getChangedRegisters(&oldXmmHandler, &newXmmHandler, cellsData, cellIndex)
		selectedCellRegisters := joinWithPriority(&requestedCellRegisters, &changedCellRegisters)

		oldXmmHandler = newXmmHandler

		res.CellRegs = append(res.CellRegs, selectedCellRegisters)
		cellIndex++
		ptrace.Cont(pid, 0)
		syscall.Wait4(pid, &ws, syscall.WALL, nil)
	}

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())
	res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())

	return res
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	//Testing JSON Request

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	var dat map[string]interface{}
	err := json.Unmarshal(bodyBytes, &dat)

	if err != nil {
		panic(err)
	}

	jsonData, _ := json.Marshal(dat)

	fmt.Println(string(jsonData))

	// fmt.Println(bodyBytes)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//End of Testing JSON Request

	cellsData := cellshandler.NewCellsData()

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	if decodeErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't read data from the client properly."})
		return
	}

	var fileText string
	cellsData.HandleCellsData()
	fileText = cellsData.CellsData2SourceCode()
	fileErr := ioutil.WriteFile("output.asm", []byte(fileText), 0644)
	if fileErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't create file properly."})
		return
	}

	nasmPath, nasmPathErr := exec.LookPath("nasm")

	if nasmPathErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't find nasm executable path"})
		return
	}

	nasmCmd := exec.Command(nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil {
		response(&w, ResponseObj{ConsoleOut: stderr.String()})
		return
	}

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't find linker executable path"})
		return
	}

	linkingCmd := exec.Command(linkerPath, "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil {
		response(&w, ResponseObj{ConsoleOut: stderr.String()})
		return
	}
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		response(&w, ResponseObj{ConsoleOut: "Could't find server path"})
		return
	}
	fullPath := path.Join(path.Dir(filename), "output")
	exeCmd := exec.Command(fullPath)
	exeCmd.Stderr = os.Stderr
	exeCmd.Stdin = os.Stdin
	exeCmd.Stdout = os.Stdout
	exeCmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}

	startErr := exeCmd.Start()

	if startErr != nil {
		panic(startErr)
	}

	exeCmd.Wait()

	pid := exeCmd.Process.Pid

	responseObj := cellsLoop(&cellsData, pid)
	response(&w, responseObj)

}

func main() {
	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
