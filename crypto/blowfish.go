package crypto

import (
	"crypto/rand"
	"github.com/ameise84/pi_net/crypto/blowfish"
)

var _gBlowfish blowfishCryptos

type blowfishCryptos struct {
	crypto Cryptos
	key    []byte
}

func init() {
	key := make([]byte, 56)
	n, err2 := rand.Read(key)
	if err2 != nil {
		panic(err2.Error())
	}
	if n < 1 {
		panic("init blowfish key n<1")
	}
	_gBlowfish.key = key[:n]
	_gBlowfish.crypto = blowfish.NewCipher(_gBlowfish.key)
}
