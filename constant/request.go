package constant

// required fields for api endpoints
var (
	AppointmentBookRequiredFields                     = []string{"appointment_slot_id", "date", "time"}
	AppointmentRescheduleRequiredFields               = []string{"date", "time"}
	AppointmentRatingAddRequiredFields                = []string{"appointment_id", "rating", "client_id", "counsellor_id"}
	CounsellorOrderCreateRequiredFields               = []string{"client_id", "counsellor_id", "date", "time", "no_session"}
	ClientEventOrderCreateRequiredFields              = []string{"client_id", "event_order_id"}
	CounsellorEventOrderCreateRequiredFields          = []string{"counsellor_id", "title", "description", "topic_id", "date", "time", "duration", "price"}
	ListenerOrderCreateRequiredFields                 = []string{"client_id", "listener_id", "date", "time"}
	CounsellorOrderPaymentCompleteRequiredFields      = []string{"order_id", "payment_method", "payment_id"}
	ClientEventOrderPaymentCompleteRequiredFields     = []string{"order_id", "payment_method", "payment_id"}
	CounsellorEventOrderPaymentCompleteRequiredFields = []string{"order_id", "payment_method", "payment_id"}
	ListenerOrderPaymentCompleteRequiredFields        = []string{"order_id"}
	ClientProfileAddRequiredFields                    = []string{"first_name", "phone", "email"}
	CounsellorProfileAddRequiredFields                = []string{"first_name", "phone", "email", "price"}
	ListenerProfileAddRequiredFields                  = []string{"first_name", "phone", "email"}
	TherapistOrderCreateRequiredFields                = []string{"client_id", "therapist_id", "date", "time", "no_session"}
	TherapistEventOrderCreateRequiredFields           = []string{"therapist_id", "title", "description", "topic_id", "date", "time", "duration", "price"}
	TherapistOrderPaymentCompleteRequiredFields       = []string{"order_id", "payment_method", "payment_id"}
	TherapistEventOrderPaymentCompleteRequiredFields  = []string{"order_id", "payment_method", "payment_id"}
	TherapistProfileAddRequiredFields                 = []string{"first_name", "phone", "email", "price"}
)
