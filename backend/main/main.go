package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/cellshandler"
	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/xmmhandler"
)

const (
	//MAXCPUTIME is the maximum time in seconds the process can be scheduled
	MAXCPUTIME uint64 = 1

	//MAXPROCESSTIME is the maximum wall time in seconds the process can run in the server
	MAXPROCESSTIME time.Duration = 2
)

// FPRegs represents a user_fpregs_struct in /usr/include/x86_64-linux-gnu/sys/user.h.
type FPRegs struct {
	Cwd      uint16     // Control Word
	Swd      uint16     // Status Word
	Ftw      uint16     // Tag Word
	Fop      uint16     // Last Instruction Opcode
	Rip      uint64     // Instruction Pointer
	Rdp      uint64     // Data Pointer
	Mxcsr    uint32     // MXCSR Register State
	MxcrMask uint32     // MXCR Mask
	StSpace  [32]uint32 // 8*16 bytes for each FP-reg = 128 bytes
	XMMSpace [256]byte  // 16*16 bytes for each XMM-reg = 256 bytes
	_        [24]uint32 // padding
}

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

func runningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
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
		xmmData := XMMData{
			XmmID:       request.XmmID,
			XmmValues:   xmmHandler.GetXMMData(request.XmmNumber, request.DataFormat, request.PrintFormat),
			PrintFormat: request.PrintFormat}

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
				xmmData := XMMData{
					XmmID:       xmmString,
					XmmValues:   newXmmHandler.GetXMMData(index, cellsData.DefaultDataFormat[index], cellsData.DefaultPrintingFormat[index]),
					PrintFormat: cellsData.DefaultPrintingFormat[index]}
				cellRegisters = append(cellRegisters, xmmData)
			}
		}
	}

	return cellRegisters
}

// func getXMMRegs(pid int) (xmmhandler.XMMHandler, error) {
// 	var unixRegs unix.PtraceRegs

// 	err := ptrace.GetFPRegs(pid, &unixRegs)
// 	if err != nil {
// 		//Something went wrong, time to kill
// 		return xmmhandler.XMMHandler{}, err
// 	}
// 	fmt.Printf("\nAddress: %p\n", &unixRegs)

// 	fpPointer := (*FPRegs)(unsafe.Pointer(&unixRegs))
// 	xmmSlice := fpPointer.XMMSpace[:]

// 	return xmmhandler.NewXMMHandler(&xmmSlice), err
// }

func getXMMRegs(pid int) (xmmhandler.XMMHandler, error) {
	var fpRegs FPRegs

	err := getFPRegs(pid, &fpRegs)
	fmt.Printf("\nAddress fp: %p\n", &fpRegs)

	if err != nil {

		fmt.Println(getProcessStatus(pid, "State:\t"))
		return xmmhandler.XMMHandler{}, err
	}
	xmmSlice := fpRegs.XMMSpace[:]

	return xmmhandler.NewXMMHandler(&xmmSlice), err
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

// func ptrace(request int, pid int, addr uintptr, data uintptr) (err error) {
// 	_, _, e1 := Syscall6(SYS_PTRACE, uintptr(request), uintptr(pid), uintptr(addr), uintptr(data), 0, 0)
// 	if e1 != 0 {
// 		err = errnoErr(e1)
// 	}
// 	return
// }

func getFPRegs(pid int, data *FPRegs) error {
	sysPtrace := 101
	_, _, errno := syscall.RawSyscall6(uintptr(sysPtrace),
		uintptr(syscall.PTRACE_GETFPREGS),
		uintptr(pid),
		uintptr(0),
		uintptr(unsafe.Pointer(data)),
		0,
		0)

	var err error
	if errno != 0 {
		err = errno
		return err
	}
	return nil
}

func prLimit(pid int, limit uintptr, rlimit *syscall.Rlimit) error {
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRLIMIT64,
		uintptr(pid),
		limit,
		uintptr(unsafe.Pointer(rlimit)),
		0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
		return err
	}
	return nil
}

func limitCPUTime(pid int, maxTime uint64) {
	var rlimit syscall.Rlimit

	rlimit.Cur = maxTime
	rlimit.Max = maxTime
	prLimit(pid, syscall.RLIMIT_CPU, &rlimit)
}

func stateStopped(state string) bool {
	stopedStates := []string{"t", "T", "Z", "X"}

	for _, stopState := range stopedStates {
		if state == stopState {
			return true
		}
	}
	return false
}

func checkTimeout(timeoutCH chan error, resCH chan bool, pid int) {
	select {
	case <-time.After(2 * time.Second):
		// fmt.Println("Is ", pid, " present: ", pidExists(pid))
		// fmt.Println("Killing process")
		// syscall.Kill(pid, syscall.SYS_KILL)
		// path := "/proc/" + strconv.Itoa(pid) + "/status"
		// echoExe := exec.Command("ls", path)
		// echoExe.Stdout = os.Stdout
		// echoExe
		// fmt.Println("Is ", pid, " present: ", pidExists(pid))
		resCH <- true

		return
	case <-timeoutCH:
		fmt.Println("Process stopped before timeout")
		resCH <- false
		return
	}
}

func waitingTime(pid int, ws *syscall.WaitStatus, timeoutCH chan error) {
	_, waitErr := syscall.Wait4(pid, ws, syscall.WALL, nil)
	timeoutCH <- waitErr

}

func cellsLoop(cellsData *cellshandler.CellsData, pid int, cmd *exec.Cmd) ResponseObj {

	res := ResponseObj{CellRegs: make([]CellRegisters, 0)}

	cellIndex := 0

	if cellsData.HasDataCell {
		res.CellRegs = append(res.CellRegs, CellRegisters{})
		cellIndex++
	}

	oldXmmHandler, getErr := getXMMRegs(pid)
	if getErr != nil {
		syscall.Kill(pid, syscall.SYS_KILL)
		return ResponseObj{ConsoleOut: "Could not get XMM registers"}
	}

	timeoutChannel := make(chan error, 1)
	resChannel := make(chan bool, 1)
	// ptrace.Cont(pid, 0)
	syscall.PtraceCont(pid, 0)

	var ws syscall.WaitStatus

	go checkTimeout(timeoutChannel, resChannel, pid)
	go waitingTime(pid, &ws, timeoutChannel)

	if <-resChannel == true {
		fmt.Println("Is ", pid, " present: ", pidExists(pid))
		fmt.Println("Killing process")
		syscall.Kill(pid, syscall.SYS_KILL)
		fmt.Println("Is ", pid, " present: ", pidExists(pid))

		return ResponseObj{ConsoleOut: "Execution timeout"}
	}

	for cellIndex < len(cellsData.Data) {
		newXmmHandler, getErr := getXMMRegs(pid)
		if getErr != nil {
			syscall.Kill(pid, syscall.SYS_KILL)
			return ResponseObj{ConsoleOut: "Could not get XMM registers"}
		}

		updatePrintFormat(cellsData, cellIndex)
		requestedCellRegisters := getRequestedRegisters(cellsData, &newXmmHandler, cellIndex)
		changedCellRegisters := getChangedRegisters(&oldXmmHandler, &newXmmHandler, cellsData, cellIndex)
		selectedCellRegisters := joinWithPriority(&requestedCellRegisters, &changedCellRegisters)

		oldXmmHandler = newXmmHandler

		res.CellRegs = append(res.CellRegs, selectedCellRegisters)
		cellIndex++
		// ptrace.Cont(pid, 0)
		syscall.PtraceCont(pid, 0)
		go checkTimeout(timeoutChannel, resChannel, pid)
		go waitingTime(pid, &ws, timeoutChannel)

		if <-resChannel == true {
			fmt.Println("Is ", pid, " present: ", pidExists(pid))
			fmt.Println("Killing process")
			syscall.Kill(pid, syscall.SYS_KILL)
			fmt.Println("Is ", pid, " present: ", pidExists(pid))

			return ResponseObj{ConsoleOut: "Execution timeout"}
		}

	}

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())
	res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())
	fmt.Println(res)

	return res
}

func printLibraries(pid int) {
	content, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/maps")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}

func getProcessStatus(pid int, paramString string) string {

	content, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		panic(err)
	}
	contentStr := string(content)
	index := strings.Index(contentStr, paramString)
	return string(contentStr[index+len(paramString)])
}

func pidExists(pid int) bool {

	_, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		return false
	}
	return true
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	//Testing JSON Request

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	fmt.Println(string(bodyBytes))

	var jsonMap map[string]interface{}
	err := json.Unmarshal(bodyBytes, &jsonMap)

	if err != nil {
		panic(err)
	}

	jsonData, _ := json.MarshalIndent(jsonMap, "", "\t")

	fmt.Println(string(jsonData))

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

	nasmCmd := exec.Command(nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm", "-o", "output.o")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Start()

	if nasmErr != nil {
		response(&w, ResponseObj{ConsoleOut: stderr.String()})
		return
	}

	nasmPid := nasmCmd.Process.Pid

	fmt.Println("\n\nNASM")
	printLibraries(nasmPid)

	nasmCmd.Wait()

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't find linker executable path"})
		return
	}

	linkingCmd := exec.Command(linkerPath, "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Start()

	if linkingErr != nil {
		response(&w, ResponseObj{ConsoleOut: stderr.String()})
		return
	}

	ldPid := linkingCmd.Process.Pid

	fmt.Println("\n\nLINKER")
	printLibraries(ldPid)

	linkingCmd.Wait()

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
		response(&w, ResponseObj{ConsoleOut: stderr.String()})
		return
	}

	pid := exeCmd.Process.Pid

	fmt.Println("\n\nEXEC")
	fmt.Println(pid)
	printLibraries(pid)

	exeCmd.Wait()

	limitCPUTime(pid, MAXCPUTIME)

	responseObj := cellsLoop(&cellsData, pid, exeCmd)
	response(&w, responseObj)
	fmt.Printf("\n\nfinished\n\n")

}

func main() {

	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
