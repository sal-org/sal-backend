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
	CounsellorBlocked     = "3"
)

// listener status
const (
	ListenerNotApproved = "0"
	ListenerActive      = "1"
	ListenerInactive    = "2"
	ListenerBlocked     = "3"
)

// order status
const (
	OrderWaiting    = "0"
	OrderInProgress = "1"
	OrderCompleted  = "2"
)

// appointment status
const (
	AppointmentToBeStarted = "1"
	AppointmentStarted     = "2"
	AppointmentCompleted   = "3"
)

// appointment slots status
const (
	AppointmentSlotsActive = "1"
)

// invoice status
const (
	InvoiceInProgress      = "1"
	InvoiceRefundInitiated = "2"
	InvoiceRefunded        = "3"
	InvoiceCompleted       = "4"
)

// coupon types
const (
	CouponFlatType       = "1"
	CouponPercentageType = "2"
)

// coupon status
const (
	CouponInactive = "0"
	CouponActive   = "1"
)

// appointment slots types
const (
	SlotUnavailable = "0"
	SlotAvailable   = "1"
	SlotBooked      = "2"
)

// event status
const (
	EventWaiting     = "0" // payment not done by counsellor,
	EventToBeStarted = "1"
	EventStarted     = "2"
	EventCompleted   = "3"
)
