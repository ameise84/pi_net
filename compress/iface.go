package compress

type Compressor interface {
	Zip(src, dst []byte) ([]byte, error)
	Unzip(src, dst []byte) ([]byte, error)
}
