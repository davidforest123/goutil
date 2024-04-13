package gfs

import (
	"testing"

	"github.com/davidforest123/goutil/basic/gtest"
)

func TestDirJoinFile(t *testing.T) {
	cl := gtest.NewCaseList()

	cl.New().Input("/usrs/tony").Input("README.md").Expect("/usrs/tony/README.md").Expect(true)
	cl.New().Input("/usrs/tony").Input("../README.md").Expect("/usrs/README.md").Expect(true)
	cl.New().Input("/usrs/tony/").Input("README.md").Expect("/usrs/tony/README.md").Expect(true)
	cl.New().Input("/usrs/tony/").Input("../README.md").Expect("/usrs/README.md").Expect(true)

	cl.New().Input("C:\\Temp").Input("README.md").Expect("C:\\Temp\\README.md").Expect(true)
	cl.New().Input("C:\\Temp").Input("..\\README.md").Expect("C:\\README.md").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("README.md").Expect("C:\\Temp\\README.md").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("..\\README.md").Expect("C:\\README.md").Expect(true)

	cl.New().Input("C:\\Temp").Input("C:\\README.md").Expect("").Expect(false)
	cl.New().Input("C:\\Temp\\").Input("C:\\README.md").Expect("").Expect(false)
	cl.New().Input("/root").Input("/boot/README.md").Expect("").Expect(false)
	cl.New().Input("/root/").Input("/boot/README.md").Expect("").Expect(false)
	cl.New().Input("/root").Input("C:\\README.md").Expect("").Expect(false)
	cl.New().Input("/root/").Input("C:\\README.md").Expect("").Expect(false)
	cl.New().Input("C:\\Temp\\").Input("../README.md").Expect("").Expect(false)
	cl.New().Input("C:\\Temp").Input("../README.md").Expect("").Expect(false)

	for _, v := range cl.Get() {
		src := v.Inputs[0].(string)
		dst := v.Inputs[1].(string)
		expect := v.Expects[0].(string)
		expectOk := v.Expects[1].(bool)
		result, resultErr := DirJoinFile(src, dst)
		if result != expect {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect %s but got (%s, %s)", src, dst, expect, result, resultErr.Error())
		}
		if expectOk && resultErr != nil {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect ok but got error %s", src, dst, resultErr.Error())
		}
		if !expectOk && resultErr == nil {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect err but got ok", src, dst)
		}
	}
}

func TestDirJoinDir(t *testing.T) {
	cl := gtest.NewCaseList()

	cl.New().Input("/usrs/tony").Input("DOCS").Expect("/usrs/tony/DOCS/").Expect(true)
	cl.New().Input("/usrs/tony").Input("../DOCS").Expect("/usrs/DOCS/").Expect(true)
	cl.New().Input("/usrs/tony/").Input("DOCS").Expect("/usrs/tony/DOCS/").Expect(true)
	cl.New().Input("/usrs/tony/").Input("../DOCS").Expect("/usrs/DOCS/").Expect(true)
	cl.New().Input("/usrs/tony/").Input("DOCS/ENG").Expect("/usrs/tony/DOCS/ENG/").Expect(true)
	cl.New().Input("/usrs/tony/").Input("../DOCS/ENG").Expect("/usrs/DOCS/ENG/").Expect(true)

	cl.New().Input("C:\\Temp").Input("DOCS").Expect("C:\\Temp\\DOCS\\").Expect(true)
	cl.New().Input("C:\\Temp").Input("..\\DOCS").Expect("C:\\DOCS\\").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("DOCS").Expect("C:\\Temp\\DOCS\\").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("..\\DOCS").Expect("C:\\DOCS\\").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("DOCS\\").Expect("C:\\Temp\\DOCS\\").Expect(true)
	cl.New().Input("C:\\Temp\\").Input("..\\DOCS\\").Expect("C:\\DOCS\\").Expect(true)

	cl.New().Input("C:\\Temp").Input("C:\\DOCS").Expect("").Expect(false)
	cl.New().Input("C:\\Temp\\").Input("C:\\DOCS").Expect("").Expect(false)
	cl.New().Input("/root").Input("/boot/DOCS").Expect("").Expect(false)
	cl.New().Input("/root/").Input("/boot/DOCS").Expect("").Expect(false)
	cl.New().Input("/root").Input("C:\\DOCS").Expect("").Expect(false)
	cl.New().Input("/root/").Input("C:\\DOCS").Expect("").Expect(false)
	cl.New().Input("C:\\Temp\\").Input("../DOCS").Expect("").Expect(false)
	cl.New().Input("C:\\Temp").Input("../DOCS").Expect("").Expect(false)

	for _, v := range cl.Get() {
		src := v.Inputs[0].(string)
		dst := v.Inputs[1].(string)
		expect := v.Expects[0].(string)
		expectOk := v.Expects[1].(bool)
		result, resultErr := DirJoinDir(src, dst)
		if result != expect {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect %s but got (%s, %s)", src, dst, expect, result, resultErr.Error())
		}
		if expectOk && resultErr != nil {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect ok but got error %s", src, dst, resultErr.Error())
		}
		if !expectOk && resultErr == nil {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect err but got ok", src, dst)
		}
	}
}
