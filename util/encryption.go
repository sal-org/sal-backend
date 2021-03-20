package util

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func getSHA512(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func getMD5Hash(savedFileName string) string {
	openedFile, _ := os.Open(savedFileName)
	defer openedFile.Close()

	h := md5.New()
	if _, err := io.Copy(h, openedFile); err != nil {
		fmt.Println("getMD5Hash", err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
