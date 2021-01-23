package xmmhandler

import (
	"encoding/binary"
	"fmt"
	"math"
)

//XmmBytes is the number of bytes inside a XMM register
const XmmBytes = 16

//SizeOfInt16 is the number of bytes inside an int16
const SizeOfInt16 = 2

//SizeOfInt32 is the number of bytes inside an int16
const SizeOfInt32 = 4

//SizeOfInt64 is the number of bytes inside an int16
const SizeOfInt64 = 8

//XmmRegisters is the number of xmm registers in intel x86
const XmmRegisters = 16

//XMM ...
type XMM struct {
	Values []byte
}

//NewXMM creates a new XMM
func NewXMM(p *[]byte) XMM {
	var resXMM = XMM{Values: make([]byte, XmmBytes)}
	slice := *p

	for i := 0; i < XmmBytes; i++ {
		resXMM.Values[i] = slice[i]
	}

	return resXMM
}

//Print prints the values in the xmm register as bytes.
func (xmm XMM) Print() {
	fmt.Println(xmm.Values)
}

//PrintAs prints the values in the xmm register as bytes, words,
//double words or quad words depending in the string received.
//Posible entries: int8, int16, int32, int64, float32, float64
func (xmm XMM) PrintAs(format string) {
	switch format {
	case "int8":
		fmt.Printf("%v", xmm.Values)

	case "int16":
		data := make([]int16, len(xmm.Values)/SizeOfInt16)
		for i := range data {
			data[i] = int16(binary.LittleEndian.Uint16(xmm.Values[i*SizeOfInt16 : (i+1)*SizeOfInt16]))
		}
		fmt.Printf("%v", data)

	case "int32":
		data := make([]int32, len(xmm.Values)/SizeOfInt32)
		for i := range data {
			data[i] = int32(binary.LittleEndian.Uint32(xmm.Values[i*SizeOfInt32 : (i+1)*SizeOfInt32]))
		}
		fmt.Printf("%v", data)

	case "int64":
		data := make([]int64, len(xmm.Values)/SizeOfInt64)
		for i := range data {
			data[i] = int64(binary.LittleEndian.Uint64(xmm.Values[i*SizeOfInt64 : (i+1)*SizeOfInt64]))
		}
		fmt.Printf("%v", data)

	case "float32":
		data := make([]float32, len(xmm.Values)/SizeOfInt32)
		for i := range data {
			data[i] = math.Float32frombits(binary.LittleEndian.Uint32(xmm.Values[i*SizeOfInt32 : (i+1)*SizeOfInt32]))
		}
		fmt.Printf("%v", data)

	case "float64":
		data := make([]float64, len(xmm.Values)/SizeOfInt64)
		for i := range data {
			data[i] = math.Float64frombits(binary.LittleEndian.Uint64(xmm.Values[i*SizeOfInt64 : (i+1)*SizeOfInt64]))
		}
		fmt.Printf("%v", data)

	}

}

//AsInt8 returns a slice with the values in the xmm register as bytes.
func (xmm XMM) AsInt8() []byte {
	return xmm.Values
}

//AsInt16 returns a slice with the values in the xmm register as words.
func (xmm XMM) AsInt16() []int16 {
	data := make([]int16, len(xmm.Values)/SizeOfInt16)
	for i := range data {
		data[i] = int16(binary.LittleEndian.Uint16(xmm.Values[i*SizeOfInt16 : (i+1)*SizeOfInt16]))
	}
	return data
}

//AsInt32 returns a slice with the values in the xmm register as double words.
func (xmm XMM) AsInt32() []int32 {
	data := make([]int32, len(xmm.Values)/SizeOfInt32)
	for i := range data {
		data[i] = int32(binary.LittleEndian.Uint32(xmm.Values[i*SizeOfInt32 : (i+1)*SizeOfInt32]))
	}
	return data
}

//AsInt64 returns a slice with the values in the xmm register as quad words.
func (xmm XMM) AsInt64() []int64 {
	data := make([]int64, len(xmm.Values)/SizeOfInt64)
	for i := range data {
		data[i] = int64(binary.LittleEndian.Uint64(xmm.Values[i*SizeOfInt64 : (i+1)*SizeOfInt64]))
	}
	return data
}

//AsFloat32 returns a slice with the values in the xmm register as simple precision numbers.
func (xmm XMM) AsFloat32() []float32 {
	data := make([]float32, len(xmm.Values)/SizeOfInt32)
	for i := range data {
		data[i] = math.Float32frombits(binary.LittleEndian.Uint32(xmm.Values[i*SizeOfInt32 : (i+1)*SizeOfInt32]))
	}
	return data
}

//AsFloat64 returns a slice with the values in the xmm register as double precision numbers.
func (xmm XMM) AsFloat64() []float64 {
	data := make([]float64, len(xmm.Values)/SizeOfInt64)
	for i := range data {
		data[i] = math.Float64frombits(binary.LittleEndian.Uint64(xmm.Values[i*SizeOfInt64 : (i+1)*SizeOfInt64]))
	}
	return data
}

//XMMHandler has all 16 XMM registers and is created with a pointer to the start of XMM Space.
type XMMHandler struct {
	Xmm []XMM
}

//NewXMMHandler creates a new XMMHandler
func NewXMMHandler(p *[]byte) XMMHandler {
	handlerRes := XMMHandler{Xmm: make([]XMM, XmmRegisters)}
	slice := *p

	for i, _ := range handlerRes.Xmm {
		xmmSlice := slice[16*i : 16*(i+1)]
		handlerRes.Xmm[i] = NewXMM(&xmmSlice)
	}

	return handlerRes
}

//PrintAs print all XMM registers as the type passed by parameter.
func (handler XMMHandler) PrintAs(format string) {

	for i, xmm := range handler.Xmm {
		fmt.Printf("\nXMM%-2v:  ", i)
		xmm.PrintAs(format)
	}
	fmt.Printf("\n")
}

// func main() {
// 	var valores [32]byte

// 	for i := 0; i < 32; i++ {
// 		valores[i] = 2 * byte(i)
// 	}

// 	slice := valores[:]
// 	xmm1 := NewXMM(&slice)
// 	xmm1.PrintAs("int8")
// 	fmt.Println(xmm1.AsInt8())
// 	xmm1.PrintAs("int16")
// 	fmt.Println(xmm1.AsInt16())
// 	xmm1.PrintAs("int32")
// 	fmt.Println(xmm1.AsInt32())
// 	xmm1.PrintAs("int64")
// 	fmt.Println(xmm1.AsInt64())
// 	xmm1.PrintAs("float32")
// 	fmt.Println(xmm1.AsFloat32())
// 	xmm1.PrintAs("float64")
// 	fmt.Println(xmm1.AsFloat64())

// }
