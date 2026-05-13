package b2c

import (
	"encoding/json"
	"log"
	"net/http"
	CONSTANT "salbackend/constant"
	UTIL "salbackend/util"

	"github.com/gorilla/mux"
)

// LoadAdminRoutes - load all web routes with web prefix
func LoadWebB2CRoutes(router *mux.Router) {
	webRoutes := router.PathPrefix("/web").Subrouter()

	webRoutes.HandleFunc("", CreateUserHandler).Methods("POST")

	// client login
	// webClientRoutes.HandleFunc("/access_code", CheckAccessCode).Methods("POST")
	// webClientRoutes.HandleFunc("/send_otp", SendOTPWithCorporateEmail).Methods("POST")
	// webClientRoutes.HandleFunc("/verify_otp", VerifyOTPWithCorporateEmail).Methods("POST")

	// // client profile
	// webClientRoutes.HandleFunc("/profile/get", ProfileGet).Methods("POST")
	// webClientRoutes.HandleFunc("/profile/add", ProfileAdd).Methods("POST")
	// webClientRoutes.HandleFunc("/profile/update", ProfileUpdate).Methods("POST")

	// // client search
	// webClientRoutes.HandleFunc("/search", ListSearch).Methods("POST")

	// // client appointment
	// webClientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/detail", AppointmentDetail).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/reschedule", AppointmentReschedule).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/cancel", AppointmentCancel).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/rate", AppointmentRatingAdd).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/start", AppointmentStart).Methods("POST")
	// webClientRoutes.HandleFunc("/appointment/end", AppointmentEnd).Methods("POST")

	// // client agora token
	// webClientRoutes.HandleFunc("/agora/token", GenerateAgoraToken).Methods("POST")

	// // client therapist
	// webClientRoutes.HandleFunc("/therapist/get", TherapistProfile).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/slots", TherapistSlots).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/order", CorporateCounsellorOrderCreate).Methods("POST")
	// webClientRoutes.HandleFunc("/therapist/complete", CorporateCounsellorOrderPaymentComplete).Methods("POST")

}

// CreateUserHandler handles user creation with detailed error logging
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	// Recover from any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming payload
	var encryptedPayload struct {
		Payload string `json:"payload"`
	}

	// Decode JSON with error logging
	err := json.NewDecoder(r.Body).Decode(&encryptedPayload)
	if err != nil {
		log.Printf("JSON decoding error: %v", err)
		http.Error(w, "Invalid payload format", http.StatusBadRequest)
		return
	}

	decrypted, _ := UTIL.DecryptPayload(encryptedPayload.Payload, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB_PROD, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB_PROD)

	switch decrypted["path"] {
	case "/book_demo":
		ClientBookDemo(w, r, decrypted)
	case "/access_code":
		CheckAccessCode(w, r, decrypted)
	case "/send_otp":
		SendOTPWithCorporateEmail(w, r, decrypted)
	case "/verify_otp":
		VerifyOTPWithCorporateEmail(w, r, decrypted)
	case "/profile/get":
		ProfileGet(w, r, decrypted)
	case "/profile/add":
		ProfileAdd(w, r, decrypted)
	case "/profile/update":
		ProfileUpdate(w, r, decrypted)
	case "/search":
		ListSearch(w, r, decrypted)
	case "/appointment/upcoming":
		AppointmentsUpcoming(w, r, decrypted)
	case "/appointment/past":
		AppointmentsPast(w, r, decrypted)
	case "/appointment/detail":
		AppointmentDetail(w, r, decrypted)
	case "/appointment/reschedule":
		AppointmentReschedule(w, r, decrypted)
	case "/appointment/cancel":
		AppointmentCancel(w, r, decrypted)
	case "/appointment/rate":
		AppointmentRatingAdd(w, r, decrypted)
	case "/appointment/start":
		AppointmentStart(w, r, decrypted)
	case "/appointment/end":
		AppointmentEnd(w, r, decrypted)
	case "/agora/token":
		GenerateAgoraToken(w, r, decrypted)
	case "/therapist/get":
		TherapistProfile(w, r, decrypted)
	case "/therapist/slots":
		TherapistSlots(w, r, decrypted)
	case "/therapist/order":
		CorporateCounsellorOrderCreate(w, r, decrypted)
	case "/therapist/complete":
		CorporateCounsellorOrderPaymentComplete(w, r, decrypted)
	default:
		w.Header().Set("Status", "200")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Page Not Found")
	}

}
