package b2c

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	CONSTANT "salbackend/constant"

	"github.com/gorilla/mux"
)

// LoadAdminRoutes - load all web routes with web prefix
func LoadWebB2CRoutes(router *mux.Router) {
	webRoutes := router.PathPrefix("/web").Subrouter()

	webRoutes.HandleFunc("", CreateUserHandler).Methods("POST")

	// client login
	// webClientRoutes.HandleFunc("/access_code", CheckAccessCode).Methods("POST")
	// webClientRoutes.HandleFunc("/send_otp", SendOTPWithCorporateEmail).Methods("POST")
	// webClientRoutes.HandleFunc("/verify_otp", VerifyOTPWithCorporateEmail).Methods("POST")

	// // client profile
	// webClientRoutes.HandleFunc("/profile/get", ProfileGet).Methods("POST")
	// webClientRoutes.HandleFunc("/profile/add", ProfileAdd).Methods("POST")
	// webClientRoutes.HandleFunc("/profile/update", ProfileUpdate).Methods("POST")

	// // client search
	// webClientRoutes.HandleFunc("/search", ListSearch).Methods("POST")

	// // client appointment
	// webClientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/detail", AppointmentDetail).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/reschedule", AppointmentReschedule).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/cancel", AppointmentCancel).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/rate", AppointmentRatingAdd).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/start", AppointmentStart).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/end", AppointmentEnd).Methods("POST")

	// // client agora token
	// webClientRoutes.HandleFunc("/agora/token", GenerateAgoraToken).Methods("POST")

	// // client therapist
	// webClientRoutes.HandleFunc("/therapist/get", TherapistProfile).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/slots", TherapistSlots).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/order", CorporateCounsellorOrderCreate).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/complete", CorporateCounsellorOrderPaymentComplete).Methods("POST")

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

// CreateUserHandler handles user creation with detailed error logging
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	// Recover from any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming payload
	var encryptedPayload struct {
		Payload string `json:"payload"`
	}

	// Decode JSON with error logging
	err := json.NewDecoder(r.Body).Decode(&encryptedPayload)
	if err != nil {
		log.Printf("JSON decoding error: %v", err)
		http.Error(w, "Invalid payload format", http.StatusBadRequest)
		return
	}

	decrypted, _ := DecryptPayload(encryptedPayload.Payload, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)

	switch decrypted["path"] {
	case "/book_demo":
		ClientBookDemo(w, r, decrypted)
	default:
		w.Header().Set("Status", "200")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Page Not Found")
	}

}
