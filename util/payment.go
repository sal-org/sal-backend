package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	MODEL "salbackend/model"
)

// Verify Payment Signature
func GenerateSignature(signature, orderID, paymentID string) bool {
	data := "" + orderID + "|" + paymentID + ""

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(CONFIG.AWSSecretKey))

	// Write Data to it
	_, err := h.Write([]byte(data))

	if err != nil {
		return false
	}

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return false
	}

	return true

}

// func CreateOrderID(orderID, user_id string, amount float64) {

// 	createOrderBodyBytes, _ := json.Marshal(map[string]interface{}{
// 		"amount":   amount,
// 		"currency": "INR",
// 		"receipt":  "OrderId_" + orderID,
// 		"notes": map[string]interface{}{
// 			"userId": user_id, "paymentOrder": "Created",
// 		},
// 	})

// 	req, _ := http.NewRequest("POST", CONSTANT.RazorPayURL+"/orders", bytes.NewBuffer(createOrderBodyBytes))
// 	req.Header.Add("Authorization", CONFIG.RazorpayAuth)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return MODEL.RazorPayTransaction{}
// 	}

// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	razorPayTransaction := MODEL.RazorPayTransaction{}
// 	json.Unmarshal(body, &razorPayTransaction)

// 	return razorPayTransaction

// }

// GetRazorpayPayment - get razorpay transaction details
func GetRazorpayPayment(transactionID string) MODEL.RazorPayTransaction {
	req, _ := http.NewRequest("GET", CONSTANT.RazorPayURL+"/payments/"+transactionID, nil)
	req.Header.Add("Authorization", CONFIG.RazorpayAuth)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MODEL.RazorPayTransaction{}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	razorPayTransaction := MODEL.RazorPayTransaction{}
	json.Unmarshal(body, &razorPayTransaction)

	return razorPayTransaction
}

// CaptureRazorpayPayment - capture razorpay transaction
func CaptureRazorpayPayment(transactionID string, amount float64) {
	refundBodyBytes, _ := json.Marshal(map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
	})

	req, _ := http.NewRequest("POST", CONSTANT.RazorPayURL+"/payments/"+transactionID+"/capture", bytes.NewBuffer(refundBodyBytes))
	req.Header.Add("Authorization", CONFIG.RazorpayAuth)
	req.Header.Add("Content-Type", "application/json")

	http.DefaultClient.Do(req)
}

// RefundRazorpayPayment - refund amount from razorpay transaction
func RefundRazorpayPayment(transactionID string, amount float64) {
	refundBodyBytes, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	req, _ := http.NewRequest("POST", CONSTANT.RazorPayURL+"/payments/"+transactionID+"/refund", bytes.NewBuffer(refundBodyBytes))
	req.Header.Add("Authorization", CONFIG.RazorpayAuth)
	req.Header.Add("Content-Type", "application/json")

	http.DefaultClient.Do(req)
}
