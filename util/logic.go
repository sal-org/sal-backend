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

// CheckIfAppointmentSlotAvailable - for both counsellor and listener, check if the specfied slot is available - date (2021-01-12), time (0-47 slots in IST)
func CheckIfAppointmentSlotAvailable(counsellorID, date, time string) bool {
	data, _, _ := DB.SelectSQL(CONSTANT.SlotsTable, []string{"1"}, map[string]string{"counsellor_id": counsellorID, "date": date, time: "1"}) // if the date time data is 1 in database
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

// FilterAvailableSlots - show only available slots and dates
func FilterAvailableSlots(slots []map[string]string) []map[string]string {
	// remove dates with no availability
	filteredSlots := []map[string]string{}
	for _, slot := range slots {
		filteredSlot := map[string]string{}
		startSlot := 0
		if strings.EqualFold(GetCurrentTime().Format("2006-01-02"), slot["date"]) {
			// use from next hour and multiply by 2 to get 30 min slots
			startSlot = (GetCurrentTime().Hour() + 1) * 2 // use next slot for removing expired time for today
		}

		for i := startSlot; i < 48; i++ { // 48 - 30 min slots
			// show only times with availability
			if strings.EqualFold(slot[strconv.Itoa(i)], "1") {
				filteredSlot[strconv.Itoa(i)] = "1"
			}
		}
		if len(filteredSlot) > 0 { // atleast 1 slot is available
			filteredSlot["date"] = slot["date"]
			filteredSlots = append(filteredSlots, filteredSlot)
		}
	}

	return filteredSlots
}
