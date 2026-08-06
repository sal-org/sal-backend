package miscellaneous

import (
	"fmt"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// UploadFile godoc
// @Tags Miscellaneous
// @Summary Upload files like photos, certificates, aadhar etc
// @Router /upload [post]
// @Param file formData file true "File to be uploaded"
// @Param type formData string true "1(counsellor)/2(listener)/3(client)/4(therapist)"
// @Accept multipart/form-data
// @Security JWTAuth
// @Produce json
// @Success 200
func UploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	s3Path := CONSTANT.MiscellaneousS3Path
	switch r.FormValue("type") {
	case CONSTANT.CounsellorType:
		s3Path = CONSTANT.CounsellorS3Path
	case CONSTANT.ListenerType:
		s3Path = CONSTANT.ListenerS3Path
	case CONSTANT.ClientType:
		s3Path = CONSTANT.ClientS3Path
	case CONSTANT.TherapistType:
		s3Path = CONSTANT.TherapistS3Path
	}

	var fileName string
	// file upload
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("UploadFile", err)
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	if file != nil {
		defer file.Close()

		name, uploaded := UTIL.UploadToS3(CONFIG.S3Bucket, s3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(handler.Filename), CONSTANT.S3PublicRead, file)
		if !uploaded {
			fmt.Println("UploadFile", err)
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name
	}

	urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, fileName, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
	fileName = endPointURL

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	s3Path := CONSTANT.MiscellaneousS3Path
	switch r.FormValue("type") {
	case CONSTANT.CounsellorType:
		s3Path = CONSTANT.CounsellorS3Path
	case CONSTANT.ListenerType:
		s3Path = CONSTANT.ListenerS3Path
	case CONSTANT.ClientType:
		s3Path = CONSTANT.ClientS3Path
	case CONSTANT.TherapistType:
		s3Path = CONSTANT.TherapistS3Path
	}

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, s3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
