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

const (
	NewEventWaitingForApprovalTitle = "New SAL Cafe Event has been Created!"
)

const (
	NewEventWaitingForApprovalBody = "New Event , Counsellor Id : ###counsellor_id### , Type 1 means Counsellor or Type 2 means Therapists ###type###, Event Title : ###title###, Event Description : ###description###, Photo : ###photo###, Topic id : ###topic_id### , Date : ###date### , Time : ###time### , Duration : ###duration### , Event Price : ###price###"
)

const (
	CounsellorProfileWaitingForApprovalTitle = "New Profile Of Counsellor has been Created!"
	CounsellorProfileWaitingForApprovalBody  = "Counsellor Name is : ###first_name### ###last_name### , Gender : ###gender### , Phone Number : ###phone### , Photo : ###photo### , Email-Id : ###email### , Education : ###education### , Experience : ###experience### , About : ###about### , Resume : ###resume### , Ceritificate : ###certificate### , Aadhar : ###aadhar### , Linkedin : ###linkedin### , Status : ###status###"
)
