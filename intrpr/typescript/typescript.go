package typescript

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/clarkmcc/go-typescript/utils"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/grand"
	"github.com/dop251/goja"
	"strings"
)

//go:embed typescript_source/v4.7.2.js
var tsSrcV4x7x2 string

var (
	allTs       = map[string]*goja.Program{}
	latestTsVer = "4.7.2"
)

func init() {
	programV4x7x2, err := goja.Compile("", tsSrcV4x7x2, true)
	if err != nil {
		panic(fmt.Errorf("compiling registered source for tag '%s': %w", "4.7.2", err))
	}
	allTs["4.7.2"] = programV4x7x2
}

func TranspileToJavaScript(ts string, tsVer string) (string, error) {
	if tsVer == "" {
		tsVer = latestTsVer
	}
	tsProg, exist := allTs[tsVer]
	if !exist {
		return "", gerrors.New("typescript version %s not found", tsVer)
	}
	if tsProg == nil {
		return "", gerrors.New("typescript program %s is nil", tsVer)
	}

	rt := goja.New()
	_, err := rt.RunProgram(tsProg)
	if err != nil {
		return "", fmt.Errorf("running typescript compiler: %w", err)
	}
	compileOptions := map[string]interface{}{
		"module": "none",
	}
	optionBytes, err := json.Marshal(compileOptions)
	if err != nil {
		return "", fmt.Errorf("marshalling compile options: %w", err)
	}
	decoderName := grand.RandomString(24)
	rt.Set(decoderName, utils.ErrorWrapper(rt, func(call goja.FunctionCall) (interface{}, error) {
		bs, err := base64.StdEncoding.DecodeString(call.Argument(0).String())
		if err != nil {
			return nil, err
		}
		return string(bs), nil
	}))
	transpileCmd := fmt.Sprintf("ts.transpile(%s('%s'), %s, /*fileName*/ undefined, /*diagnostics*/ undefined, /*moduleName*/ \"%s\")",
		decoderName, base64.StdEncoding.EncodeToString([]byte(ts)), optionBytes, "default")
	value, err := rt.RunString(transpileCmd)
	if err != nil {
		return "", fmt.Errorf("running compiler: %w", err)
	}
	return strings.TrimSuffix(value.String(), "\r\n"), nil
}
