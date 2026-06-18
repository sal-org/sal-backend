package util

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	CONFIG "salbackend/config"
	MODEL "salbackend/model"
	"strings"
)

// GetPayUPayment - get razorpay transaction details
func GetPayUPayment(transactionID string) bool {

	isRecordVaild := false

	url := CONFIG.PayUURL

	data := CONFIG.PayUKey + "|" + "verify_payment" + "|" + transactionID + "|" + CONFIG.PayUSalt
	hash := sha512.New()
	hash.Write([]byte(data))
	hashInString := hex.EncodeToString(hash.Sum(nil))

	payload := strings.NewReader("key=" + CONFIG.PayUKey + "&command=verify_payment&var1=" + transactionID + "&hash=" + hashInString)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return false
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	payUTransaction := MODEL.PaymentVerify{}
	json.Unmarshal(body, &payUTransaction)

	if payUTransaction.Status == 1 {
		if payUTransaction.TransactionDetails[transactionID]["status"] == "success" {
			isRecordVaild = true
		}
	}

	return isRecordVaild
}
