package constant

// required fields for api endpoints
var (
	AdminProfileAddRequiredFields                = []string{"username", "password", "type"}
	AppointmentBookRequiredFields                = []string{"appointment_slot_id", "date", "time"}
	AppointmentRescheduleRequiredFields          = []string{"date", "time"}
	AppointmentRatingAddRequiredFields           = []string{"appointment_id", "rating", "client_id", "counsellor_id"}
	ContentAddRequiredFields                     = []string{"title", "photo", "content", "type"}
	CounsellorOrderCreateRequiredFields          = []string{"client_id", "counsellor_id", "date", "time", "no_session"}
	CouponAddRequiredFields                      = []string{"coupon_code", "discount", "minimum_order_value", "type", "start_by", "end_by"}
	EventOrderCreateRequiredFields               = []string{"user_id", "event_order_id"}
	EventBlockOrderCreateRequiredFields          = []string{"counsellor_id", "title", "description", "topic_id", "date", "time", "price"}
	EventOrderPaymentCompleteRequiredFields      = []string{"order_id", "payment_method", "payment_id"}
	EventBlockOrderPaymentCompleteRequiredFields = []string{"order_id", "payment_method", "payment_id"}
	ListenerOrderCreateRequiredFields            = []string{"client_id", "listener_id", "date", "time"}
	CounsellorOrderPaymentCompleteRequiredFields = []string{"order_id", "payment_method", "payment_id"}
	ListenerOrderPaymentCompleteRequiredFields   = []string{"order_id"}
	ClientProfileAddRequiredFields               = []string{"first_name", "phone", "email"}
	CounsellorProfileAddRequiredFields           = []string{"first_name", "phone", "price"}
	MoodAddRequiredFields                        = []string{"mood_id", "date"}
	ListenerProfileAddRequiredFields             = []string{"first_name", "phone"}
	QuoteAddRequiredFields                       = []string{"quote"}
	TherapistOrderCreateRequiredFields           = []string{"client_id", "therapist_id", "date", "time", "no_session"}
	TherapistOrderPaymentCompleteRequiredFields  = []string{"order_id", "payment_method", "payment_id"}
	TherapistProfileAddRequiredFields            = []string{"first_name", "phone", "price"}
	CancellationUpdateRequiredFields             = []string{"cancellation_reason"}
)
