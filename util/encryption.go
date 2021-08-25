package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GetStringMD5Hash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

func getFileMD5Hash(savedFileName string) string {
	openedFile, _ := os.Open(savedFileName)
	defer openedFile.Close()

	h := md5.New()
	if _, err := io.Copy(h, openedFile); err != nil {
		fmt.Println("getMD5Hash", err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
