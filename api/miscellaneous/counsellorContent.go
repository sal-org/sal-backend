package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// Counsellor Content godoc
// @Tags Miscellaneous
// @Summary Get contents
// @Router /counsellor-content [get]
// @Param therapist_id query string true "Logged in therapists ID (counsellor_id/therapist_id)"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListCounsellorContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where counsellor_id = ? and training = 0 and status = 1 order by created_at desc limit 20", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellor_content"] = counsellor

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
