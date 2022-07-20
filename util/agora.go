package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	Model "salbackend/model"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

func GenerateAgoraRTCToken(channelName string, roleStr string, uidStr string, expireTime uint32) (result string, err error) {

	// var appID, appCertificate string
	var role rtctokenbuilder.Role

	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else if roleStr == "subscriber" {
		role = rtctokenbuilder.RoleSubscriber
	} else {
		role = rtctokenbuilder.RoleAttendee
	}

	// appID = CONSTANT.AGORA_APP_ID
	// appCertificate = CONSTANT.AGORA_APP_CERTIFICATE
	result, err = rtctokenbuilder.BuildTokenWithUserAccount(CONFIG.AGORA_APP_ID, CONFIG.AGORA_APP_CERTIFICATE, channelName, uidStr, role, expireTime)

	return result, err

}

func BasicAuthorization(channelName, uid string) (string, error) {
	customerKey := CONFIG.AGORA_Customer_Key
	// Customer secret
	customerSecret := CONFIG.AGORA_Customer_Secret

	// Concatenate customer key and customer secret and use base64 to encode the concatenated string
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	body := Model.PostRequestForAgora{
		CName: channelName,
		Uid:   uid,
		ClientRequest: Model.ClientRequestS{
			Region:              "AP",
			ResourceExpiredHour: 24,
			Scene:               0,
		},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", CONSTANT.AgoraURL+"/apps/"+CONFIG.AGORA_APP_ID+"/cloud_recording/acquire", payloadBuf)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	bodyy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return "", err
	}

	return string(mapp["resourceId"]), nil

}

func AgoraRecordingCallStart(uid, channelName, token, resourceid string) (string, error) {
	customerKey := CONFIG.AGORA_Customer_Key
	// Customer secret
	customerSecret := CONFIG.AGORA_Customer_Secret

	// Concatenate customer key and customer secret and use base64 to encode the concatenated string
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	body := Model.AgoraCallStartModel{
		Uid:   uid,
		CName: channelName,
		ClientRequest: Model.ClientRequestForStartCall{
			Token: token,
			RecordingConfig: Model.RecordingConfigModel{
				MaxIdleTime: 660,
				// StreamMode:        "standard",
				StreamTypes: 2,
				ChannelType: 0,
				// SubscribeUidGroup: 0,
			},
			RecordingFileConfig: Model.RecordingFileConfigModel{
				AvFileType: []string{"hls", "mp4"},
			},
			StorageConfig: Model.StorageConfigModel{
				AccessKey:      CONFIG.AWSAccesKey,
				Region:         14,
				Bucket:         "clove-cloud-recording",
				SecretKey:      CONFIG.AWSSecretKey,
				Vendor:         1,
				FileNamePrefix: []string{"recordingfile"},
				ExtensionParams: Model.ExtensionParamsModel{
					Tag: "public",
				},
			},
		},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", CONSTANT.AgoraURL+"/apps/"+CONFIG.AGORA_APP_ID+"/cloud_recording/resourceid/"+resourceid+"/mode/mix/start", payloadBuf)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	bodyy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return "", err
	}

	return string(mapp["sid"]), nil

}

func AgoraRecordingCallStop(uid, channelName, resourceid, sid string) (string, string, error) {
	customerKey := CONFIG.AGORA_Customer_Key
	// Customer secret
	customerSecret := CONFIG.AGORA_Customer_Secret

	// Concatenate customer key and customer secret and use base64 to encode the concatenated string
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	body := Model.AgoraCallStopModel{
		CName: channelName,
		Uid:   uid,
		ClientRequest: Model.ClientRequestForStopCall{
			Async_stop: false,
		},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", CONSTANT.AgoraURL+"/apps/"+CONFIG.AGORA_APP_ID+"/cloud_recording/resourceid/"+resourceid+"/sid/"+sid+"/mode/mix/stop", payloadBuf)

	if err != nil {
		return "", "", err
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	bodyy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var mapp Model.AgoraFileNameModel
	var fileNameInMP4, fileNameInM3U8 string

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		fmt.Println("eerror")
		return "", "", err
	}

	if mapp.Code == 200 {
		fileNameInMP4 = mapp.Body.ServerResponse.FileList[0].FileName
		fileNameInM3U8 = mapp.Body.ServerResponse.FileList[1].FileName
	}

	return fileNameInMP4, fileNameInM3U8, nil
}
