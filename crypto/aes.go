package crypto

import (
	"crypto/rand"
	"github.com/ameise84/pi_net/crypto/aes"
)

var _gDefaultAES aesCryptos

type aesCryptos struct {
	crypto Cryptos
	key    []byte
	iv     []byte
}

func init() {
	_gDefaultAES.key = make([]byte, 16)
	n, err := rand.Read(_gDefaultAES.key)
	if err != nil {
		panic(err.Error())
	}
	if n != 16 {
		panic("init aes key err")
	}

	_gDefaultAES.iv = make([]byte, 16)
	n, err = rand.Read(_gDefaultAES.iv)
	if err != nil {
		panic(err.Error())
	}
	if n != 16 {
		panic("init aes iv err")
	}
	_gDefaultAES.crypto = aes.NewCipher(_gDefaultAES.key, _gDefaultAES.iv)

}
