package util

import (
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"
)

// SendNotification - send notification using onesignal
func SendNotification(heading, content, personID, personType string) {
	if strings.Contains(content, "###") { // check if notification variables are replaced
		return
	}

	// add data to notifications
	notification := map[string]string{}
	notification["person_id"] = personID
	notification["title"] = heading
	notification["body"] = content
	notification["status"] = CONSTANT.NotificationActive
	notification["onesignal_id"] = GetNotificationID(personID, personType)
	if len(notification["onesignal_id"]) > 0 {
		notification["notification_status"] = CONSTANT.NotificationInProgress
	} else {
		// set notification sent status as sent if no onesignal id is available
		notification["notification_status"] = CONSTANT.NotificationSent
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
