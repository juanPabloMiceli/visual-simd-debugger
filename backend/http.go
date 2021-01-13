package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"runtime"
)

type cellData struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Output string `json:"output"`
}

type cellsData struct {
	CellsData []cellData `json:"CellsData"`
}

type responseObj struct {
	ConsoleOut string
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

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	_cellsData := cellsData{}

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&_cellsData)

	if decodeErr != nil {
		response(&w, responseObj{"Could't read data from the server properly."})
		return
	}

	fileErr := ioutil.WriteFile("output.asm", []byte(_cellsData.CellsData[0].Code), 0644)
	if fileErr != nil {
		response(&w, responseObj{"Could't create file properly."})
		return
	}

	nasmPath, nasmPathErr := exec.LookPath("nasm")

	if nasmPathErr != nil {
		response(&w, responseObj{"Could't find nasm executable path"})
		fmt.Println("Could't find nasm executable path")
		return
	}

	nasmCmd := exec.Command(nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil {
		response(&w, responseObj{stderr.String()})
		fmt.Println("this:" + fmt.Sprint(nasmErr) + ": " + stderr.String())
		return
	}

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		response(&w, responseObj{"Could't find linker executable path"})
		fmt.Println("Could't find linker executable path")
		return
	}

	linkingCmd := exec.Command(linkerPath, "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil {
		response(&w, responseObj{stderr.String()})
		fmt.Println(fmt.Sprint(linkingErr) + ": " + stderr.String())
		return
	}
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		response(&w, responseObj{"Could't find server path"})
		fmt.Println("Could't find server path")
		return
	}
	fullPath := path.Join(path.Dir(filename), "output")
	exeCmd := exec.Command(fullPath)
	var out bytes.Buffer
	exeCmd.Stdout = &out
	exeCmd.Stderr = &stderr
	exeErr := exeCmd.Run()

	if exeErr != nil {
		response(&w, responseObj{stderr.String()})
		fmt.Println(fmt.Sprint(exeErr) + ": " + stderr.String())
		return
	}

	response(&w, responseObj{out.String()})

	fmt.Println(out.String())
}

func main() {
	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
