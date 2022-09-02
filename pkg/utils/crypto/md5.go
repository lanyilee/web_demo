package crypto

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strconv"
)

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func hexToByte(hex string) []byte {
	length := len(hex) / 2
	slice := make([]byte, length)
	rs := []rune(hex)

	for i := 0; i < length; i++ {
		s := string(rs[i*2 : i*2+2])
		value, _ := strconv.ParseInt(s, 16, 10)
		slice[i] = byte(value & 0xFF)
	}
	return slice
}

func NewMD5Hash(s string) string {
	return "{MD5}" + base64.StdEncoding.EncodeToString(hexToByte(MD5(s)))
}
