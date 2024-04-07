package gconfig

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/crypto/gbase"
	"github.com/davidforest123/goutil/crypto/gencrypt"
	"github.com/davidforest123/goutil/crypto/ghash"
)

const (
	head = "1024"
)

func configMartolodEncrypt(plain []byte, userSecret, saltSecret string) ([]byte, error) {
	if len(plain) == 0 {
		return nil, gerrors.New("empty plain")
	}
	if userSecret == "" {
		return nil, gerrors.New("empty user secret")
	}
	if saltSecret == "" {
		return nil, gerrors.New("empty salt secret")
	}

	secretHMAC := []byte(userSecret)
	for i := 0; i < 6; i++ {
		secretHMAC = ghash.GetHMAC(ghash.HashTypeSHA256, secretHMAC, []byte(head+saltSecret))
	}

	cipher, err := gencrypt.NewAesGcm256().Encrypt(plain, secretHMAC, true)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func configMartolodDecrypt(cipher []byte, userSecret, saltSecret string) ([]byte, error) {
	if len(cipher) == 0 {
		return nil, gerrors.New("empty cipher")
	}
	if userSecret == "" {
		return nil, gerrors.New("empty user secret")
	}
	if saltSecret == "" {
		return nil, gerrors.New("empty salt secret")
	}

	secretHMAC := []byte(userSecret)
	for i := 0; i < 6; i++ {
		secretHMAC = ghash.GetHMAC(ghash.HashTypeSHA256, secretHMAC, []byte(head+saltSecret))
	}

	plain, err := gencrypt.NewAesGcm256().Decrypt(cipher, secretHMAC, true)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func configTriMartolodEncrypt(plain string, userSecret, saltSecret string) (string, error) {
	buf, err := configMartolodEncrypt([]byte(plain), userSecret, saltSecret)
	if err != nil {
		return "", err
	}
	return gbase.Base64Encode(buf), nil
}

func configTriMartolodDecrypt(cipher string, userSecret, saltSecret string) (string, error) {
	buf, err := gbase.Base64Decode(cipher)
	if err != nil {
		return "", err
	}
	res, err := configMartolodDecrypt(buf, userSecret, saltSecret)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
