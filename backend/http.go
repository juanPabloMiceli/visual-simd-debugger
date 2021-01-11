package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type codeData struct {
	CodeText string `json:"CodeText"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func codeSave(w http.ResponseWriter, req *http.Request) {

	enableCors(&w)

	_codeData := codeData{}

	dec := json.NewDecoder(req.Body)

	dec.DisallowUnknownFields()

	decodeErr := dec.Decode(&_codeData)

	if decodeErr != nil {
		fmt.Println("Could't read data from the server properly.")
		w.WriteHeader(500)
		return
	}

	fileErr := ioutil.WriteFile("output.asm", []byte(_codeData.CodeText), 0644)
	if fileErr != nil {
		println("Could't create file properly.")
		w.WriteHeader(500)
		return
	}

	nasmCmd := exec.Command("/usr/bin/nasm", "-f", "elf64", "-g", "-F", "DWARF", "output.asm")
	var stderr bytes.Buffer
	nasmCmd.Stderr = &stderr
	nasmErr := nasmCmd.Run()

	if nasmErr != nil {
		fmt.Println(fmt.Sprint(nasmErr) + ": " + stderr.String())
		w.WriteHeader(500)
		return
	}

	linkingCmd := exec.Command("/usr/bin/ld", "-o", "output", "output.o")
	linkingCmd.Stderr = &stderr
	linkingErr := linkingCmd.Run()

	if linkingErr != nil {
		fmt.Println(fmt.Sprint(linkingErr) + ": " + stderr.String())
		w.WriteHeader(500)
		return
	}

	exeCmd := exec.Command("/home/juampi/probando_go/backend/output")
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
