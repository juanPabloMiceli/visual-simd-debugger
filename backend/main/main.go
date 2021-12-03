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

	"../cellshandler"
	"../models"
	"../utils"
	"../xmmhandler"
)

const (
	//MAXCPUTIME is the maximum time in seconds the process can be scheduled
	MAXCPUTIME uint64 = 2

	//CHARS is a string containing all possible characters in filename
	CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-0123456789."

	//RANDOMBYTES is the number of bytes that will be encoded to the filename in base64.
	//The filename length will be RANDOMBYTES*8/6
	RANDOMBYTES = 12

	//MAXBYTES is the maximum bytes the asm file can use
	MAXBYTES = 30720 //30KiBytes
)

func getRequestedRegisters(requests *cellshandler.XmmRequests, xmmHandler *xmmhandler.XMMHandler, xmmFormat *cellshandler.XMMFormat) models.CellRegisters {
	cellRegisters := models.CellRegisters{}

	for _, request := range *requests {
		// fmt.Println("Request: ", request.PrintFormat)
		if request.PrintFormat == "" {
			request.PrintFormat = xmmFormat.DefaultPrintingFormat[request.XmmNumber]
		}
		// fmt.Println("Final request: ", request.PrintFormat)
		xmmData := models.XMMData{
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

func getChangedRegisters(hiddenRegs *cellshandler.HiddenInCell, oldXmmHandler *xmmhandler.XMMHandler, newXmmHandler *xmmhandler.XMMHandler, xmmFormat *cellshandler.XMMFormat) models.CellRegisters {
	cellRegisters := models.CellRegisters{}

	for index := range oldXmmHandler.Xmm {
		if !containsInt(index, *hiddenRegs) {
			oldXmm := oldXmmHandler.Xmm[index]
			newXmm := newXmmHandler.Xmm[index]
			if !oldXmm.Equals(newXmm) {
				xmmString := "XMM" + strconv.Itoa(index)
				xmmData := models.XMMData{
					XmmID:     xmmString,
					XmmValues: newXmmHandler.GetXMMData(index, xmmFormat.DefaultDataFormat[index], xmmFormat.DefaultPrintingFormat[index])}
				cellRegisters = append(cellRegisters, xmmData)
			}
		}
	}

	return cellRegisters
}

func getXMMRegs(pid int) (xmmhandler.XMMHandler, error) {
	var fpRegs models.FPRegs

	err := utils.GetFPRegs(pid, &fpRegs)
	fmt.Printf("\nAddress fp: %p\n", &fpRegs)

	if err != nil {
		fmt.Println(err)
		return xmmhandler.XMMHandler{}, err
	}
	xmmSlice := fpRegs.XMMSpace[:]

	return xmmhandler.NewXMMHandler(&xmmSlice), err
}

func joinWithPriority(cellRegs1 *models.CellRegisters, cellRegs2 *models.CellRegisters) models.CellRegisters {

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

func cellsLoop(cellsData *cellshandler.CellsData, pid int, xmmFormat *cellshandler.XMMFormat) models.ResponseObj {

	res := models.ResponseObj{CellRegs: make([]models.CellRegisters, 0)}
	cellIndex := 0

	oldXmmHandler, getErr := getXMMRegs(pid)
	if getErr != nil {
		return utils.KillProcess(pid, "Could not get XMM registers.")
	}
	var ws syscall.WaitStatus

	for cellIndex < len(cellsData.Data) {
		newXmmHandler, getErr := getXMMRegs(pid)
		if getErr != nil {
			return utils.KillProcess(pid, "Could not get XMM registers.")
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
			return utils.KillProcess(pid, execErr.Error())
		}

		_, waitErr := syscall.Wait4(pid, &ws, syscall.WALL, nil)

		if waitErr != nil {
			return utils.KillProcess(pid, waitErr.Error())
		}
		if !utils.PidExists(pid) && cellIndex < len(cellsData.Data)-1 {
			return utils.KillProcess(pid, "Something stopped the program.\n")
		}

	}

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())

	if utils.PidExists(pid) {
		aux := utils.KillProcess(pid, "Something went wrong, program did not reach the end.")
		res.ConsoleOut = aux.ConsoleOut
	} else {
		res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())
	}

	return res
}

func getCellsData(req *http.Request) (cellshandler.CellsData, error) {
	cellsData := cellshandler.NewCellsData()

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	return cellsData, decodeErr
}

func printJSONInput(req *http.Request) error {
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = ioutil.ReadAll(req.Body)
	}

	if err != nil {
		return err
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonMap)

	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(jsonMap, "", "\t")

	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
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

	utils.LimitFileSize(syscall.Getpid(), MAXBYTES)

	utils.EnableCors(&w)

	err := printJSONInput(req)

	if err != nil {
		utils.Response(&w, models.ResponseObj{ConsoleOut: "The json received was mega rancio."})
		return
	}
	xmmFormat := cellshandler.NewXMMFormat()
	cellsData, err := getCellsData(req)
	if err != nil {
		utils.Response(&w, models.ResponseObj{ConsoleOut: "Could't read data from the client properly."})
		return
	}
	if cellsData.HandleCellsData(&xmmFormat) {

		utils.Response(&w, models.ResponseObj{ConsoleOut: "Please insert some code."})
		return
	}

	var fileText = cellsData.CellsData2SourceCode()

	execMap, missingPaths := checkExecutables("nasm", "minijail0", "ld", "microjail")

	if len(missingPaths) > 0 {
		responseObj := models.ResponseObj{ConsoleOut: "Could't find next executable paths:"}
		for _, path := range missingPaths {
			responseObj.ConsoleOut += "\n* " + path
		}
		utils.Response(&w, responseObj)
		return
	}

	randomID, randErr := randomString(RANDOMBYTES)

	if randErr != nil {
		utils.Response(&w, models.ResponseObj{ConsoleOut: "Could't create file name properly. (Server error please notify)"})
		fmt.Println("Error creating random file: ", randErr)
		return
	}

	randomFolder := path.Join("/clients", randomID) //Cada usuario tiene una carpeta propia lo cual le va a dar aislamiento
	randomFile := path.Join(randomFolder, randomID)
	os.Mkdir(randomFolder, os.FileMode(0777))
	fileErr := ioutil.WriteFile(randomFile+".asm", []byte(fileText), 0644)

	if fileErr != nil {
		utils.Response(&w, models.ResponseObj{ConsoleOut: "Could't create file properly. Maybe the file is greater than 30Kb."})
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
		res := models.ResponseObj{ConsoleOut: "Text file must not be greater than 30Kb."}
		fmt.Println("File is larger than 30Kb. Aborting.")
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	//NASM con namespace de todo tipo y color
	nasmCmd := exec.Command(execMap["minijail0"], //Path a minijail
		"-p",                            //PID namespace
		"-n",                            //No new priviligies
		"-S", "../policies/nasm.policy", //Setea las policies para NASM
		"-v",               //Vamos a crear un nuevo VFS
		"-P", "/var/empty", //Hacemos un pivot_root a /var/empty
		"-b", fmt.Sprintf("%s,,1", randomFolder), // Bindeamos la carpeta del cliente a si misma con permiso de escritura
		"-b", "/usr/bin/nasm", //Bindiamos esto para tener el binario a NASM
		"-b", "/proc", //Bindiamos /proc
		"-r",                                                                     //Remonta /proc a readonly
		execMap["nasm"], "-f", "elf64", randomFile+".asm", "-o", randomFile+".o") //Comando ejecutador de NASM
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmCmd.Stdout = os.Stdout
	nasmErr := nasmCmd.Start()

	if nasmErr != nil {
		fmt.Println("Error starting NASM: ", nasmErr.Error())
		res := models.ResponseObj{ConsoleOut: nasmErr.Error()}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	nasmPID := nasmCmd.Process.Pid
	utils.LimitCPUTime(nasmPID, MAXCPUTIME)

	nasmErr = nasmCmd.Wait()

	if nasmErr != nil {
		fmt.Println("Error executing nasm: ", nasmErr)
		fmt.Println(stderr.String())
		errorString := strings.ReplaceAll(stderr.String(), randomFile, "output")
		res := models.ResponseObj{ConsoleOut: errorString}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	if !utils.FileExists(randomFile + ".o") {
		fmt.Println("NASM execution finished correctly but didn't create expected file: " + randomFile + ".o")
		res := models.ResponseObj{ConsoleOut: "NASM execution failed."}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
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
		res := models.ResponseObj{ConsoleOut: linkingErr.Error()}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	linkerPID := linkingCmd.Process.Pid
	utils.LimitCPUTime(linkerPID, MAXCPUTIME)

	linkingErr = linkingCmd.Wait()

	if linkingErr != nil {
		fmt.Println("Error executing LD: ", linkingErr)
		fmt.Println(stderr.String())
		errorString := strings.ReplaceAll(stderr.String(), randomFile, "output")
		res := models.ResponseObj{ConsoleOut: errorString}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	if !utils.FileExists(randomFile + ".o") {
		fmt.Println("LD execution finished correctly but didn't create expected file: " + randomFile)
		res := models.ResponseObj{ConsoleOut: "LD execution failed."}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	err = os.Chmod(randomFile, 0111)
	if err != nil {
		fmt.Println("Could not change permissions to executable: " + randomFile)
		res := models.ResponseObj{ConsoleOut: "LD execution failed."}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
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
		res := models.ResponseObj{ConsoleOut: startErr.Error()}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	microjailPID := exeCmd.Process.Pid
	utils.LimitCPUTime(microjailPID, MAXCPUTIME)

	exeCmd.Wait()

	optErr := syscall.PtraceSetOptions(microjailPID, 0x100000|syscall.PTRACE_O_TRACEEXEC) //0x100000 = PTRACE_O_EXITKILL, 0x200000 = PTRACE_O_SUSPEND_SECCOMP

	if optErr != nil {
		res := utils.KillProcess(microjailPID, optErr.Error())
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	//One continue such that the C execve is made
	execErr := syscall.PtraceCont(microjailPID, 0)
	if execErr != nil {
		res := utils.KillProcess(microjailPID, execErr.Error())
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	var ws syscall.WaitStatus
	_, waitErr := syscall.Wait4(microjailPID, &ws, syscall.WALL, nil)
	if waitErr != nil {
		res := utils.KillProcess(microjailPID, waitErr.Error())
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	if !utils.PidExists(microjailPID) {
		res := models.ResponseObj{ConsoleOut: "Microjail error."}
		utils.DeleteFiles(randomFolder, &res)
		utils.Response(&w, res)
		return
	}

	responseObj := cellsLoop(&cellsData, microjailPID, &xmmFormat)

	runtime.UnlockOSThread()
	utils.DeleteFiles(randomFolder, &responseObj)
	fmt.Println(responseObj)
	utils.Response(&w, responseObj)

}

func main() {
	runtime.GOMAXPROCS(1)

	http.HandleFunc("/codeSave", codeSave)

	fmt.Println("===============================")
	fmt.Println("=Server listening on port 8080=")
	fmt.Println("===============================")
	http.ListenAndServe(":8080", nil)

}
