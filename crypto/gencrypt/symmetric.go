package gencrypt

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/wumansgy/goEncrypt/aes"
)

// TODO：使用https://github.com/forgoer/openssl等库验证测试当前的包

func checkKey(ea Alg, key []byte) error {
	if ea == ALG_AES_128_CBC || ea == ALG_AES_128_CTR {
		if len(key) != 16 {
			return gerrors.Errorf("AES-128 requires 16 bytes (128 bits) key")
		}
	}
	if ea == ALG_AES_256_CBC || ea == ALG_AES_256_CTR {
		if len(key) != 32 {
			return gerrors.Errorf("AES-256 requires 32 bytes (256 bits) key")
		}
	}
	return nil
}

func SymmetricEncrypt(ea Alg, plainText, key []byte) (cipherText []byte, err error) {
	if err := checkKey(ea, key); err != nil {
		return nil, err
	}

	switch ea {
	case ALG_AES_128_CBC:
		return aes.AesCbcEncrypt(plainText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_128_CTR:
		return aes.AesCtrEncrypt(plainText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_256_CBC:
		return aes.AesCbcEncrypt(plainText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_256_CTR:
		return aes.AesCtrEncrypt(plainText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	}
	return nil, gerrors.Errorf("unknown encrypt algorithm %s", string(ea))
}

func SymmetricDecrypt(ea Alg, cipherText, key []byte) (plainText []byte, err error) {
	if err := checkKey(ea, key); err != nil {
		return nil, err
	}

	switch ea {
	case ALG_AES_128_CBC:
		return aes.AesCbcDecrypt(cipherText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_128_CTR:
		return aes.AesCtrDecrypt(cipherText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_256_CBC:
		return aes.AesCbcDecrypt(cipherText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	case ALG_AES_256_CTR:
		return aes.AesCtrDecrypt(cipherText, key, nil) // TODO: 确定最后一个参数是否影响输出结果
	}
	return nil, gerrors.Errorf("unknown decrypt algorithm %s", string(ea))
}
