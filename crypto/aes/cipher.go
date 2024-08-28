package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/ameise84/pi_common/bytes_buffer"
	"github.com/ameise84/pi_net/crypto/internal"
)

// NewCipher key要求16, 24, 32位, iv 的长度必须等于block的长度
func NewCipher(key, iv []byte) *CipherAES {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	r := &CipherAES{
		key:   key,
		iv:    iv,
		block: block,
	}
	return r
}

type CipherAES struct {
	key   []byte
	iv    []byte
	block cipher.Block
}

func (a *CipherAES) BlockSize() int {
	return a.block.BlockSize()
}

func (a *CipherAES) Encrypt(src, dst bytes_buffer.ShiftBuffer) error {
	srcBytes, _ := src.Peek()
	origData := internal.PaddingBlock(srcBytes, a.block.BlockSize())
	crypto := make([]byte, len(origData))
	bm := cipher.NewCBCEncrypter(a.block, a.iv)
	bm.CryptBlocks(crypto, origData) //加密
	//base64.StdEncoding.EncodeToString(crypto)
	return dst.AssignBytes(nil)
}

func (a *CipherAES) Decrypt(src, dst bytes_buffer.ShiftBuffer) error {
	//crypto:=base64.StdEncoding.DecodeString(src.Peek())
	crypto, _ := src.Peek()
	bm := cipher.NewCBCEncrypter(a.block, a.iv)
	origData := make([]byte, len(crypto))
	bm.CryptBlocks(origData, crypto) //解密
	origData = internal.TrimPaddedBlock(origData)
	return dst.AssignBytes(nil)
}
