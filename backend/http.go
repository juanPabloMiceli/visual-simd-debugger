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
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"faa"

	"github.com/zchee/kube-timeleap/pkg/ptrace"
	"golang.org/x/sys/unix"
)

//CellData ...
type CellData struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Output string `json:"output"`
}

//CellsData ...
type CellsData struct {
	Data []CellData `json:"CellsData"`
}

//ResponseObj ...
type ResponseObj struct {
	ConsoleOut string
	CellRegs   []unix.PtraceRegs
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func response(w *http.ResponseWriter, obj interface{}) {
	faa.Foo()

	responseJSON, err := json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(responseJSON)
}

func notLastCell(currentIndex int, size int) bool {
	return currentIndex < size-1
}

func notDataCell(index int, hasDataCell bool) bool {
	return (index == 0 && !hasDataCell) || index != 0
}

func handleCellsData(obj *CellsData) (string, bool) {
	var resText string = ""
	var hasDataCell bool = false
	var startText string = "global _start\n"
	startText += "section .text\n"
	startText += "_start:\n"

	for index, cellData := range obj.Data {
		if cellData.Code != "" {
			if index == 0 {
				if strings.Contains(cellData.Code, ";data") {
					cellData.Code = strings.Replace(cellData.Code, ";data", "section .data", 1)
					hasDataCell = true
				} else {
					cellData.Code = startText + cellData.Code
				}
			}

			if strings.Contains(cellData.Code, "int 3") {
				cellData.Code = strings.ReplaceAll(cellData.Code, "int 3", ";int 3")
			}

			if index == 1 && hasDataCell {
				cellData.Code = startText + cellData.Code
			}

			resText += cellData.Code + "\n"

			if notLastCell(index, len(obj.Data)) && notDataCell(index, hasDataCell) {
				resText += "int 3\n"
			}

		}
	}

	return resText, hasDataCell
}

func cellsLoop(cellsData *CellsData, pid int, hasDataCell bool) ResponseObj {

	res := ResponseObj{}

	if hasDataCell {
		res.CellRegs = append(res.CellRegs, unix.PtraceRegs{})
	}

	ptrace.Cont(pid, 0)

	var ws syscall.WaitStatus

	syscall.Wait4(pid, &ws, syscall.WALL, nil)
	for !ws.Exited() {
		var fpregs unix.PtraceRegs
		ptrace.GetRegs(pid, &fpregs)
		res.CellRegs = append(res.CellRegs, fpregs)
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

	cellsData := CellsData{}

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	if decodeErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't read data from the server properly."})
		return
	}

	var fileText string
	var hasDataCell bool
	fileText, hasDataCell = handleCellsData(&cellsData)
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

	//////

	// ptrace.Cont(pid, 0)

	// var ws syscall.WaitStatus

	// syscall.Wait4(pid, &ws, syscall.WALL, nil)

	// var reg unix.PtraceRegs
	// ptrace.GetFPRegs(pid, &reg)
	// fmt.Printf("\n%p\n", &reg)
	// var xmmReg *ptrace.FPRegs

	// xmmReg = (*ptrace.FPRegs)(unsafe.Pointer(&reg))
	// var f float32
	// f = *(*float32)(unsafe.Pointer(&xmmReg.XMMSpace[0]))
	// fmt.Printf("\nBuenas buenas %v\n", f)

	responseObj := cellsLoop(&cellsData, pid, hasDataCell)
	// var responseObj ResponseObj
	// hasDataCell = hasDataCell
	response(&w, responseObj)

}

func main() {
	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
