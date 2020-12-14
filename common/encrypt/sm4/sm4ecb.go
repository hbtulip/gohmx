/*
// Copyright 2017 cetc-30. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package china crypto algorithm implements the sm2, sm3, sm4 algorithms
*/
package sm4

import (
	"bytes"
	"errors"
)


const (
	ENC cryptMode = iota
	DEC
)

type cryptMode int

/*
const BlockSize = 16
var fk = [4]uint32{0xA3B1BAC6, 0x56AA3350, 0x677D9197, 0xB27022DC}
var ck = [32]uint32{0x00070e15, 0x1c232a31, 0x383f464d, 0x545b6269, 0x70777e85,
	0x8c939aa1, 0xa8afb6bd, 0xc4cbd2d9, 0xe0e7eef5, 0xfc030a11, 0x181f262d,
	0x343b4249, 0x50575e65, 0x6c737a81, 0x888f969d, 0xa4abb2b9, 0xc0c7ced5,
	0xdce3eaf1, 0xf8ff060d, 0x141b2229, 0x30373e45, 0x4c535a61, 0x686f767d,
	0x848b9299, 0xa0a7aeb5, 0xbcc3cad1, 0xd8dfe6ed, 0xf4fb0209, 0x10171e25,
	0x2c333a41, 0x484f565d, 0x646b7279}
*/
var sboxt = [16][16]byte{
	{0xd6, 0x90, 0xe9, 0xfe, 0xcc, 0xe1, 0x3d, 0xb7, 0x16, 0xb6, 0x14, 0xc2, 0x28, 0xfb, 0x2c, 0x05},
	{0x2b, 0x67, 0x9a, 0x76, 0x2a, 0xbe, 0x04, 0xc3, 0xaa, 0x44, 0x13, 0x26, 0x49, 0x86, 0x06, 0x99},
	{0x9c, 0x42, 0x50, 0xf4, 0x91, 0xef, 0x98, 0x7a, 0x33, 0x54, 0x0b, 0x43, 0xed, 0xcf, 0xac, 0x62},
	{0xe4, 0xb3, 0x1c, 0xa9, 0xc9, 0x08, 0xe8, 0x95, 0x80, 0xdf, 0x94, 0xfa, 0x75, 0x8f, 0x3f, 0xa6},
	{0x47, 0x07, 0xa7, 0xfc, 0xf3, 0x73, 0x17, 0xba, 0x83, 0x59, 0x3c, 0x19, 0xe6, 0x85, 0x4f, 0xa8},
	{0x68, 0x6b, 0x81, 0xb2, 0x71, 0x64, 0xda, 0x8b, 0xf8, 0xeb, 0x0f, 0x4b, 0x70, 0x56, 0x9d, 0x35},
	{0x1e, 0x24, 0x0e, 0x5e, 0x63, 0x58, 0xd1, 0xa2, 0x25, 0x22, 0x7c, 0x3b, 0x01, 0x21, 0x78, 0x87},
	{0xd4, 0x00, 0x46, 0x57, 0x9f, 0xd3, 0x27, 0x52, 0x4c, 0x36, 0x02, 0xe7, 0xa0, 0xc4, 0xc8, 0x9e},
	{0xea, 0xbf, 0x8a, 0xd2, 0x40, 0xc7, 0x38, 0xb5, 0xa3, 0xf7, 0xf2, 0xce, 0xf9, 0x61, 0x15, 0xa1},
	{0xe0, 0xae, 0x5d, 0xa4, 0x9b, 0x34, 0x1a, 0x55, 0xad, 0x93, 0x32, 0x30, 0xf5, 0x8c, 0xb1, 0xe3},
	{0x1d, 0xf6, 0xe2, 0x2e, 0x82, 0x66, 0xca, 0x60, 0xc0, 0x29, 0x23, 0xab, 0x0d, 0x53, 0x4e, 0x6f},
	{0xd5, 0xdb, 0x37, 0x45, 0xde, 0xfd, 0x8e, 0x2f, 0x03, 0xff, 0x6a, 0x72, 0x6d, 0x6c, 0x5b, 0x51},
	{0x8d, 0x1b, 0xaf, 0x92, 0xbb, 0xdd, 0xbc, 0x7f, 0x11, 0xd9, 0x5c, 0x41, 0x1f, 0x10, 0x5a, 0xd8},
	{0x0a, 0xc1, 0x31, 0x88, 0xa5, 0xcd, 0x7b, 0xbd, 0x2d, 0x74, 0xd0, 0x12, 0xb8, 0xe5, 0xb4, 0xb0},
	{0x89, 0x69, 0x97, 0x4a, 0x0c, 0x96, 0x77, 0x7e, 0x65, 0xb9, 0xf1, 0x09, 0xc5, 0x6e, 0xc6, 0x84},
	{0x18, 0xf0, 0x7d, 0xec, 0x3a, 0xdc, 0x4d, 0x20, 0x79, 0xee, 0x5f, 0x3e, 0xd7, 0xcb, 0x39, 0x48},
}

func scSbox(in byte) byte {
	var x, y int
	x = (int)(in >> 4 & 0x0f)
	y = (int)(in & 0x0f)
	return sboxt[x][y]
}

func tt(in uint32) uint32 {
	var tmp [4]byte
	var re uint32
	tmp[0] = byte(in>>24) & 0xff
	tmp[1] = byte(in>>16) & 0xff
	tmp[2] = byte(in>>8) & 0xff
	tmp[3] = byte(in) & 0xff
	re = (uint32(scSbox(tmp[3])) |
		(uint32(scSbox(tmp[2])) << 8) |
		(uint32(scSbox(tmp[1])) << 16) |
		(uint32(scSbox(tmp[0])) << 24))
	return re
}

func l(in uint32) uint32 {
	return in ^ leftRotate(in, 2) ^ leftRotate(in, 10) ^ leftRotate(in, 18) ^ leftRotate(in, 24)
}

func key_l(in uint32) uint32 {
	return in ^ leftRotate(in, 13) ^ leftRotate(in, 23)
}

func t3(in uint32) uint32 {
	return l(tt(in))
}
func key_t(in uint32) uint32 {
	return key_l(tt(in))
}

func keyExp(key [4]uint32) [32]uint32 {
	var k [36]uint32
	var rk [32]uint32
	for i := 0; i < 4; i++ {
		k[i] = key[i] ^ fk[i]
	}
	for i := 0; i < 32; i++ {
		k[i+4] = k[i] ^ key_t(k[i+1]^k[i+2]^k[i+3]^ck[i])
		rk[i] = k[i+4]
	}
	return rk
}

func encrypt_oneround(rk [32]uint32, msg []byte) []byte {
	var x [36]uint32
	var cipher = make([]byte, 16)
	for i := 0; i < 4; i++ {
		x[i] = (uint32(msg[i*4+3])) |
			(uint32(msg[i*4+2]) << 8) |
			(uint32(msg[i*4+1]) << 16) |
			(uint32(msg[i*4]) << 24)
	}
	for i := 0; i < 32; i++ {
		x[i+4] = x[i] ^ t3(x[i+1]^x[i+2]^x[i+3]^rk[i])
	}
	for i := 0; i < 4; i++ {
		cipher[i*4] = byte(x[35-i]>>24) & 0xff
		cipher[i*4+1] = byte(x[35-i]>>16) & 0xff
		cipher[i*4+2] = byte(x[35-i]>>8) & 0xff
		cipher[i*4+3] = byte(x[35-i]) & 0xff
	}
	return cipher
}

func rk_swap(rk [32]uint32) [32]uint32 {
	var tmp uint32
	for i := 0; i < 16; i++ {
		tmp = rk[i]
		rk[i] = rk[31-i]
		rk[31-i] = tmp
	}
	return rk
}
/*
	SM4 PKCS7Padding填充模式
     原理：与16的倍数进行相比，缺少多少位就填充多少位的位数值
	 例如：test字节数为4，填充12个0x0c；12345678字节数为8，填充8个8
*/

func pkcs7Padding(src []byte) []byte {
	padding := BlockSize - len(src)%BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > BlockSize || unpadding == 0 {
		return nil, errors.New("Invalid pkcs7 padding (unpadding > BlockSize || unpadding == 0)")
	}

	pad := src[len(src)-unpadding:]
	for i := 0; i < unpadding; i++ {
		if pad[i] != byte(unpadding) {
			return nil, errors.New("Invalid pkcs7 padding (pad[i] != unpadding)")
		}
	}

	return src[:(length - unpadding)], nil
}

func Sm4Ecb(key []byte, msg []byte, mode cryptMode) []byte {
	var inData []byte
	if mode == ENC {
		inData = pkcs7Padding(msg)
	} else {
		inData = msg
	}
	var key_u32 [4]uint32
	cipher := make([]byte, len(inData))
	for i := 0; i < 4; i++ {
		key_u32[i] = (uint32(key[i*4+3])) |
			(uint32(key[i*4+2]) << 8) |
			(uint32(key[i*4+1]) << 16) |
			(uint32(key[i*4]) << 24)
	}
	var rk [32]uint32
	rk = keyExp(key_u32)
	if mode == DEC {
		rk = rk_swap(rk)
	}
	for i := 0; i < len(inData)/16; i++ {
		msg_tmp := inData[i*16 : i*16+16]
		cipher_tmp := encrypt_oneround(rk, msg_tmp)
		copy(cipher[i*16:i*16+16], cipher_tmp)
	}
	if mode == DEC {
		cipher, _ = pkcs7UnPadding(cipher)
	}
	return cipher
}

func leftRotate(x uint32, r int) uint32 {
	var rr uint32 = uint32(r)
	return ((x << rr) | (x >> (32 - rr))) & 0xffffffff
}
