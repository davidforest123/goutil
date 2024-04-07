package gcryptfmt

import (
	"bytes"
	"encoding/hex"
	"goutil/basic/gerrors"
	"math/big"
)

const btcBase58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type (
	KeyEncode func(keyBytes []byte) (string, error)
	KeyDecode func(keyString string) ([]byte, error)
)

func BtcKeyEncode(data []byte) (string, error) {
	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(data)
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(btcBase58Alphabet[mod.Int64()]) + encoded
	}

	return encoded, nil
}

func BtcKeyDecode(data string) ([]byte, error) {
	decimalData := new(big.Int)
	alphabetBytes := []byte(btcBase58Alphabet)
	multiplier := big.NewInt(58)

	for _, value := range data {
		pos := bytes.IndexByte(alphabetBytes, byte(value))
		if pos == -1 {
			return nil, gerrors.New("Character `%v` not found in alphabet", byte(value))
		}
		decimalData.Mul(decimalData, multiplier)
		decimalData.Add(decimalData, big.NewInt(int64(pos)))
	}

	return decimalData.Bytes(), nil
}

func EthKeyEncode(keyBytes []byte) (string, error) {
	return hex.EncodeToString(keyBytes), nil
}

func EthKeyDecode(keyString string) ([]byte, error) {
	return hex.DecodeString(keyString)
}
