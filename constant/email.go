package constant

const (
	SameerEmailID = "sameer.littlemagix@gmail.com"
	AnandEmailID  = "anand.shah@sal.foundation"
	ShivamEmailID = "tiwarisomprakash.2000@gmail.com"
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
	ClientAppointmentCancelClientBody                   = "you have cancelled your scheduled session. Refund, if any, will be processed shortly as per the Cancellation & Refund Policy. We wish you good health and know that you can always come back to us for any support"
	ClientSignupClientEmailBody                         = "Your profile is ready! Many features like audios, articles, mood diaries, workshops on CLOVE Cafe and booking counselling sessions are now unlocked."
	ClientAppointmentBookCounsellorEmailBody            = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details. Please note that the call may be recorded for quality and training purposes in accordance with the privacy policy."
	ClientAppointmentRescheduleCounsellorEmailBody      = "Your Client has rescheduled a current booking. Please click on the CLOVE Mobile app and view the change in the 'Booking section' for details"
	ClientAppointmentCancelCounsellorEmailBody          = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client."
	ClientAppointmentBulkCancelClientEmailBody          = "You have cancelled all your pending sessions. Refund, if any, will be processed shortly as per the Cancellation & Refund Policy. We wish you good health and know that you can always come back to us for any support"
	ClientAppointmentRescheduleClientEmailBody          = "Your request for a session reschedule has been confirmed for ###date_time### with ###therpists_name###, Click 'My Sessions' for more details"
	ClientAppointmentBookClientEmailBody                = "Your appointment for private & confidential talk with ###therpist_name### is confirmed for ###date_time###. Please note that the call may be recorded for quality and training purposes in accordance with the privacy policy."
	CounsellorAppointmentCancelCounsellorBodyEmailBody  = "You have now cancelled your scheduled session. Cancellation charges, if any, will be processed as per the Cancellation & Refund Policy."
	ClientAppointmentFollowUpSessionClientEmailBody     = "Your follow-up appointment for one to one chat with ###therpist_name### is confirmed  for ###date_time###."
	ClientAppointmentFollowUpSessionCounsellorEmailBody = "You have a followup counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details."
	CounsellorAppointmentCancelClientBodyEmailBody      = "We are sorry that your session has been cancelled by your ###therapist_name### for ###date_time### due to a personal emergency. Please reschedule your session within the next 7 days or request for a credit refund."
	ClientPaymentSucessClientBody                       = "Hi ###client_name###, successful payment of Rs. ###paid_amount### has been received for your consultation booking. Manage your account anytime, anywhere from your SAL Mobile app on your phone. Click to view your Transaction ID #  and receipt"
	CounsellorAccountSignupCounsellorEmailBody          = "Thank you for successfully completing the onboarding form. Our offline team will contact you shortly for Agreement signup."
)

// Event Approval For Sal Team to Send a Email

const (
	NewEventWaitingForApprovalTitle = "New SAL Cafe Event has been Created!"
)

// Same as Profile Approval SAL Team to Send A Email

const (
	CounsellorProfileWaitingForApprovalTitle = "Clove successful sign in"
)

// For Receipt type
const (
	AppointmentSessionsTypeForReceipt = "Counselling Sessions"
	SalCafeTypeForReceipt             = "SAL Café"
	SalCafeQty                        = "01"
)
