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
		break
	case CONSTANT.ListenerType:
		s3Path = CONSTANT.ListenerS3Path
		break
	case CONSTANT.ClientType:
		s3Path = CONSTANT.ClientS3Path
		break
	case CONSTANT.TherapistType:
		s3Path = CONSTANT.TherapistS3Path
		break
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

		name, uploaded := UTIL.UploadToS3(CONFIG.S3Bucket, s3Path, CONFIG.S3AccesKey, CONFIG.S3SecretKey, CONFIG.S3Region, filepath.Ext(handler.Filename), CONSTANT.S3PublicRead, file)
		if !uploaded {
			fmt.Println("UploadFile", err)
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name
	}

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
