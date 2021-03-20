package util

import "strconv"

// GetPayUPaymentObject - all the details required by payu, to create payment
func GetPayUPaymentObject(merchantKey, salt, orderID, productInfo, firstName, lastName, email, phone, sURL, fURL string, amount float64) map[string]interface{} {
	payment := map[string]interface{}{}

	payment["key"] = merchantKey
	payment["txnid"] = orderID
	payment["amount"] = amount
	payment["productinfo"] = productInfo
	payment["firstname"] = firstName
	payment["lastname"] = lastName
	payment["email"] = email
	payment["phone"] = phone
	payment["surl"] = sURL
	payment["furl"] = fURL
	payment["hash"] = getSHA512(merchantKey + "|" + orderID + "|" + strconv.FormatFloat(amount, 'f', 2, 64) + "|" + productInfo + "|" + firstName + "|" + email + "|||||||||||" + salt + "")

	return payment
}
