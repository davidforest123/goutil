package gconv

import (
	"fmt"
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"goutil/container/gnum"
	"goutil/sys/gtime"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// convert between built-in types: string, bytes, number......

// BytesToBinString convert bytes to binary string
// Used to verify Net Mask
// delimiter: separator between bytes.
// Sample: convert []byte{254, 255, 255, 0} to "11111110111111111111111100000000"
func BytesToBinString(b []byte, delimiter string) string {
	var result string
	for i := 0; i < len(b); i++ {
		for j := 1; j <= 8; j++ {
			if b[i]<<(j-1)>>7 == 1 {
				result += "1"
			} else {
				result += "0"
			}
		}
		if i < len(b)-1 {
			result += delimiter
		}
	}
	return result
}

// BytesToDecString convert bytes to decimal string
// Used to format IP, Net Mask, MAC
// delimiter: separator between bytes.
// Sample: convert []byte{254, 255, 255, 0} to "254255255000"
func BytesToDecString(b []byte, delimiter string, makeUpZeroes int) string {
	var result string
	for i := 0; i < len(b); i++ {
		if makeUpZeroes == 3 {
			result += fmt.Sprintf("%03d", b[i])
		} else if makeUpZeroes == 2 {
			result += fmt.Sprintf("%02d", b[i])
		} else {
			result += fmt.Sprintf("%d", b[i])
		}
		if i < len(b)-1 {
			result += delimiter
		}
	}
	return result
}

// BytesToHexString convert bytes to hex string
// Used to format IP, Net Mask, MAC
// delimiter: separator between bytes.
// Sample: convert []byte{254, 255, 255, 0} to "feffff00"
func BytesToHexString(b []byte, delimiter string, makeUpZeroes bool) string {
	var result string
	for i := 0; i < len(b); i++ {
		if makeUpZeroes {
			result += fmt.Sprintf("%02x", b[i])
		} else {
			result += fmt.Sprintf("%x", b[i])
		}
		if i < len(b)-1 {
			result += delimiter
		}
	}
	return result
}

func AnyToDecimal(val any) (*gnum.Decimal, error) {
	return gnum.NewDecimalFromAny(val)
}

// TODO
func AnyToString(val any) (string, error) {
	return "", nil
}

func NumToString(num any) string {
	switch v := num.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case gnum.Uint24:
		return strconv.FormatUint(uint64(v.Uint32()), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 6, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64)
	case big.Int:
		return v.String()
	case big.Float:
		return v.String()
	case big.Rat:
		return v.String()
	default:
		return fmt.Sprintf("Unsupported NumToString(%s): !(%T=%+v)", gany.Type(num), num, num)
	}
	return ""
}

func StringToAny(s string, convSample any) (any, error) {
	resType := gany.Type(convSample)
	var res any
	err := error(nil)
	switch resType {
	case "bool":
		if strings.ToLower(s) == "true" {
			res = true
		} else if strings.ToLower(s) == "false" {
			res = false
		} else {
			return nil, gerrors.New("Can't convert string(%s) to bool in StringToAny", s)
		}
	case gany.Type(time.Time{}):
		res, err = gtime.ParseTimeStringStrict(s)
	case "rune":
		if len(s) != 4 {
			return nil, gerrors.New("Can't convert string(%s) to rune in StringToAny", s)
		}
		res = strings.Map(func(r rune) rune {
			return r
		}, s)
	case "byte":
		if len(s) != 1 {
			return nil, gerrors.New("Can't convert string(%s) to byte in StringToAny", s)
		}
		res = s[0]
	case "uint", "uint8", "uint16", gany.Type(gnum.NewUint24(0)), "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"float32", "float64", gany.Type(big.Int{}), gany.Type(big.Float{}), gany.Type(gnum.Decimal{}):
		dcml, err := gnum.NewDecimalFromString(s)
		if err != nil {
			return nil, err
		}
		return dcml.Conv(convSample, false)
	case "string":
		res = s
	default:
		return nil, gerrors.New("unsupported convSample type %s in StringToAny", resType)
	}

	return res, err
}
