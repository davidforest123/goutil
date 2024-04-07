package gnet

import (
	"goutil/basic/gtest"
	"testing"
)

func TestParseIPMask(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("255.255.255.0").Expect(true)
	cl.New().Input("255.255.254.0").Expect(true)
	cl.New().Input("255.255.253.0").Expect(false)
	cl.New().Input("255.255.252.0").Expect(true)
	cl.New().Input("254.255.255.0").Expect(false)

	for _, v := range cl.Get() {
		s := v.Inputs[0].(string)
		expect := v.Expects[0].(bool)
		_, err := ParseIPMask(s)
		if expect != (err == nil) {
			gtest.PrintlnExit(t, "test fail for ParseIPMask(%s), expect %v but got %v", s, expect, err)
		}
	}
}
