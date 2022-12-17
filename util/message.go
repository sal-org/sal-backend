package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"
)

func RemoveMessage(tagID, phone string) {
	// delete any previous message, if any
	DB.DeleteSQL(CONSTANT.MessagesTable, map[string]string{"tag_id": tagID, "phone": phone})
}

// SendMessage - send text message using message provider. now : true - send now without background workers
func SendMessage(text, route, phone, sent_at, tagID string, now bool) {
	if strings.Contains(text, "###") { // check if message variables are replaced
		return
	}

	// add data to messages
	message := map[string]string{}
	message["text"] = text
	message["route"] = route
	message["phone"] = phone
	message["tag_id"] = tagID
	message["send_at"] = sent_at
	if now {
		// set message sent status as sent if now is true
		message["status"] = CONSTANT.MessageSent
		message["message_status"] = CONSTANT.MessageSent
		sendMessageCoreFactor(text, route, phone)
	} else {
		message["status"] = CONSTANT.MessageInProgress
		message["message_status"] = CONSTANT.MessageInProgress
	}
	message["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.MessagesTable, CONSTANT.MessagesDigits, message, "message_id")

}

func sendMessageCoreFactor(text, route, phone string) {
	resp, err := http.Get(buildMessageURL(text, route, phone))
	if err != nil {
		fmt.Println("sendMessageCoreFactor", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("sendMessageCoreFactor", err)
		return
	}

	fmt.Println("sendMessageCoreFactor", string(body))
}

func buildMessageURL(text, route, phone string) string {
	u, _ := url.Parse(CONSTANT.CorefactorsSendSMSEndpoint)

	v := url.Values{}
	v.Add("text", text)
	v.Add("key", CONSTANT.CorefactorsAPIKey)
	v.Add("to", phone)
	v.Add("route", route)
	v.Add("from", CONSTANT.TextMessageFrom)

	u.RawQuery = v.Encode()

	fmt.Println(u.String())
	return u.String()
}
