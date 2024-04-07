package intrpr

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"testing"
)

func TestVm_SetVal(t *testing.T) {
	vms := map[string]Vm{}
	vm1, err := NewVM("goja")
	if err != nil {
		gtest.PrintlnExit(t, err.Error())
	}
	vms["goja"] = vm1

	vm2, err := NewVM("yaegi")
	if err != nil {
		gtest.PrintlnExit(t, err.Error())
	}
	vms["yaegi"] = vm2

	cl := gtest.NewCaseList()
	cl.New().Input("goja").Input(LangJavaScript).Input(`
	function scriptSum() {
		return a + b;
	}`).Input("a").Input(40).Input("b").Input(2).Expect("42")

	cl.New().Input("yaegi").Input(LangGo).Input(`
	func scriptSum() int {
		return a + b
	}`).Input("a").Input(40).Input("b").Input(2).Expect("42")

	cl.New().Input("goja").Input(LangJavaScript).Input(`
	function scriptSum() {
		return a + b;
	}`).Input("a").Input(50).Input("b").Input(3).Expect("53")

	for _, cl := range cl.Get() {
		engine := cl.Inputs[0].(string)
		lang := cl.Inputs[1].(Lang)
		script := cl.Inputs[2].(string)
		n1 := cl.Inputs[3].(string)
		v1 := cl.Inputs[4].(int)
		n2 := cl.Inputs[5].(string)
		v2 := cl.Inputs[6].(int)
		expect := cl.Expects[0].(string)

		vm, exist := vms[engine]
		if !exist {
			gtest.PrintlnExit(t, "vm engine %s not found", engine)
		}

		// set function
		err = vm.SetVal(n1, v1)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}
		err = vm.SetVal(n2, v2)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// load script and run
		err = vm.LoadScript(lang, script)
		if err != nil {
			gtest.PrintlnExit(t, "LoadScript for engine %s error: %s", engine, err.Error())
		}
		err = vm.Run()
		if err != nil {
			gtest.PrintlnExit(t, "Run for engine %s error: %s", engine, err.Error())
		}

		// register script function
		scriptSum, err := vm.GetFunc("scriptSum")
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// call script function
		res, err := scriptSum()
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		numStr := res[0].String()
		if numStr != expect {
			gtest.PrintlnExit(t, "engine %s result should be %s but not %s", engine, expect, numStr)
		}
	}
}

func TestVm_SetFunc(t *testing.T) {
	goSum := func(a, b int) int {
		return a + b
	}

	cl := gtest.NewCaseList()
	cl.New().Input("goja").Input(LangJavaScript).Input(`
	function scriptSum(a, b) {
		return goSum(a, b);
	}`).Input(40).Input(2).Expect("42")

	cl.New().Input("yaegi").Input(LangGo).Input(`
	func scriptSum(a int, b int) int {
		return goSum(a, b)
	}`).Input(40).Input(2).Expect("42")

	for _, cl := range cl.Get() {
		engine := cl.Inputs[0].(string)
		lang := cl.Inputs[1].(Lang)
		script := cl.Inputs[2].(string)
		n1 := cl.Inputs[3].(int)
		n2 := cl.Inputs[4].(int)
		expect := cl.Expects[0].(string)

		vm, err := NewVM(engine)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// set function
		err = vm.SetFunc("goSum", goSum)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// load script and run
		err = vm.LoadScript(lang, script)
		if err != nil {
			gtest.PrintlnExit(t, "LoadScript for engine %s error: %s", engine, err.Error())
		}
		err = vm.Run()
		if err != nil {
			gtest.PrintlnExit(t, "Run for engine %s error: %s", engine, err.Error())
		}

		// register script function
		scriptSum, err := vm.GetFunc("scriptSum")
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// call script function
		res, err := scriptSum(n1, n2)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		numStr := res[0].String()
		if numStr != expect {
			gtest.PrintlnExit(t, "engine %s result should be %s but not %s", engine, expect, numStr)
		}
	}
}

type (
	ctx struct{}

	bar struct{ Field int }
)

func (f *ctx) Bar(field int) *bar {
	return &bar{Field: field}
}

func (b *bar) Double() int {
	return b.Field * 2
}

func TestVm_GetFunc(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("goja").Input(LangJavaScript).Input(`
	function onReply(a, b) {
		return ctx.Bar(a + b).Double();
	}`).Input(40).Input(2).Expect("84")

	cl.New().Input("yaegi").Input(LangGo).Input(`
	func onReply(a int, b int) int {
		return ctx.Bar(a + b).Double()
	}`).Input(40).Input(2).Expect("84")

	for _, cl := range cl.Get() {
		engine := cl.Inputs[0].(string)
		lang := cl.Inputs[1].(Lang)
		script := cl.Inputs[2].(string)
		n1 := cl.Inputs[3].(int)
		n2 := cl.Inputs[4].(int)
		expect := cl.Expects[0].(string)

		vm, err := NewVM(engine)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// set context
		err = vm.SetVal("ctx", &ctx{})
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// load script and run
		err = vm.LoadScript(lang, script)
		if err != nil {
			gtest.PrintlnExit(t, "LoadScript for engine %s error: %s", engine, err.Error())
		}
		err = vm.Run()
		if err != nil {
			gtest.PrintlnExit(t, "Run for engine %s error: %s", engine, err.Error())
		}

		// register script function
		onReply, err := vm.GetFunc("onReply")
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// call script function
		res, err := onReply(n1, n2)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		numStr := res[0].String()
		if numStr != expect {
			gtest.PrintlnExit(t, "engine %s result should be %s but not %s", engine, expect, numStr)
		}
	}
}

func TestNewVm_RunTypeScript(t *testing.T) {
	vm, err := NewVM("goja")
	gtest.Assert(t, err)
	script := `
	class Student {
	name:string;
	age:number;
	}`
	fmt.Println(vm.LoadScript(LangTypeScript, script))
	fmt.Println(vm.Run())
}
