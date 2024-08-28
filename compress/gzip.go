package compress

import (
	"github.com/pierrec/lz4/v4"
	"sync"
)

var (
	_gzip     Compressor
	_gzipPool sync.Pool
)

func init() {
	_glz4 = lz4Compress{}
	_glz4Pool = sync.Pool{New: func() any {
		return &lz4.Compressor{}
	}}
}
