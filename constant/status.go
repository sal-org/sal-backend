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

// content types
const (
	VideoContentType   = "1"
	AudioContentType   = "2"
	ArticleContentType = "3"
)

// content status
const (
	ContentActive  = "1"
	ContentDeleted = "2"
)

// appointment or event order
const (
	OrderAppointmentType = "1"
	OrderEventBookType   = "2"
	OrderEventBlockType  = "3"
)

// counsellor or listener
const (
	CounsellorType = "1"
	ListenerType   = "2"
	ClientType     = "3"
	TherapistType  = "4"
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

// therapist status
const (
	TherapistNotApproved = "0"
	TherapistActive      = "1"
	TherapistInactive    = "2"
	TherapistBlocked     = "3"
)

// admin status
const (
	AdminActive  = "1"
	AdminBlocked = "2"
)

// payment status
const (
	PaymentValid   = "1"
	PaymentInvalid = "2"
)

// notification status
const (
	NotificationInactive = "0"
	NotificationActive   = "1"
	NotificationRead     = "2"
	NotificationDeleted  = "3"
)

// email status
const (
	EmailInProgress = "1"
	EmailSent       = "2"
)

// message status
const (
	MessageInProgress = "1"
	MessageSent       = "2"
)

// notification sent status
const (
	NotificationInProgress = "1"
	NotificationSent       = "2"
	NotificationInvalid    = "3"
)

// order status
const (
	OrderWaiting    = "0"
	OrderInProgress = "1"
	OrderCompleted  = "2"
)

// appointment status
const (
	AppointmentToBeDuplicate                  = "0"
	AppointmentToBeStarted                    = "1"
	AppointmentStarted                        = "2"
	AppointmentCompleted                      = "3"
	AppointmentUserCancelled                  = "4"
	AppointmentUserCancelledWithin4Hour       = "12"
	AppointmentCounsellorCancelledWithin4Hour = "13"
	AppointmentInTheReview                    = "14"
	AppointmentCounsellorCancelled            = "5"
	AppointmentAdminCancelled                 = "6"
	AppointmentNoShowClient                   = "7"
	AppointmentNoShowCounsellor               = "8"
	AppointmentNoShowBoth                     = "9"
	AppointmentAdminRefunds                   = "10"
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

// refund status
const (
	RefundInProgress = "1"
	RefundCompleted  = "2"
)

// refund type
const (
	RefundCancellationType = "1"
	RefundNoShowType       = "2"
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

// partner status
const (
	PartnerInactive = "0"
	PartnerActive   = "1"
)

// partner address status
const (
	PartnerAddressInactive = "0"
	PartnerAddressActive   = "1"
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
	EventNoShow      = "4"
)

// assessment status
const (
	AssessmentInactive = "0"
	AssessmentActive   = "1"
)

// assessment question status
const (
	AssessmentQuestionInactive = "0"
	AssessmentQuestionActive   = "1"
)

// assessment question option status
const (
	AssessmentQuestionOptionInactive = "0"
	AssessmentQuestionOptionActive   = "1"
)

// Appointment Request Option status
const (
	AppointmentRequestProgress  = "1"
	AppointmentRequestCompleted = "2"
)

// assessment result status
const (
	AssessmentResultInactive = "0"
	AssessmentResultActive   = "1"
)

// payment  status
const (
	PaymentActive        = "1"
	PaymentActiveDeleted = "2"
)

// mood result status
const (
	MoodResultInactive = "0"
	MoodResultActive   = "1"
)

// bulk notification type
const (
	BulkNotificationAll      = "1"
	BulkNotificationSpecific = "2"
)

// Agora status
const (
	AgoraResourceID  = "1" // agora resource id gerenated
	AgoraResourceID2 = "2"
	AgoraCallStart1  = "3" // agora call start
	AgoraCallStart2  = "4"
	AgoraCallStop1   = "5" // agora call stop
	AgoraCallStop2   = "6"
)

const (
	QualityCheckDataInsert = "1" // data insert in qualitycheckdetails table
	QualityCheckLinkInsert = "2" // for link
)
