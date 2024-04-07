package intrpr

import (
	"context"
	"github.com/d5/tengo/v2"
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"reflect"
)

// call go function from tengo script: https://play.golang.org/p/Zb8hfBAf-WI

type (
	vmTengo struct {
		setValues map[string]any
		prog      *tengo.Compiled
	}
)

func valTengo2Comm(tengoVal reflect.Value) gany.Val {
	return gany.NewVal(tengoVal)
}

func newVMTengo() (*vmTengo, error) {
	res := &vmTengo{}
	res.setValues = map[string]any{}
	return res, nil
}

func (vm *vmTengo) LoadScript(lang Lang, script string) error {
	if lang != LangTengo {
		return gerrors.New("engine yaegj doesn't support language %s", lang)
	}

	s := tengo.NewScript([]byte(script))
	err := error(nil)
	vm.prog, err = s.Compile()
	return err
}

func (vm *vmTengo) LoadScriptAndRun(script string) (gany.Val, error) {
	s := tengo.NewScript([]byte(script))
	for k, v := range vm.setValues {
		if err := s.Add(k, v); err != nil {
			return gany.ValNil, err
		}
	}
	_, err := s.RunContext(context.Background())
	return gany.ValNil, err
}

func (vm *vmTengo) SetVal(name string, value any) error {
	vm.setValues[name] = value
	return nil
}

func (vm *vmTengo) SetFunc(name string, fn any) error {
	vm.setValues[name] = fn
	return nil
}

func (vm *vmTengo) Run() error {
	if vm.prog == nil {
		return gerrors.New("tengo prog not initialized")
	}
	for k, v := range vm.setValues {
		if err := vm.prog.Set(k, v); err != nil {
			return err
		}
	}
	return vm.prog.Run()
}

func (vm *vmTengo) GetFunc(name string) (Callable, error) {
	return nil, gerrors.ErrNotSupport // Note: for now, tengo doesn't support call script function from go
}

func (vm *vmTengo) GetVal(name string) (gany.Val, error) {
	vrb := vm.prog.Get(name)
	return gany.NewVal(vrb.Value()), nil
}
