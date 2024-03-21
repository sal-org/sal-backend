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

func RemoveNotification(tagID, userID string) {
	// delete any previous notifications, if any
	DB.DeleteSQL(CONSTANT.NotificationsTable, map[string]string{"tag_id": tagID, "user_id": userID})
}

// SendNotification - send notification using onesignal
func SendNotification(heading, content, userID, personType, sendAt, status, tagID string) {
	if strings.Contains(content, "###") { // check if notification variables are replaced
		return
	}

	// add data to notifications
	notification := map[string]string{}
	notification["user_id"] = userID
	notification["title"] = heading
	notification["body"] = content
	notification["send_at"] = sendAt
	notification["tag_id"] = tagID
	notification["type"] = personType
	notification["status"] = CONSTANT.NotificationActive
	notification["onesignal_id"] = GetNotificationID(userID, personType)

	notification_status := CheckNotificationEnableORDisable(userID, personType)

	if len(notification["onesignal_id"]) > 0 && notification_status == "0" {
		notification["notification_status"] = CONSTANT.NotificationInvalid
	} else if CONSTANT.NotificationInProgress == status {
		notification["notification_status"] = CONSTANT.NotificationInProgress
	} else {
		// set notification sent status as sent if no onesignal id is available
		notification["notification_status"] = CONSTANT.NotificationSent
		sendNotification(heading, content, notification["onesignal_id"], sendAt, personType)

	}
	notification["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.NotificationsTable, CONSTANT.NotificationsDigits, notification, "notification_id")

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

func CheckNotificationEnableORDisable(id, idType string) string {
	switch idType {
	case CONSTANT.CounsellorType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", id)
	case CONSTANT.ListenerType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.ListenersTable+" where listener_id = ?", id)
	case CONSTANT.ClientType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.ClientsTable+" where client_id = ?", id)
	case CONSTANT.TherapistType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.TherapistsTable+" where therapist_id = ?", id)
	}
	return ""
}

func sendNotification(heading, content, notificationID, sentAt, usertype string) {
	// sent to onesignal
	var app_id, apiKey string
	var byteData []byte

	if usertype == "3" {
		app_id = CONFIG.OneSignalAppIDForClient
		apiKey = CONFIG.OneSignalApiKeyForClient

		// data := MODEL.OneSignalNotificationData{
		// 	AppID:            app_id,
		// 	Headings:         map[string]string{"en": heading},
		// 	Contents:         map[string]string{"en": content},
		// 	IncludePlayerIDs: []string{notificationID},
		// 	Data:             map[string]string{},
		// }
		// byteData, _ = json.Marshal(data)
	} else {
		app_id = CONFIG.OneSignalAppIDForTherapist
		apiKey = CONFIG.OneSignalApiKeyForTherapist

	}

	if strings.Contains(notificationID, "-") {

		data := MODEL.OneSignalNotificationData{
			AppID:            app_id,
			Headings:         map[string]string{"en": heading},
			Contents:         map[string]string{"en": content},
			IncludePlayerIDs: []string{notificationID},
			Data:             map[string]string{},
		}
		byteData, _ = json.Marshal(data)

	} else {
		data := MODEL.OneSignalNotificatnData{
			AppID:          app_id,
			Headings:       map[string]string{"en": heading},
			Contents:       map[string]string{"en": content},
			IncludeAliases: MODEL.IncludeAliase{ExternalID: []string{notificationID}},
			Channels:       []string{"push"},
			Data:           map[string]string{},
		}
		byteData, _ = json.Marshal(data)
	}

	// resp, err := http.Post("https://onesignal.com/api/v1/notifications", "application/json", bytes.NewBuffer(byteData))
	// if err != nil {
	// 	fmt.Println("sendNotification", err)
	// 	return
	// }
	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(byteData))
	req.Header.Add("Authorization", "Basic "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("sendNotification", err)
		return
	}

	fmt.Println(string(body))
}

func SendBulkNotification(heading, content string) {

	data := MODEL.OneSignalNotificationBulkData{
		AppID:            CONFIG.OneSignalAppIDForClient, // change according to client : OneSignalAppIDForClient , therpists : OneSignalAppIDForTherapist required
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		IncludedSegments: []string{"Active Users", "Inactive Users"},
		Data:             map[string]string{},
	}

	byteData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(byteData))
	req.Header.Add("Authorization", "Basic "+CONFIG.OneSignalApiKeyForClient) // change according to client : ZDMxNGU3NTYtM2RkNS00NmMzLWJhMjMtYWUwYTAzYzg3Nzdk , therpists: N2RmZGRlNTMtYTM1MC00YmZmLTg3MjEtNzNkMDViMGZlNGEz required
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
