package gnum

import (
	"goutil/basic/gerrors"
	"math/big"
	"strconv"
)

func Format(num any, base int) (string, error) {
	switch v := num.(type) {
	case int:
		return strconv.FormatInt(int64(v), base), nil
	case int8:
		return strconv.FormatInt(int64(v), base), nil
	case int16:
		return strconv.FormatInt(int64(v), base), nil
	case int32:
		return strconv.FormatInt(int64(v), base), nil
	case int64:
		return strconv.FormatInt(v, base), nil
	case uint:
		return strconv.FormatUint(uint64(v), base), nil
	case uint8:
		return strconv.FormatUint(uint64(v), base), nil
	case uint16:
		return strconv.FormatUint(uint64(v), base), nil
	case uint32:
		return strconv.FormatUint(uint64(v), base), nil
	case uint64:
		return strconv.FormatUint(v, base), nil
	case float32: // untested
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64: // untested
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case big.Int:
		return BaseConvert(v.String(), 10, base)
	default:
		return "", gerrors.New("Unsupported type")
	}
	return "", gerrors.New("Unknown error")
}

func FormatInt(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FormatUint8(u uint8) string {
	return strconv.FormatUint(uint64(u), 10)
}

func FormatUint16(u uint16) string {
	return strconv.FormatUint(uint64(u), 10)
}

func FormatUint32(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

func FormatUint64(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func FormatFloat64(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}
