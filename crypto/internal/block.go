package internal

import "bytes"

func PaddingBlock(ciphertext []byte, blockSize int) []byte {
	size := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(size)}, size)
	return append(ciphertext, padText...)
}

func TrimPaddedBlock(origData []byte) []byte {
	length := len(origData)
	size := int(origData[length-1])
	return origData[:(length - size)]
}
