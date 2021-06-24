package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	"strings"
)

// SendNotification - send notification using onesignal
func SendNotification(heading, content, personID, personType string) {

	// add data to notifications
	notification := map[string]string{}
	notification["person_id"] = personID
	notification["title"] = heading
	notification["body"] = content
	notification["status"] = CONSTANT.NotificationActive
	notification["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.NotificationsTable, CONSTANT.NotificationsDigits, notification, "notification_id")

	// get notification id by person type
	notificationID := GetNotificationID(personID, personType)
	if len(notificationID) == 0 || strings.Contains(content, "###") { // check if notification id is available and notification variables are replaced
		return
	}

	// sent to onesignal
	data := MODEL.OneSignalNotificationData{
		AppID:            CONFIG.OneSignalAppID,
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		IncludePlayerIDs: []string{notificationID},
		Data:             map[string]string{},
	}
	byteData, _ := json.Marshal(data)
	resp, err := http.Post("https://onesignal.com/api/v1/notifications", "application/json", bytes.NewBuffer(byteData))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data, string(body))
}

// GetNotificationID - get notification ID of client/counselors
func GetNotificationID(id string, idType string) string {
	switch idType {
	case CONSTANT.CounsellorType:
		return DB.QueryRowSQL("select device_id from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", id)
	case CONSTANT.ListenerType:
		return DB.QueryRowSQL("select device_id from "+CONSTANT.ListenersTable+" where listener_id = ?", id)
	case CONSTANT.ClientType:
		return DB.QueryRowSQL("select device_id from "+CONSTANT.ClientsTable+" where client_id = ?", id)
	case CONSTANT.TherapistType:
		return DB.QueryRowSQL("select device_id from "+CONSTANT.TherapistsTable+" where therapist_id = ?", id)
	}
	return ""
}
