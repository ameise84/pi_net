package crypto

import (
	"github.com/ameise84/pi_common/bytes_buffer"
)

//名称	密钥长度			运算速度	安全性	资源消耗
//DES	56位			较快		低		中
//3DES	112位或168位		慢		中		高
//AES	128、192、256位	快		高		低

//名称	成熟度	安全性(取决于密钥长度)	运算速度	资源消耗
//RSA	高		高					慢		高
//DSA	高		高					慢		只能用于数字签名
//ECC	低		高					快		低(计算量小,存储空间占用小,带宽要求低)

func Blowfish() Cryptos {
	return _gBlowfish.crypto
}

func AES() Cryptos {
	return _gDefaultAES.crypto
}

func Encrypt(src, dst bytes_buffer.ShiftBuffer) error {
	return _gBlowfish.crypto.Encrypt(src, dst)
}

func Decrypt(src, dst bytes_buffer.ShiftBuffer) error {
	return _gBlowfish.crypto.Decrypt(src, dst)
}
