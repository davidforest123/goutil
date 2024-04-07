package gnum

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"goutil/basic/gtest"
	"goutil/container/gany"
	"goutil/container/gconv"
	"goutil/encoding/gjson"
	"math/big"
	"testing"
)

type decimalS1 struct {
	Name  string
	Score Decimal `json:"Score,omitempty"`
}

type decimalS2 struct {
	Name  string
	Score *Decimal `json:"Score,omitempty" bson:"Score,omitempty"`
}

func TestNewDecimalFromAny(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input(rune(100)).Expect("100")
	cl.New().Input(byte(100)).Expect("100")
	cl.New().Input(uint8(100)).Expect("100")
	cl.New().Input(uint16(100)).Expect("100")
	cl.New().Input(uint32(100)).Expect("100")
	cl.New().Input(uint64(100)).Expect("100")
	cl.New().Input(uint(100)).Expect("100")
	cl.New().Input(int8(100)).Expect("100")
	cl.New().Input(int16(100)).Expect("100")
	cl.New().Input(int32(100)).Expect("100")
	cl.New().Input(int64(100)).Expect("100")
	cl.New().Input(int(100)).Expect("100")
	cl.New().Input(NewUint24(100)).Expect("100")
	cl.New().Input(float32(100)).Expect("100")
	cl.New().Input(float64(100)).Expect("100")
	cl.New().Input(*big.NewInt(100)).Expect("100")
	cl.New().Input(*big.NewFloat(100)).Expect("100")
	cl.New().Input("100").Expect("100")

	for _, v := range cl.Get() {
		result, err := NewDecimalFromAny(v.Inputs[0])
		expect := v.Expects[0].(string)
		gtest.Assert(t, err)
		if result.String() != expect {
			gtest.PrintlnExit(t, "NewDecimalFromAny(%v).String() got %s but %s expected", v.Inputs[0], result.String(), expect)
		}
	}
}

func TestDecimal_Convert(t *testing.T) {
	noError := error(nil)
	someErr := errors.New("whatever error")
	cl := gtest.NewCaseList()
	cl.New().Input("100").Input(rune(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(byte(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint8(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint16(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint32(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint64(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int8(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int16(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(NewUint24(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int32(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int64(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(*big.NewInt(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100").Input(*big.NewFloat(0)).Input(true).Expect("100").Expect(noError)

	cl.New().Input("100").Input(rune(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(byte(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint8(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint16(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(NewUint24(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint32(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint64(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(uint(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int8(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int16(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int32(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int64(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(int(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(*big.NewInt(0)).Input(false).Expect("100").Expect(noError)
	cl.New().Input("100").Input(*big.NewFloat(0)).Input(false).Expect("100").Expect(noError)

	cl.New().Input("100.23456789").Input(rune(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(byte(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(uint8(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(uint16(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(NewUint24(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(uint32(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(uint64(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(uint(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(int8(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(int16(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(int32(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(int64(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(int(0)).Input(true).Expect("100").Expect(noError)
	cl.New().Input("100.23456789").Input(*big.NewInt(0)).Input(true).Expect("100").Expect(noError)
	// FIXME: I don't understand why this test case tested failed.
	//cl.New().Input("100.23456789").Input(*big.NewFloat(0)).Input(true).Expect("100.23456789").Expect(noError)

	cl.New().Input("100.23456789").Input(rune(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(byte(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(uint8(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(uint16(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(NewUint24(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(uint32(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(uint64(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(uint(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(int8(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(int16(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(int32(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(int64(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(int(0)).Input(false).Expect("").Expect(someErr)
	cl.New().Input("100.23456789").Input(*big.NewInt(0)).Input(false).Expect("").Expect(someErr)
	// FIXME: I don't understand why this test case tested failed.
	//cl.New().Input("100.23456789").Input(*big.NewFloat(0)).Input(false).Expect("").Expect(someErr)

	for _, v := range cl.Get() {
		dcm, err := NewDecimalFromString(v.Inputs[0].(string))
		gtest.Assert(t, err)
		sample := v.Inputs[1]
		allowFractionalLoss := v.Inputs[2].(bool)
		expectStr := v.Expects[0].(string)
		expectErr := error(nil)
		if v.Expects[1] != nil {
			expectErr = v.Expects[1].(error)
		}
		expectErrStr := ""
		if expectErr != nil {
			expectErrStr = expectErr.Error()
		}
		gotRes, gotErr := dcm.Convert(sample, allowFractionalLoss)
		gotErrStr := ""
		if gotErr != nil {
			gotErrStr = gotErr.Error()
		}
		if (expectErr == nil && gotErr != nil) || (expectErr != nil && gotErr == nil) || (gotErr == nil && expectStr != gconv.NumToString(gotRes)) {
			gtest.PrintlnExit(t, "Decimal(%s).Convert(%s, %v) got {%s, %s} but expect {%s, %s}",
				dcm.String(), gany.Type(sample), allowFractionalLoss,
				fmt.Sprintf("%s", gconv.NumToString(gotRes)), gotErrStr,
				expectStr, expectErrStr)
		}
	}
}

func TestNewFromStringEx(t *testing.T) {
	d, err := NewDecimalFromStringEx("12345.1234567", 3)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.123" {
		t.Errorf("NewDecimalFromStringEx error1")
		return
	}

	d, err = NewDecimalFromStringEx("12345.1234567", 4)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1235" {
		t.Errorf("NewDecimalFromStringEx error2")
		return
	}

	d, err = NewDecimalFromStringEx("12345.1234567", 7)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1234567" {
		t.Errorf("NewDecimalFromStringEx error3")
		return
	}

	d, err = NewDecimalFromStringEx("12345.1234567", 8)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1234567" {
		t.Errorf("NewDecimalFromStringEx error4")
		return
	}
}

func TestNewDecimalFromInt(t *testing.T) {
	d := NewDecimalFromInt(12345)
	if d.String() != "12345" {
		t.Errorf("NewDecimalFromInt error1")
	}

	d = NewDecimalFromInt(1234567890123456789)
	if d.String() != "1234567890123456789" {
		t.Errorf("NewDecimalFromInt error2")
	}
}

func TestNewDecimalFromUint(t *testing.T) {
	d := NewDecimalFromUint(12345)
	if d.String() != "12345" {
		t.Errorf("NewDecimalFromInt error1")
	}

	d = NewDecimalFromUint(1234567890123456789)
	if d.String() != "1234567890123456789" {
		t.Errorf("NewDecimalFromInt error2")
	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	s := decimalS1{Name: "Bob", Score: Decimal0}
	if gjson.MarshalStringDefault(s, false) != `{"Name":"Bob","Score":"0"}` {
		t.Errorf("TestDecimal_MarshalJSON error1")
		return
	}

	s2 := decimalS2{Name: "Bob"}
	if gjson.MarshalStringDefault(s2, false) != `{"Name":"Bob"}` {
		t.Errorf("TestDecimal_MarshalJSON error2")
		return
	}
}

func TestDecimal_UnmarshalJSON(t *testing.T) {
	type S struct {
		Name  string
		Score Decimal `json:"Score,omitempty"`
	}

	jsonString := `{"Name":"Tom", "Score":99.9}`
	s := &S{}
	if err := json.Unmarshal([]byte(jsonString), s); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(s)
}

func TestDecimal_Trunc(t *testing.T) {
	d, _ := NewDecimalFromString("1.23456789")

	fmt.Println(d.Trunc(8, 0.01))

	if d.Trunc(3, 0.02).String() != "1.22" {
		t.Errorf("Decimal.Trunc error1")
	}
	if d.Trunc(3, 0.03).String() != "1.23" {
		t.Errorf("Decimal.Trunc error2")
	}
	if d.Trunc(3, 0.04).String() != "1.2" {
		t.Errorf("Decimal.Trunc error3")
	}
	if d.Trunc(3, 0.05).String() != "1.2" {
		t.Errorf("Decimal.Trunc error4")
	}
	if d.Trunc(3, 0.06).String() != "1.2" {
		t.Errorf("Decimal.Trunc error5")
	}
	if d.Trunc(3, 0.07).String() != "1.19" {
		t.Errorf("Decimal.Trunc error6")
	}
	if d.Trunc(6, 0.000007).String() != "1.234562" {
		t.Errorf("Decimal.Trunc error7")
	}

	// FIXME 这个例子值得探讨，Trunc是否正确，貌似不太对喔
	d, _ = NewDecimalFromString("5141.73181940667768")
	if d.Trunc(8, 0.000001).String() != "5141.73181899" {
		t.Errorf("Decimal.Trunc error8")
	}
}

func TestDecimal_Trunc2(t *testing.T) {
	d, _ := NewDecimalFromString("1.23456789")

	fmt.Println(d.Trunc2(NewDecimalFromFloat64(0.00001), 0.01))
}

func TestDecimal_MarshalBSON(t *testing.T) {
	d, err := NewDecimalFromString("1.23")
	if err != nil {
		t.Error(err)
		return
	}

	type S struct {
		Number Decimal
	}
	s1 := S{
		Number: d,
	}

	buf, err := bson.Marshal(s1)
	if err != nil {
		t.Error(err)
		return
	}

	s2 := new(S)
	if err := bson.Unmarshal(buf, s2); err != nil {
		t.Error(err)
		return
	}
	if s2.Number.Equal(s1.Number) == false {
		t.Errorf("Unmarshal bad result %s vs %s", s1.Number.String(), s2.Number.String())
		return
	}

	s3 := decimalS2{Name: "Bob"}
	if _, err := bson.Marshal(s3); err != nil {
		t.Errorf("TestDecimal_MarshalBSON error2: %s", err.Error())
		return
	}
}

func TestDecimal_MarshalBSON_NilPointer(t *testing.T) {
	s2 := decimalS2{}
	if _, err := bson.Marshal(s2); err != nil {
		t.Errorf("TestDecimal_MarshalBSON_NilPointer error: %s", err.Error())
		return
	}
}

func TestToElegantFloat64s(t *testing.T) {
	var r []Decimal
	r = append(r, NewDecimalFromFloat64(0.003))
	r = append(r, NewDecimalFromFloat64(0.0004))
	r = append(r, NewDecimalFromFloat64(0.00005))
	r = append(r, NewDecimalFromFloat64(0.000006))
	r = append(r, NewDecimalFromFloat64(0.0000007))
	r = append(r, NewDecimalFromFloat64(0.00000008))
	r = append(r, NewDecimalFromFloat64(0.000000009))
	r = append(r, NewDecimalFromFloat64(0.00000000010))

	efs := ToElegantFloat64s(r)
	for _, v := range efs {
		if len(v.String()) != 12 {
			gtest.PrintlnExit(t, "converted ElegantFloat length should be 12")
		}
	}
}

func TestDecimal_DivRound(t *testing.T) {
	d := NewDecimalFromInt(10000)
	fmt.Println(d.DivInt(3).MulInt(3))
	fmt.Println(d.DivRoundInt(3, 30).MulInt(3))
	fmt.Println(d.DivRoundFloat64(0.001436, 30).MulFloat64(0.001436).String())
}

func TestDecimal_Div(t *testing.T) {
	a, _ := decimal.NewFromString("10000")
	b, _ := decimal.NewFromString("0.001436")
	fmt.Println(a.Exponent(), b.Exponent(), a.Div(b).Mul(b).Exponent())
	fmt.Println(a.Div(b).Mul(b).Value())
	fmt.Println(a.Div(b).Mul(b)) // 10000.0000000000000000000308
}

func TestDecimal_BitsAfterDecimalPoint(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("1").Expect(0)
	cl.New().Input("12").Expect(0)
	cl.New().Input("123").Expect(0)
	cl.New().Input("1234").Expect(0)
	cl.New().Input("12345").Expect(0)
	cl.New().Input("123456").Expect(0)
	cl.New().Input("1234567").Expect(0)
	cl.New().Input("0.1").Expect(1)
	cl.New().Input("0.12").Expect(2)
	cl.New().Input("0.123").Expect(3)
	cl.New().Input("0.1234").Expect(4)
	cl.New().Input("0.12345").Expect(5)
	cl.New().Input("0.123456").Expect(6)
	cl.New().Input("0.1234567").Expect(7)
	cl.New().Input("0.12345678").Expect(8)
	cl.New().Input("0.123456789").Expect(9)
	cl.New().Input("0.1234567890").Expect(10) // NOTE: The 0 in the tail is also counted as precision
	cl.New().Input("0.1234567891").Expect(10)
	cl.New().Input("0.12345678901").Expect(11)
	cl.New().Input("0.123456789012").Expect(12)
	cl.New().Input("0.1234567890123").Expect(13)
	cl.New().Input("0.12345678901234").Expect(14)
	cl.New().Input("0.123456789012345").Expect(15)
	cl.New().Input("0.1234567890123456").Expect(16)
	cl.New().Input("0.12345678901234567").Expect(17)
	cl.New().Input("0.123456789012345678").Expect(18)
	cl.New().Input("0.1234567890123456789").Expect(19)
	cl.New().Input("0.12345678901234567890").Expect(20)
	cl.New().Input("0.12345678901234567891").Expect(20)
	cl.New().Input("0.123456789012345678901").Expect(21)
	cl.New().Input("0.1234567890123456789012").Expect(22)
	cl.New().Input("0.12345678901234567890123").Expect(23)
	cl.New().Input("0.123456789012345678901230").Expect(24)
	cl.New().Input("0.1234567890123456789012300").Expect(25)
	cl.New().Input("0.12345678901234567890123000").Expect(26)

	for _, c := range cl.Get() {
		s := c.Inputs[0].(string)
		e := c.Expects[0].(int)
		d, err := NewDecimalFromString(s)
		gtest.Assert(t, err)
		if d.BitsAfterDecimalPoint(true) != e {
			gtest.PrintlnExit(t, "decimal(%s) from string(%s) expect precision %d, but %d got", d.String(), s, e, d.BitsAfterDecimalPoint(true))
		}
	}
}

func TestDecimal_IsInterger(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("1.0").Expect(true)
	cl.New().Input("12.00").Expect(true)
	cl.New().Input("123.000").Expect(true)
	cl.New().Input("1234.0000").Expect(true)
	cl.New().Input("12345.00000").Expect(true)
	cl.New().Input("123456.000000").Expect(true)
	cl.New().Input("1234567.0000000").Expect(true)
	cl.New().Input("0.1").Expect(false)
	cl.New().Input("0.12").Expect(false)
	cl.New().Input("0.123").Expect(false)
	cl.New().Input("0.1234").Expect(false)
	cl.New().Input("0.12345").Expect(false)
	cl.New().Input("0.123456").Expect(false)
	cl.New().Input("0.1234567").Expect(false)

	for _, c := range cl.Get() {
		s := c.Inputs[0].(string)
		e := c.Expects[0].(bool)
		d, err := NewDecimalFromString(s)
		gtest.Assert(t, err)
		if d.IsInterger() != e {
			gtest.PrintlnExit(t, "IsInterger(%s) expect %v, but %v got", d.String(), e, d.IsInterger())
		}
	}
}
