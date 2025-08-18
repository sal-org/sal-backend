package util

import (
	"strconv"
	"strings"
	"time"

	CONSTANT "salbackend/constant"
	DB "salbackend/database"
)

// GetBillingDetails - calculate tax, paid amount
func GetBillingDetails(price, discount string) map[string]string {
	billing := map[string]string{}

	paidAmount, _ := strconv.ParseFloat(price, 64)
	discountAmount, _ := strconv.ParseFloat(discount, 64)
	paidAmount -= discountAmount
	if paidAmount < 0 { // if amount becomes negative after discount
		paidAmount = 0
	}
	tax := (float64(paidAmount) / float64((100 + CONSTANT.GSTPercent))) * float64(CONSTANT.GSTPercent)
	actualAmount := float64(paidAmount) - tax
	cgst, sgst := tax/2, tax/2

	billing["paid_amount"] = strconv.FormatFloat(paidAmount, 'f', 2, 64)
	billing["discount"] = discount
	billing["tax"] = strconv.FormatFloat(tax, 'f', 2, 64)
	billing["actual_amount"] = strconv.FormatFloat(actualAmount, 'f', 2, 64)
	billing["cgst"] = strconv.FormatFloat(cgst, 'f', 2, 64)
	billing["sgst"] = strconv.FormatFloat(sgst, 'f', 2, 64)

	return billing
}

func CalculateAge(birthdateStr string) (int, error) {
	// Define the format of the birthdate string (e.g., "YYYY-MM-DD")
	const dateFormat = "2006-01-02"

	// Parse the birthdate string to time.Time
	birthdate, err := time.Parse(dateFormat, birthdateStr)
	if err != nil {
		return 0, err
	}

	// Get the current date
	currentDate := time.Now()

	// Calculate the age based on the year difference
	age := currentDate.Year() - birthdate.Year()

	// Adjust if the current date is before the birthday this year
	if currentDate.YearDay() < birthdate.YearDay() {
		age--
	}

	return age, nil
}

func EncodeEmailID(email string) string {
	// Find the position of the '@' symbol
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		// Return the email as is if '@' is not found (invalid email format)
		return email
	}

	// Split the email into the local part and domain part
	localPart := email[:atIndex]
	domainPart := email[atIndex:]

	// Check if the local part has at least 1 characters
	if len(localPart) <= 1 {
		// Replace the first 1 characters with 'x'
		localPart = strings.Repeat("x", len(localPart))
	} else if len(localPart) == 2 {
		// Replace the first 2 characters with 'x'
		localPart = localPart[:1] + "x"
	} else if len(localPart) == 3 {
		// Replace the first 3 characters with 'x'
		localPart = localPart[:1] + "x" + localPart[2:]
	} else if len(localPart) > 3 && len(localPart) <= 6 {
		localPart = localPart[:1] + "xx" + localPart[3:]
	} else if len(localPart) > 6 && len(localPart) <= 9 {
		localPart = localPart[:2] + "xxx" + localPart[5:]
	} else if len(localPart) > 9 && len(localPart) <= 12 {
		localPart = localPart[:3] + "xxxx" + localPart[8:]
	} else if len(localPart) > 12 && len(localPart) <= 15 {
		localPart = localPart[:3] + "xxxxx" + localPart[8:]
	} else {
		localPart = localPart[:5] + "xxxxxxx" + localPart[12:]
	}

	// Reconstruct the email with modified local part
	return localPart + domainPart
}

// GetBillingDetails - calculate tax, paid amount
func GetDiscount(price, discount string) map[string]string {
	billing := map[string]string{}

	paidAmount, _ := strconv.ParseFloat(price, 64)
	discountAmount, _ := strconv.ParseFloat(discount, 64)
	paidAmount -= discountAmount
	if paidAmount < 0 { // if amount becomes negative after discount
		paidAmount = 0
	}

	billing["paid_amount"] = strconv.FormatFloat(paidAmount, 'f', 2, 64)
	billing["discount"] = discount

	return billing
}

func AvgRatingFromula(rating []map[string]string, totalCount string, paramName string) string {

	totalCnt, _ := strconv.ParseFloat(totalCount, 32)

	var sum float64
	for i := 0; i < len(rating); i++ {
		rating, _ := strconv.ParseFloat(rating[i][paramName], 64)
		sum = sum + rating
	}

	avg := sum / totalCnt
	avgInString := strconv.FormatFloat(avg, 'f', 1, 64)

	return avgInString
}

// CheckIfAppointmentSlotAvailable - for both counsellor and listener, check if the specfied slot is available - date (2021-01-12), time (0-47 slots in IST)
func CheckIfAppointmentSlotAvailable(counsellorID, date, time string) bool {
	data, _, _ := DB.SelectSQL(CONSTANT.SlotsTable, []string{"1"}, map[string]string{"counsellor_id": counsellorID, "date": date, time: CONSTANT.SlotAvailable}) // if the date time data is 1 in database
	return len(data) > 0
}

// CheckIfAppointmentSlotAvailable - for both counsellor and listener, check if the specfied slot is available - date (2021-01-12), time (0-47 slots in IST)
func CheckIfAppointmentzInPersonSlotAvailable(counsellorID, date, time string) bool {
	data, _, _ := DB.SelectSQL(CONSTANT.InPersonSLotsTable, []string{"1"}, map[string]string{"counsellor_id": counsellorID, "date": date, time: CONSTANT.SlotAvailable}) // if the date time data is 1 in database
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

	filteredSlots := []map[string]string{}
	for _, slot := range slots {
		filteredSlot := map[string]string{}
		startSlot := 0
		if strings.EqualFold(GetCurrentTime().Format("2006-01-02"), slot["date"]) {
			// use from next hour and multiply by 2 to get 30 min slots
			startSlot = (GetCurrentTime().Add(330*time.Minute).Hour()+1)*2 + 2 // use next slot for removing expired time for today
		}

		for i := startSlot; i < 48; i++ { // 48 - 30 min slots
			// show only times with availability
			if strings.EqualFold(slot[strconv.Itoa(i)], "1") {
				filteredSlot[strconv.Itoa(i)] = "1"
			}
		}

		if len(filteredSlot) > 0 { // atleast 1 slot is available
			//filteredSlot["date"] = slot["date"]
			filteredSlot["date"] = slot["date"]
			filteredSlots = append(filteredSlots, filteredSlot)
		}
	}

	// Time Zone Conversion

	// counsellorTimeZone := "330" // IST

	// clientTimeZone := "-300" // IST

	// counsellorTimeInInt, _ := strconv.Atoi(counsellorTimeZone)

	// clientTimeInInt, _ := strconv.Atoi(clientTimeZone) // IST

	// counsellorTimeInInt = counsellorTimeInInt / 30
	// clientTimeInInt = clientTimeInInt / 30

	// //remove dates with no availability
	// filteredSlots := []map[string]string{}
	// for _, slot := range slots {
	// 	filteredSlot := map[string]string{}
	// 	privousSlot := map[string]string{}
	// 	startSlot := 0
	// 	if strings.EqualFold(GetCurrentTime().Format("2006-01-02"), slot["date"]) {
	// 		// use from next hour and multiply by 2 to get 30 min slots
	// 		startSlot = (GetCurrentTime().Add(330*time.Minute).Hour()+1)*2 + 2 // use next slot for removing expired time for today
	// 	}

	// 	for i := startSlot; i < 48; i++ { // 48 - 30 min slots
	// 		// show only times with availability
	// 		if strings.EqualFold(slot[strconv.Itoa(i)], "1") {
	// 			index := i
	// 			index = index - counsellorTimeInInt
	// 			index = index + clientTimeInInt
	// 			if index < 0 {
	// 				index = index + 47
	// 				privousSlot[strconv.Itoa(index)] = "1"
	// 			} else {
	// 				filteredSlot[strconv.Itoa(index)] = "1"
	// 			}
	// 		}
	// 	}

	// 	if len(filteredSlot) > 0 { // atleast 1 slot is available
	// 		//filteredSlot["date"] = slot["date"]
	// 		filteredSlot["date"] = slot["date"]
	// 		filteredSlots = append(filteredSlots, filteredSlot)
	// 	}

	// 	if len(privousSlot) > 0 { // atleast 1 slot is available
	// 		//filteredSlot["date"] = slot["date"]
	// 		date, _ := time.Parse("2006-01-02", slot["date"])

	// 		previousDate := date.AddDate(0, 0, -1)

	// 		// Format the resulting date back to a string
	// 		previousDateStr := previousDate.Format("2006-01-02")
	// 		privousSlot["date"] = previousDateStr
	// 		filteredSlots = append(filteredSlots, privousSlot)
	// 	}
	// }

	// // Step 1: Merge by date
	// mergedMap := make(map[string]map[string]string)

	// for _, entry := range filteredSlots {
	// 	date := entry["date"]

	// 	if _, exists := mergedMap[date]; !exists {
	// 		mergedMap[date] = make(map[string]string)
	// 		mergedMap[date]["date"] = date
	// 	}

	// 	for k, v := range entry {
	// 		if k != "date" {
	// 			mergedMap[date][k] = v
	// 		}
	// 	}
	// }

	// // Step 2: Collect and sort dates
	// var sortedDates []string
	// for date := range mergedMap {
	// 	sortedDates = append(sortedDates, date)
	// }
	// sort.Strings(sortedDates) // sorts in ascending order

	// // Step 3: Build result using sorted dates
	// var result []map[string]string
	// for _, date := range sortedDates {
	// 	result = append(result, mergedMap[date])
	// }

	return filteredSlots
}

// FilterAvailableSlots - show only available slots and dates
func FilterAvailableForInPersonSlots(slots []map[string]string) []map[string]string {
	// remove dates with no availability
	filteredSlots := []map[string]string{}
	for _, slot := range slots {
		filteredSlot := map[string]string{}
		startSlot := 0
		if strings.EqualFold(GetCurrentTime().Format("2006-01-02"), slot["date"]) {
			// use from next hour and multiply by 2 to get 30 min slots
			startSlot = (GetCurrentTime().Add(330 * time.Minute).Hour()) * 2 // use next slot for removing expired time for today

			if (GetCurrentTime().Add(330 * time.Minute).Minute()) >= 30 {
				startSlot = startSlot + 1
			}
		}

		for i := startSlot; i < 48; i++ { // 48 - 30 min slots
			// show only times with availability
			if strings.EqualFold(slot[strconv.Itoa(i)], "1") {
				filteredSlot[strconv.Itoa(i)] = "1"
			}
		}

		if len(filteredSlot) > 0 { // atleast 1 slot is available
			//filteredSlot["date"] = slot["date"]
			filteredSlot["date"] = slot["date"]
			filteredSlot["company_name"] = slot["company_name"]
			filteredSlot["company_location"] = slot["company_location"]
			filteredSlot["counsellor_id"] = slot["counsellor_id"]
			filteredSlots = append(filteredSlots, filteredSlot)
		}
	}

	return filteredSlots
}

// check if schedule available at a particular time slot
func CheckIfScheduleAvailable(schedules []map[string]string, time string) string {
	for _, schedule := range schedules {
		if strings.EqualFold(schedule["availability_status"], "1") && strings.EqualFold(schedule["status"], "1") && strings.EqualFold(schedule[time], CONSTANT.SlotAvailable) {
			return CONSTANT.SlotAvailable
		}
	}
	return CONSTANT.SlotUnavailable
}

func CalculateExperience(startDate, nowDate, gapYears, gapMonth string) string {
	layout := "2006-01-02"

	start, _ := time.Parse(layout, startDate)

	end, _ := time.Parse(layout, nowDate)

	gapyears, _ := strconv.Atoi(gapYears)
	gapmonths, _ := strconv.Atoi(gapMonth)

	years := end.Year() - start.Year()
	years = years - gapyears
	months := int(end.Month()) - int(start.Month())
	months = months - gapmonths
	days := end.Day() - start.Day()

	if days < 0 {
		months--
		days += 30 // Approximation, adjust as needed
	}
	if months < 0 {
		years--
		months += 12
	}

	exprience := ""

	if months >= 6 {
		exprience = strconv.Itoa(years) + ".5"
	} else {
		exprience = strconv.Itoa(years)
	}

	return exprience

}

func ConvertTimeZoneClientToSystem(dateInClient, timeInClient, counsellorTimeZone, clientTimeZone string) (timeZoneDate string, timeZoneTime string) {

	// Timezone conversion

	counsellorTimeInInt, _ := strconv.Atoi(counsellorTimeZone)

	clientTimeInInt, _ := strconv.Atoi(clientTimeZone) // IST

	counsellorTimeInInt = counsellorTimeInInt / 30
	clientTimeInInt = clientTimeInInt / 30

	index, _ := strconv.Atoi(timeInClient)
	index = index - clientTimeInInt
	index = index + counsellorTimeInInt

	if index > 47 {
		index = index - 47
		date, _ := time.Parse("2006-01-02", dateInClient)

		lastestDate := date.AddDate(0, 0, 1)

		// Format the resulting date back to a string
		lastestDateStr := lastestDate.Format("2006-01-02")
		timeZoneDate = lastestDateStr
		timeZoneTime = strconv.Itoa(index)
	} else {
		timeZoneDate = dateInClient
		timeZoneTime = strconv.Itoa(index)
	}

	return timeZoneDate, timeZoneTime
}

func ConvertTimeZoneSystemToClient(dateInClient, timeInClient, counsellorTimeZone, clientTimeZone string) (timeZoneDate string, timeZoneTime string) {

	// Timezone conversion

	counsellorTimeInInt, _ := strconv.Atoi(counsellorTimeZone)

	clientTimeInInt, _ := strconv.Atoi(clientTimeZone) // IST

	counsellorTimeInInt = counsellorTimeInInt / 30
	clientTimeInInt = clientTimeInInt / 30

	index, _ := strconv.Atoi(timeInClient)
	index = index - counsellorTimeInInt
	index = index + clientTimeInInt

	if index < 0 {
		index = index + 47
		date, _ := time.Parse("2006-01-02", dateInClient)

		privousDate := date.AddDate(0, 0, -1)

		// Format the resulting date back to a string
		privousDateStr := privousDate.Format("2006-01-02")
		timeZoneDate = privousDateStr
		timeZoneTime = strconv.Itoa(index)
	} else {
		timeZoneDate = dateInClient
		timeZoneTime = strconv.Itoa(index)
	}

	return timeZoneDate, timeZoneTime
}
