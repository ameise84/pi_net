package crypto

import "github.com/ameise84/pi_common/bytes_buffer"

type Cryptos interface {
	BlockSize() int
	Encrypt(src, dst bytes_buffer.ShiftBuffer) error
	Decrypt(src, dst bytes_buffer.ShiftBuffer) error
}
