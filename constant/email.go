package constant

const (
	InstantSendEmailMessage = true
	LaterSendEmailMessage   = false
)

// format event-Target
// notification headings
const (
	// client
	ClientAppointmentCancelClientTitle     = "Booking Cancellation"
	CounsellorAppointmentCancelClientTitle = "Appointment Cancellation"
	ClientPaymentSucessClientTitle         = "Payment Confirmation"
)

// notification contents
const (
	// client
	ClientAppointmentCancelClientBody     = "Hi ###client_name###, you have now cancelled your scheduled session. Refund, if any, will be processed shortly as per the Cancellation & Refund Policy. You can also re-book a session at any time."
	CounsellorAppointmentCancelClientBody = "Hi ###client_name###, Your scheduled session is cancelled by the counselor. You may reschedule your appointment within the next 7 days or request for a session credit refund"
	ClientPaymentSucessClientBody         = "Hi ###client_name###, successful payment of Rs. ###paid_amount### has been received for your consultation booking. Manage your account anytime, anywhere from your SAL Mobile app on your phone. Click to view your Transaction ID #  and receipt"
)
