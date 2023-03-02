package constant

// database tables
const (
	AdminsTable                    = "admins"
	UsersPermissionTable           = "usersPermission"
	RolesTable                     = "roles"
	QualityCheckTable              = "qualitycheck"
	QualityCheckDetailsTable       = "qualitycheck_details"
	AppointmentsTable              = "appointments"
	AppointmentSlotsTable          = "appointment_slots"
	AssessmentsTable               = "assessments"
	AssessmentQuestionsTable       = "assessment_questions"
	AssessmentQuestionOptionsTable = "assessment_question_options"
	AssessmentScoresTable          = "assessment_scores"
	AssessmentResultsTable         = "assessment_results"
	AssessmentResultDetailsTable   = "assessment_result_details"
	ClientsTable                   = "clients"
	ContentsTable                  = "contents"
	ContentCategoriesTable         = "content_categories"
	ContentLikesTable              = "content_likes"
	CouponsTable                   = "coupons"
	CounsellorsTable               = "counsellors"
	CounsellorLanguagesTable       = "counsellor_languages"
	CounsellorTopicsTable          = "counsellor_topics"
	EmailsTable                    = "emails"
	QualityCheckEmailTable         = "qualitycheck_email"
	InvoicesTable                  = "invoices"
	LanguagesTable                 = "languages"
	ListenersTable                 = "listeners"
	MessagesTable                  = "messages"
	MoodsTable                     = "moods"
	MoodResultsTable               = "mood_results"
	NotificationsTable             = "notifications"
	NotificationsBulkTable         = "notifications_bulk"
	OrderClientAppointmentTable    = "order_client_appointments"
	OrderEventTable                = "order_events"
	OrderCounsellorEventTable      = "order_counsellor_events"
	PaymentsTable                  = "payments"
	PhoneOTPVerifiedTable          = "phone_otp_verified"
	QuotesTable                    = "quotes"
	RatingTypesTable               = "rating_types"
	RefundsTable                   = "refunds"
	SchedulesTable                 = "schedules"
	SlotsTable                     = "slots"
	TherapistsTable                = "therapists"
	TopicsTable                    = "topics"
	ReceiptTable                   = "receipts"
	AssessmentPdfTable             = "assessment_pdf"
	AgoraTable                     = "agora"
)

// NumberOfTimesUniqueInserts - number of times insert statement should get executed for unqiue id
const NumberOfTimesUniqueInserts = 10

// RandomIDDigits - random unqiue ID, for generating unique random id
const RandomIDDigits = "abcdefghijklmnopqrstuvwxyz0123456789"

// length of unqiue digits to be generated for each table
const (
	AdminDigits             = 4
	AppointmentDigits       = 12
	AppointmentSlotDigits   = 11
	AssessmentResultsDigits = 16
	MoodResultsDigits       = 17
	ClientDigits            = 13
	ContentDigits           = 10
	CounsellorDigits        = 6
	ListenerDigits          = 9
	MessagesDigits          = 10
	EmailsDigits            = 13
	NotificationsDigits     = 15
	EventDigits             = 7
	InvoiceDigits           = 8
	ReceiptDigits           = 10
	PaymentsDigits          = 10
	OrderDigits             = 10
	OrderEventDigits        = 11
	PaymentDigits           = 10
	RefundDigits            = 9
	TherapistDigits         = 5
	AgoraDigits             = 20
)
