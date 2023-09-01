package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type HashFunc func(text string) string

type HashesToInsertType struct {
	Id     int
	Input  string
	Output string
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func SHA1Hash(text string) string {
	hash := sha1.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func SHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CheckErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func GetHashFuncFromName(algorithmName string) (HashFunc, error) {
	formattedName := strings.ReplaceAll(strings.ToLower(algorithmName), "-", "")
	switch formattedName {
	case "sha256":
		return SHA256Hash, nil
	case "sha1":
		return SHA1Hash, nil
	case "md5":
		return MD5Hash, nil
	default:
		return nil, fmt.Errorf("algorithm name not supported: %s", algorithmName)
	}
}
