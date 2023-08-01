package util

import (
	"crypto/md5"
	"encoding/hex"
)

type HashesToInsertType struct {
	Id     int
	Input  string
	Output string
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CheckErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
