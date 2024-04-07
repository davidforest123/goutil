package gmnemonic

import (
	"fmt"
	"goutil/basic/gerrors"
	"goutil/basic/gtest"
	"goutil/crypto/gcryptfmt"
	"goutil/crypto/gcrypto"
	"strings"
	"testing"
)

func TestGenerateMasterKey(t *testing.T) {
	bipTestMnemonic := "someone prosper toward lunar human dolphin bulb motion cricket suit favorite fuel all coach sauce cloud dune home peasant trick general fold area meat"
	exceptMasterPrivateKey := "xprv9s21ZrQH143K315v9tUybKQDzDtA5TkMf4Vyj5e6vAgi97zw6AaJYJg1zSgszyccYkRAV1KWtcNC7kh5kL8Fy1dmXXKRnwEYRaPFQQVkGh6" // generated by https://iancoleman.io/bip39/

	masterPrivateKey, masterPublicKey, err := GenerateBIP32MasterKey(bipTestMnemonic, "abc")
	gtest.Assert(t, err)
	fmt.Println("Master public key:", masterPublicKey)

	if masterPrivateKey != exceptMasterPrivateKey {
		gtest.Assert(t, gerrors.New("TestGenerateMasterKey failed"))
	}
}

func TestGenerateChildKeyWithPath(t *testing.T) {
	bipTestMnemonic := "someone prosper toward lunar human dolphin bulb motion cricket suit favorite fuel all coach sauce cloud dune home peasant trick general fold area meat"

	cl := gtest.NewCaseList()

	// except address generated by MyEtherWallet.com
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/0").Expect("0x3137e77647f56d89abce30c048b05f2d8cfb815e")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/1").Expect("0x3cb83dbbc71af2613d724b4f56a8d4174984daeb")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/2").Expect("0x385d80fbf24b47c4a44d982e627bb2706e11739d")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/3").Expect("0xcda97fdd38f643aa40dd6127e241eb8e55147751")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/4").Expect("0x4fa19ea20ee39744a7bb9d39675f7a1e163a058d")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/5").Expect("0xbc2f1dd36eabd84c67d71a9bb46d6eea4637346f")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/6").Expect("0xbf7e56c1f5683174df52683277a2a08491d5f9c5")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/7").Expect("0xf93dcfe5a6a22b55ab9e27c833a0e21dc1d8ad37")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/8").Expect("0xa35078c25419b558d1b7361efa07401002fb72cb")
	cl.New().Input(bipTestMnemonic).Input("abc").Input("m/44'/60'/0'/0/9").Expect("0xf3548f860e7ed291b28018915a6961eefb49b8e3")

	// except address generated by MyEtherWallet.com
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/0").Expect("0xd9a0cedf8792fe65ff5d59427822c587a58aaf70")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/1").Expect("0xa73c4f233c7702063e3a328bfb075202d3cce13c")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/2").Expect("0xd6676e5a914872958676ad8340a31bb03229806b")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/3").Expect("0x9e1c2c27ec1a05cf2eb8c9d4516fda62cdd36f39")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/4").Expect("0x82c6fa41ec0ebbdcfe56f905390058f6573e8184")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/5").Expect("0x080fe54c82cb152318bee78a58a4ec0d405b1da1")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/6").Expect("0xd78e663f4a6a58f2b9753f686c97ecee9e2ef7f4")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/7").Expect("0x04e126b040b48f142fc86dad51d54921bcd4c7d7")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/8").Expect("0x7835f0f5879049fefbdfd49c01a21a3416de3dda")
	cl.New().Input(bipTestMnemonic).Input("").Input("m/44'/60'/0'/0/9").Expect("0x2b9c8ae416eaa1fce7c1891398fa5f70ed72c1c8")

	ec := gcrypto.NewEthCrypt()
	for i := 0; i < len(cl.Get()); i++ {
		testMnemonic := cl.Get()[i].Inputs[0].(string)
		testPassword := cl.Get()[i].Inputs[1].(string)
		testPath := cl.Get()[i].Inputs[2].(string)
		expectAddr := cl.Get()[i].Expects[0].(string)

		childKey, err := GenerateChildKeyWithPath(testMnemonic, testPassword, testPath)
		gtest.Assert(t, err)
		childPrivateKey, err := childKey.PrivateKey.String(gcryptfmt.EthKeyEncode)
		gtest.Assert(t, err)
		childPublicKey, err := childKey.PublicKey.String(gcryptfmt.EthKeyEncode)
		gtest.Assert(t, err)
		gotAddr, err := ec.PublicKeyToAddress(childKey.PublicKey)
		gtest.Assert(t, err)

		if strings.ToLower(gotAddr) != strings.ToLower(expectAddr) {
			gtest.PrintlnExit(t, "password(%s) path(%s) except addr %s but got %s", testPassword, testPath, expectAddr, gotAddr)
		}
		fmt.Println(fmt.Sprintf("pwd %s path %s privateKey %s publicKey %s", testPassword, testPath, childPrivateKey, childPublicKey))
	}
}
