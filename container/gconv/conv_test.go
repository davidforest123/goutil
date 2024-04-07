package gconv

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"testing"
)

func TestBytesToBinString(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input([]byte{254, 255, 255, 0}).Input("").Expect("11111110111111111111111100000000")

	for _, v := range cl.Get() {
		b := v.Inputs[0].([]byte)
		delimiter := v.Inputs[1].(string)
		expect := v.Expects[0].(string)
		result := BytesToBinString(b, delimiter)
		if result != expect {
			gtest.PrintlnExit(t, "BytesToBinString(%v, %s) expect %s but got %s", b, delimiter, expect, result)
		}
	}
}

func TestBytesToDecString(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input([]byte{254, 255, 255, 0}).Input(".").Input(0).Expect("254.255.255.0")

	for _, v := range cl.Get() {
		b := v.Inputs[0].([]byte)
		delimiter := v.Inputs[1].(string)
		makeUpZeros := v.Inputs[2].(int)
		expect := v.Expects[0].(string)
		result := BytesToDecString(b, delimiter, makeUpZeros)
		if result != expect {
			gtest.PrintlnExit(t, "BytesToBinString(%v, %s, %d) expect %s but got %s", b, delimiter, makeUpZeros, expect, result)
		}
	}
}

func TestBytesToHexString(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input([]byte{254, 255, 255, 0}).Input(".").Input(true).Expect("fe.ff.ff.00")

	for _, v := range cl.Get() {
		b := v.Inputs[0].([]byte)
		delimiter := v.Inputs[1].(string)
		makeUpZeros := v.Inputs[2].(bool)
		expect := v.Expects[0].(string)
		result := BytesToHexString(b, delimiter, makeUpZeros)
		if result != expect {
			gtest.PrintlnExit(t, "BytesToBinString(%v, %s, %v) expect %s but got %s", b, delimiter, makeUpZeros, expect, result)
		}
	}
}

func TestStringToAny(t *testing.T) {
	fmt.Println(StringToAny("123", string("")))
}
