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
	"strings"
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

func notLastCell(currentIndex int, size int) bool {
	return currentIndex < size-1
}

func notInData(index int, hasDataCell bool) bool {
	return (index == 0 && !hasDataCell) || index != 0
}

func handleCellsData(obj *CellsData) string {
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

			if index == 1 && hasDataCell {
				cellData.Code = startText + cellData.Code
			}

			resText += cellData.Code + "\n"

			if notLastCell(index, len(obj.Data)) && notInData(index, hasDataCell) {
				resText += "int 3\n"
			}

		}
	}

	return resText
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	cellsData := CellsData{}

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&cellsData)

	if decodeErr != nil {
		response(&w, responseObj{"Could't read data from the server properly."})
		return
	}

	var fileText string = handleCellsData(&cellsData)

	fileErr := ioutil.WriteFile("output.asm", []byte(fileText), 0644)
	if fileErr != nil {
		response(&w, responseObj{"Could't create file properly."})
		return
	}

	nasmPath, nasmPathErr := exec.LookPath("nasm")

	if nasmPathErr != nil {
		response(&w, responseObj{"Could't find nasm executable path"})
		return
	}

	nasmCmd := exec.Command(nasmPath, "-f", "elf64", "-g", "-F", "DWARF", "output.asm")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil {
		response(&w, responseObj{stderr.String()})
		return
	}

	linkerPath, linkerPathErr := exec.LookPath("ld")

	if linkerPathErr != nil {
		response(&w, responseObj{"Could't find linker executable path"})
		return
	}

	linkingCmd := exec.Command(linkerPath, "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil {
		response(&w, responseObj{stderr.String()})
		return
	}
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		response(&w, responseObj{"Could't find server path"})
		return
	}
	fullPath := path.Join(path.Dir(filename), "output")
	exeCmd := exec.Command(fullPath)
	var out bytes.Buffer
	exeCmd.Stdout = &out
	exeCmd.Stderr = &stderr
	exeErr := exeCmd.Run()

	var responseString string = out.String()
	if responseString != "" {
		responseString += "\n"
	}

	if exeErr != nil {
		responseString += fmt.Sprint(exeErr)
	} else {
		responseString += "exit status 0"
	}

	response(&w, responseObj{responseString})

}

func main() {
	http.HandleFunc("/codeSave", codeSave)

	http.ListenAndServe(":8080", nil)
}
