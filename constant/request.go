package constant

// required fields for api endpoints
var (
	CounsellorOrderCreateRequiredFields          = []string{"client_id", "counsellor_id", "date", "time", "no_session"}
	ListenerOrderCreateRequiredFields            = []string{"client_id", "listener_id", "date", "time"}
	CounsellorOrderPaymentCompleteRequiredFields = []string{"order_id", "payment_method", "payment_id"}
	ListenerOrderPaymentCompleteRequiredFields   = []string{"order_id"}
)
