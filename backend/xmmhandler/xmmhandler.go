package xmmhandler

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

const (
	//XMMBYTES is the number of bytes inside a XMM register
	XMMBYTES = 16

	//XMMREGISTERS is the number of xmm registers in intel x86
	XMMREGISTERS = 16

	//SIZEOFINT16 is the number of bytes inside an int16
	SIZEOFINT16 = 2

	//SIZEOFINT32 is the number of bytes inside an int16
	SIZEOFINT32 = 4

	//SIZEOFINT64 is the number of bytes inside an int16
	SIZEOFINT64 = 8
)

const (
	//INT8STRING is the string that will define we want to work with the xmm values as int 8
	INT8STRING = "v16_int8"

	//INT16STRING is the string that will define we want to work with the xmm values as int 16
	INT16STRING = "v8_int16"

	//INT32STRING is the string that will define we want to work with the xmm values as int 32
	INT32STRING = "v4_int32"

	//INT64STRING is the string that will define we want to work with the xmm values as int 64
	INT64STRING = "v2_int64"

	//FLOAT32STRING is the string that will define we want to work with the xmm values as float 32
	FLOAT32STRING = "v4_float"

	//FLOAT64STRING is the string that will define we want to work with the xmm values as float 64
	FLOAT64STRING = "v2_double"
)

//XMM register is represented as a 16 bytes slice with the corresponding values
type XMM []byte

//NewXMM creates a new XMM
func NewXMM(p *[]byte) XMM {
	var resXMM = XMM{}
	slice := *p

	for i := 0; i < XMMBYTES; i++ {
		resXMM = append(resXMM, slice[i])
	}

	return resXMM
}

//Equals compares two XMM
func (xmm XMM) Equals(newXmm XMM) bool {

	for index := range xmm {
		if xmm[index] != newXmm[index] {
			return false
		}
	}

	return true

}

//Print prints the values in the xmm register as bytes.
func (xmm XMM) Print() {
	fmt.Println(xmm)
}

//PrintAs prints the values in the xmm register as bytes, words,
//double words or quad words depending in the string received.
//Posible entries: int8, int16, int32, int64, float32, float64
func (xmm XMM) PrintAs(format string) {
	switch format {
	case "int8":
		fmt.Printf("%v", xmm)

	case "int16":
		data := make([]int16, len(xmm)/SIZEOFINT16)
		for i := range data {
			data[i] = int16(binary.LittleEndian.Uint16(xmm[i*SIZEOFINT16 : (i+1)*SIZEOFINT16]))
		}
		fmt.Printf("%v", data)

	case "int32":
		data := make([]int32, len(xmm)/SIZEOFINT32)
		for i := range data {
			data[i] = int32(binary.LittleEndian.Uint32(xmm[i*SIZEOFINT32 : (i+1)*SIZEOFINT32]))
		}
		fmt.Printf("%v", data)

	case "int64":
		data := make([]int64, len(xmm)/SIZEOFINT64)
		for i := range data {
			data[i] = int64(binary.LittleEndian.Uint64(xmm[i*SIZEOFINT64 : (i+1)*SIZEOFINT64]))
		}
		fmt.Printf("%v", data)

	case "float32":
		data := make([]float32, len(xmm)/SIZEOFINT32)
		for i := range data {
			data[i] = math.Float32frombits(binary.LittleEndian.Uint32(xmm[i*SIZEOFINT32 : (i+1)*SIZEOFINT32]))
		}
		fmt.Printf("%v", data)

	case "float64":
		data := make([]float64, len(xmm)/SIZEOFINT64)
		for i := range data {
			data[i] = math.Float64frombits(binary.LittleEndian.Uint64(xmm[i*SIZEOFINT64 : (i+1)*SIZEOFINT64]))
		}
		fmt.Printf("%v", data)

	}

}

//AsInt8 returns a slice with the values in the xmm register as bytes.
//Must convert values to int16 because javascript won't recognize bytes as numbers.
func (xmm XMM) AsInt8() []int16 {
	data := make([]int16, len(xmm))
	for i := range data {
		data[i] = int16(xmm[i])
		data[i] = 2
	}
	reverseSlice(data)
	return data
}

//AsInt16 returns a slice with the values in the xmm register as words.
func (xmm XMM) AsInt16() []int16 {
	data := make([]int16, len(xmm)/SIZEOFINT16)
	for i := range data {
		data[i] = int16(binary.LittleEndian.Uint16(xmm[i*SIZEOFINT16 : (i+1)*SIZEOFINT16]))
	}
	reverseSlice(data)
	return data
}

//AsInt32 returns a slice with the values in the xmm register as double words.
func (xmm XMM) AsInt32() []int32 {
	data := make([]int32, len(xmm)/SIZEOFINT32)
	for i := range data {
		data[i] = int32(binary.LittleEndian.Uint32(xmm[i*SIZEOFINT32 : (i+1)*SIZEOFINT32]))
	}
	reverseSlice(data)
	return data
}

//AsInt64 returns a slice with the values in the xmm register as quad words.
func (xmm XMM) AsInt64() []int64 {
	data := make([]int64, len(xmm)/SIZEOFINT64)
	for i := range data {
		data[i] = int64(binary.LittleEndian.Uint64(xmm[i*SIZEOFINT64 : (i+1)*SIZEOFINT64]))
	}
	reverseSlice(data)
	return data
}

//AsFloat32 returns a slice with the values in the xmm register as simple precision numbers.
func (xmm XMM) AsFloat32() []float32 {
	data := make([]float32, len(xmm)/SIZEOFINT32)
	for i := range data {
		data[i] = math.Float32frombits(binary.LittleEndian.Uint32(xmm[i*SIZEOFINT32 : (i+1)*SIZEOFINT32]))
	}
	reverseSlice(data)
	return data
}

//AsFloat64 returns a slice with the values in the xmm register as double precision numbers.
func (xmm XMM) AsFloat64() []float64 {
	data := make([]float64, len(xmm)/SIZEOFINT64)
	for i := range data {
		data[i] = math.Float64frombits(binary.LittleEndian.Uint64(xmm[i*SIZEOFINT64 : (i+1)*SIZEOFINT64]))
	}
	reverseSlice(data)
	return data
}

//XMMHandler has all 16 XMM registers and is created with a pointer to the start of XMM Space.
type XMMHandler struct {
	Xmm []XMM
}

//NewXMMHandler creates a new XMMHandler
func NewXMMHandler(p *[]byte) XMMHandler {
	handlerRes := XMMHandler{Xmm: make([]XMM, XMMREGISTERS)}
	slice := *p

	for i := range handlerRes.Xmm {
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

//GetXMMData will call the corresponding As<format> function given the xmmNumber and the data format desired.
func (handler *XMMHandler) GetXMMData(xmmNumber int, dataFormat string) interface{} {
	switch dataFormat {
	case INT8STRING:
		return handler.Xmm[xmmNumber].AsInt8()
	case INT16STRING:
		return handler.Xmm[xmmNumber].AsInt16()
	case INT32STRING:
		return handler.Xmm[xmmNumber].AsInt32()
	case INT64STRING:
		return handler.Xmm[xmmNumber].AsInt64()
	case FLOAT32STRING:
		return handler.Xmm[xmmNumber].AsFloat32()
	case FLOAT64STRING:
		return handler.Xmm[xmmNumber].AsFloat64()
	default:
		panic("The XMM format is invalid")
	}
}

func reverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	valueLen := value.Len()

	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

// func main() {
// 	var valores [32]byte

// 	for i := 0; i < 32; i++ {
// 		valores[i] = 2 * byte(i)
// 	}

// 	slice := valores[:]
// 	xmm1 := NewXMM(&slice)
// 	fmt.Println(xmm1.AsInt8())
// 	// xmm1.PrintAs(INT8STRING)
// 	// fmt.Println(xmm1.AsInt8())
// 	// xmm1.PrintAs(INT16STRING)
// 	// fmt.Println(xmm1.AsInt16())
// 	// xmm1.PrintAs(INT32STRING)
// 	// fmt.Println(xmm1.AsInt32())
// 	// xmm1.PrintAs(INT64STRING)
// 	// fmt.Println(xmm1.AsInt64())
// 	// xmm1.PrintAs(FLOAT32STRING)
// 	// fmt.Println(xmm1.AsFloat32())
// 	// xmm1.PrintAs(FLOAT64STRING)
// 	// fmt.Println(xmm1.AsFloat64())

// }
