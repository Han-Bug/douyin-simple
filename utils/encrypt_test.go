package utils

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("12345678")
	encrypt, err := EncryptDes("12", key)
	if err != nil {
		return
	}
	fmt.Println("encrypt:", encrypt)
	decrypt, err := DecryptDes(encrypt, key)
	if err != nil {
		return
	}
	fmt.Println("decrypt:", decrypt)

}
