package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

func ListenerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get listeners
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " (first_name like '%%"+val[0]+"%%' or last_name like '%%"+val[0]+"%%') ")
			}
		case "phone":
			if len(val[0]) > 0 {
				wheres = append(wheres, " phone = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "email":
			if len(val[0]) > 0 {
				wheres = append(wheres, " email = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "listener_id":
			wheres = append(wheres, " listener_id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	listeners, status, ok := DB.SelectProcess("select * from "+CONSTANT.ListenersTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of listeners
	listenersCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.ListenersTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["listeners"] = listeners
	response["listeners_count"] = listenersCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(listenersCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ListenerUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	listenerID, ok := UTIL.Required(r.FormValue("listener_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, listenerID, CONSTANT.ShowDialog, response)
		return
	}

	body := MODEL.UpdateListenerProfileRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if order_id exists
	if !DB.CheckIfExists(CONSTANT.ListenersTable, map[string]string{"listener_id": listenerID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add listener
	listener := map[string]string{}
	listener["first_name"] = body.FirstName
	listener["last_name"] = body.LastName
	listener["phone"] = body.Phone
	listener["email"] = body.Email
	listener["gender"] = body.Gender
	listener["occupation"] = body.Occupation
	listener["age_group"] = body.AgeGroup
	listener["about"] = body.About
	listener["status"] = body.Status
	listener["modified_by"] = id
	listener["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("listener_id")}, listener)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
