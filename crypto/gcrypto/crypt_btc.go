package gcrypto

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/crypto/gbase"
	"github.com/davidforest123/goutil/crypto/gecdsa"
	"github.com/davidforest123/goutil/crypto/ghash"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"strings"
)

type (
	btcCrypt struct{}
)

const (
	btcPrivateKeyBytesLen = 32
)

/*
// convert private/public binary key to HEX format
func (bk btcKey) String() string {
	result := ""
	for i := 0; i < len(bk.data); i++ {
		result += fmt.Sprintf("%02X", bk.data[i])
	}
	return result
}*/

func NewBtcCrypt() *btcCrypt {
	return &btcCrypt{}
}

/*
// 把src数组转换成指定长度的数组，长度不够则添加0
// :param size: 要返回的数组长度
//
//	:param dst: byte类型的切片，需要返回的切片
//	:param src: 原byte数组
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}


// GenerateKeyPair generates key pair by a random string.
// randomKey: random string or mnemonic words, length must greater than 36
func (bc *btcCrypt) GenerateKeyPair(randomKey string) (*KeyPair, error) {
	var curve elliptic.Curve

	// verify randomKey length
	rkLen := len(randomKey)
	if rkLen > 521/8+8 { // > 73
		curve = elliptic.P521()
	} else if rkLen > 384/8+8 { // > 56
		curve = elliptic.P384()
	} else if rkLen > 256/8+8 { // > 40
		curve = elliptic.P256()
	} else if rkLen > 224/8+8 { // > 36
		curve = elliptic.P224()
	} else { // <= 36
		err := gerrors.New("randomKey length %d is too short. It mast be longer than 36 bytes.", rkLen)
		return nil, err
	}

	// generate key
	key, err := ecdsa.GenerateKey(curve, strings.NewReader(randomKey))
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0, btcPrivateKeyBytesLen)
	privateKey := paddedAppend(btcPrivateKeyBytesLen, b, key.D.Bytes())
	publicKey := append(key.PublicKey.X.Bytes(), key.Y.Bytes()...)

	return &KeyPair{
		PrivateKey: NewKey(privateKey),
		PublicKey:  NewKey(publicKey),
	}, nil
}*/

// PrivateKeyToPublicKey generates hex string public key from hex string private key.
func (bc *btcCrypt) PrivateKeyToPublicKey(privateKey gecdsa.Key) (publicKey gecdsa.Key, err error) {
	key, err := loadECDSAFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	publicKeyBuf := append(key.PublicKey.X.Bytes(), key.Y.Bytes()...)
	return gecdsa.NewKey(publicKeyBuf), nil
}

// NOTE: PublicKeyToAddress2 is copied from gist, test required
func (bc *btcCrypt) PublicKeyToAddress2(publicKey gecdsa.Key) (addr string, err error) {
	mainNetAddr, err := btcutil.NewAddressPubKey(publicKey.Bytes(), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return mainNetAddr.EncodeAddress(), nil
}

// PublicKeyToAddress returns address from public key
func (bc *btcCrypt) PublicKeyToAddress(publicKey gecdsa.Key, typ string) (addr string, err error) {
	if strings.ToLower(typ) != "legacy" {
		return "", gerrors.New("only legacy type supported for now")
	}

	// step 1 - sha256 hash
	pubHash := ghash.GetSHA256(publicKey.Bytes()) // 对公钥进行hash256运算

	// step 2 - ripemd160 hash
	ripemd160 := ripemd160.New()
	ripemd160.Write(pubHash)
	pubHash = ripemd160.Sum(nil)

	// step 3 - insert version byte
	verPubHash := append([]byte{0x00} /*version 0x00*/, pubHash...)

	// step 4 - sha256 hash
	checksumHash := ghash.GetSHA256(verPubHash)

	// step 5 - sha256 hash
	checksumHash = ghash.GetSHA256(checksumHash)

	// step 6 - check sum is first 4 bytes
	checksum := checksumHash[:4]

	// step 7 - encode address
	return gbase.Base58BtcStyleEncode(append(verPubHash, checksum...)), nil
}

func (bc *btcCrypt) Sign(privateKey gecdsa.Key, originData []byte) (sign *gecdsa.Signature, err error) {
	pk, err := loadECDSAFromPrivateKey(privateKey)

	r, s, err := ecdsa.Sign(rand.Reader, pk, originData)
	if err != nil {
		return nil, err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	st, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return nil, err
	}
	w.Flush()

	return &gecdsa.Signature{Data: b.Bytes()}, nil
}

// decodeSign decodes signature into r & s
func decodeSign(sign *gecdsa.Signature) (r, s *big.Int, err error) {
	rd, err := gzip.NewReader(bytes.NewBuffer(sign.Bytes()))
	if err != nil {
		err = gerrors.New("decode error," + err.Error())
		return
	}
	defer rd.Close()

	buf := make([]byte, 1024)
	count, err := rd.Read(buf)
	if err != nil {
		return nil, nil, gerrors.New("decode read error," + err.Error())
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		return nil, nil, gerrors.New("decode fail")
	}

	resultR := big.NewInt(0)
	resultS := big.NewInt(0)
	err = resultR.UnmarshalText([]byte(rs[0]))
	if err != nil {
		return nil, nil, gerrors.New("decrypt r fail, " + err.Error())
	}
	err = resultS.UnmarshalText([]byte(rs[1]))
	if err != nil {
		return nil, nil, gerrors.New("decrypt s fail, " + err.Error())
	}
	return resultR, resultS, nil
}

func (bc *btcCrypt) Verify(originData []byte, sign *gecdsa.Signature, pubKey *ecdsa.PublicKey) (bool, error) {
	r, s, err := decodeSign(sign)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(pubKey, originData, r, s)
	return result, nil
}

func loadECDSAFromPrivateKey(privateKey gecdsa.Key) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	pk.D = new(big.Int).SetBytes(privateKey.Bytes())
	pk.PublicKey.Curve = elliptic.P256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
}

func loadECDSAFromPrivateKeyBytes(privateKeyBytes []byte) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	pk.D = new(big.Int).SetBytes(privateKeyBytes)
	pk.PublicKey.Curve = elliptic.P256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
}

func loadECDSAFromPrivateKeyHexString(privateKeyHexString string) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	pk.D, _ = new(big.Int).SetString(privateKeyHexString, 16)
	pk.PublicKey.Curve = elliptic.P256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
}

func loadSignatureFromHexString(sign string) (*gecdsa.Signature, error) {
	buf, err := hex.DecodeString(sign)
	if err != nil {
		return nil, err
	}
	return &gecdsa.Signature{Data: buf}, nil
}
