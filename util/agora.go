package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
				MaxIdleTime:     660,
				StreamTypes:     2,
				ChannelType:     0,
				VideoStreamType: 0,
				TranscodingConfigs: Model.TranscodingConfig{
					Height:           640,
					Width:            360,
					Fps:              15,
					MixedVideoLayout: 1,
				},
			},
			RecordingFileConfig: Model.RecordingFileConfigModel{
				AvFileType: []string{"hls", "mp4"},
			},
			StorageConfig: Model.StorageConfigModel{
				AccessKey:      CONFIG.AWSAccesKey,
				Region:         14,
				Bucket:         CONFIG.S3BUCKETAGORA,
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
	bodyy, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(bodyy))
	var result Model.AgoraCallStartResponse
	json.NewDecoder(res.Body).Decode(&result)

	// bodyy, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return "", err
	}

	fmt.Println(result)

	return mapp["sid"], nil

}

func AgoraUpateRecordingSettingCallStart(uid, channelName, sid, resourceid, uid1 string) (string, string, error) {
	customerKey := CONFIG.AGORA_Customer_Key
	// Customer secret
	customerSecret := CONFIG.AGORA_Customer_Secret

	var allUsers = []string{uid + "123", uid1 + "321"}

	fmt.Println(allUsers)

	// Concatenate customer key and customer secret and use base64 to encode the concatenated string
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	body := Model.ClientRequestForUpdateStartCall{
		UID:           uid,
		Cname:         channelName,
		ClientRequest: Model.ClientRequest{StreamSubscribe: Model.StreamSubscribe{AudioUIDList: Model.AudioUIDList{SubscribeAudioUids: allUsers}, VideoUIDList: Model.VideoUIDList{SubscribeVideoUids: allUsers}}},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	//https://api.agora.io/v1/apps/<appid>/cloud_recording/resourceid/<resourceid>/sid/<sid>/mode/<mode>/update
	req, err := http.NewRequest("POST", CONSTANT.AgoraURL+"/apps/"+CONFIG.AGORA_APP_ID+"/cloud_recording/resourceid/"+resourceid+"/sid/"+sid+"/mode/mix/update", payloadBuf)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	defer res.Body.Close()

	// var result Model.AgoraCallStartResponse
	// json.NewDecoder(res.Body).Decode(&result)

	bodyy, _ := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }
	fmt.Println(string(bodyy))
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return "", "", err
	}
	fmt.Println(mapp)
	return mapp["sid"], mapp["resourceId"], nil
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
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result Model.AgoraCallStopResponseModel
	json.NewDecoder(resp.Body).Decode(&result)

	codeStatus, _ := CallStatus(resourceid, sid)
	fmt.Println(codeStatus)
	fmt.Println(result)

	var fileNameInMP4, fileNameInM3U8 string

	if result.Code == 200 {

		// fmt.Println(result)
		// b, _ := json.Marshal(result)
		// fmt.Println(string(b))

		bodyy, _ := io.ReadAll(resp.Body)
		// if err != nil {
		// 	return "", "", err
		// }
		fmt.Println(string(bodyy))

		// var mapp Model.AgoraCallStopResponseModel

		// err = json.Unmarshal(b, &mapp)
		// if err != nil {
		// 	fmt.Println("eerror")
		// 	return "", "", err
		// }

		// fmt.Println(mapp)

		if len(result.Body.ServerResponse.FileList[0].FileName) == 0 {
			fileNameInM3U8 = ""
			fileNameInMP4 = ""
		} else {
			fileNameInMP4 = result.Body.ServerResponse.FileList[0].FileName
			// fmt.Println(fileNameInMP4)
			fileNameInM3U8 = result.Body.ServerResponse.FileList[1].FileName
			// fmt.Println(fileNameInM3U8)
		}

	} else {
		fileNameInM3U8 = ""
		fileNameInMP4 = ""
	}

	return fileNameInMP4, fileNameInM3U8, nil
}

func CallStatus(resourceid string, sid string) (Model.AgoraCallStatus, error) {

	customerKey := CONFIG.AGORA_Customer_Key
	// Customer secret
	customerSecret := CONFIG.AGORA_Customer_Secret

	// Concatenate customer key and customer secret and use base64 to encode the concatenated string
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	client := &http.Client{}
	req, err := http.NewRequest("GET", CONSTANT.AgoraURL+"/apps/"+CONFIG.AGORA_APP_ID+"/cloud_recording/resourceid/"+resourceid+"/sid/"+sid+"/mode/mix/query", nil)

	if err != nil {
		return Model.AgoraCallStatus{}, err
	}
	// Add Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return Model.AgoraCallStatus{}, err
	}
	defer resp.Body.Close()

	var result Model.AgoraCallStatus
	json.NewDecoder(resp.Body).Decode(&result)
	// // b, _ := json.Marshal(result)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	return result, nil
}
