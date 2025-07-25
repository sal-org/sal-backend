package constant

// database tables
const (
	AdminsTable                                 = "admins"
	UsersPermissionTable                        = "usersPermission"
	RolesTable                                  = "roles"
	QualityCheckTable                           = "qualitycheck"
	QualityCheckDetailsTable                    = "qualitycheck_details"
	AppointmentsTable                           = "appointments"
	ClientCounsellingLimitTable                 = "client_counselling_limit"
	ClientCounsellingUnLimitTable               = "client_counselling_unlimit"
	InPersonAppointmentsTable                   = "appointments_inperson"
	AppointmentRequestTable                     = "request_appointment"
	EventInPersonRequestTable                   = "request_inperson_cafe"
	InPersonAppointmentRequestTable             = "request_appointment_inperson"
	AppointmentSlotsTable                       = "appointment_slots"
	AssessmentsTable                            = "assessments"
	AssessmentQuestionsTable                    = "assessment_questions"
	AssessmentQuestionOptionsTable              = "assessment_question_options"
	AssessmentScoresTable                       = "assessment_scores"
	AssessmentResultsTable                      = "assessment_results"
	AssessmentResultDetailsTable                = "assessment_result_details"
	ClientsTable                                = "clients"
	FeedbackTable                               = "feedback"
	WebinarsTable                               = "webinars"
	WebinarsBookTable                           = "webinars_book"
	B2B2CAppointmentTransitionsTable            = "b2b2c_appointment_transition"
	WebB2BBookDemoTable                         = "web_b2b_book_demo"
	AppInfoTable                                = "app_info"
	CounsellorDocumentListTable                 = "counsellor_documnets_list"
	CorporatePartnersTable                      = "corporate_partners"
	CounsellorMyTimeSheetTable                  = "counsellor_time_sheet"
	CorporatePartnersAddressTable               = "partners_address"
	ContentsTable                               = "contents"
	ContentUserViewsTable                       = "content_user_views"
	InPersonCounsellorConnectWithCorporateTable = "inperson_connect_with_counsellor"
	ContentCategoriesTable                      = "content_categories"
	ContentLikesTable                           = "content_likes"
	CouponsTable                                = "coupons"
	CounsellorsTable                            = "counsellors"
	CounsellorRecordsTable                      = "counsellor_record"
	CounsellorRecordFormCategoryTable           = "counsellor_record_form_category"
	CounsellorRecordFormSubCategoryTable        = "counsellor_record_form_sub_category"
	CounsellorRecordFormEmotionalState          = "counsellor_record_form_emotional_state"
	CounsellorLanguagesTable                    = "counsellor_languages"
	CounsellorTopicsTable                       = "counsellor_topics"
	EmailsTable                                 = "emails"
	QualityCheckEmailTable                      = "qualitycheck_email"
	InvoicesTable                               = "invoices"
	LanguagesTable                              = "languages"
	ListenersTable                              = "listeners"
	MessagesTable                               = "messages"
	MoodsTable                                  = "moods"
	AdsContentTable                             = "adsContent"
	MoodResultsTable                            = "mood_results"
	NotificationsTable                          = "notifications"
	NotificationsBulkTable                      = "notifications_bulk"
	OrderClientAppointmentTable                 = "order_client_appointments"
	InPersonOrderClientAppointmentTable         = "order_client_appointments_inperson"
	OrderEventTable                             = "order_events"
	OrderEventInPersonTable                     = "order_events_inperson"
	OrderCounsellorEventTable                   = "order_counsellor_events"
	OrderCounsellorEventInPersonTable           = "order_counsellor_events_inperson"
	PaymentsTable                               = "payments"
	PhoneOTPVerifiedTable                       = "phone_otp_verified"
	QuotesTable                                 = "quotes"
	RatingTypesTable                            = "rating_types"
	RefundsTable                                = "refunds"
	SchedulesTable                              = "schedules"
	SchedulesDatesTable                         = "schedules_dates"
	SlotsTable                                  = "slots"
	InPersonSLotsTable                          = "in_person_slots"
	InPersonSLotsScheduleTable                  = "in_person_schedules"
	TherapistsTable                             = "therapists"
	TopicsTable                                 = "topics"
	ReceiptTable                                = "receipts"
	AssessmentPdfTable                          = "assessment_pdf"
	AgoraTable                                  = "agora"
	CorporateClientFamilyAccessControlTable     = "corporate_client_family_access_control"
	CompanyAccessControlTable                   = "company_access_control"
	ClientAccessControlTable                    = "client_access_control"
)

// NumberOfTimesUniqueInserts - number of times insert statement should get executed for unqiue id
const NumberOfTimesUniqueInserts = 10

// RandomIDDigits - random unqiue ID, for generating unique random id
const RandomIDDigits = "abcdefghijklmnopqrstuvwxyz0123456789"

// length of unqiue digits to be generated for each table
const (
	AdminDigits                    = 4
	AppointmentDigits              = 12
	LimitAppointmentDigits         = 12
	AppointmentRequestDigits       = 17
	AppointmentSlotDigits          = 11
	AssessmentResultsDigits        = 16
	AssessmentDigits               = 12
	AssessmentQuestionDigits       = 14
	AssessmentQuestionOptionDigits = 18
	MoodResultsDigits              = 17
	ClientDigits                   = 13
	FeedbackDigits                 = 17
	ContentDigits                  = 10
	CorporateDigits                = 8
	CorporateAddressDigits         = 10
	CounsellorDigits               = 6
	CounsellorRecordDigits         = 16
	ListenerDigits                 = 9
	MessagesDigits                 = 10
	EmailsDigits                   = 13
	NotificationsDigits            = 15
	EventDigits                    = 7
	WebinarsDigits                 = 10
	InvoiceDigits                  = 8
	ReceiptDigits                  = 10
	PaymentsDigits                 = 10
	OrderDigits                    = 10
	OrderEventDigits               = 11
	PaymentDigits                  = 10
	RefundDigits                   = 9
	TherapistDigits                = 5
	AgoraDigits                    = 20
)
