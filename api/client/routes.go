package client

import "github.com/gorilla/mux"

// LoadClientRoutes - load all client routes with client prefix
func LoadClientRoutes(router *mux.Router) {
	clientRoutes := router.PathPrefix("/client").Subrouter()

	// appointment
	clientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/inperson_appointment/upcoming", InPersonAppointmentsUpcoming).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/appointment/slots", AppointmentSlotsUnused).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/inperson_appointment/past", InPersonAppointmentsPast).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/appointment", AppointmentDetail).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/inperson_appointment", InPersonAppointmentDetail).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/appointment", AppointmentBook).Methods("POST")
	clientRoutes.HandleFunc("/appointment", AppointmentReschedule).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/inperson_appointment", InPersonAppointmentReschedule).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/appointment", AppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")

	clientRoutes.HandleFunc("/inperson_appointment", InPersonAppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")

	clientRoutes.HandleFunc("/appointment/bulk", AppointmentBulkCancel).Queries(
		"appointment_slot_id", "{appointment_slot_id}",
	).Methods("DELETE")
	clientRoutes.HandleFunc("/appointment/rate", AppointmentRatingAdd).Methods("POST")

	clientRoutes.HandleFunc("/inperson_appointment/rate", InPersonAppointmentRatingAdd).Methods("POST")

	clientRoutes.HandleFunc("/appointment/download", DownloadReceipt).Queries(
		"invoice_id", "{invoice_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/appointment/cancellationreason", CancellationReason).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/inperson_appointment/cancellationreason", InPersonCancellationReason).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/inperson_appointment/no_show", InPersonAppointmentNoShow).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/appointment/agoratoken", GenerateAgoraToken).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/request", AppointmentRequest).Methods("POST")

	clientRoutes.HandleFunc("/inperson_appointment/request", InPersonAppointmentRequest).Methods("POST")

	clientRoutes.HandleFunc("/appointment/request", GetAppointmentRequest).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/inperson_appointment/request", GetInPersonAppointmentRequest).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/appointment/start", AppointmentStart).Queries(
		"appointment_id", "{appointment_id}",
		"uid", "{uid}",
	).Methods("PUT")
	clientRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
		"uid", "{uid}",
	).Methods("PUT")

	clientRoutes.HandleFunc("/appointment/coupon", CouponGet).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// assessment
	clientRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	clientRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	clientRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/assessment/download", AssessmentDownload).Queries(
		"assessment_result_id", "{assessment_result_id}",
	).Methods("GET")

	// counsellor
	clientRoutes.HandleFunc("/counsellor", CounsellorProfile).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/slots", CounsellorSlots).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// In Person slots
	clientRoutes.HandleFunc("/in_person_counsellor/slots", InPersonCounsellorSlots).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	clientRoutes.HandleFunc("/counsellor/order", CounsellorOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/counsellor/paymentcomplete", CounsellorOrderPaymentComplete).Methods("POST")

	// event
	clientRoutes.HandleFunc("/events", EventsList).Methods("GET")
	clientRoutes.HandleFunc("/events_inperson", InPersonEventsList).Methods("GET")
	clientRoutes.HandleFunc("/event", EventDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson", EventInPersonDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event/booked", EventsBooked).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson/booked", EventsBookedInPerson).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson/past", PastEventsInPerson).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson/request", EventsInPersonRequest).Methods("POST")
	clientRoutes.HandleFunc("/event_inperson/request", GetEventInPersonRequest).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson/cancel", EventsInPersonCancel).Methods("PUT")
	clientRoutes.HandleFunc("/event_inperson/rate", GetEventsInPersonRate).Queries(
		"user_id", "{user_id}",
		"order_id", "{order_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event_inperson/rate", EventsInPersonRate).Methods("POST")
	clientRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/event_inperson/order", EventOrderInPersonCreate).Methods("POST")
	clientRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// home
	clientRoutes.HandleFunc("/home", Home).Methods("GET")

	// cor login with email
	clientRoutes.HandleFunc("/coremail/sendotp", SendOTPWithCorporateEmail).Queries(
		"cor_email", "{cor_email}",
	).Methods("GET")

	// cor register with  email
	clientRoutes.HandleFunc("/coremail/send", SendOTPWithCorporateEmailForRegister).Queries(
		"cor_email", "{cor_email}",
	).Methods("GET")

	// check access code
	clientRoutes.HandleFunc("/check_code", CheckAccessCode).Queries(
		"access_code", "{access_code}",
	).Methods("GET")

	// get address
	clientRoutes.HandleFunc("/get_address", GetAddressForCorporateClient).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// check access token
	clientRoutes.HandleFunc("/check_access_token", CheckIfAccessTokenExpired).Methods("GET")

	// check phone and  email
	clientRoutes.HandleFunc("/check_uniqueid", CheckEmailANDPhone).Methods("GET")

	clientRoutes.HandleFunc("/coremail/verifyotp", VerifyOTPWithCorporateEmail).Queries(
		"cor_email", "{cor_email}",
		"otp", "{otp}",
		"device_id", "{device_id}",
	).Methods("GET")

	// depandent client send otp
	clientRoutes.HandleFunc("/depandent_client/sendotp", GetDenpendantClientOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")

	clientRoutes.HandleFunc("/depandent_client/verifyotp", VerifyOTPWithDependantClientEmail).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
		"device_id", "{device_id}",
	).Methods("GET")

	// phone number verification for corporate client family member
	clientRoutes.HandleFunc("/family_member/sendotp", SendOTPForForFamilyRegister).Queries(
		"family_phone_no", "{family_phone_no}",
		"family_email_id", "{family_email_id}",
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/family_member/verifyotp", VerifyOTPForRegisterFamilyMember).Queries(
		"family_phone_no", "{family_phone_no}",
		"otp", "{otp}",
	).Methods("GET")

	// login
	clientRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	clientRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
		"device_id", "{device_id}",
	).Methods("GET")
	clientRoutes.Path("/refresh-token").Queries(
		"client_id", "{client_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// listener
	clientRoutes.HandleFunc("/listener", ListenerProfile).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/slots", ListenerSlots).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/order", ListenerOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/listener/paymentcomplete", ListenerOrderPaymentComplete).Methods("POST")

	// mood
	clientRoutes.HandleFunc("/mood", MoodAdd).Methods("POST")
	clientRoutes.HandleFunc("/mood/history", MoodHistory).Queries(
		"client_id", "{client_id}",
		"dates", "{dates}",
	).Methods("GET")
	clientRoutes.HandleFunc("/mood/content", ListMoodContent).Queries(
		"user_id", "{user_id}",
	).Methods("GET")

	// notification
	clientRoutes.HandleFunc("/notification", NotificationsGet).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// delete
	clientRoutes.HandleFunc("/family_member/delete", DeleteAccountForFamilyMember).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// profile
	clientRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	clientRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	clientRoutes.HandleFunc("/corporate", ProfileAddForCor).Methods("POST")
	clientRoutes.HandleFunc("", ProfileUpdate).Queries(
		"client_id", "{client_id}",
	).Methods("PUT")
	clientRoutes.HandleFunc("/corporate_get_relation", RelativeProfileAdd).Methods("POST")
	clientRoutes.HandleFunc("/corporate_get_relation", GetRelativeProfile).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// search
	clientRoutes.HandleFunc("/search", ListSearch).Methods("GET")

	// search
	clientRoutes.HandleFunc("/webinar", WebinarList).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	// search
	clientRoutes.HandleFunc("/webinar", WebinarOrderCreate).Methods("POST")

	// corporate search
	clientRoutes.HandleFunc("/corporate_search", ListSearchForCorporate).Methods("GET")

	// corporate search
	clientRoutes.HandleFunc("/corporate/available-slots", ListSearchForCorporateForAvailableSlots).Methods("GET")

	// corporate in person search
	clientRoutes.HandleFunc("/corporate_inperson_search", ListSearchForCorporateInPerson).Methods("GET")

	// therapist
	clientRoutes.HandleFunc("/therapist", TherapistProfile).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/therapist/slots", TherapistSlots).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// In Person slots
	clientRoutes.HandleFunc("/in_person_therapist/slots", InPersonTherapistSlots).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/therapist/order", TherapistOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/therapist/paymentcomplete", TherapistOrderPaymentComplete).Methods("POST")

	clientRoutes.HandleFunc("/therapist/getHashData", GenerateHashForPayment).Queries(
		"hashData", "{hashData}",
	).Methods("GET")

	clientRoutes.HandleFunc("/restore-user-account", RestoreUserProfile).Methods("PUT")

	// corporate client

	clientRoutes.HandleFunc("/corporateCounsellor/order", CorporateCounsellorOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/corporateCounsellor/paymentcomplete", CorporateCounsellorOrderPaymentComplete).Methods("POST")

	// in person corporate clients

	clientRoutes.HandleFunc("/inperson_corporateCounsellor/order", InPersonCorporateCounsellorOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/inperson_corporateCounsellor/paymentcomplete", InPersonCorporateCounsellorOrderPaymentComplete).Methods("POST")

}
