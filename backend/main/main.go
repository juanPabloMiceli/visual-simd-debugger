package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
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

	//CHARS is a string containing all possible characters in filename
	CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-0123456789."

	//FILENAMELEN is the created filename length
	FILENAMELEN = 10

	//MAXBYTES is the maximum bytes the asm file can use
	MAXBYTES = 30720 //30KBytes
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

//XMMData contains the data that has to be delivered to the frontend for each XMM register
type XMMData struct {
	XmmID     string
	XmmValues []string
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

//Estoy adentro de docker?
func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

//ResponseObj is the object send to the client as a JSON.
//This contains the console error and the info of every register to print.
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

func getRequestedRegisters(requests *cellshandler.XmmRequests, xmmHandler *xmmhandler.XMMHandler, xmmFormat *cellshandler.XMMFormat) CellRegisters {
	cellRegisters := CellRegisters{}

	for _, request := range *requests {
		// fmt.Println("Request: ", request.PrintFormat)
		if request.PrintFormat == "" {
			request.PrintFormat = xmmFormat.DefaultPrintingFormat[request.XmmNumber]
		}
		// fmt.Println("Final request: ", request.PrintFormat)
		xmmData := XMMData{
			XmmID:     request.XmmID,
			XmmValues: xmmHandler.GetXMMData(request.XmmNumber, request.DataFormat, request.PrintFormat)}

		cellRegisters = append(cellRegisters, xmmData)
	}

	return cellRegisters
}

func containsInt(elem int, s []int) bool {
	for _, current := range s {
		if current == elem {
			return true
		}
	}
	return false
}

func getChangedRegisters(hiddenRegs *cellshandler.HiddenInCell, oldXmmHandler *xmmhandler.XMMHandler, newXmmHandler *xmmhandler.XMMHandler, xmmFormat *cellshandler.XMMFormat) CellRegisters {
	cellRegisters := CellRegisters{}

	for index := range oldXmmHandler.Xmm {
		if !containsInt(index, *hiddenRegs) {
			oldXmm := oldXmmHandler.Xmm[index]
			newXmm := newXmmHandler.Xmm[index]
			if !oldXmm.Equals(newXmm) {
				xmmString := "XMM" + strconv.Itoa(index)
				xmmData := XMMData{
					XmmID:     xmmString,
					XmmValues: newXmmHandler.GetXMMData(index, xmmFormat.DefaultDataFormat[index], xmmFormat.DefaultPrintingFormat[index])}
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
		fmt.Println(err)
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

func setDefaultDataFormat(xmmFormat *cellshandler.XMMFormat, newDataFormat string) {
	for i := range xmmFormat.DefaultDataFormat {
		xmmFormat.DefaultDataFormat[i] = newDataFormat
	}
}

func setDefaultPrintFormat(xmmFormat *cellshandler.XMMFormat, newPrintFormat string) {
	for i := range xmmFormat.DefaultPrintingFormat {
		xmmFormat.DefaultPrintingFormat[i] = newPrintFormat
	}
}

func updatePrintFormat(cellsData *cellshandler.CellsData, cellIndex int, xmmFormat *cellshandler.XMMFormat) {
	r := regexp.MustCompile(`(( |\t)+)?;(( |\t)+)?(print|p)(( |\t)+)?(?P<printFormat>\/(d|x|t|u))?(( |\t)+)?(?P<xmmID>xmm([0-9]|1[0-5])?)\.(?P<dataFormat>v16_int8|v8_int16|v4_int32|v2_int64|v4_float|v2_double)`)
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
					xmmFormat.DefaultPrintingFormat[xmmNumber] = values["printFormat"]
				}
				xmmFormat.DefaultDataFormat[xmmNumber] = values["dataFormat"]

			} else {
				//I want to change all defaultValues

				setDefaultDataFormat(xmmFormat, values["dataFormat"])
				setDefaultPrintFormat(xmmFormat, values["printFormat"])

			}
		}
	}

}

func getFPRegs(pid int, data *FPRegs) error {
	_, _, errno := syscall.RawSyscall6(uintptr(syscall.SYS_PTRACE),
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

func limitFileSize(pid int, maxSize uint64) {
	var rlimit syscall.Rlimit

	rlimit.Cur = maxSize
	rlimit.Max = maxSize
	prLimit(pid, syscall.RLIMIT_FSIZE, &rlimit)
}

func limitCPUTime(pid int, maxTime uint64) {
	var rlimit syscall.Rlimit

	rlimit.Cur = maxTime
	rlimit.Max = maxTime
	prLimit(pid, syscall.RLIMIT_CPU, &rlimit)
}

func killProcess(pid int, err string) ResponseObj {
	fmt.Println("Killing Process")
	var ws syscall.WaitStatus
	_, _, killErr := syscall.RawSyscall6(syscall.SYS_KILL,
		uintptr(pid),
		uintptr(syscall.SIGKILL),
		0, 0, 0, 0)
	syscall.Wait4(pid, &ws, syscall.WALL, nil)
	if pidExists(pid) {
		return ResponseObj{ConsoleOut: err + "\nCould not kill process: " + strconv.Itoa(pid) + "\nError: " + killErr.Error()}

	}
	fmt.Println("Process killed succesfully.")

	return ResponseObj{ConsoleOut: err + "\nProcess killed succesfully."}
}

func cellsLoop(cellsData *cellshandler.CellsData, pid int, xmmFormat *cellshandler.XMMFormat) ResponseObj {

	res := ResponseObj{CellRegs: make([]CellRegisters, 0)}
	cellIndex := 0

	oldXmmHandler, getErr := getXMMRegs(pid)
	if getErr != nil {
		return killProcess(pid, "Could not get XMM registers.")
	}
	var ws syscall.WaitStatus

	for cellIndex < len(cellsData.Data) {
		newXmmHandler, getErr := getXMMRegs(pid)
		if getErr != nil {
			return killProcess(pid, "Could not get XMM registers.")
		}

		if cellIndex != 0 {
			updatePrintFormat(cellsData, cellIndex, xmmFormat)

		}
		requestedCellRegisters := getRequestedRegisters(&cellsData.Requests[cellIndex], &newXmmHandler, xmmFormat)
		changedCellRegisters := getChangedRegisters(&cellsData.HiddenRegs[cellIndex], &oldXmmHandler, &newXmmHandler, xmmFormat)
		selectedCellRegisters := joinWithPriority(&requestedCellRegisters, &changedCellRegisters)

		oldXmmHandler = newXmmHandler
		// fmt.Println(oldXmmHandler)

		res.CellRegs = append(res.CellRegs, selectedCellRegisters)
		cellIndex++
		fmt.Println(cellIndex)

		execErr := syscall.PtraceCont(pid, 0)
		if execErr != nil {
			return killProcess(pid, execErr.Error())
		}

		_, waitErr := syscall.Wait4(pid, &ws, syscall.WALL, nil)

		if waitErr != nil {
			return killProcess(pid, waitErr.Error())
		}
		if !pidExists(pid) && cellIndex < len(cellsData.Data)-1 {
			return killProcess(pid, "Something stopped the program.\n")
		}

	}

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())

	if pidExists(pid) {
		aux := killProcess(pid, "Something went wrong, program did not reach the end.")
		res.ConsoleOut = aux.ConsoleOut
	} else {
		res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())
	}

	return res
}

func pidExists(pid int) bool {
	_, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/status")
	return err == nil
}

func deleteFile(filePath string) error {
	if fileExists(filePath) {
		delExe := exec.Command("rm", filePath)
		delErr := delExe.Run()
		return delErr
	}

	return nil

}

//deleteFiles randomFolder 3 files created "output.asm", "output.o" and "output"
//So that next request is clean
func deleteFiles(folderPath string, res *ResponseObj) {
	err := os.RemoveAll(folderPath)
	// err1 := deleteFile(path.Join(filesPath, fileName))
	// err2 := deleteFile(path.Join(filesPath, fileName+".o"))
	// err3 := deleteFile(path.Join(filesPath, fileName+".asm"))

	if err != nil {
		fmt.Printf("Could not remove folder %s. Error: %s\n", folderPath, err.Error())
		res.ConsoleOut += "\nCould not remove your files from server, please notify. Error: " + err.Error()
	}
	// if err1 != nil {
	// 	res.ConsoleOut += "\nCould not remove exeecutable from server. Error: " + err1.Error()
	// }
	// if err2 != nil {
	// 	res.ConsoleOut += "\nCould not remove binary from server. Error: " + err2.Error()
	// }
	// if err3 != nil {
	// 	res.ConsoleOut += "\nCould not remove text file from server. Error: " + err3.Error()
	// }
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func getCellsData(req *http.Request) (cellshandler.CellsData, error) {
	cellsData := cellshandler.NewCellsData()

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	return cellsData, decodeErr
}

func printJSONInput(req *http.Request) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	var jsonMap map[string]interface{}
	err := json.Unmarshal(bodyBytes, &jsonMap)

	if err != nil {
		panic(err)
	}

	jsonData, _ := json.MarshalIndent(jsonMap, "", "\t")

	fmt.Println(string(jsonData))

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}

func checkExecutables(paths ...string) (map[string]string, []string) {
	resMap := make(map[string]string)
	var missingPaths []string

	for _, path := range paths {
		execPath, execErr := exec.LookPath(path)
		if execErr != nil {
			missingPaths = append(missingPaths, path)
		} else {
			resMap[path] = execPath
		}
	}
	return resMap, missingPaths

}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	safeString := base64.URLEncoding.EncodeToString(b)
	if err != nil {
		return "", err
	}
	return safeString, err
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	limitFileSize(syscall.Getpid(), MAXBYTES)

	enableCors(&w)

	printJSONInput(req)

	xmmFormat := cellshandler.NewXMMFormat()
	cellsData, decodeErr := getCellsData(req)
	if decodeErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't read data from the client properly."})
		return
	}
	if cellsData.HandleCellsData(&xmmFormat) {

		response(&w, ResponseObj{ConsoleOut: "Please insert some code."})
		return
	}

	var fileText = cellsData.CellsData2SourceCode()

	execMap, missingPaths := checkExecutables("nasm", "minijail0", "ld", "microjail")

	if len(missingPaths) > 0 {
		responseObj := ResponseObj{ConsoleOut: "Could't find next executable paths:"}
		for _, path := range missingPaths {
			responseObj.ConsoleOut += "\n* " + path
		}
		response(&w, responseObj)
		return
	}

	randomID, randErr := randomString(FILENAMELEN)

	if randErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't create file name properly. (Server error please notify)"})
		fmt.Println("Error creating random file: ", randErr)
		return
	}

	randomFolder := path.Join("/clients", randomID) //Cada usuario tiene una carpeta propia lo cual le va a dar aislamiento
	randomFile := path.Join(randomFolder, randomID) //Cada usuario tiene una carpeta propia lo cual le va a dar aislamiento
	os.Mkdir(randomFolder, os.FileMode(0777))
	fileErr := ioutil.WriteFile(randomFile+".asm", []byte(fileText), 0644)

	if fileErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't create file properly. Maybe the file is greater than 30Kb."})
		fmt.Println("Error creating asm file. Maybe the file is greater than 30Kb. Error: ", fileErr)
		return
	}

	fileInfo, newFileErr := os.Stat(randomFile + ".asm")
	if newFileErr != nil {
		fmt.Println("Error obtaining file stats. Error: ", newFileErr)
		panic(newFileErr)
	}
	fileSize := fileInfo.Size()

	if fileSize > MAXBYTES {
		res := ResponseObj{ConsoleOut: "Text file must not be greater than 30Kb."}
		fmt.Println("File is larger than 30Kb. Aborting.")
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	/*
		////////////////////////////////////////////////////////
		///////////inicio codigo para informe///////////////////
		////////////////////////////////////////////////////////

		//Empezamos corriendo PS sin namespace
		fmt.Println("")
		fmt.Println("=============================================================")
		fmt.Println("=Procesos visibles para NASM/LD sin aplicar ningun namespace=")
		fmt.Println("=============================================================")
		fmt.Println("$ ps -aux")
		informeCmd := exec.Command("/bin/ps", "-aux") //Ejecutamos PS

		informeCmd.Stdout = os.Stdout
		informeErr := informeCmd.Run()
		if informeErr != nil {
			fmt.Println("Error: ", informeErr)
		}

		//Ahora corremos findmnt sin namespace
		fmt.Println("")
		fmt.Println("================================================================")
		fmt.Println("=Filesystems visibles para NASM/LD sin aplicar ningun namespace=")
		fmt.Println("================================================================")
		fmt.Println("$ findmnt --output TARGET,SOURCE,VFS-OPTIONS")

		informeCmd = exec.Command("/bin/findmnt", "--output", "TARGET,SOURCE,VFS-OPTIONS") //Ejecutamos findmnt

		informeCmd.Stdout = os.Stdout
		informeErr = informeCmd.Run()
		if informeErr != nil {
			fmt.Println("Error: ", informeErr)
		}

		//Corremos PS con el namespace
		fmt.Println("")
		fmt.Println("===================================================================================")
		fmt.Println("=Procesos visibles para NASM/LD luego de aplicar minijail con PID y VFS namespaces=")
		fmt.Println("===================================================================================")
		fmt.Println("$ ps -aux")
		informeCmd = exec.Command(execMap["minijail0"], //Path a minijail
			"-p",               //PID namespace
			"-v",               //Vamos a crear un nuevo VFS
			"-P", "/var/empty", //Hacemos un pivot_root a /var/empty
			"-b", fmt.Sprintf("%s,,1", randomFolder), // Bindeamos la carpeta del cliente con permiso de escritura
			"-b", "/bin/ps", //Bindiamos el binario de PS
			"-b", "/proc", //Bindiamos /proc porque PS lo necesita
			"-r",              //Remonta /proc a readonly
			"/bin/ps", "-aux") //Ejecutamos PS

		informeCmd.Stdout = os.Stdout
		informeErr = informeCmd.Run()
		if informeErr != nil {
			fmt.Println("Error: ", informeErr)
		}

		//Corremos findmnt con el namespace
		fmt.Println("")
		fmt.Println("======================================================================================")
		fmt.Println("=Filesystems visibles para NASM/LD luego de aplicar minijail con PID y VFS namespaces=")
		fmt.Println("======================================================================================")
		fmt.Println("$ findmnt --output TARGET,SOURCE,VFS-OPTIONS")

		informeCmd = exec.Command(execMap["minijail0"], //Path a minijail
			"-p",               //PID namespace
			"-v",               //Vamos a crear un nuevo VFS
			"-P", "/var/empty", //Hacemos un pivot_root a /var/empty
			"-b", fmt.Sprintf("%s,,1", randomFolder), // Bindeamos la carpeta del cliente con permiso de escritura
			"-b", "/bin/findmnt", //Bindiamos el binario de findmnt
			"-b", "/proc", //Bindiamos /proc porque findmnt la necesita
			"-r",                                                    //Remonta /proc a readonly
			"/bin/findmnt", "--output", "TARGET,SOURCE,VFS-OPTIONS") //Ejecutamos findmnt

		informeCmd.Stdout = os.Stdout
		informeErr = informeCmd.Run()
		if informeErr != nil {
			fmt.Println("Error: ", informeErr)
		}

		///////////////////////////////////////////////////////
		////////////fin codigo para informe////////////////////
		///////////////////////////////////////////////////////
	*/
	//NASM con namespace de todo tipo y color
	nasmCmd := exec.Command(execMap["minijail0"], //Path a minijail
		"-p",               //PID namespace
		"-v",               //Vamos a crear un nuevo VFS
		"-P", "/var/empty", //Hacemos un pivot_root a /var/empty
		"-b", fmt.Sprintf("%s,,1", randomFolder), // Bindeamos la carpeta del cliente a si misma con permiso de escritura
		"-b", "/usr/bin/nasm", //Bindiamos esto para tener el binario a NASM
		"-b", "/proc", //Bindiamos /proc
		"-r",                                                                                          //Remonta /proc a readonly
		execMap["nasm"], "-f", "elf64", "-g", "-F", "DWARF", randomFile+".asm", "-o", randomFile+".o") //Comando ejecutador de NASM

	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmCmd.Stdout = os.Stdout
	nasmErr := nasmCmd.Start()

	if nasmErr != nil {
		fmt.Println("Error starting NASM: ", nasmErr.Error())
		res := ResponseObj{ConsoleOut: nasmErr.Error()}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	nasmPID := nasmCmd.Process.Pid
	limitCPUTime(nasmPID, MAXCPUTIME)

	nasmErr = nasmCmd.Wait()

	if nasmErr != nil {
		fmt.Println("Error executing nasm: ", nasmErr)
		fmt.Println(stderr.String())
		errorString := strings.ReplaceAll(stderr.String(), randomFile, "output")
		res := ResponseObj{ConsoleOut: errorString}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	if !fileExists(randomFile + ".o") {
		fmt.Println("NASM execution finished correctly but didn't create expected file: " + randomFile + ".o")
		res := ResponseObj{ConsoleOut: "NASM execution failed."}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	fmt.Println("")
	fmt.Println("Program compiled correctly")

	//LD con namespace de todo tipo y color
	linkingCmd := exec.Command(execMap["minijail0"], //Path a minijail
		"-p",                          //PID namespace
		"-n",                          //No new priviligies
		"-S", "../policies/ld.policy", //Setea las policies para LD
		"-v",               //Vamos a crear un nuevo VFS
		"-P", "/var/empty", //Hacemos un pivot_root a /var/empty
		"-b", fmt.Sprintf("%s,,1", randomFolder), // Bindeamos la carpeta del cliente a si misma con permiso de escritura
		"-b", "/usr/bin/ld", //Bindiamos esto para tener el binario a LD
		"-b", "/proc", //Bindiamos /proc
		"-r",                                                                     //Remonta /proc a readonly
		execMap["ld"], "-nostdlib", "-static", "-o", randomFile, randomFile+".o") //Comando ejecutador de LD

	linkingCmd.Stderr = &stderr
	linkingCmd.Stdout = os.Stdout

	linkingErr := linkingCmd.Start()

	if linkingErr != nil {
		fmt.Println("Error starting LD: ", linkingErr.Error())
		res := ResponseObj{ConsoleOut: linkingErr.Error()}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	linkerPID := linkingCmd.Process.Pid
	limitCPUTime(linkerPID, MAXCPUTIME)

	linkingErr = linkingCmd.Wait()

	if linkingErr != nil {
		fmt.Println("Error executing LD: ", linkingErr)
		fmt.Println(stderr.String())
		errorString := strings.ReplaceAll(stderr.String(), randomFile, "output")
		res := ResponseObj{ConsoleOut: errorString}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	if !fileExists(randomFile + ".o") {
		fmt.Println("LD execution finished correctly but didn't create expected file: " + randomFile)
		res := ResponseObj{ConsoleOut: "LD execution failed."}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	fmt.Println("Program linked correctly")

	exeCmd := exec.Command(execMap["microjail"], randomFile)

	exeCmd.Stderr = os.Stderr
	exeCmd.Stdin = os.Stdin
	exeCmd.Stdout = os.Stdout
	exeCmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}
	runtime.LockOSThread()

	startErr := exeCmd.Start()

	if startErr != nil {
		fmt.Println("Error starting microjail: ", startErr.Error())
		res := ResponseObj{ConsoleOut: startErr.Error()}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	microjailPID := exeCmd.Process.Pid
	limitCPUTime(microjailPID, MAXCPUTIME)

	exeCmd.Wait()

	optErr := syscall.PtraceSetOptions(microjailPID, 0x100000|syscall.PTRACE_O_TRACEEXEC) //0x100000 = PTRACE_O_EXITKILL, 0x200000 = PTRACE_O_SUSPEND_SECCOMP

	if optErr != nil {
		res := killProcess(microjailPID, optErr.Error())
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	//One continue such that the C execve is made
	execErr := syscall.PtraceCont(microjailPID, 0)
	if execErr != nil {
		res := killProcess(microjailPID, execErr.Error())
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	var ws syscall.WaitStatus
	_, waitErr := syscall.Wait4(microjailPID, &ws, syscall.WALL, nil)
	if waitErr != nil {
		res := killProcess(microjailPID, waitErr.Error())
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	if !pidExists(microjailPID) {
		res := ResponseObj{ConsoleOut: "Microjail error."}
		deleteFiles(randomFolder, &res)
		response(&w, res)
		return
	}

	responseObj := cellsLoop(&cellsData, microjailPID, &xmmFormat)

	runtime.UnlockOSThread()
	deleteFiles(randomFolder, &responseObj)
	fmt.Println(responseObj)
	response(&w, responseObj)

}

func main() {
	runtime.GOMAXPROCS(1)

	http.HandleFunc("/codeSave", codeSave)

	fmt.Println("===============================")
	fmt.Println("=Server listening on port 8080=")
	fmt.Println("===============================")
	http.ListenAndServe(":8080", nil)

}
