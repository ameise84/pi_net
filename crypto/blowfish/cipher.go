// 参考blowfish 加密 ,减少了 s 数量

package blowfish

import (
	"encoding/binary"
	"github.com/ameise84/pi_common/bytes_buffer"
)

const BlockSize = 8 //!!!不要修改

func NewCipher(key []byte) *Cipher {
	var result Cipher
	if k := len(key); k < 1 || k > 56 {
		panic("crypto/blowfish2: invalid key size ,must be [1,56]")
	}
	return initCipher(key, &result)
}

type Cipher struct {
	p      [18]uint32
	s0, s1 [256]uint32
}

func (c *Cipher) BlockSize() int {
	return BlockSize
}

func (c *Cipher) Encrypt(src, dst bytes_buffer.ShiftBuffer) error {
	baseLen := src.GetDataSize()
	padding := BlockSize - baseLen%BlockSize

	var lenBytes [BlockSize]byte
	if padding != BlockSize {
		_ = src.AppendBytes(lenBytes[:padding])
	}

	binary.LittleEndian.PutUint64(lenBytes[:], uint64(baseLen))
	_ = src.AppendBytes(lenBytes[:])

	baseLen = src.GetDataSize()
	if dst.GetDataSize() != baseLen {
		dst.ResetLen(baseLen)
	}

	in, _ := src.Peek()
	out, _ := dst.Peek()
	ts := baseLen / BlockSize
	for i := 0; i < ts; i++ {
		idx := i * BlockSize
		encryptBlock(c, in[idx:], out[idx:])
	}
	return nil
}

func (c *Cipher) Decrypt(src, dst bytes_buffer.ShiftBuffer) error {
	baseLen := src.GetDataSize()
	if dst.GetDataSize() != baseLen {
		dst.ResetLen(baseLen)
	}
	in, _ := src.Peek()
	out, _ := dst.Peek()
	ts := baseLen / BlockSize
	for i := 0; i < ts; i++ {
		idx := i * BlockSize
		decryptBlock(c, in[idx:], out[idx:])
	}

	lenBytes := make([]byte, BlockSize)
	copy(lenBytes, out[baseLen-BlockSize:])
	baseLen = int(binary.LittleEndian.Uint64(lenBytes))
	dst.ResetLen(baseLen)
	return nil
}

func initCipher(key []byte, c *Cipher) *Cipher {
	copy(c.p[0:], p[0:])
	copy(c.s0[0:], s0[0:])
	copy(c.s1[0:], s1[0:])
	initKey(key, c)
	return c
}

func initKey(key []byte, c *Cipher) {
	j := 0
	for i := 0; i < 18; i++ {
		var d uint32
		for k := 0; k < 4; k++ {
			d = d<<8 | uint32(key[j])
			j++
			if j >= len(key) {
				j = 0
			}
		}
		c.p[i] ^= d
	}

	out := make([]byte, 8)
	for i := 0; i < 18; i += 2 {
		c.p[i], c.p[i+1] = encryptBlock(c, out, out)
	}

	for i := 0; i < 256; i += 2 {
		c.s0[i], c.s0[i+1] = encryptBlock(c, out, out)
	}
	for i := 0; i < 256; i += 2 {
		c.s1[i], c.s1[i+1] = encryptBlock(c, out, out)
	}
}
