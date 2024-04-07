package gnum

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gany"
	"math"
)

func Min(numTypeSample any) (any, error) {
	ntsType := gany.Type(numTypeSample)
	switch ntsType {
	case "uint8":
		return uint8(0), nil
	case "uint16":
		return uint16(0), nil
	case gany.Type(NewUint24(0)):
		return NewUint24(0), nil
	case "uint32":
		return uint32(0), nil
	case "uint64":
		return uint64(0), nil
	case "uint":
		return uint(0), nil
	case "int8":
		return int8(math.MinInt8), nil
	case "int16":
		return int16(math.MinInt16), nil
	case "int32":
		return int32(math.MinInt32), nil
	case "int64":
		return int64(math.MinInt64), nil
	case "int":
		return math.MinInt, nil
	default:
		return 0, gerrors.New("Unsupported Min type %s", ntsType)
	}
}

func Max(numTypeSample any) (any, error) {
	ntsType := gany.Type(numTypeSample)
	switch ntsType {
	case "uint8":
		return uint8(math.MaxUint8), nil
	case "uint16":
		return uint16(math.MaxUint16), nil
	case gany.Type(NewUint24(0)):
		return MaxUint24, nil
	case "uint32":
		return uint32(math.MaxUint32), nil
	case "uint64":
		return uint64(math.MaxUint64), nil
	case "uint":
		return uint(math.MaxUint), nil
	case "int8":
		return int8(math.MaxInt8), nil
	case "int16":
		return int16(math.MaxInt16), nil
	case "int32":
		return int32(math.MaxInt32), nil
	case "int64":
		return int64(math.MaxInt64), nil
	case "int":
		return math.MaxInt, nil
	case "float32":
		return math.MaxFloat32, nil
	case "float64":
		return math.MaxFloat64, nil
	default:
		return 0, gerrors.New("Unsupported Max type %s", ntsType)
	}
}

func MustMin(numTypeSample any) any {
	result, err := Min(numTypeSample)
	if err != nil {
		panic(gerrors.Wrap(err, "MustMin()"))
	}
	return result
}

func MustMax(numTypeSample any) any {
	result, err := Max(numTypeSample)
	if err != nil {
		panic(gerrors.Wrap(err, "MustMax()"))
	}
	return result
}
