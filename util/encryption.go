package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	MODEL "salbackend/model"
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

func getMD5Hash(savedFileName string) string {
	h := md5.New()
	return hex.EncodeToString(h.Sum(nil))
}

// EncryptPayload encrypts data using AES-CBC with PKCS7 padding
func EncryptPayload(plainText map[string]interface{}, key string, iv string) (string, error) {

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(plainText)
	if err != nil {
		log.Fatalf("Error marshaling map: %v", err)
	}

	plainTextBytes := []byte(string(jsonData))
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	// Validate parameters
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("key length must be 16, 24, or 32 bytes")
	}

	if len(ivBytes) != 16 {
		return "", errors.New("IV must be exactly 16 bytes")
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Add PKCS7 padding
	plainTextBytes = addPKCS7Padding(plainTextBytes, aes.BlockSize)

	// Create the CBC encrypter
	mode := cipher.NewCBCEncrypter(block, ivBytes)

	// Allocate space for ciphertext
	ciphertext := make([]byte, len(plainTextBytes))

	// Encrypt
	mode.CryptBlocks(ciphertext, plainTextBytes)

	// Base64 encode for transmission
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// GenerateRandomIV creates a cryptographically secure random IV
func GenerateRandomIV() ([]byte, error) {
	iv := make([]byte, 16) // AES block size is always 16 bytes
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// addPKCS7Padding adds PKCS7 padding to plaintext
func addPKCS7Padding(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := make([]byte, padLen)
	for i := 0; i < padLen; i++ {
		padding[i] = byte(padLen)
	}
	return append(data, padding...)
}

func DecryptPayload(encryptedData string, key string, iv string) (map[string]string, error) {
	// Decode the base64 encrypted string from CryptoJS
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return map[string]string{}, errors.New("base64 decode error: " + err.Error())
	}

	// Convert key and IV to byte slices
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return map[string]string{}, errors.New("cipher creation error: " + err.Error())
	}

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	// Decrypt the data
	decrypted := make([]byte, len(decodedData))
	mode.CryptBlocks(decrypted, decodedData)

	// Remove PKCS7 padding
	unpaddedData, err := removePKCS7Padding(decrypted)
	if err != nil {
		return map[string]string{}, errors.New("padding removal error: " + err.Error())
	}

	// Declare a map to store the unmarshalled data
	result := make(map[string]string)

	// Unmarshal the JSON string into the map
	err = json.Unmarshal(unpaddedData, &result)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return result, nil
}


func DecryptPayloadForAssessment(encryptedData string, key string, iv string) (MODEL.AssessmentAddRequest, error) {
	// Decode the base64 encrypted string from CryptoJS
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return MODEL.AssessmentAddRequest{}, errors.New("base64 decode error: " + err.Error())
	}

	// Convert key and IV to byte slices
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return MODEL.AssessmentAddRequest{}, errors.New("cipher creation error: " + err.Error())
	}

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	// Decrypt the data
	decrypted := make([]byte, len(decodedData))
	mode.CryptBlocks(decrypted, decodedData)

	// Remove PKCS7 padding
	unpaddedData, err := removePKCS7Padding(decrypted)
	if err != nil {
		return MODEL.AssessmentAddRequest{}, errors.New("padding removal error: " + err.Error())
	}

	// Declare a struct to store the unmarshalled data
	result := MODEL.AssessmentAddRequest{}

	// Unmarshal the JSON string into the struct
	err = json.Unmarshal(unpaddedData, &result)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return result, nil
}

func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty input")
	}

	padLength := int(data[len(data)-1])
	if padLength > len(data) || padLength == 0 {
		return nil, errors.New("invalid padding")
	}

	// Verify padding
	for i := 1; i <= padLength; i++ {
		if data[len(data)-i] != byte(padLength) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-padLength], nil
}
