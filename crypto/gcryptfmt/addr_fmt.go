package gcryptfmt

import (
	"encoding/hex"
	"github.com/davidforest123/goutil/basic/gerrors"
	"golang.org/x/crypto/sha3"
	"strconv"
	"strings"
)

/**
BTC addr types
a) Legacy/P2PKH/PayToPubKeyHash: starts with '1'
b) SegWitNested/NestedSegWit/SegWit/P2SH/PayToScriptHash: starts with '3'
c) SegWitBech32/NativeSegWit/Bech32: starts with 'bc1'
*/

// IMPORTANT:
// https://github.com/modood/btckeygen 参考这个项目，应该是把BTC的地址玩明白了
// BTC地址测试用例可以从下列地址获取：
// https://kimbatt.github.io/btc-address-generator/?page=mnemonic-seed
// https://iancoleman.io/bip39/
// https://iancoleman.io/bitcoin-key-compression/

type (
	AddrEncode func(addrBytes []byte) (string, error)
	AddrDecode func(addrString string) ([]byte, error)
)

func ethChecksumAddress(address string) (string, error) {
	address = strings.ToLower(address)
	address = strings.Replace(address, "0x", "", 1)
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(address))
	hash := sha.Sum(nil)
	hashStr := hex.EncodeToString(hash)
	if len(hashStr) < len(address) {
		return "", gerrors.New("address(%s) is invalid ETH address", address)
	}
	var result []string
	for i, v := range address {
		res, err := strconv.ParseInt(string(hashStr[i]), 16, 64)
		if err != nil {
			return "", err
		}
		if res > 7 {
			result = append(result, strings.ToUpper(string(v)))
			continue
		}
		result = append(result, string(v))
	}
	return strings.Join(result, ""), nil
}

func EthAddrEncode(keyBytes []byte) (string, error) {
	s := hex.EncodeToString(keyBytes)
	if len(s) == 40 { // is a ETH address
		err := error(nil)
		if s, err = ethChecksumAddress(s); err != nil {
			return "", err
		}
	}
	return "0x" + s, nil
}

func EthAddrDecode(keyString string) ([]byte, error) {
	return hex.DecodeString(keyString)
}
