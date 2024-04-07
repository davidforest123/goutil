package intrpr

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"reflect"
)

type (
	vmYaegi struct {
		vmYaegi   *interp.Interpreter
		customLib map[string]map[string]reflect.Value
		prog      *interp.Program
		retVal    gany.Val
	}
)

func valYaegi2Comm(yaegiVal reflect.Value) gany.Val {
	return gany.NewVal(yaegiVal)
}

func newVMYaegi() (*vmYaegi, error) {
	res := &vmYaegi{vmYaegi: interp.New(interp.Options{})}
	res.customLib = make(map[string]map[string]reflect.Value)
	res.customLib["custom/custom"] = make(map[string]reflect.Value)

	if err := res.vmYaegi.Use(stdlib.Symbols); err != nil {
		return nil, err
	}

	return res, nil
}

func (vm *vmYaegi) LoadScript(lang Lang, script string) error {
	if lang != LangGo {
		return gerrors.New("engine yaegi doesn't support language %s", lang)
	}
	err := error(nil)
	vm.prog, err = vm.vmYaegi.Compile(script)
	return err
}

func (vm *vmYaegi) LoadScriptAndRun(script string) (gany.Val, error) {
	val, err := vm.vmYaegi.Eval(script)
	if err != nil {
		return gany.ValNil, err
	}
	return valYaegi2Comm(val), nil
}

// TODO 确定重复调用这个接口，比如映射两个不同的name、value对， 是可以正常执行的
// 引入的指针和相关的自定义类型（非基础类型）都可以存取，但是默认不可以声明
// 如果需要在脚本中主动声明Go Runtime中的自定义类型，目前知道的方法是在脚本中import该类型所在的包，
func (vm *vmYaegi) SetVal(name string, value any) error {
	vm.customLib["custom/custom"][name] = reflect.ValueOf(value)
	if err := vm.vmYaegi.Use(vm.customLib); err != nil {
		return err
	}

	vm.customLib["custom/custom"]["ctx"] = reflect.ValueOf(value)
	if err := vm.vmYaegi.Use(vm.customLib); err != nil {
		return err
	}

	if _, err := vm.vmYaegi.Eval(`import . "custom"`); err != nil {
		return err
	}
	return nil
}

func (vm *vmYaegi) SetFunc(name string, fn any) error {
	return vm.SetVal(name, fn)
}

func (vm *vmYaegi) Run() error {
	if vm.vmYaegi == nil {
		return gerrors.New("yaegi vm not initialized")
	}
	if vm.prog == nil {
		return gerrors.New("yaegi prog not initialized")
	}

	vm.retVal = gany.ValNil
	retVal, err := vm.vmYaegi.Execute(vm.prog)
	if err == nil {
		vm.retVal = valYaegi2Comm(retVal)
	}
	return err
}

func (vm *vmYaegi) GetFunc(name string) (Callable, error) {
	fn, err := vm.vmYaegi.Eval(name)
	if err != nil {
		return nil, err
	}

	return func(args ...any) ([]gany.Val, error) {
		var argsSlice []reflect.Value
		for _, item := range args {
			argsSlice = append(argsSlice, reflect.ValueOf(item))
		}

		resVals := fn.Call(argsSlice)
		var res []gany.Val
		for _, item := range resVals {
			res = append(res, gany.NewVal(item))
		}
		return res, nil
	}, nil
}

func (vm *vmYaegi) GetVal(name string) (gany.Val, error) {
	val, exsit := vm.vmYaegi.Globals()[name]
	if !exsit {
		return gany.ValNil, gerrors.New("value %s doesn't exist", name)
	}
	return valYaegi2Comm(val), nil
}
