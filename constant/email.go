package constant

const (
	SameerEmailID = "sameer.littlemagix@gmail.com"
	AnandEmailID  = "anand.shah@sal.foundation"
	AkshayEmailID = "akshay.gandhi@clovemind.com"
	ShivamEmailID = "shivam.tiwari@clovemind.com"
)

const (
	InstantSendEmailMessage = true
	LaterSendEmailMessage   = false
)

// format event-Target
// notification headings
const (
	// client
	ClientSignupProfileTitle                        = "Congratulations on your Clove registration"
	ClientAppointmentBookCounsellorTitle            = "Clove: You have a new booking"
	ClientAppointmentBookClientTitle                = "Clove: Your appointment has been confirmed"
	ClientAppointmentRescheduleClientTitle          = "Clove: Appointment has been rescheduled"
	ClientAppointmentCancelClientTitle              = "Clove: Your appointment has been cancelled successfully"
	ClientAppointmentBulkCancelClientTitle          = "Clove: All your scheduled appointment sessions have been cancelled successfully"
	ClientCorLoginOTPTitle                          = "Clove: Your Mobile App OTP is ###otp###"
	ClientAppointmentCancelCounsellorTitle          = "Clove: Your session has been cancelled"
	CounsellorAppointmentCancelCounsellorTitle      = "Clove: Your appointment has been cancelled successfully"
	ClientAppointmentFollowUpSessionCounsellorTitle = "Clove: Followup session booked!"
	ClientAppointmentFollowUpSessionClientTitle     = "Clove: Your follow-up appointment has been confirmed"
	CounsellorAppointmentCancelClientTitle          = "Clove: Your session has been cancelled"
	ClientPaymentSucessClientTitle                  = "Payment Confirmation"
)

// notification contents
const (
	// client
	ClientAppointmentCancelClientBody                   = "Your scheduled session with ###therapist_name### on ###date### at ###time### has been cancelled successfully."
	ClientCorLoginOTPBody                               = "###otp### is the OTP to login/register yourself on the Clove Mobile App. Please do not share it with anyone."
	ClientSignupClientEmailBody                         = "Welcome to Clove mobile app! You can now access self-care audios, relevant articles, daily journaling, self-assessments and book your sessions seamlessly."
	ClientAppointmentBookCounsellorEmailBody            = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details. Please note that the call may be recorded for quality and training purposes in accordance with the privacy policy."
	ClientAppointmentRescheduleCounsellorEmailBody      = "Your client, ###first_name### has rescheduled the appointment to ###date### on ###time###. Please check the Upcoming Sessions Section for further details."
	ClientAppointmentCancelCounsellorEmailBody          = "Your upcoming session with ###client_name### for ###date### on ###time### has been cancelled by the client."
	ClientAppointmentBulkCancelClientEmailBody          = "You have cancelled all your pending sessions. Refund, if any, will be processed shortly as per the Cancellation & Refund Policy. We wish you good health and know that you can always come back to us for any support"
	ClientAppointmentRescheduleClientEmailBody          = "Your request for a session reschedule has been confirmed for ###time### on ###date### with ###therpists_name###, Click 'My Sessions' for more details"
	ClientAppointmentBookClientEmailBody                = "Your appointment for private & confidential talk with ###therpist_name### is confirmed for ###date### at ###time###. Please note that the call may be recorded for quality and training purposes in accordance with the privacy policy."
	CounsellorAppointmentCancelCounsellorBodyEmailBody  = "You have now cancelled your scheduled session. Cancellation charges, if any, will be processed as per the Cancellation & Refund Policy."
	ClientAppointmentFollowUpSessionClientEmailBody     = "Your follow-up appointment for one to one chat with ###therpist_name### is confirmed  for ###date_time###."
	ClientAppointmentFollowUpSessionCounsellorEmailBody = "You have a followup counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details."
	CounsellorAppointmentCancelClientBodyEmailBody      = "We are sorry that your session has been cancelled by ###therapist_name### for ###date_time### due to a personal emergency. Please reschedule your session within the next 7 days or request for a credit refund."
	ClientPaymentSucessClientBody                       = "Hi ###client_name###, successful payment of Rs. ###paid_amount### has been received for your consultation booking. Manage your account anytime, anywhere from your SAL Mobile app on your phone. Click to view your Transaction ID #  and receipt"
	CounsellorAccountSignupCounsellorEmailBody          = "Thank you for successfully completing the onboarding form. Our offline team will contact you shortly for Agreement signup."
	AdminRefundAmonutForClientEmailBody                 = "Your refund of Rs. ###amount### has been initiated. This is for your appointment dated ###date### at ###time### and is being processed to your card/bank account. It will reflect in 7-10 working days."
)

// Event Approval For Sal Team to Send a Email

const (
	NewEventWaitingForApprovalTitle = "New SAL Cafe Event has been Created!"
)

// Same as Profile Approval SAL Team to Send A Email

const (
	CounsellorProfileWaitingForApprovalTitle = "Clove successful sign in"
)

// counsellor - client record
const (
	CounsellorRecordForClientTitle = "Clove: New Client Record"
)

// counsellor - client record
const (
	CounselloRecordClientEmergencyCaseTitle = "Clove: New Emergency Case"
)

// counsellor - client record
const (
	CounsellorVisitForClientTitle = "Clove: New Client Visit Record"
)

// For Receipt type
const (
	AppointmentSessionsTypeForReceipt = "Counselling Sessions"
	SalCafeTypeForReceipt             = "SAL Café"
	SalCafeQty                        = "01"
)
