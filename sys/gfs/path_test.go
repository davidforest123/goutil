package gfs

import (
	"testing"

	"github.com/davidforest123/goutil/basic/gtest"
)

func TestDirJoinFile(t *testing.T) {
	cl := gtest.NewCaseList()
	//cl.New().Input("/root").Input("C:\\README.md").Expect("C:\\README.md")
	//cl.New().Input("/root/").Input("C:\\README.md").Expect("C:\\README.md")
	cl.New().Input("/usrs/tony").Input("README.md").Expect("/usrs/tony/README.md")
	cl.New().Input("/usrs/tony").Input("../README.md").Expect("/usrs/README.md")
	cl.New().Input("/usrs/tony/").Input("README.md").Expect("/usrs/tony/README.md")
	cl.New().Input("/usrs/tony/").Input("../README.md").Expect("/usrs/README.md")

	//cl.New().Input("C:\\Temp").Input("C:\\README.md").Expect("C:\\README.md")
	//cl.New().Input("C:\\Temp\\").Input("C:\\README.md").Expect("C:\\README.md")
	//cl.New().Input("C:\\Temp").Input("README.md").Expect("C:\\Temp\\README.md")
	//cl.New().Input("C:\\Temp").Input("../README.md").Expect("C:\\README.md")
	//cl.New().Input("C:\\Temp\\").Input("README.md").Expect("C:\\Temp\\README.md")
	//cl.New().Input("C:\\Temp\\").Input("../README.md").Expect("C:\\README.md")

	for _, v := range cl.Get() {
		src := v.Inputs[0].(string)
		dst := v.Inputs[1].(string)
		expect := v.Expects[0].(string)
		result := DirJoinFile(src, dst)
		if result != expect {
			gtest.PrintlnExit(t, "DirJoinFile(%v, %s) expect %s but got %s", src, dst, expect, result)
		}
	}
}
