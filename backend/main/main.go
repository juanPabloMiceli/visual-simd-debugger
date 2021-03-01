package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

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
	MAXCPUTIME uint64 = 2
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

func getXMMRegs(pid int) (xmmhandler.XMMHandler, error) {
	var fpRegs FPRegs

	err := getFPRegs(pid, &fpRegs)
	fmt.Printf("\nAddress fp: %p\n", &fpRegs)

	if err != nil {

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

func killProcess(pid int, err string) ResponseObj {
	killErr := syscall.Kill(pid, syscall.SYS_KILL)
	if pidExists(pid) {
		fmt.Println("err: ", err)
		fmt.Println("pid: ", pid)
		fmt.Println("killErr: ", killErr.Error())
		return ResponseObj{ConsoleOut: err + "\nCould not kill process: " + strconv.Itoa(pid) + "\nError: " + killErr.Error()}

	}
	return ResponseObj{ConsoleOut: err + "\nProcess killed succesfully."}
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
		return killProcess(pid, "Could not get XMM registers.")
	}

	execErr := syscall.PtraceCont(pid, 0)

	if execErr != nil {
		return killProcess(pid, execErr.Error())
	}

	var ws syscall.WaitStatus

	_, waitErr := syscall.Wait4(pid, &ws, syscall.WALL, nil)
	if waitErr != nil {
		return killProcess(pid, waitErr.Error())
	}

	for cellIndex < len(cellsData.Data) {
		newXmmHandler, getErr := getXMMRegs(pid)
		if getErr != nil {
			return killProcess(pid, "Could not get XMM registers.")
		}

		updatePrintFormat(cellsData, cellIndex)
		requestedCellRegisters := getRequestedRegisters(cellsData, &newXmmHandler, cellIndex)
		changedCellRegisters := getChangedRegisters(&oldXmmHandler, &newXmmHandler, cellsData, cellIndex)
		selectedCellRegisters := joinWithPriority(&requestedCellRegisters, &changedCellRegisters)

		oldXmmHandler = newXmmHandler
		fmt.Println(oldXmmHandler)

		res.CellRegs = append(res.CellRegs, selectedCellRegisters)
		cellIndex++
		execErr = syscall.PtraceCont(pid, 0)
		if execErr != nil {
			return killProcess(pid, execErr.Error())
		}

		_, waitErr = syscall.Wait4(pid, &ws, syscall.WALL, nil)
		if waitErr != nil {
			return killProcess(pid, waitErr.Error())
		}

		fmt.Println(pidExists(pid))

	}

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())
	res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())

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

func deleteFile(filePath string) error {
	if fileExists(filePath) {
		delExe := exec.Command("rm", filePath)
		delErr := delExe.Run()
		return delErr
	}

	return nil

}

//deleteFiles removes the 3 files created "output.asm", "output.o" and "output"
//So that next request is clean
func deleteFiles(filesPath string) error {

	err := deleteFile(path.Join(filesPath, "output"))
	if err != nil {
		return err
	}
	err = deleteFile(path.Join(filesPath, "output.o"))
	if err != nil {
		return err
	}
	err = deleteFile(path.Join(filesPath, "output.asm"))
	if err != nil {
		return err
	}

	return err
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		response(&w, ResponseObj{ConsoleOut: "Could't find server path"})
		return
	}

	filepath := path.Dir(filename)
	//Deleting files just in case previous execution failed to do that
	deleteFiles(filepath)
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

	mjPath, mjPathError := exec.LookPath("minijail0")

	if mjPathError != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't find minijail executable path"})
		return
	}

	nasmCmd := exec.Command(mjPath, "-n", "-S", "../policies/nasm.policy", nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm", "-o", "output.o")

	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil || !fileExists(path.Join(filepath, "output.o")) {
		if stderr.String() == "" {
			stderr.WriteString("NASM execution failed")
		}
		res := ResponseObj{ConsoleOut: stderr.String()}
		delErr := deleteFiles(filepath)
		if delErr != nil {
			res.ConsoleOut += "\n Could not remove server files. Error: " + delErr.Error()
			fmt.Println("\n Could not remove server files. Error: " + delErr.Error())
		}
		response(&w, res)
		return
	}
	fmt.Println("Program compiled")

	// nasmPid := nasmCmd.Process.Pid

	// fmt.Println("\n\nNASM")
	// printLibraries(nasmPid)

	// nasmErr = nasmCmd.Wait()

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't find linker executable path"})
		return
	}

	linkingCmd := exec.Command(mjPath, "-n", "-S", "../policies/ld.policy", linkerPath, "-o", "output", "output.o")

	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil || !fileExists(path.Join(filepath, "output")) {
		if stderr.String() == "" {
			stderr.WriteString("Linker execution failed")
		}
		res := ResponseObj{ConsoleOut: stderr.String()}
		delErr := deleteFiles(filepath)
		if delErr != nil {
			res.ConsoleOut += "\n Could not remove server files. Error: " + delErr.Error()
			fmt.Println("\n Could not remove server files. Error: " + delErr.Error())
		}
		response(&w, res)
		return
	}
	fmt.Println("Program Linked")

	// ldPid := linkingCmd.Process.Pid

	// fmt.Println("\n\nLINKER")
	// printLibraries(ldPid)

	// linkingCmd.Wait()

	fullPath := path.Join(filepath, "output")

	exeCmd := exec.Command(fullPath)

	// exeCmd := exec.Command("minijail0", "-n", "-S", "../policies/exec.policy", fullPath)
	exeCmd.Stderr = os.Stderr
	exeCmd.Stdin = os.Stdin
	exeCmd.Stdout = os.Stdout
	exeCmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}
	runtime.LockOSThread()

	fmt.Println("Starting...")
	startErr := exeCmd.Start()

	if startErr != nil {
		res := ResponseObj{ConsoleOut: startErr.Error()}
		delErr := deleteFiles(filepath)
		if delErr != nil {
			res.ConsoleOut += "\n Could not remove server files. Error: " + delErr.Error()
			fmt.Println("\n Could not remove server files. Error: " + delErr.Error())
		}
		response(&w, res)
		return
	}

	pid := exeCmd.Process.Pid
	limitCPUTime(pid, MAXCPUTIME)

	// fmt.Println("\n\nEXEC")
	// fmt.Println(pid)
	// printLibraries(pid)
	fmt.Println("Waiting process ", pid)
	exeCmd.Wait()

	responseObj := cellsLoop(&cellsData, pid, exeCmd)

	runtime.UnlockOSThread()
	delErr := deleteFiles(filepath)
	if delErr != nil {
		responseObj.ConsoleOut += "\n Could not remove server files. Error: " + delErr.Error()
		fmt.Println("\n Could not remove server files. Error: " + delErr.Error())
	}
	response(&w, responseObj)

}

func main() {
	runtime.GOMAXPROCS(1)

	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
