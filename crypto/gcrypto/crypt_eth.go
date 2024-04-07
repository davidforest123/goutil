package gcrypto

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
	"goutil/crypto/gcryptfmt"
	"goutil/crypto/gecdsa"
)

type (
	ethCrypt struct{}
)

func NewEthCrypt() *ethCrypt {
	return &ethCrypt{}
}

// GenerateRandomPrivateKey generates private key in hex string format.
func (ec *ethCrypt) GenerateRandomPrivateKey() ([]byte, error) {
	// generate random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return crypto.FromECDSA(privateKey), nil // SAVE BUT DO NOT SHARE THIS (Private Key)
}

// PrivateKeyToPublicKey generates hex string public key from hex string private key.
func (ec *ethCrypt) PrivateKeyToPublicKey(privateKey gecdsa.Key) (gecdsa.Key, error) {
	// load private key from bytes
	pk, err := gecdsa.LoadPrivateKeyFromBytes(privateKey.Bytes())
	if err != nil {
		return nil, err
	}

	keyPair, err := gecdsa.PrivateKeyToKeyPair(pk)
	if err != nil {
		return nil, err
	}
	return keyPair.PublicKey, nil
}

func (ec *ethCrypt) PublicKeyToAddress(publicKey gecdsa.Key) (addr string, err error) {
	// load public key
	pk, err := crypto.UnmarshalPubkey(publicKey.Bytes())
	if err != nil {
		return "", err
	}

	// generate account address
	return gcryptfmt.EthAddrEncode(crypto.PubkeyToAddress(*pk).Bytes())
}

func (ec *ethCrypt) Sign(privateKey gecdsa.Key, originData []byte) ([]byte, error) {
	// load private key
	pk, err := crypto.ToECDSA(privateKey.Bytes())
	if err != nil {
		return nil, err
	}

	// get 32 bytes length hash of `data` to make sure to sign message is not too long to sign.
	hash := crypto.Keccak256Hash(originData)
	return crypto.Sign(hash.Bytes(), pk)
}

func (ec *ethCrypt) Verify(publicKey gecdsa.Key, originData, signature []byte) (bool, error) {
	// implement 1
	hash := crypto.Keccak256Hash(originData)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, err
	}
	return bytes.Equal(sigPublicKey, publicKey.Bytes()), nil

	// implement 2
	/*sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return false, err
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	return bytes.Equal(sigPublicKeyBytes, publicKeyBytes), nil*/

	// implement 3
	/*signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	return crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID), nil*/
}
