package constant

// database tables
const (
	AppointmentsTable        = "appointments"
	ClientsTable             = "clients"
	CouponsTable             = "coupons"
	CounsellorsTable         = "counsellors"
	CounsellorLanguagesTable = "counsellor_languages"
	CounsellorTopicsTable    = "counsellor_topics"
	InvoicesTable            = "invoices"
	LanguagesTable           = "languages"
	ListenersTable           = "listeners"
	OrdersTable              = "orders"
	SlotsTable               = "slots"
	TopicsTable              = "topics"
)

// NumberOfTimesUniqueInserts - number of times insert statement should get executed for unqiue id
const NumberOfTimesUniqueInserts = 10

// RandomIDDigits - random unqiue ID, for generating unique random id
const RandomIDDigits = "abcdefghijklmnopqrstuvwxyz0123456789"

// length of unqiue digits to be generated for each table
const (
	InvoiceDigits = 8
	OrderDigits   = 10
)
