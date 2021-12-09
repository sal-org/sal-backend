package util

import (
	"strconv"
	"time"
)

// GetCurrentTime - current time in utc
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// BuildDateTime - build UTC time from given inputs
func BuildDateTime(date string, timeSlot string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", date+" "+GetTimeFromTimeSlot(timeSlot))
	return t
}

// ConvertToPersonTimezone - convert UTC to person timezone
func ConvertTimezone(date time.Time, timezone string) time.Time {
	timezoneMinutes, _ := strconv.Atoi(timezone)
	return date.Add(time.Duration(timezoneMinutes) * time.Minute)
}

// GetTimeFromTimeSlot -
func GetTimeFromTimeSlot(timeSlot string) string {
	slot, _ := strconv.Atoi(timeSlot)
	switch slot {
	case 0:
		return "00:00:00"
	case 1:
		return "00:30:00"
	case 2:
		return "01:00:00"
	case 3:
		return "01:30:00"
	case 4:
		return "02:00:00"
	case 5:
		return "02:30:00"
	case 6:
		return "03:00:00"
	case 7:
		return "03:30:00"
	case 8:
		return "04:00:00"
	case 9:
		return "04:30:00"

	case 10:
		return "05:00:00"
	case 11:
		return "05:30:00"
	case 12:
		return "06:00:00"
	case 13:
		return "06:30:00"
	case 14:
		return "07:00:00"
	case 15:
		return "07:30:00"
	case 16:
		return "08:00:00"
	case 17:
		return "08:30:00"
	case 18:
		return "09:00:00"
	case 19:
		return "09:30:00"

	case 20:
		return "10:00:00"
	case 21:
		return "10:30:00"
	case 22:
		return "11:00:00"
	case 23:
		return "11:30:00"
	case 24:
		return "12:00:00"
	case 25:
		return "12:30:00"
	case 26:
		return "13:00:00"
	case 27:
		return "13:30:00"
	case 28:
		return "14:00:00"
	case 29:
		return "14:30:00"

	case 30:
		return "15:00:00"
	case 31:
		return "15:30:00"
	case 32:
		return "16:00:00"
	case 33:
		return "16:30:00"
	case 34:
		return "17:00:00"
	case 35:
		return "17:30:00"
	case 36:
		return "18:00:00"
	case 37:
		return "18:30:00"
	case 38:
		return "19:00:00"
	case 39:
		return "19:30:00"

	case 40:
		return "20:00:00"
	case 41:
		return "20:30:00"
	case 42:
		return "21:00:00"
	case 43:
		return "21:30:00"
	case 44:
		return "22:00:00"
	case 45:
		return "22:30:00"
	case 46:
		return "23:00:00"
	case 47:
		return "23:30:00"

	}
	return "00:00:00"
}

/*func BuildDateTimeMintues(date string, time string, mintues string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", date+" "+time+":"+GetTimeFromMintuesSlot(mintues))
	return t
}

func GetTimeFromMintuesSlot(mintues string) string {
	slot, _ := strconv.Atoi(mintues)
	switch slot {
	case 00:
		return "00:00"
	case 05:
		return "05:00"
	case 10:
		return "10:00"
	case 15:
		return "15:00"
	case 20:
		return "20:00"
	case 25:
		return "25:00"
	case 30:
		return "30:00"
	case 35:
		return "35:00"
	case 40:
		return "40:00"
	case 45:
		return "45:00"
	case 50:
		return "50:00"
	case 55:
		return "55:00"
	case 60:
		return "60:00"
	}
	return "00:00"

}
*/
