package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

	UTIL "salbackend/util"
)

func AvailabilityGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get counsellor availability hours
	availability, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsScheduleTable+" where counsellor_id = ? ", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["availability"] = availability
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AvailabilityUpdate godoc
// @Tags Counsellor Availability
// @Summary Update counsellor availability hours
// @Router /counsellor/availability [put]
// @Param counsellor_id query string true "Counsellor ID to update availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var newValue = false

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBodyInListMap(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, day := range body {

		if strings.EqualFold(day["status"], "0") {

			availabilityDates, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InPersonSLotsScheduleTable+" where counsellor_id = ? and date = ? order by date", r.FormValue("counsellor_id"), day["date"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			DB.DeleteSQL(CONSTANT.InPersonSLotsScheduleTable, map[string]string{"id": day["id"]})

			availabileDates, status, ok := DB.SelectProcess("select id from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and `date` = ?", r.FormValue("counsellor_id"), day["date"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			slots := map[string]string{}

			if availabilityDates[0]["ctn"] == "1" {

				slots["available"] = "0"
				slots["0"] = "0"
				slots["1"] = "0"
				slots["2"] = "0"
				slots["3"] = "0"
				slots["4"] = "0"
				slots["5"] = "0"
				slots["6"] = "0"
				slots["7"] = "0"
				slots["8"] = "0"
				slots["9"] = "0"
				slots["10"] = "0"
				slots["11"] = "0"
				slots["12"] = "0"
				slots["13"] = "0"
				slots["14"] = "0"
				slots["15"] = "0"
				slots["16"] = "0"
				slots["17"] = "0"
				slots["18"] = "0"
				slots["19"] = "0"
				slots["20"] = "0"
				slots["21"] = "0"
				slots["22"] = "0"
				slots["23"] = "0"
				slots["24"] = "0"
				slots["25"] = "0"
				slots["26"] = "0"
				slots["27"] = "0"
				slots["28"] = "0"
				slots["29"] = "0"
				slots["30"] = "0"
				slots["31"] = "0"
				slots["32"] = "0"
				slots["33"] = "0"
				slots["34"] = "0"
				slots["35"] = "0"
				slots["36"] = "0"
				slots["37"] = "0"
				slots["38"] = "0"
				slots["39"] = "0"
				slots["40"] = "0"
				slots["41"] = "0"
				slots["42"] = "0"
				slots["43"] = "0"
				slots["44"] = "0"
				slots["45"] = "0"
				slots["46"] = "0"
				slots["47"] = "0"

			} else {

				slots["available"] = "1"

				if day["0"] == "1" {
					slots["0"] = "0"
				}

				if day["1"] == "1" {
					slots["1"] = "0"
				}

				if day["2"] == "1" {
					slots["2"] = "0"
				}

				if day["3"] == "1" {
					slots["3"] = "0"
				}

				if day["4"] == "1" {
					slots["4"] = "0"
				}

				if day["5"] == "1" {
					slots["5"] = "0"
				}

				if day["6"] == "1" {
					slots["6"] = "0"
				}

				if day["7"] == "1" {
					slots["7"] = "0"
				}

				if day["8"] == "1" {
					slots["8"] = "0"
				}

				if day["9"] == "1" {
					slots["9"] = "0"
				}

				if day["10"] == "1" {
					slots["10"] = "0"
				}

				if day["11"] == "1" {
					slots["11"] = "0"
				}

				if day["12"] == "1" {
					slots["12"] = "0"
				}

				if day["13"] == "1" {
					slots["13"] = "0"
				}

				if day["14"] == "1" {
					slots["14"] = "0"
				}

				if day["15"] == "1" {
					slots["15"] = "0"
				}

				if day["16"] == "1" {
					slots["16"] = "0"
				}

				if day["17"] == "1" {
					slots["17"] = "0"
				}

				if day["18"] == "1" {
					slots["18"] = "0"
				}

				if day["19"] == "1" {
					slots["19"] = "0"
				}

				if day["20"] == "1" {
					slots["20"] = "0"
				}

				if day["21"] == "1" {
					slots["21"] = "0"
				}

				if day["22"] == "1" {
					slots["22"] = "0"
				}

				if day["23"] == "1" {
					slots["23"] = "0"
				}

				if day["24"] == "1" {
					slots["24"] = "0"
				}

				if day["25"] == "1" {
					slots["25"] = "0"
				}

				if day["26"] == "1" {
					slots["26"] = "0"
				}

				if day["27"] == "1" {
					slots["27"] = "0"
				}

				if day["28"] == "1" {
					slots["28"] = "0"
				}

				if day["29"] == "1" {
					slots["29"] = "0"
				}

				if day["30"] == "1" {
					slots["30"] = "0"
				}

				if day["31"] == "1" {
					slots["31"] = "0"
				}

				if day["32"] == "1" {
					slots["32"] = "0"
				}

				if day["33"] == "1" {
					slots["33"] = "0"
				}

				if day["34"] == "1" {
					slots["34"] = "0"
				}

				if day["35"] == "1" {
					slots["35"] = "0"
				}

				if day["36"] == "1" {
					slots["36"] = "0"
				}

				if day["37"] == "1" {
					slots["37"] = "0"
				}

				if day["38"] == "1" {
					slots["38"] = "0"
				}

				if day["39"] == "1" {
					slots["39"] = "0"
				}

				if day["40"] == "1" {
					slots["40"] = "0"
				}

				if day["41"] == "1" {
					slots["41"] = "0"
				}

				if day["42"] == "1" {
					slots["42"] = "0"
				}

				if day["43"] == "1" {
					slots["43"] = "0"
				}

				if day["44"] == "1" {
					slots["44"] = "0"
				}

				if day["45"] == "1" {
					slots["45"] = "0"
				}

				if day["46"] == "1" {
					slots["46"] = "0"
				}

				if day["47"] == "1" {
					slots["47"] = "0"
				}

			}

			for key, val := range slots {

				DB.ExecuteSQL("update "+CONSTANT.InPersonSLotsTable+" set `"+key+"` = "+val+" where id = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", availabileDates[0]["id"]) // dont update already booked slots

			}
		} else {
			availabileDates, status, ok := DB.SelectProcess("select date,id from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and `date` = ?", r.FormValue("counsellor_id"), day["date"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			if len(availabileDates) == 0 {
				newValue = true
			}

			slotsDate := map[string]string{}

			slotsDate["availability_status"] = day["availability_status"]
			slotsDate["counsellor_id"] = day["counsellor_id"]
			slotsDate["date"] = day["date"]
			slotsDate["format"] = day["format"]
			slotsDate["break"] = day["break"]
			slotsDate["partner_name"] = day["partner_name"]
			slotsDate["partner_location"] = day["partner_location"]
			slotsDate["counselling_address"] = day["counselling_address"]
			slotsDate["status"] = day["status"]
			slotsDate["0"] = day["0"]
			slotsDate["1"] = day["1"]
			slotsDate["2"] = day["2"]
			slotsDate["3"] = day["3"]
			slotsDate["4"] = day["4"]
			slotsDate["5"] = day["5"]
			slotsDate["6"] = day["6"]
			slotsDate["7"] = day["7"]
			slotsDate["8"] = day["8"]
			slotsDate["9"] = day["9"]
			slotsDate["10"] = day["10"]
			slotsDate["11"] = day["11"]
			slotsDate["12"] = day["12"]
			slotsDate["13"] = day["13"]
			slotsDate["14"] = day["14"]
			slotsDate["15"] = day["15"]
			slotsDate["16"] = day["16"]
			slotsDate["17"] = day["17"]
			slotsDate["18"] = day["18"]
			slotsDate["19"] = day["19"]
			slotsDate["20"] = day["20"]
			slotsDate["21"] = day["21"]
			slotsDate["22"] = day["22"]
			slotsDate["23"] = day["23"]
			slotsDate["24"] = day["24"]
			slotsDate["25"] = day["25"]
			slotsDate["26"] = day["26"]
			slotsDate["27"] = day["27"]
			slotsDate["28"] = day["28"]
			slotsDate["29"] = day["29"]
			slotsDate["30"] = day["30"]
			slotsDate["31"] = day["31"]
			slotsDate["32"] = day["32"]
			slotsDate["33"] = day["33"]
			slotsDate["34"] = day["34"]
			slotsDate["35"] = day["35"]
			slotsDate["36"] = day["36"]
			slotsDate["37"] = day["37"]
			slotsDate["38"] = day["38"]
			slotsDate["39"] = day["39"]
			slotsDate["40"] = day["40"]
			slotsDate["41"] = day["41"]
			slotsDate["42"] = day["42"]
			slotsDate["43"] = day["43"]
			slotsDate["44"] = day["44"]
			slotsDate["45"] = day["45"]
			slotsDate["46"] = day["46"]
			slotsDate["47"] = day["47"]

			slots := map[string]string{}

			if len(day["id"]) > 0 {
				DB.UpdateSQL(CONSTANT.InPersonSLotsScheduleTable, map[string]string{"id": day["id"]}, slotsDate)

				availabilityDates, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InPersonSLotsScheduleTable+" where counsellor_id = ? and date = ? order by date", r.FormValue("counsellor_id"), day["date"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				if availabilityDates[0]["ctn"] == "1" {

					slots["available"] = day["availability_status"]
					slots["0"] = day["0"]
					slots["1"] = day["1"]
					slots["2"] = day["2"]
					slots["3"] = day["3"]
					slots["4"] = day["4"]
					slots["5"] = day["5"]
					slots["6"] = day["6"]
					slots["7"] = day["7"]
					slots["8"] = day["8"]
					slots["9"] = day["9"]
					slots["10"] = day["10"]
					slots["11"] = day["11"]
					slots["12"] = day["12"]
					slots["13"] = day["13"]
					slots["14"] = day["14"]
					slots["15"] = day["15"]
					slots["16"] = day["16"]
					slots["17"] = day["17"]
					slots["18"] = day["18"]
					slots["19"] = day["19"]
					slots["20"] = day["20"]
					slots["21"] = day["21"]
					slots["22"] = day["22"]
					slots["23"] = day["23"]
					slots["24"] = day["24"]
					slots["25"] = day["25"]
					slots["26"] = day["26"]
					slots["27"] = day["27"]
					slots["28"] = day["28"]
					slots["29"] = day["29"]
					slots["30"] = day["30"]
					slots["31"] = day["31"]
					slots["32"] = day["32"]
					slots["33"] = day["33"]
					slots["34"] = day["34"]
					slots["35"] = day["35"]
					slots["36"] = day["36"]
					slots["37"] = day["37"]
					slots["38"] = day["38"]
					slots["39"] = day["39"]
					slots["40"] = day["40"]
					slots["41"] = day["41"]
					slots["42"] = day["42"]
					slots["43"] = day["43"]
					slots["44"] = day["44"]
					slots["45"] = day["45"]
					slots["46"] = day["46"]
					slots["47"] = day["47"]

				} else {

					if day["availability_status"] == "0" {

						slots["available"] = "1"

						if day["0"] == "1" {
							slots["0"] = "0"
						}

						if day["1"] == "1" {
							slots["1"] = "0"
						}

						if day["2"] == "1" {
							slots["2"] = "0"
						}

						if day["3"] == "1" {
							slots["3"] = "0"
						}

						if day["4"] == "1" {
							slots["4"] = "0"
						}

						if day["5"] == "1" {
							slots["5"] = "0"
						}

						if day["6"] == "1" {
							slots["6"] = "0"
						}

						if day["7"] == "1" {
							slots["7"] = "0"
						}

						if day["8"] == "1" {
							slots["8"] = "0"
						}

						if day["9"] == "1" {
							slots["9"] = "0"
						}

						if day["10"] == "1" {
							slots["10"] = "0"
						}

						if day["11"] == "1" {
							slots["11"] = "0"
						}

						if day["12"] == "1" {
							slots["12"] = "0"
						}

						if day["13"] == "1" {
							slots["13"] = "0"
						}

						if day["14"] == "1" {
							slots["14"] = "0"
						}

						if day["15"] == "1" {
							slots["15"] = "0"
						}

						if day["16"] == "1" {
							slots["16"] = "0"
						}

						if day["17"] == "1" {
							slots["17"] = "0"
						}

						if day["18"] == "1" {
							slots["18"] = "0"
						}

						if day["19"] == "1" {
							slots["19"] = "0"
						}

						if day["20"] == "1" {
							slots["20"] = "0"
						}

						if day["21"] == "1" {
							slots["21"] = "0"
						}

						if day["22"] == "1" {
							slots["22"] = "0"
						}

						if day["23"] == "1" {
							slots["23"] = "0"
						}

						if day["24"] == "1" {
							slots["24"] = "0"
						}

						if day["25"] == "1" {
							slots["25"] = "0"
						}

						if day["26"] == "1" {
							slots["26"] = "0"
						}

						if day["27"] == "1" {
							slots["27"] = "0"
						}

						if day["28"] == "1" {
							slots["28"] = "0"
						}

						if day["29"] == "1" {
							slots["29"] = "0"
						}

						if day["30"] == "1" {
							slots["30"] = "0"
						}

						if day["31"] == "1" {
							slots["31"] = "0"
						}

						if day["32"] == "1" {
							slots["32"] = "0"
						}

						if day["33"] == "1" {
							slots["33"] = "0"
						}

						if day["34"] == "1" {
							slots["34"] = "0"
						}

						if day["35"] == "1" {
							slots["35"] = "0"
						}

						if day["36"] == "1" {
							slots["36"] = "0"
						}

						if day["37"] == "1" {
							slots["37"] = "0"
						}

						if day["38"] == "1" {
							slots["38"] = "0"
						}

						if day["39"] == "1" {
							slots["39"] = "0"
						}

						if day["40"] == "1" {
							slots["40"] = "0"
						}

						if day["41"] == "1" {
							slots["41"] = "0"
						}

						if day["42"] == "1" {
							slots["42"] = "0"
						}

						if day["43"] == "1" {
							slots["43"] = "0"
						}

						if day["44"] == "1" {
							slots["44"] = "0"
						}

						if day["45"] == "1" {
							slots["45"] = "0"
						}

						if day["46"] == "1" {
							slots["46"] = "0"
						}

						if day["47"] == "1" {
							slots["47"] = "0"
						}

					} else {

						slots["available"] = "1"

						if day["0"] == "1" {
							slots["0"] = "1"
						}

						if day["1"] == "1" {
							slots["1"] = "1"
						}

						if day["2"] == "1" {
							slots["2"] = "1"
						}

						if day["3"] == "1" {
							slots["3"] = "1"
						}

						if day["4"] == "1" {
							slots["4"] = "1"
						}

						if day["5"] == "1" {
							slots["5"] = "1"
						}

						if day["6"] == "1" {
							slots["6"] = "1"
						}

						if day["7"] == "1" {
							slots["7"] = "1"
						}

						if day["8"] == "1" {
							slots["8"] = "1"
						}

						if day["9"] == "1" {
							slots["9"] = "1"
						}

						if day["10"] == "1" {
							slots["10"] = "1"
						}

						if day["11"] == "1" {
							slots["11"] = "1"
						}

						if day["12"] == "1" {
							slots["12"] = "1"
						}

						if day["13"] == "1" {
							slots["13"] = "1"
						}

						if day["14"] == "1" {
							slots["14"] = "1"
						}

						if day["15"] == "1" {
							slots["15"] = "1"
						}

						if day["16"] == "1" {
							slots["16"] = "1"
						}

						if day["17"] == "1" {
							slots["17"] = "1"
						}

						if day["18"] == "1" {
							slots["18"] = "1"
						}

						if day["19"] == "1" {
							slots["19"] = "1"
						}

						if day["20"] == "1" {
							slots["20"] = "1"
						}

						if day["21"] == "1" {
							slots["21"] = "1"
						}

						if day["22"] == "1" {
							slots["22"] = "1"
						}

						if day["23"] == "1" {
							slots["23"] = "1"
						}

						if day["24"] == "1" {
							slots["24"] = "1"
						}

						if day["25"] == "1" {
							slots["25"] = "1"
						}

						if day["26"] == "1" {
							slots["26"] = "1"
						}

						if day["27"] == "1" {
							slots["27"] = "1"
						}

						if day["28"] == "1" {
							slots["28"] = "1"
						}

						if day["29"] == "1" {
							slots["29"] = "1"
						}

						if day["30"] == "1" {
							slots["30"] = "1"
						}

						if day["31"] == "1" {
							slots["31"] = "1"
						}

						if day["32"] == "1" {
							slots["32"] = "1"
						}

						if day["33"] == "1" {
							slots["33"] = "1"
						}

						if day["34"] == "1" {
							slots["34"] = "1"
						}

						if day["35"] == "1" {
							slots["35"] = "1"
						}

						if day["36"] == "1" {
							slots["36"] = "1"
						}

						if day["37"] == "1" {
							slots["37"] = "1"
						}

						if day["38"] == "1" {
							slots["38"] = "1"
						}

						if day["39"] == "1" {
							slots["39"] = "1"
						}

						if day["40"] == "1" {
							slots["40"] = "1"
						}

						if day["41"] == "1" {
							slots["41"] = "1"
						}

						if day["42"] == "1" {
							slots["42"] = "1"
						}

						if day["43"] == "1" {
							slots["43"] = "1"
						}

						if day["44"] == "1" {
							slots["44"] = "1"
						}

						if day["45"] == "1" {
							slots["45"] = "1"
						}

						if day["46"] == "1" {
							slots["46"] = "1"
						}

						if day["47"] == "1" {
							slots["47"] = "1"
						}

					}

				}

			} else {
				// newly added schedule
				DB.InsertSQL(CONSTANT.InPersonSLotsScheduleTable, slotsDate)

				availabilityDates, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InPersonSLotsScheduleTable+" where counsellor_id = ? and date = ? order by date", r.FormValue("counsellor_id"), day["date"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				if availabilityDates[0]["ctn"] == "1" {

					slots["available"] = day["availability_status"]
					slots["0"] = day["0"]
					slots["1"] = day["1"]
					slots["2"] = day["2"]
					slots["3"] = day["3"]
					slots["4"] = day["4"]
					slots["5"] = day["5"]
					slots["6"] = day["6"]
					slots["7"] = day["7"]
					slots["8"] = day["8"]
					slots["9"] = day["9"]
					slots["10"] = day["10"]
					slots["11"] = day["11"]
					slots["12"] = day["12"]
					slots["13"] = day["13"]
					slots["14"] = day["14"]
					slots["15"] = day["15"]
					slots["16"] = day["16"]
					slots["17"] = day["17"]
					slots["18"] = day["18"]
					slots["19"] = day["19"]
					slots["20"] = day["20"]
					slots["21"] = day["21"]
					slots["22"] = day["22"]
					slots["23"] = day["23"]
					slots["24"] = day["24"]
					slots["25"] = day["25"]
					slots["26"] = day["26"]
					slots["27"] = day["27"]
					slots["28"] = day["28"]
					slots["29"] = day["29"]
					slots["30"] = day["30"]
					slots["31"] = day["31"]
					slots["32"] = day["32"]
					slots["33"] = day["33"]
					slots["34"] = day["34"]
					slots["35"] = day["35"]
					slots["36"] = day["36"]
					slots["37"] = day["37"]
					slots["38"] = day["38"]
					slots["39"] = day["39"]
					slots["40"] = day["40"]
					slots["41"] = day["41"]
					slots["42"] = day["42"]
					slots["43"] = day["43"]
					slots["44"] = day["44"]
					slots["45"] = day["45"]
					slots["46"] = day["46"]
					slots["47"] = day["47"]

				} else {

					slots["available"] = "1"

					if day["0"] == "1" {
						slots["0"] = "1"
					}

					if day["1"] == "1" {
						slots["1"] = "1"
					}

					if day["2"] == "1" {
						slots["2"] = "1"
					}

					if day["3"] == "1" {
						slots["3"] = "1"
					}

					if day["4"] == "1" {
						slots["4"] = "1"
					}

					if day["5"] == "1" {
						slots["5"] = "1"
					}

					if day["6"] == "1" {
						slots["6"] = "1"
					}

					if day["7"] == "1" {
						slots["7"] = "1"
					}

					if day["8"] == "1" {
						slots["8"] = "1"
					}

					if day["9"] == "1" {
						slots["9"] = "1"
					}

					if day["10"] == "1" {
						slots["10"] = "1"
					}

					if day["11"] == "1" {
						slots["11"] = "1"
					}

					if day["12"] == "1" {
						slots["12"] = "1"
					}

					if day["13"] == "1" {
						slots["13"] = "1"
					}

					if day["14"] == "1" {
						slots["14"] = "1"
					}

					if day["15"] == "1" {
						slots["15"] = "1"
					}

					if day["16"] == "1" {
						slots["16"] = "1"
					}

					if day["17"] == "1" {
						slots["17"] = "1"
					}

					if day["18"] == "1" {
						slots["18"] = "1"
					}

					if day["19"] == "1" {
						slots["19"] = "1"
					}

					if day["20"] == "1" {
						slots["20"] = "1"
					}

					if day["21"] == "1" {
						slots["21"] = "1"
					}

					if day["22"] == "1" {
						slots["22"] = "1"
					}

					if day["23"] == "1" {
						slots["23"] = "1"
					}

					if day["24"] == "1" {
						slots["24"] = "1"
					}

					if day["25"] == "1" {
						slots["25"] = "1"
					}

					if day["26"] == "1" {
						slots["26"] = "1"
					}

					if day["27"] == "1" {
						slots["27"] = "1"
					}

					if day["28"] == "1" {
						slots["28"] = "1"
					}

					if day["29"] == "1" {
						slots["29"] = "1"
					}

					if day["30"] == "1" {
						slots["30"] = "1"
					}

					if day["31"] == "1" {
						slots["31"] = "1"
					}

					if day["32"] == "1" {
						slots["32"] = "1"
					}

					if day["33"] == "1" {
						slots["33"] = "1"
					}

					if day["34"] == "1" {
						slots["34"] = "1"
					}

					if day["35"] == "1" {
						slots["35"] = "1"
					}

					if day["36"] == "1" {
						slots["36"] = "1"
					}

					if day["37"] == "1" {
						slots["37"] = "1"
					}

					if day["38"] == "1" {
						slots["38"] = "1"
					}

					if day["39"] == "1" {
						slots["39"] = "1"
					}

					if day["40"] == "1" {
						slots["40"] = "1"
					}

					if day["41"] == "1" {
						slots["41"] = "1"
					}

					if day["42"] == "1" {
						slots["42"] = "1"
					}

					if day["43"] == "1" {
						slots["43"] = "1"
					}

					if day["44"] == "1" {
						slots["44"] = "1"
					}

					if day["45"] == "1" {
						slots["45"] = "1"
					}

					if day["46"] == "1" {
						slots["46"] = "1"
					}

					if day["47"] == "1" {
						slots["47"] = "1"
					}

				}

			}

			if newValue {
				slots["date"] = day["date"]
				slots["counsellor_id"] = r.FormValue("counsellor_id")
				DB.InsertSQL(CONSTANT.InPersonSLotsTable, slots)

			} else {
				for key, val := range slots {

					DB.ExecuteSQL("update "+CONSTANT.InPersonSLotsTable+" set `"+key+"` = "+val+" where id = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", availabileDates[0]["id"]) // dont update already booked slots

				}
			}

		}

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
