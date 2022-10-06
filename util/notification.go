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
func SendNotification(heading, content, userID, personType, sendAt, tagID string) {
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
	notification["status"] = CONSTANT.NotificationActive
	notification["onesignal_id"] = GetNotificationID(userID, personType)

	notification_status, usertype := CheckNotificationEnableORDisable(userID, personType)

	if len(notification["onesignal_id"]) > 0 && notification_status == "0" {
		notification["notification_status"] = CONSTANT.NotificationInProgress
	} else {
		// set notification sent status as sent if no onesignal id is available
		notification["notification_status"] = CONSTANT.NotificationSent
		sendNotification(heading, content, notification["onesignal_id"], sendAt, usertype)

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

func CheckNotificationEnableORDisable(id, idType string) (string, string) {
	switch idType {
	case CONSTANT.CounsellorType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", id), "1"
	case CONSTANT.ListenerType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.ListenersTable+" where listener_id = ?", id), "1"
	case CONSTANT.ClientType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.ClientsTable+" where client_id = ?", id), "0"
	case CONSTANT.TherapistType:
		return DB.QueryRowSQL("select notification_status from "+CONSTANT.TherapistsTable+" where therapist_id = ?", id), "1"
	}
	return "", ""
}

func sendNotification(heading, content, notificationID, sentAt, usertype string) {
	// sent to onesignal
	var app_id string

	if usertype == "1" {
		app_id = CONFIG.OneSignalAppIDForTherapist
	} else {
		app_id = CONFIG.OneSignalAppIDForClient
	}

	data := MODEL.OneSignalNotificationData{
		AppID:            app_id,
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		IncludePlayerIDs: []string{notificationID},
		Data:             map[string]string{},
		SendAfter:        sentAt,
	}
	byteData, _ := json.Marshal(data)
	resp, err := http.Post("https://onesignal.com/api/v1/notifications", "application/json", bytes.NewBuffer(byteData))
	if err != nil {
		fmt.Println("sendNotification", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("sendNotification", err)
		return
	}

	fmt.Println(data, string(body))
}
