package compress

func Zip(src []byte, out []byte) ([]byte, error) {
	return _glz4.Zip(src, out)
}

func Unzip(src []byte, out []byte) ([]byte, error) {
	return _glz4.Unzip(src, out)
}
