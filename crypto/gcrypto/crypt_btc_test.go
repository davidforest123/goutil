package gcrypto

import (
	"encoding/hex"
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"github.com/davidforest123/goutil/crypto/gcryptfmt"
	"testing"
)

// https://github.com/modood/btckeygen 这里有多中地址类型的测试用例

func TestBtcCrypt_Verify(t *testing.T) {
	cl := gtest.NewCaseList()

	// addr type Bech32, generated by kimbatt.github.io/btc-address-generator/
	cl.New().Input("Kx4AqA7oBmU5AgWZnT29bcMj7tgcubDJx4baQcrBqYwo6kPYJGS3").Input("bc1qy4vagt63cqx50zfhy39a527xyfus4f6yn6uuay")

	// addr type Legacy, generated by kimbatt.github.io/btc-address-generator/
	cl.New().Input("KxEinRi8ykepw5RdSYfBht6zS2SpgAX5pWhA9EHs36ogJew218ZK").Input("136chcVixvjAK7XR741PMwyiiJ9Q3w5Gfj")

	// addr type Segwit, generated by kimbatt.github.io/btc-address-generator/
	cl.New().Input("L1FaiAuUD9dBNR999gY4RBvNZt1L4VRA5W9n8M5rtavq8R98Nv17").Input("3CbfBKa3u9R1g1wRdGdAriy18KmrEmAp49")

	// generated by bitaddress.org
	cl.New().Input("L2nWh2np6GjMK9zvbvb3Ywq1ZnKSrYCZK9563y5VUfst6oeakSxh").Input("1HdebB6WmX16fKWjQ4W5momCFhUQCQdwKT")

	bc := NewBtcCrypt()
	data, err := gcryptfmt.BtcKeyDecode("L2nWh2np6GjMK9zvbvb3Ywq1ZnKSrYCZK9563y5VUfst6oeakSxh")
	gtest.Assert(t, err)
	public, err := bc.PrivateKeyToPublicKey(data)
	fmt.Println(hex.EncodeToString(public.Bytes()))
	fmt.Println(bc.PublicKeyToAddress(public, "legacy"))
}
