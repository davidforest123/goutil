package gmnemonic

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"goutil/crypto/gecdsa"
)

// NewRandMnemonic generates a string consisting of the mnemonic words for the given entropy.
// The Mnemonic specification is defined by BIP39.
// entropy:
//
//	 128 bits -> 12 words
//		160 bits -> 15 words
//		192 bits -> 18 words
//		224 bits -> 21 words
//		256 bits -> 24 words
func NewRandMnemonic(entropy int) (string, error) {
	et, err := bip39.NewEntropy(entropy)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(et)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// GenerateBIP32MasterKey generates BIP32 root/master key with mnemonic and passphrase.
// mnemonic: mnemonic word list
// passphrase: salt
func GenerateBIP32MasterKey(mnemonic, passphrase string) (masterPrivateKey, masterPublicKey string, err error) {
	// Generate a BIP32 HD wallet seed for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, passphrase)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", err
	}

	// .String() implements Base58(Version+Depth+FingerPrint+ChildNumber+ChainCode+KeyBytes+Checksum)
	return masterKey.String(), masterKey.PublicKey().String(), nil
}

// GenerateChildKeyWithPath generates child key with mnemonic, passphrase and derivation path.
// The derivationPath specification is defined by BIP32,BIP39,BIP44,BIP49,BIP84....
// m / purpose' / coin_type' / account' / change / address_index
// MEW Ethereum derivation path: m/44'/60'/1'/0/0
func GenerateChildKeyWithPath(mnemonic, passphrase, derivationPath string) (*gecdsa.KeyPair, error) {
	// Generate a BIP32 HD wallet seed for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, passphrase)

	path, err := accounts.ParseDerivationPath(derivationPath)
	if err != nil {
		return nil, err
	}

	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	keyEC, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	privateKey := keyEC.ToECDSA()
	return gecdsa.PrivateKeyToKeyPair(privateKey)
}

/*
// Note: Can't find other wallets to verify whether this function works fine.
// GenerateChildKeyWithIndex generates child key with mnemonic, passphrase and child index.
// BIP32.
func GenerateChildKeyWithIndex(mnemonic, passphrase string, childIdx uint32) (*gecdsa.KeyPair, error) {
	// Generate a BIP32 HD wallet seed for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, passphrase)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	childKey, err := masterKey.NewChildKey(childIdx)
	if err != nil {
		return nil, err
	}
	return &gecdsa.KeyPair{
		PrivateKey: gecdsa.NewKey(childKey.Key),
		PublicKey:  gecdsa.NewKey(childKey.PublicKey().Key),
	}, nil
}*/
