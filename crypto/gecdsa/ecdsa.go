package gecdsa

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/crypto/gcryptfmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
)

type (
	Signature struct {
		Data []byte
	}

	KeyPair struct {
		PrivateKey Key
		PublicKey  Key
	}

	Key []byte

	bitcoinKeyFmt struct{}

	Ecdsa interface {
		PrivateKeyToPublicKey(privateKey Key) (publicKey Key, err error)
		PublicKeyToAddress(publicKey Key) (addr Key, err error)
		Sign(privateKey Key, originData []byte) (sign *Signature, err error)
		Verify(publicKey Key, originData []byte, sign *Signature) (bool, error)
	}
)

var (
	BitcoinKeyFmt *bitcoinKeyFmt
)

func (s *Signature) Bytes() []byte {
	return s.Data
}

func (s *Signature) String() string {
	return hex.EncodeToString(s.Data)
}

func NewKey(key []byte) Key {
	return key
}

func (dk Key) Bytes() []byte {
	return dk
}

func (dk Key) String(kf gcryptfmt.KeyEncode) (string, error) {
	return kf(dk)
}

func (bkf *bitcoinKeyFmt) Bytes(k Key) []byte {
	return k
}

func (bkf *bitcoinKeyFmt) String(k Key) string {
	return base58.Encode(k)
}

func LoadPrivateKeyFromBytes(privateKeyBytes []byte) (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(privateKeyBytes)
}

func PrivateKeyToKeyPair(key *ecdsa.PrivateKey) (*KeyPair, error) {
	// get private key bytes
	privateKeyBuf := crypto.FromECDSA(key)

	// get public key bytes
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, gerrors.New("error casting public key to ECDSA")
	}
	publicKeyBuf := crypto.FromECDSAPub(publicKeyECDSA)

	return &KeyPair{
		PrivateKey: NewKey(privateKeyBuf),
		PublicKey:  NewKey(publicKeyBuf),
	}, nil
}
