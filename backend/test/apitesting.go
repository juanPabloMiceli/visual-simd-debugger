package main

import (
	"fmt"
)

func main() {
	fmt.Print("Tests no implementados")
}

// //CODESAVEURL is the url where the submitCode request is handled
// const CODESAVEURL = "http://localhost:8080/codeSave"

// //XMMData ...
// type XMMData struct {
// 	XmmID     string
// 	XmmValues interface{}
// }

// //CellRegisters contains the different XMMData in a cell.
// type CellRegisters []XMMData

// //ResponseObj ...
// type ResponseObj struct {
// 	ConsoleOut string
// 	CellRegs   []CellRegisters //Could be a slice of any of int or float types
// }

// //Publish send post requests with a json input to url input
// func Publish(url string, bodyString string) (*http.Response, error) {
// 	rawIn := json.RawMessage(bodyString)
// 	body, _ := rawIn.MarshalJSON()
// 	fmt.Println(string(body))
// 	resp, err := http.Post(url, "application/json", bytes.NewReader(body))

// 	if err != nil {
// 		return resp, err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return resp, fmt.Errorf("server didn’t respond 200 OK: %s", resp.Status)
// 	}

// 	return resp, nil
// }

// func main() {
// 	var bodyString string

// 	bodyString = fmt.Sprintf(`
// 	{
// 		"CellsData": [
// 			{
// 				"code": ";data\np: dd 11, 12, 12, 4",
// 				"id": 0,
// 				"output": []
// 			},
// 			{
// 				"code": "movdqu xmm3, [p]\n;print xmm3.v4_int32",
// 				"id": 1,
// 				"output": []
// 			}
// 		]
// 	}
// 	`)

// 	resp, _ := Publish(CODESAVEURL, bodyString)

// 	// buf := new(bytes.Buffer)
// 	// buf.ReadFrom(resp.Body)
// 	// str := buf.String()

// 	// fmt.Println([]byte(str))

// 	// var jsonMap map[string]interface{}
// 	// json.Unmarshal([]byte(str), &jsonMap)

// 	// fmt.Println(str)
// 	// // fmt.Println(jsonMap["CellRegs"])
// 	// hm := jsonMap["CellRegs"].([] interface{})
// 	// fmt.Println(reflect.TypeOf(jsonMap["CellRegs"]))

// 	var respObj ResponseObj

// 	dec := json.NewDecoder(resp.Body)

// 	dec.DisallowUnknownFields()

// 	decodeErr := dec.Decode(&respObj)

// 	if decodeErr != nil {
// 		panic(decodeErr)
// 	}

// 	fmt.Println(respObj.CellRegs[1][0].XmmID)
// }

// func TestOneRegisterInBytes(t *testing.T) {

// 	expectedConsoleOut := "Exited status: 0"
// 	expectedXmmID := "XMM0"
// 	// expectedXmmValues := make([]int, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
// 	expectedConsoleOut = expectedConsoleOut
// 	expectedXmmID = expectedXmmID
// 	// expectedXmmValues = expectedXmmValues

// 	var bodyString string

// 	bodyString = fmt.Sprintf(`
// 	{
// 		"CellsData": [
// 			{
// 				"code": ";data\np: db 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16",
// 				"id": 0,
// 				"output": []
// 			},
// 			{
// 				"code": "movdqu xmm0, [p]",
// 				"id": 1,
// 				"output": []
// 			}
// 		]
// 	}
// 	`)

// 	resp, _ := Publish(CODESAVEURL, bodyString)

// 	var respObj ResponseObj

// 	dec := json.NewDecoder(resp.Body)

// 	dec.DisallowUnknownFields()

// 	decodeErr := dec.Decode(&respObj)

// 	if decodeErr != nil {
// 		t.Errorf("actual decodeErr = %v; expected decodeErr = <nil>")
// 	}

// 	if respObj.ConsoleOut != "" {
// 		a := 0
// 		a = a + 1
// 	}

// 	fmt.Println(respObj.CellRegs[1][0].XmmID)

// }
