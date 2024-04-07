package gstring

import (
	"goutil/basic/gtest"
	"testing"
)

func TestRegexMatch(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("www.google.com").Input("google.com?").Expect(true)
	cl.New().Input("www.google.com").Input("google?").Expect(true)
	cl.New().Input("www.yahoo.com").Input("google?").Expect(false)
	cl.New().Input("whois.org").Input(".org$").Expect(true)
	cl.New().Input("www.google.com").Input("google.com$").Expect(true)
	cl.New().Input("www.google.com.hk").Input("google.com$").Expect(false)

	for _, v := range cl.Get() {
		s := v.Inputs[0].(string)
		exp := v.Inputs[1].(string)
		expect := v.Expects[0].(bool)
		result, err := RegexMatch(s, exp)
		if err != nil {
			gtest.PrintlnExit(t, "RegexMatch(%s, %s) error: %s", s, exp, err.Error())
		}
		if expect != result {
			gtest.PrintlnExit(t, "RegexMatch(%s, %s) expect %v but get %v", s, exp, expect, result)
		}
	}
}
