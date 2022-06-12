package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncryptMd5 md5加密
func EncryptMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
