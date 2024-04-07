package gbase

import (
	"encoding/base64"
	"math/big"
	"strings"
)

const (
// btcAlphabet is the modified base58 alphabet used by Bitcoin.
 btcAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

 flickrAlphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var (
 bigRadix = big.NewInt(58)
 bigZero = big.NewInt(0)
)

func Base64Encode(message []byte) string {
	return base64.RawStdEncoding.EncodeToString(message)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}

func IsBase58BtcStyle(s string) bool {
	for _, v := range s {
		if !strings.Contains(btcAlphabet, string(v)) {
			return false
		}
	}
	return true
}

// Base58BtcStyleEncode encodes a byte slice to a modified base58 string, using BTCAlphabet
func Base58BtcStyleEncode(message []byte) string {
	return encodeAlphabet(message, btcAlphabet)
}

// Base58BtcStyleDecode decodes a modified base58 string to a byte slice, using BTCAlphabet
func Base58BtcStyleDecode(s string) ([]byte, error) {
	return decodeAlphabet(s, btcAlphabet), nil
}

// DecodeAlphabet decodes a modified base58 string to a byte slice, using alphabet.
func decodeAlphabet(b, alphabet string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	for i := len(b) - 1; i >= 0; i-- {
		tmp := strings.IndexAny(alphabet, string(b[i]))
		if tmp == -1 {
			return []byte("")
		}
		idx := big.NewInt(int64(tmp))
		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		answer.Add(answer, tmp1)
		j.Mul(j, bigRadix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabet[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// encodeAlphabet encodes a byte slice to a modified base58 string, using alphabet
func encodeAlphabet(b []byte, alphabet string) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabet[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}
