package util

import (
	"strconv"
	"strings"

	CONSTANT "salbackend/constant"
	DB "salbackend/database"
)

// GetBillingDetails - calculate tax, paid amount
func GetBillingDetails(price, discount string) map[string]string {
	billing := map[string]string{}

	actualAmount, _ := strconv.ParseFloat(price, 64)
	discountAmount, _ := strconv.ParseFloat(discount, 64)
	actualAmount -= discountAmount
	if actualAmount < 0 { // if amount becomes negative after discount
		actualAmount = 0
	}
	tax := actualAmount * CONSTANT.GSTPercent
	paidAmount := actualAmount + tax

	billing["actual_amount"] = price
	billing["discount"] = discount
	billing["tax"] = strconv.FormatFloat(tax, 'f', 2, 64)
	billing["paid_amount"] = strconv.FormatFloat(paidAmount, 'f', 2, 64)

	return billing
}

// CheckIfAppointmentSlotAvailable - for both counsellor and listener, check if the specfied slot is available - date (2021-01-12), time (0-23 slots in IST)
func CheckIfAppointmentSlotAvailable(counsellorID, date, time string) bool {
	data, _, _ := DB.SelectSQL(CONSTANT.SlotsTable, []string{date}, map[string]string{"counsellor_id": counsellorID, "date": date, time: "1"}) // if the date time data is 1 in database
	return len(data) > 0
}

// AssociateLanguagesAndTopics - add/update languages and topics for counsellor/listener
func AssociateLanguagesAndTopics(topicIDs, languageIDs, id string) {
	if len(topicIDs) > 0 {
		// first delete all and add topics to listener - to update
		DB.DeleteSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": id})
		for _, topicID := range strings.Split(topicIDs, ",") {
			DB.InsertSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": id, "topic_id": topicID})
		}
	}

	if len(languageIDs) > 0 {
		// first delete all and add languages to listener - to update
		DB.DeleteSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": id})
		for _, languageID := range strings.Split(languageIDs, ",") {
			DB.InsertSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": id, "language_id": languageID})
		}
	}
}
