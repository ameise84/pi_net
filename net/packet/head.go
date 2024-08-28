package packet

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

const HeadSize = 8
const (
	cryptoBit   byte = 1 << 0 //bool
	compressBit byte = 1 << 1 //bool
	bigPacket   byte = 1 << 2 //bool
)

func NewHead() *Head {
	return &Head{
		Mark: &[HeadSize]byte{},
	}
}

func NewHeadWarp(b []byte) *Head {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return &Head{
		Mark: (*[HeadSize]byte)((unsafe.Pointer)(bh.Data)),
	}
}

type Head struct {
	Mark *[HeadSize]byte
}

func (h *Head) Init() {
	for i := 0; i < HeadSize; i++ {
		h.Mark[i] = 0
	}
}

func (h *Head) ToBytes() []byte {
	return h.Mark[:]
}

func (h *Head) FromBytes(b []byte) {
	copy(h.Mark[:], b)
}

func (h *Head) Warp(b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Mark = (*[HeadSize]byte)((unsafe.Pointer)(bh.Data))
}

// Id 消息包ID
func (h *Head) Id() uint32 {
	return binary.BigEndian.Uint32(h.Mark[0:4])
}

// SetId 设置消息包 ID
func (h *Head) SetId(id uint32) {
	binary.BigEndian.PutUint32(h.Mark[0:4], id)
}

// IsCompress 消息包是否压缩
func (h *Head) IsCompress() bool {
	return h.Mark[4]&compressBit != 0
}

// SetCompress 设置压缩
func (h *Head) SetCompress(do bool) {
	h.Mark[4] &= ^compressBit
	if do {
		h.Mark[4] |= compressBit
	}
}

// IsCrypto 消息包是否加密
func (h *Head) IsCrypto() bool {
	return h.Mark[4]&cryptoBit != 0
}

// SetCrypto 设置加密
func (h *Head) SetCrypto(do bool) {
	h.Mark[4] &= ^cryptoBit
	if do {
		h.Mark[4] |= cryptoBit
	}
}

// BodySize 消息包长度
func (h *Head) BodySize() uint16 {
	return binary.BigEndian.Uint16(h.Mark[6:8])
}

func (h *Head) SetBodySize(size uint16) {
	binary.BigEndian.PutUint16(h.Mark[6:8], size)
}
