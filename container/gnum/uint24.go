package gnum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

const (
	// MaxUint24 maximum value of uint24 variable
	MaxUint24 = 1<<24 - 1
)

// Uint24 type
type Uint24 [3]uint8

// NewUint24 creates a new uint24 type
func NewUint24(val uint32) Uint24 {
	if val > MaxUint24 {
		panic(fmt.Sprintf("val %d > uint24 max limit %d", val, MaxUint24))
	}

	n := Uint24{}
	n[0] = uint8(val & 0xFF)
	n[1] = uint8((val >> 8) & 0xFF)
	n[2] = uint8((val >> 16) & 0xFF)
	return n
}

func NewUint24FromBytes(b [3]byte) Uint24 {
	n := Uint24{}
	n[0] = b[0]
	n[1] = b[1]
	n[2] = b[2]
	return n
}

// Uint32 converts uint24 to uint32
func (u Uint24) Uint32() uint32 {
	return uint32(u[0]) | uint32(u[1])<<8 | uint32(u[2])<<16
}

// String converts uint24 to string
func (u Uint24) String() string {
	return strconv.Itoa(int(u.Uint32()))
}

// ToBytes converts uint24 to bytes array
func (u Uint24) ToBytes() []byte {
	var buf = &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, u); err != nil {
		return nil
	}
	return buf.Bytes()
}