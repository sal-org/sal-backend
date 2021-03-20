package constant

// database tables
const (
	AppointmentsTable        = "appointments"
	AppointmentSlotsTable    = "appointment_slots"
	ClientsTable             = "clients"
	CouponsTable             = "coupons"
	CounsellorsTable         = "counsellors"
	CounsellorLanguagesTable = "counsellor_languages"
	CounsellorTopicsTable    = "counsellor_topics"
	InvoicesTable            = "invoices"
	EventsTable              = "events"
	LanguagesTable           = "languages"
	ListenersTable           = "listeners"
	OrdersTable              = "orders"
	PhoneOTPVerifiedTable    = "phone_otp_verified"
	SchedulesTable           = "schedules"
	SlotsTable               = "slots"
	TopicsTable              = "topics"
)

// NumberOfTimesUniqueInserts - number of times insert statement should get executed for unqiue id
const NumberOfTimesUniqueInserts = 10

// RandomIDDigits - random unqiue ID, for generating unique random id
const RandomIDDigits = "abcdefghijklmnopqrstuvwxyz0123456789"

// length of unqiue digits to be generated for each table
const (
	AppointmentDigits     = 12
	AppointmentSlotDigits = 11
	ClientDigits          = 13
	CounsellorDigits      = 6
	ListenerDigits        = 9
	EventDigits           = 7
	InvoiceDigits         = 8
	OrderDigits           = 10
)
