package compress

import (
	"github.com/pierrec/lz4/v4"
	"sync"
)

var (
	_glz4     Compressor
	_glz4Pool sync.Pool
)

func init() {
	_glz4 = lz4Compress{}
	_glz4Pool = sync.Pool{New: func() any {
		return &lz4.Compressor{}
	}}
}

type lz4Compress struct{}

func (lz4Compress) Zip(src []byte, dst []byte) ([]byte, error) {
	c := _glz4Pool.Get().(*lz4.Compressor)
	defer _glz4Pool.Put(c)
	n, err := c.CompressBlock(src, dst)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return src, nil
	}
	return dst[:n], nil
}

func (lz4Compress) Unzip(src []byte, dst []byte) ([]byte, error) {
	n, err := lz4.UncompressBlock(src, dst)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}
