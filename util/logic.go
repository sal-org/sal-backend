package util

import (
	"fmt"
	"math"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
)

func IsStringInSlice(target string, list []string) bool {
	return slices.Contains(list, target)
}

// GetBillingDetails - calculate tax, paid amount
func GetBillingDetails(price, discount string) map[string]string {
	billing := map[string]string{}

	paidAmount, _ := strconv.ParseFloat(price, 64)
	discountAmount, _ := strconv.ParseFloat(discount, 64)
	paidAmount -= discountAmount
	if paidAmount < 0 { // if amount becomes negative after discount
		paidAmount = 0
	}
	paidAmount = math.Round(paidAmount)
	tax := (float64(paidAmount) / float64((100 + CONSTANT.GSTPercent))) * float64(CONSTANT.GSTPercent)
	actualAmount := float64(paidAmount) - tax
	cgst, sgst := tax/2, tax/2

	billing["paid_amount"] = strconv.FormatFloat(paidAmount, 'f', 0, 64)
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

func GetBaseURLAndEndpointFromURL(fullURL string) (string, string) {
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}

	baseURL := u.Scheme + "://" + u.Host
	endpoint := strings.TrimPrefix(u.RequestURI(), "/")
	return baseURL, endpoint
}

func GetEndpointFromURL(fullURL string) string {

	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}

	u.RawQuery = ""

	return u.String()
}

// CapitalizeFirst - capitalize first letter of the string
func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
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

func getAvailabilityPriority(date time.Time, today time.Time, tomorrow time.Time) int {

	if date.Equal(today) {
		return 1
	}

	if date.Equal(tomorrow) {
		return 2
	}

	return 3
}

func SortTherapists(therapists []MODEL.TherapistResponse) []MODEL.TherapistResponse {

	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	tomorrow := today.AddDate(0, 0, 1)

	sort.SliceStable(
		therapists,
		func(i, j int) bool {

			a := therapists[i]
			b := therapists[j]

			// --------------------------------
			// No availability
			// --------------------------------

			if a.NextAvailable == nil &&
				b.NextAvailable == nil {
				return false
			}

			if a.NextAvailable == nil {
				return false
			}

			if b.NextAvailable == nil {
				return true
			}

			// --------------------------------
			// Parse dates
			// --------------------------------

			aDate, aErr := time.ParseInLocation(
				"2006-01-02",
				a.NextAvailable.Date,
				now.Location(),
			)

			bDate, bErr := time.ParseInLocation(
				"2006-01-02",
				b.NextAvailable.Date,
				now.Location(),
			)

			if aErr != nil && bErr != nil {
				return false
			}

			if aErr != nil {
				return false
			}

			if bErr != nil {
				return true
			}

			// --------------------------------
			// Priority:
			//
			// Today
			// Tomorrow
			// Later
			// No Available
			// --------------------------------

			aPriority := getAvailabilityPriority(
				aDate,
				today,
				tomorrow,
			)

			bPriority := getAvailabilityPriority(
				bDate,
				today,
				tomorrow,
			)

			if aPriority != bPriority {
				return aPriority < bPriority
			}

			// --------------------------------
			// Same group:
			// sort by date
			// --------------------------------

			if !aDate.Equal(bDate) {
				return aDate.Before(bDate)
			}

			// --------------------------------
			// Same date:
			// sort by slot
			// --------------------------------

			aSlot, _ := strconv.Atoi(
				a.NextAvailable.Slot,
			)

			bSlot, _ := strconv.Atoi(
				b.NextAvailable.Slot,
			)

			return aSlot < bSlot
		},
	)

	return therapists
}

func RotateInHouseTherapist(page string, pageSize int, therapists []MODEL.TherapistResponse) []MODEL.TherapistResponse {
	var inHouse []MODEL.TherapistResponse

	pageInt := 1

	if len(page) > 0 {
		pageInt, _ = strconv.Atoi(page)
	}

	if pageInt < 1 {
		pageInt = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (pageInt - 1) * pageSize

	for _, therapist := range therapists {
		if therapist.InHouseTherapist == "1" {
			inHouse = append(inHouse, therapist)
		}
	}

	if len(inHouse) == 0 {
		return therapists
	}

	weekIndex := GetWeekIndex()

	selectedIndex := weekIndex % len(inHouse)

	selectedTherapist := inHouse[selectedIndex]

	result := make([]MODEL.TherapistResponse, 0, len(therapists))

	result = append(result, selectedTherapist)

	for _, therapist := range therapists {
		if therapist.ID != selectedTherapist.ID {
			result = append(result, therapist)
		}
	}

	result = SortTherapists(result)

	// Apply pagination AFTER rotation.
	if offset >= len(result) {
		return []MODEL.TherapistResponse{}
	}

	end := offset + pageSize

	if end > len(result) {
		end = len(result)
	}

	return result[offset:end]
}

func BuildTherapistAvailability(rows []map[string]string, UserType string) []MODEL.TherapistResponse {

	therapistMap := make(map[string]*MODEL.TherapistResponse)
	order := make([]string, 0)

	// loc, _ := time.LoadLocation("Asia/Kolkata")
	now := GetCurrentTimeInIndia()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	for _, row := range rows {

		id := row["therapist_id"]

		if id == "" {
			continue
		}

		// Create therapist
		if _, exists := therapistMap[id]; !exists {

			therapistMap[id] = &MODEL.TherapistResponse{
				ID:                  id,
				FirstName:           row["first_name"],
				LastName:            row["last_name"],
				Pronoun:             row["pronoun"],
				TotalRating:         row["total_rating"],
				AverageRating:       row["average_rating"],
				Photo:               row["photo"],
				Video:               row["video"],
				SlotType:            row["slot_type"],
				Type:                UserType,
				Price:               row["price"],
				MultipleSessions:    row["multiple_sessions"],
				Education:           row["education"],
				Experience:          row["experience"],
				TherapeuticApproach: row["therapeutic_approach"],
				About:               row["about"],
				CorporateTherapist:  row["corporate_therapist"],
				InHouseTherapist:    row["in_house_therapist"],
				Available:           "false",
				NextAvailable:       nil,
			}

			order = append(order, id)
		}

		therapist := therapistMap[id]

		// If we already found the earliest slot,
		// don't process another date.
		if therapist.NextAvailable != nil {
			continue
		}

		// Date must exist
		date, err := time.ParseInLocation(
			"2006-01-02",
			row["date"],
			now.Location(),
		)

		if err != nil {
			continue
		}

		// Don't process past dates
		if date.Before(today) {
			continue
		}

		// Date-level availability
		if row["available"] != "1" {
			continue
		}

		// ------------------------------------
		// Determine starting slot
		// ------------------------------------

		startSlot := 0

		if date.Equal(today) {

			// 30-minute slots
			// +3 slots = 1 hour 30 minutes gap
			startSlot = getCurrentSlot(now) + 3

		}

		// fmt.Println(startSlot)

		if startSlot > 47 {
			continue
		}

		// ------------------------------------
		// Find first available slot
		// ------------------------------------

		for slot := startSlot; slot < 48; slot++ {

			key := fmt.Sprintf("%d", slot)

			value := row[key]

			// 0 = unavailable
			// 1 = available
			// 2 = booked
			if value != "1" {
				continue
			}

			therapist.Available = "true"

			therapist.NextAvailable = &MODEL.NextAvailable{
				Date: row["date"],
				Slot: fmt.Sprintf("%d", slot),
			}

			break
		}
	}

	// Convert map to slice
	result := make(
		[]MODEL.TherapistResponse,
		0,
		len(order),
	)

	for _, id := range order {

		result = append(
			result,
			*therapistMap[id],
		)
	}

	return result
}

// func BuildTherapistAvailability(rows []map[string]string) []MODEL.TherapistResponse {

// 	therapistMap := make(map[string]*MODEL.TherapistResponse)
// 	order := make([]string, 0)

// 	for _, row := range rows {

// 		id := row["therapist_id"]

// 		if id == "" {
// 			continue
// 		}

// 		// Create therapist
// 		if _, exists := therapistMap[id]; !exists {

// 			therapistMap[id] = &MODEL.TherapistResponse{
// 				ID:                  id,
// 				FirstName:           row["first_name"],
// 				LastName:            row["last_name"],
// 				Pronoun:             row["pronoun"],
// 				TotalRating:         row["total_rating"],
// 				AverageRating:       row["average_rating"],
// 				Photo:               row["photo"],
// 				Price:               row["price"],
// 				MultipleSessions:    row["multiple_sessions"],
// 				Education:           row["education"],
// 				Experience:          row["experience"],
// 				TherapeuticApproach: row["therapeutic_approach"],
// 				About:               row["about"],
// 				CorporateTherapist:  row["corporate_therapist"],
// 				InHouseTherapist:    row["in_house_therapist"],
// 				Available:           "false",
// 				NextAvailable:       nil,
// 			}

// 			order = append(order, id)
// 		}

// 		therapist := therapistMap[id]

// 		// Date-level availability
// 		if row["available"] != "1" {
// 			continue
// 		}

// 		// Find first slot
// 		for slot := 0; slot < 48; slot++ {

// 			key := fmt.Sprintf(
// 				"%d",
// 				slot,
// 			)

// 			value := row[key]

// 			if value != "1" {
// 				continue
// 			}

// 			// Only replace if this is the first
// 			// available slot
// 			if therapist.NextAvailable == nil {

// 				therapist.Available = "true"

// 				therapist.NextAvailable = &MODEL.NextAvailable{
// 					Date: row["date"],
// 					Slot: fmt.Sprintf(
// 						"%d",
// 						slot,
// 					),
// 				}
// 			}

// 			break
// 		}
// 	}

// 	// Convert map to slice
// 	result := make(
// 		[]MODEL.TherapistResponse,
// 		0,
// 		len(order),
// 	)

// 	for _, id := range order {

// 		result = append(
// 			result,
// 			*therapistMap[id],
// 		)
// 	}

// 	return result
// }

// func BuildTherapistAvailability(rows []map[string]string) []MODEL.TherapistResponse {

// 	therapistMap := make(map[string]*MODEL.TherapistResponse)
// 	therapistOrder := make([]string, 0)

// 	now := time.Now()

// 	// Start of today
// 	today := time.Date(
// 		now.Year(),
// 		now.Month(),
// 		now.Day(),
// 		0,
// 		0,
// 		0,
// 		0,
// 		now.Location(),
// 	)

// 	for _, row := range rows {

// 		therapistID := row["therapist_id"]

// 		if therapistID == "" {
// 			continue
// 		}

// 		therapistFirstName := row["first_name"]
// 		therapistLastName := row["last_name"]

// 		// Create therapist only once
// 		if _, exists := therapistMap[therapistID]; !exists {

// 			therapistMap[therapistID] = &MODEL.TherapistResponse{
// 				ID:            therapistID,
// 				FirstName:     therapistFirstName,
// 				LastName:      therapistLastName,
// 				Available:     "false",
// 				NextAvailable: nil,
// 			}

// 			therapistOrder = append(
// 				therapistOrder,
// 				therapistID,
// 			)
// 		}

// 		therapist := therapistMap[therapistID]

// 		fmt.Printf(
// 			"THERAPIST=%s DATE=%s AVAILABLE=%s\n",
// 			therapistID,
// 			row["date"],
// 			row["available"],
// 		)

// 		// Already found earliest slot
// 		if therapist.NextAvailable != nil {
// 			continue
// 		}

// 		dateString := row["date"]

// 		// No date
// 		if dateString == "" {
// 			continue
// 		}

// 		// Parse date
// 		date, err := time.ParseInLocation(
// 			"2006-01-02",
// 			dateString,
// 			now.Location(),
// 		)

// 		if err != nil {
// 			continue
// 		}

// 		// Ignore previous dates
// 		if date.Before(today) {
// 			continue
// 		}

// 		// ------------------------------------
// 		// IMPORTANT:
// 		// Check availability_slots.available
// 		// ------------------------------------

// 		if row["available"] != "1" {

// 			fmt.Printf(
// 				"SKIP %s because available=%q\n",
// 				therapistID,
// 				row["available"],
// 			)
// 			continue
// 		}

// 		// Determine first slot to check
// 		startSlot := 0

// 		// If today, ignore expired slots
// 		if date.Equal(today) {
// 			startSlot = getCurrentSlot(now)
// 		}

// 		// Find first available slot
// 		for slot := startSlot; slot < 48; slot++ {

// 			column := fmt.Sprintf(
// 				"`%d`",
// 				slot,
// 			)

// 			// 1 = available
// 			// 0 = unavailable
// 			// 2 = booked

// 			value := row[column]

// 			if value != "1" {
// 				continue
// 			}

// 			// Found next available slot
// 			therapist.Available = "true"

// 			therapist.NextAvailable =
// 				&MODEL.NextAvailable{
// 					Date: dateString,
// 					Slot: fmt.Sprintf("%d", slot),
// 				}

// 			break
// 		}
// 	}

// 	// Convert map to response slice
// 	result := make(
// 		[]MODEL.TherapistResponse,
// 		0,
// 		len(therapistOrder),
// 	)

// 	for _, therapistID := range therapistOrder {

// 		result = append(
// 			result,
// 			*therapistMap[therapistID],
// 		)
// 	}

// 	return result
// }

func getCurrentSlot(now time.Time) int {
	minutes := now.Hour()*60 + now.Minute()

	return minutes / 30
}

// FilterAvailableSlots - show only available slots and dates
func FilterAvailableSlots(slots []map[string]string) []map[string]string {

	filteredSlots := []map[string]string{}
	for _, slot := range slots {
		filteredSlot := map[string]string{}
		startSlot := 0
		if strings.EqualFold(GetCurrentTime().Add(330*time.Minute).Format("2006-01-02"), slot["date"]) {
			// use from next hour and multiply by 2 to get 30 min slots
			// startSlot = (GetCurrentTime().Add(330*time.Minute).Hour()+1)*2 + 2 // use next slot for removing expired time for today

			startSlot = getCurrentSlot(GetCurrentTimeInIndia()) + 3
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
