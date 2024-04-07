package gcsv

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/basic/gtest"
	"github.com/davidforest123/goutil/sys/gtime"
	"testing"
	"time"
)

type (
	user struct {
		Name      string    `csv:"name"`
		Age       int       `csv:"age,omitempty"`
		CreatedAt time.Time `csv:"createat"`
		Address   string    `csv:"address"`
		IsMale    bool      `csv:"ismale"`
	}
)

var (
	csvInput1 = []byte(`
name,age,createat,address
jack,26,2012-04-01T15:00:00Z,beijing

john,,0001-01-01T00:00:00Z,beijing

tom,23,2012-04-01T15:00:00Z,"""here""
and there"`,
	)

	csvInput2 = []byte(`
name,age,createat,address,ismale
jack,26,2012-04-01T15:00:00Z,beijing,true

john,,0001-01-01T00:00:00Z,beijing,false

tom,23,2012-04-01T15:00:00Z,"""here""
and there",true`,
	)
)

func TestMarshalBytes(t *testing.T) {
	var inUsers []user
	var outUsers []user
	inUsers = append(inUsers, user{
		Name:      "tom",
		Age:       5,
		CreatedAt: time.Now(),
		Address:   "tom hotel; tom street\n;C.A., U.S.A.",
		IsMale:    true,
	})
	inUsers = append(inUsers, user{
		Name:      "jerry",
		Age:       4,
		CreatedAt: time.Now(),
		Address:   "\"jerry hotel; jerry street\";L.A., U.S.A.",
		IsMale:    false,
	})

	checkUsersEqual := func(in, out []user) error {
		if len(in) != len(out) {
			return gerrors.New("in len %d != out len %d", len(in), len(out))
		}
		for i := range in {
			if in[i].Name != out[i].Name {
				return gerrors.New("in[%d].Name %s != out[%d].Name %s", i, in[i].Name, i, out[i].Name)
			}
			if in[i].Age != out[i].Age {
				return gerrors.New("in[%d].Age %d != out[%d].Age %d", i, in[i].Age, i, out[i].Age)
			}
			if !gtime.EqualUnixNano(in[i].CreatedAt, out[i].CreatedAt) {
				return gerrors.New("in[%d].CreatedAt %s != out[%d].CreatedAt %s", i, in[i].CreatedAt, i, out[i].CreatedAt)
			}
			if in[i].Address != out[i].Address {
				return gerrors.New("in[%d].Address %s != out[%d].Address %s", i, in[i].Address, i, out[i].Address)
			}
			if in[i].IsMale != out[i].IsMale {
				return gerrors.New("in[%d].IsMale %v != out[%d].IsMale %v", i, in[i].IsMale, i, out[i].IsMale)
			}
		}
		return nil
	}

	outUsers = nil
	buf, err := MarshalBytes(&inUsers, true)
	gtest.Assert(t, gerrors.Wrap(err, "MarshalBytes"))
	err = UnmarshalBytes(buf, &outUsers)
	gtest.Assert(t, gerrors.Wrap(err, "UnmarshalBytes"))
	gtest.Assert(t, checkUsersEqual(inUsers, outUsers))
}

func TestUnmarshalBytes(t *testing.T) {

}

func TestUnmarshalBytes_WithCsvComment(t *testing.T) {
	type (
		user struct {
			Name string `csv:"_id"`
			Age  int
		}
		ReadIdDoc struct {
			Id string `csv:"_id"`
		}
		File struct {
			Docs []ReadIdDoc
		}
		Chunk struct {
			Docs []any
		}
	)

	var users []user
	users = append(users, user{
		Name: "tom",
		Age:  5,
	})
	users = append(users, user{
		Name: "jerry",
		Age:  4,
	})

	buf, err := MarshalBytes(&users, true)
	gtest.Assert(t, err)
	var outAny File
	err = UnmarshalBytes(buf, &outAny.Docs)
	gtest.Assert(t, err)
	for _, v := range outAny.Docs {
		fmt.Println(v)
	}

	var outChunk Chunk
	err = UnmarshalBytes(csvInput2, &outChunk.Docs)
	gtest.Assert(t, err)
	for _, v := range outChunk.Docs {
		fmt.Println(v)
	}

	fmt.Println(MarshalString(outChunk.Docs, true))
}

func TestUnmarshalToTable(t *testing.T) {
	tb, err := UnmarshalToTable(csvInput1)
	gtest.Assert(t, err)
	if tb.Len() != 3 {
		gtest.PrintlnExit(t, "Table len %d != correct len 3", tb.Len())
	}
}

func TestTable_ReadRowWithTypes(t *testing.T) {
	tb, err := UnmarshalToTable(csvInput2)
	gtest.Assert(t, err)
	if tb.Len() != 3 {
		gtest.PrintlnExit(t, "Table len %d != correct len 3", tb.Len())
	}
	rowTypes := map[string]any{}
	rowTypes["name"] = ""
	rowTypes["age"] = 1
	rowTypes["createat"] = time.Now()
	rowTypes["address"] = ""
	result, err := tb.ReadRowWithTypes(2, rowTypes)
	gtest.Assert(t, err)
	fmt.Println(result)
}
