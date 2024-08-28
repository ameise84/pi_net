package net

import "github.com/ameise84/pi_common/bytes_buffer"

func newPacketBufferPool(size int) *packetBufferPool {
	grown := false
	if size <= 0 {
		size = 512
		grown = true
	}
	return &packetBufferPool{
		Pool:  bytes_buffer.NewShiftBufferPool(size, 0, grown),
		size:  size,
		grown: grown,
	}
}

type packetBufferPool struct {
	bytes_buffer.Pool[bytes_buffer.ShiftBuffer]
	size  int
	grown bool
}
