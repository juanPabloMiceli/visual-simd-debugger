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
	ID   int    `json:"id"`
	Code string `json:"code"`
}

type cellsData struct {
	CellsData []cellData `json:"CellsData"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	// req.ParseForm()
	// fmt.Println(req.Form)
	_cellsData := cellsData{}

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&_cellsData)

	if decodeErr != nil {
		fmt.Println("Could't read data from the server properly.")
		w.WriteHeader(500)
		return
	}

	fileErr := ioutil.WriteFile("output.asm", []byte(_cellsData.CellsData[0].Code), 0644)
	if fileErr != nil {
		println("Could't create file properly.")
		w.WriteHeader(500)
		return
	}

	nasmPath, nasmPathErr := exec.LookPath("nasm")

	if nasmPathErr != nil {
		fmt.Println("Could't find nasm executable path")
		w.WriteHeader(500)
		return
	}

	nasmCmd := exec.Command(nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil {
		fmt.Println(fmt.Sprint(nasmErr) + ": " + stderr.String())
		w.WriteHeader(500)
		return
	}

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		fmt.Println("Could't find linker executable path")
		w.WriteHeader(500)
		return
	}

	linkingCmd := exec.Command(linkerPath, "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil {
		fmt.Println(fmt.Sprint(linkingErr) + ": " + stderr.String())
		w.WriteHeader(500)
		return
	}
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		fmt.Println("Could't find server path")
		w.WriteHeader(500)
		return
	}
	fullPath := path.Join(path.Dir(filename), "output")
	exeCmd := exec.Command(fullPath)
	var out bytes.Buffer
	exeCmd.Stdout = &out
	exeCmd.Stderr = &stderr
	exeErr := exeCmd.Run()

	if exeErr != nil {
		fmt.Println(fmt.Sprint(exeErr) + ": " + stderr.String())
		w.WriteHeader(500)
		return
	}

	fmt.Println(out.String())
	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
