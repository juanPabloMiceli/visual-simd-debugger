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
	"syscall"
	"unsafe"

	"github.com/zchee/kube-timeleap/pkg/ptrace"
	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/cellshandler"
	"gitlab.com/juampi_miceli/visual-simd-debugger/backend/xmmhandler"
	"golang.org/x/sys/unix"
)

//ResponseObj ...
type ResponseObj struct {
	ConsoleOut string
	CellRegs   []map[string]interface{} //Could be a slice of any of int or float types
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

func cellsLoop(cellsData *cellshandler.CellsData, pid int) ResponseObj {

	res := ResponseObj{CellRegs: make([]map[string]interface{}, 0)}

	if cellsData.HasDataCell {
		res.CellRegs = append(res.CellRegs, make(map[string]interface{}))
	}

	ptrace.Cont(pid, 0)

	var ws syscall.WaitStatus

	syscall.Wait4(pid, &ws, syscall.WALL, nil)
	var unixRegs unix.PtraceRegs
	for !ws.Exited() {
		ptrace.GetFPRegs(pid, &unixRegs)

		fmt.Printf("\nAddress: %p\n", &unixRegs)

		fpPointer := (*ptrace.FPRegs)(unsafe.Pointer(&unixRegs))
		xmmSlice := fpPointer.XMMSpace[:]
		xmmHandler := xmmhandler.NewXMMHandler(&xmmSlice)

		fmt.Println("\nNueva celda")
		xmmHandler.PrintAs("float64")
		// res.CellRegs = append(res.CellRegs, xmmHandler.Xmm[0], xmmHandler.Xmm[1])
		ptrace.Cont(pid, 0)
		syscall.Wait4(pid, &ws, syscall.WALL, nil)
	}

	ptrace.GetFPRegs(pid, &unixRegs)
	fmt.Printf("\nAddress: %p\n", &unixRegs)
	fpPointer := (*ptrace.FPRegs)(unsafe.Pointer(&unixRegs))
	xmmSlice := fpPointer.XMMSpace[:]
	xmmHandler := xmmhandler.NewXMMHandler(&xmmSlice)

	fmt.Println("\nNueva celda")
	xmmHandler.PrintAs("float64")
	// res.CellRegs = append(res.CellRegs, xmmHandler.Xmm[0], xmmHandler.Xmm[1])

	fmt.Printf("Exited: %v\n", ws.Exited())
	fmt.Printf("Exited status: %v\n", ws.ExitStatus())
	res.ConsoleOut = "Exited status: " + strconv.Itoa(ws.ExitStatus())

	return res
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	cellsData := cellshandler.NewCellsData()

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	if decodeErr != nil {
		response(&w, ResponseObj{ConsoleOut: "Could't read data from the server properly."})
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
