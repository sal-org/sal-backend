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
)

// SendNotification - send notification using onesignal
func SendNotification(heading string, content string, notificationID string) {
	if len(notificationID) == 0 {
		return
	}
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

// GetClientNotificationID -
func GetClientNotificationID(clientID string) string {
	return DB.QueryRowSQL("select device_id from "+CONSTANT.ClientsTable+" where client_id = ?", clientID)
}

// GetCounsellorNotificationID -
func GetCounsellorNotificationID(counsellorID string) string {
	return DB.QueryRowSQL("select device_id from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", counsellorID)
}

// GetListenerNotificationID -
func GetListenerNotificationID(listenerID string) string {
	return DB.QueryRowSQL("select device_id from "+CONSTANT.ListenersTable+" where listener_id = ?", listenerID)
}
