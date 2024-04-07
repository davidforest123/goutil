package gany

import (
	serr "errors"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/basic/gtest"
	"reflect"
	"testing"
)

func TestTypeString(t *testing.T) {
	src := "123"
	typeString := Type(src)
	if typeString != "string" {
		t.Errorf("Parse string error, returns %s", typeString)
	}

	typeString = Type(&src)
	if typeString != "*string" {
		t.Errorf("Parse *string error, returns %s", typeString)
	}

	num := 123.456
	typeString = Type(num)
	if typeString != "float64" {
		t.Errorf("Parse float64 error, returns %s", typeString)
	}

	err := serr.New("this is a standard error")
	typeString = Type(err)
	if typeString != "*errors.errorString" {
		t.Errorf("Parse standard error type error, returns %s", typeString)
	}

	err = gerrors.Errorf("this is a extended error")
	typeString = Type(err)
	if typeString != "*gerrors.GErr" {
		t.Errorf("Parse extended error type error, returns %s", typeString)
	}

	err = gerrors.New("this is a extended error too")
	typeString = Type(err)
	if typeString != "*gerrors.GErr" {
		t.Errorf("Parse extended error too type error, returns %s", typeString)
	}

	type myStruct struct{}
	ms := myStruct{}
	typeString = Type(&ms)
	if typeString != "*gany.myStruct" {
		t.Errorf("Parse *myStruct type error, returns %s", typeString)
	}
}

func TestIsSlice(t *testing.T) {
	cl := gtest.NewCaseList()
	n1 := 1
	var a1 any
	s1 := "abc"
	var ss1 []any
	var ss2 []int
	var ss3 []string
	cl.New().Input(1).Expect(false)

	cl.New().Input(nil).Expect(false)
	cl.New().Input(n1).Expect(false)
	cl.New().Input(a1).Expect(false)
	cl.New().Input(s1).Expect(false)
	cl.New().Input(ss1).Expect(true)
	cl.New().Input(ss2).Expect(true)
	cl.New().Input(ss3).Expect(true)
	cl.New().Input(&n1).Expect(false)
	cl.New().Input(&a1).Expect(false)
	cl.New().Input(&s1).Expect(false)
	cl.New().Input(&ss1).Expect(false)
	cl.New().Input(&ss2).Expect(false)
	cl.New().Input(&ss3).Expect(false)

	for _, v := range cl.Get() {
		got := IsSlice(v.Inputs[0])
		expect := v.Expects[0].(bool)
		if got != expect {
			gtest.PrintlnExit(t, "IsSlice(type(%s)) expect %v but got %v", reflect.TypeOf(v.Inputs[0]), expect, got)
		}
	}
}

func TestIsPtr2Slice(t *testing.T) {
	cl := gtest.NewCaseList()
	n1 := 1
	var a1 any
	s1 := "abc"
	var ss1 []any
	var ss2 []int
	var ss3 []string
	cl.New().Input(1).Expect(false)
	cl.New().Input(nil).Expect(false)
	cl.New().Input(n1).Expect(false)
	cl.New().Input(a1).Expect(false)
	cl.New().Input(s1).Expect(false)
	cl.New().Input(ss1).Expect(false)
	cl.New().Input(ss2).Expect(false)
	cl.New().Input(ss3).Expect(false)
	cl.New().Input(&n1).Expect(false)
	cl.New().Input(&a1).Expect(false)
	cl.New().Input(&s1).Expect(false)
	cl.New().Input(&ss1).Expect(true)
	cl.New().Input(&ss2).Expect(true)
	cl.New().Input(&ss3).Expect(true)

	for _, v := range cl.Get() {
		got := IsPtr2Slice(v.Inputs[0])
		expect := v.Expects[0].(bool)
		if got != expect {
			gtest.PrintlnExit(t, "IsPtr2Slice(type(%s)) expect %v but got %v", reflect.TypeOf(v.Inputs[0]), expect, got)
		}
	}
}

func TestIsPtr2StructSlice(t *testing.T) {
	cl := gtest.NewCaseList()
	n1 := 1
	var a1 any
	s1 := "abc"
	var ss1 []any
	var ss2 []struct{}
	var ss3 []string

	cl.New().Input(1).Expect(false)
	cl.New().Input(nil).Expect(false)
	cl.New().Input(n1).Expect(false)
	cl.New().Input(a1).Expect(false)
	cl.New().Input(s1).Expect(false)
	cl.New().Input(ss1).Expect(false)
	cl.New().Input(ss2).Expect(false)
	cl.New().Input(ss3).Expect(false)
	cl.New().Input(&n1).Expect(false)
	cl.New().Input(&a1).Expect(false)
	cl.New().Input(&s1).Expect(false)
	cl.New().Input(&ss1).Expect(false)
	cl.New().Input(&ss2).Expect(true)
	cl.New().Input(&ss3).Expect(false)

	for _, v := range cl.Get() {
		got := IsPtr2StructSlice(v.Inputs[0])
		expect := v.Expects[0].(bool)
		if got != expect {
			gtest.PrintlnExit(t, "IsPtr2StructSlice(type(%s)) expect %v but got %v", reflect.TypeOf(v.Inputs[0]), expect, got)
		}
	}
}

func TestIsPtr2AnySlice(t *testing.T) {
	cl := gtest.NewCaseList()
	n1 := 1
	var a1 any
	s1 := "abc"
	var ss1 []any
	var ss2 []struct{}
	var ss3 []string

	cl.New().Input(1).Expect(false)
	cl.New().Input(nil).Expect(false)
	cl.New().Input(n1).Expect(false)
	cl.New().Input(a1).Expect(false)
	cl.New().Input(s1).Expect(false)
	cl.New().Input(ss1).Expect(false)
	cl.New().Input(ss2).Expect(false)
	cl.New().Input(ss3).Expect(false)
	cl.New().Input(&n1).Expect(false)
	cl.New().Input(&a1).Expect(false)
	cl.New().Input(&s1).Expect(false)
	cl.New().Input(&ss1).Expect(true)
	cl.New().Input(&ss2).Expect(false)
	cl.New().Input(&ss3).Expect(false)

	for _, v := range cl.Get() {
		got := IsPtr2AnySlice(v.Inputs[0])
		expect := v.Expects[0].(bool)
		if got != expect {
			gtest.PrintlnExit(t, "IsPtr2AnySlice(type(%s)) expect %v but got %v", reflect.TypeOf(v.Inputs[0]), expect, got)
		}
	}
}
