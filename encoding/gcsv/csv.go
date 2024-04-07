package gcsv

import (
	"bytes"
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"goutil/container/gconv"
	"io"
	"strings"
)

type (
	Table struct {
		Headers []string
		Data    [][]string
		sortIdIdx int
	}
)

// note:
// some tools export bad csv file which has invisible characters at the beginning of csv file like 0x(EFBBBF), if you find that some member can't be Unmarshal, check all members of csv file

func ReadHeaders(in []byte) ([]string, error) {
	newLineIdx := strings.IndexByte(string(in), '\n')
	if newLineIdx <= 0 {
		return nil, gerrors.New("newLineIdx %d invalid", newLineIdx)
	}
	headers := strings.Split(string(in)[:newLineIdx], ",")
	if len(headers) == 0 {
		return nil, gerrors.New("invalid nil old headers")
	}
	return headers, nil
}

func ReplaceHeaders(in []byte, newHeaders []string) ([]byte, error) {
	if len(in) == 0 {
		return nil, gerrors.New("empty param in in ReplaceHeaders")
	}
	if len(newHeaders) == 0 {
		return nil, gerrors.New("invalid nil newHeaders")
	}

	newLineIdx := strings.IndexByte(string(in), '\n')
	if newLineIdx <= 0 {
		return nil, gerrors.New("newLineIdx %d invalid", newLineIdx)
	}

	return append([]byte(strings.Join(newHeaders, ",")), in[newLineIdx:]...), nil
}

// MarshalBytes encode `v` into bytes.
// Note: \r(carriage return) before \n(newline characters) are silently removed,
// for example: string "hello\r\nworld" after csv encodiing will be "hello\nworld".
func MarshalBytes(v any, keepHeader bool) ([]byte, error) {
	if keepHeader {
		return gocsv.MarshalBytes(v)
	} else {
		buffer := bytes.NewBuffer(nil)
		if err := gocsv.MarshalWithoutHeaders(v, buffer); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
}

func MarshalString(v any, keepHeader bool) (string, error) {
	buf, err := MarshalBytes(v, keepHeader)
	return string(buf), err
}

// UnmarshalBytes decode `in` bytes into array pointer `out`.
func UnmarshalBytes(in []byte, outArrPtr any) error {
	if gany.IsPtr2StructSlice(outArrPtr) {
		return gocsv.UnmarshalBytes(in, outArrPtr)
	} else if gany.IsPtr2AnySlice(outArrPtr) {
		out := outArrPtr.(*[]any)
		table, err := UnmarshalToTable(in)
		if err != nil {
			return err
		}
		types := map[string]any{}
		for i := range table.Headers {
			types[table.Headers[i]] = ""
		}
		outMap, err := table.ReadAll(types)
		if err != nil {
			return err
		}
		for _, v := range outMap {
			*out = append(*out, v)
		}
		return nil
	} else {
		return gerrors.New("outArrPtr type only supports pointer to struct slice in UnmarshalBytes")
	}
}

func UnmarshalToTable(in []byte) (*Table, error) {
	r := csv.NewReader(bytes.NewBuffer(in))

	result := &Table{}
	err := error(nil)
	result.Headers, err = r.Read()
	if err != nil {
		if err == io.EOF {
			return result, nil
		} else {
			return nil, err
		}
	}
	for {
		lineRecords, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		result.Data = append(result.Data, lineRecords)
	}

	return result, nil
}

func (t *Table) Verify() error {
	headLen := len(t.Headers)
	for i := range t.Data {
		if len(t.Data[i]) != headLen {
			return gerrors.New("Csv Table headers len %d != len(t.Data[%d])", len(t.Headers), len(t.Data[i]))
		}
	}
	for i := range t.Data {
		_, err := gconv.AnyToDecimal(t.Data[i][0])
		if err != nil {
			return gerrors.Wrap(err, "First column of Csv Table must be a ID number, but when converting it to number, there is an error:")
		}
	}
	return nil
}

func (t *Table) SetSortIdIdx(idx int) error {
	if idx >= t.Width() {
		return gerrors.New("sortIdIdx(%d) >= width(%d)", idx, t.Width())
	}
	t.sortIdIdx = idx
	return nil
}

func (t *Table) Less(i, j int) bool {
	iID, err := gconv.AnyToDecimal(t.Data[i][t.sortIdIdx])
	panic(err)
	jID, err := gconv.AnyToDecimal(t.Data[j][t.sortIdIdx])
	panic(err)
	return iID.LessThan(*jID)
}

func (t *Table) Width() int {
	return len(t.Headers)
}

func (t *Table) Height() int {
	return len(t.Data)
}

func (t *Table) Len() int {
	return len(t.Data)
}

func (t *Table) Swap(i, j int) {
	t.Data[i], t.Data[j] = t.Data[j], t.Data[i]
}

func (t *Table) ReadRowWithTypes(rowIdx int, types map[string]any) (map[string]any, error) {
	// verify params
	if rowIdx >= len(t.Data) {
		return nil, gerrors.New("rowIdx(%d) >= len(t.Data)(%d)", rowIdx, len(t.Data))
	}
	if len(types) != len(t.Headers) {
		return nil, gerrors.New("len(types) %d != len(t.Headers) %d", len(types), len(t.Headers))
	}
	for v := range t.Headers {
		if _, ok := types[t.Headers[v]]; !ok {
			return nil, gerrors.New("header(%s) doesn't exist in types", t.Headers[v])
		}
	}

	result := map[string]any{}
	var allTypes []string
	for k := range types {
		allTypes = append(allTypes, k)
	}

	for j := range t.Data[rowIdx] {
		targetType, ok := types[t.Headers[j]]
		if !ok {
			return nil, gerrors.New("Header(%s) doesn't exist in types(%s)", t.Headers[j], strings.Join(allTypes, ","))
		}
		resItem, err := gconv.StringToAny(t.Data[rowIdx][j], targetType)
		if err != nil {
			return nil, gerrors.Wrap(err, "StringToAny")
		}
		result[t.Headers[j]] = resItem
	}

	return result, nil
}

func (t *Table) ReadAll(types map[string]any) ([]map[string]any, error) {
	result := []map[string]any{}
	for i := 0; i < t.Len(); i++ {
		row, err := t.ReadRowWithTypes(i, types)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}