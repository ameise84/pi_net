// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blowfish

import "encoding/binary"

func encryptBlock(c *Cipher, src, dst []byte) (uint32, uint32) {
	xl := binary.BigEndian.Uint32(src[:4])
	xr := binary.BigEndian.Uint32(src[4:8])
	xl ^= c.p[0]
	for i := 0; i < 16; i += 2 {
		xr ^= ((c.s0[byte(xl>>24)] + c.s1[byte(xl>>16)]) ^ c.s0[byte(xl>>8)]) + c.s1[byte(xl)] ^ c.p[i+1]
		xl ^= ((c.s0[byte(xr>>24)] + c.s1[byte(xr>>16)]) ^ c.s0[byte(xr>>8)]) + c.s1[byte(xr)] ^ c.p[i+2]
	}
	xr ^= c.p[17]
	binary.BigEndian.PutUint32(dst[:4], xr)
	binary.BigEndian.PutUint32(dst[4:8], xl)
	return xr, xl
}

func decryptBlock(c *Cipher, src, dst []byte) {
	xl := binary.BigEndian.Uint32(src[:4])
	xr := binary.BigEndian.Uint32(src[4:8])
	xl ^= c.p[17]
	for i := 0; i < 16; i += 2 {
		xr ^= ((c.s0[byte(xl>>24)] + c.s1[byte(xl>>16)]) ^ c.s0[byte(xl>>8)]) + c.s1[byte(xl)] ^ c.p[16-i]
		xl ^= ((c.s0[byte(xr>>24)] + c.s1[byte(xr>>16)]) ^ c.s0[byte(xr>>8)]) + c.s1[byte(xr)] ^ c.p[15-i]
	}
	xr ^= c.p[0]
	binary.BigEndian.PutUint32(dst[:4], xr)
	binary.BigEndian.PutUint32(dst[4:8], xl)
}
