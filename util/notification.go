package util

import (
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
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
