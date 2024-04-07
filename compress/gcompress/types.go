package gcompress

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"strings"
)

type (
	Comp    string
	Level   string
	Encrypt string
)

var (
	allComps   []Comp
	CompNone   = enrollComp("none")
	CompZip    = enrollComp("zip")
	Comp7Z     = enrollComp("7z")
	CompSnappy = enrollComp("snappy")
	CompS2     = enrollComp("s2")
	CompGzip   = enrollComp("gzip")
	CompPgZip  = enrollComp("pgzip")
	CompZStd   = enrollComp("zstd")
	CompZLib   = enrollComp("zlib")
	CompFlate  = enrollComp("flate")

	allLevels    []Level
	LevelStore   = enrollLevel("store")
	LevelDeflate = enrollLevel("deflate")

	allEncrypts   []Encrypt
	EncryptAES128 = enrollEncrypt("aes-128")
	EncryptAES192 = enrollEncrypt("aes-192")
	EncryptAES256 = enrollEncrypt("aes-256")
)

func enrollComp(algo string) Comp {
	allComps = append(allComps, Comp(algo))
	return Comp(algo)
}

func ToComp(algo string) (Comp, error) {
	for _, v := range allComps {
		if string(v) == strings.ToLower(algo) {
			return v, nil
		}
	}
	return CompNone, gerrors.New("unrecognized compress algorithm '%s'", algo)
}

func enrollLevel(level string) Level {
	allLevels = append(allLevels, Level(level))
	return Level(level)
}

func enrollEncrypt(encrypt string) Encrypt {
	allEncrypts = append(allEncrypts, Encrypt(encrypt))
	return Encrypt(encrypt)
}
