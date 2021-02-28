package constant

// server status codes
const (
	StatusCodeOk             = "200"
	StatusCodeCreated        = "201"
	StatusCodeBadRequest     = "400"
	StatusCodeForbidden      = "403"
	StatusCodeSessionExpired = "440"
	StatusCodeServerError    = "500"
	StatusCodeDuplicateEntry = "1000"
)

// type of alerts for frontend to show
const (
	NoDialog   = "0"
	ShowDialog = "1"
	ShowToast  = "2"
)

// appointment or event order
const (
	OrderAppointmentType = "1"
	OrderEventType       = "2"
)

// counsellor or listener
const (
	CounsellorType = "1"
	ListenerType   = "2"
	ClientType     = "3"
)

// counsellor status
const (
	ClientActive  = "1"
	ClientDeleted = "2"
	ClientBlocked = "3"
)

// counsellor status
const (
	CounsellorNotApproved = "0"
	CounsellorActive      = "1"
	CounsellorInactive    = "2"
	CounsellorBlocked     = "2"
)

// listener status
const (
	ListenerActive   = "1"
	ListenerInactive = "2"
	ListenerBlocked  = "2"
)

// order status
const (
	OrderWaiting    = "0"
	OrderInProgress = "1"
)

// appointment status
const (
	AppointmentWaiting   = "1"
	AppointmentStarted   = "2"
	AppointmentCompleted = "3"
)

// invoice status
const (
	InvoiceInProgress      = "1"
	InvoiceRefundInitiated = "2"
	InvoiceRefunded        = "3"
	InvoiceCompleted       = "3"
)

// coupon types
const (
	CouponFlatType       = "1"
	CouponPercentageType = "2"
)

// appointment slots types
const (
	SlotUnavailable = "0"
	SlotAvailable   = "1"
	SlotBooked      = "2"
)
