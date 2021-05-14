package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	MODEL "salbackend/model"
)

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
