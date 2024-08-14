package constant

const (
	// SameerEmailID = "sameer.littlemagix@gmail.com"
	// AnandEmailID  = "anand.shah@clovemind.com"
	// AkshayEmailID = "akshay.gandhi@clovemind.com"
	// ShivamEmailID = "shivam.tiwari@clovemind.com"
)

const (
	InstantSendEmailMessage = true
	LaterSendEmailMessage   = false
)

// format event-Target
// notification headings
const (
	// client
	ClientSignupProfileTitle                               = "Congratulations on your Clove registration"
	ClientAppointmentBookCounsellorTitle                   = "Clove: You have a new booking"
	ClientAppointmentBookClientTitle                       = "Clove: Your appointment has been confirmed"
	ClientAppointmentRescheduleClientTitle                 = "Clove: Appointment has been rescheduled"
	ClientAppointmentCancelClientTitle                     = "Clove: Your appointment has been cancelled successfully"
	ClientAppointmentBulkCancelClientTitle                 = "Clove: All your scheduled appointment sessions have been cancelled successfully"
	ClientCorLoginOTPTitle                                 = "Clove: Your Mobile App OTP is ###otp###"
	ClientAppointmentCancelCounsellorTitle                 = "Clove: Your session has been cancelled"
	CounsellorAppointmentCancelCounsellorTitle             = "Clove: Your appointment has been cancelled successfully"
	ClientAppointmentFollowUpSessionCounsellorTitle        = "Clove: Followup session booked!"
	ClientAppointmentFollowUpSessionClientTitle            = "Clove: Your follow-up appointment has been confirmed"
	CounsellorAppointmentCancelClientTitle                 = "Clove: Your session has been cancelled"
	ClientPaymentSucessClientTitle                         = "Payment Confirmation"
	CounsellorApprovedContentTitle                         = "Your content uploaded"
	ClientInPersonAppointmentBookTherapistTitle            = "Clove: You have a new in-person counselling session booking"
	ClientInPersonAppointmentBookClientTitle               = "Clove: Your in-person appointment has been confirmed"
	ClientInPersonAppointmentRescheduleTherapistTitle      = "Clove: In-person appointment with ###clientName### has been rescheduled"
	ClientInPersonAppointmentRescheduleClientTitle         = "Clove: In-person appointment has been rescheduled"
	ClientInPersonAppointmentCancellationTherapistTitle    = "Clove: Your in-person session with client has been cancelled"
	ClientInPersonAppointmentCancellationClientTitle       = "Clove: Your in-person appointment has been cancelled successfully"
	TherapistInPersonAppointmentCancellationTherapistTitle = "Clove: Your in-person appointments have been cancelled successfully"
	TherapistInPersonAppointmentCancellationClientTitle    = "Clove: Your in-person appointment has been cancelled"
)

// notification contents
const (
	// client
	ClientAppointmentCancelClientBody                     = "Your scheduled session with ###therapist_name### on ###date### at ###time### has been cancelled successfully."
	ClientCorLoginOTPBody                                 = "###otp### is the OTP to login/register yourself on the Clove Mobile App. Please do not share it with anyone."
	ClientSignupClientEmailBody                           = "Welcome to Clove mobile app! You can now access self-care audios, relevant articles, daily journaling, self-assessments and book your sessions seamlessly."
	ClientAppointmentBookCounsellorEmailBody              = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details. Please note that the call may be recorded for quality and training purposes in accordance with the privacy policy."
	ClientAppointmentRescheduleCounsellorEmailBody        = "Your client, ###first_name### has rescheduled the appointment to ###date### on ###time###. Please check the Upcoming Sessions Section for further details."
	ClientAppointmentCancelCounsellorEmailBody            = "Your upcoming session with ###client_name### for ###date### on ###time### has been cancelled by the client."
	ClientAppointmentBulkCancelClientEmailBody            = "You have cancelled all your pending sessions. Refund, if any, will be processed shortly as per the Cancellation & Refund Policy. We wish you good health and know that you can always come back to us for any support"
	ClientAppointmentRescheduleClientEmailBody            = "Your request for a session reschedule has been confirmed for ###time### on ###date### with ###therpists_name###, Click 'My Sessions' for more details"
	ClientAppointmentBookClientEmailBody                  = "Your session is booked with ###therpist_name### on ###date### at ###time###. We appreciate your commitment to prioritising your wellbeing."
	CounsellorAppointmentCancelCounsellorBodyEmailBody    = "You have now cancelled your scheduled session. Cancellation charges, if any, will be processed as per the Cancellation & Refund Policy."
	ClientAppointmentFollowUpSessionClientEmailBody       = "Your follow-up appointment for one to one chat with ###therpist_name### is confirmed  for ###date_time###."
	ClientAppointmentFollowUpSessionCounsellorEmailBody   = "You have a followup counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details."
	CounsellorAppointmentCancelClientBodyEmailBody        = "We are sorry that your session has been cancelled by ###therapist_name### for ###date_time### due to a personal emergency. Please reschedule your session."
	ClientPaymentSucessClientBody                         = "Hi ###client_name###, successful payment of Rs. ###paid_amount### has been received for your consultation booking. Manage your account anytime, anywhere from your SAL Mobile app on your phone. Click to view your Transaction ID #  and receipt"
	CounsellorAccountSignupCounsellorEmailBody            = "Thank you for successfully completing the onboarding form. Our offline team will contact you shortly for Agreement signup."
	AdminRefundAmonutForClientEmailBody                   = "Your refund of Rs. ###amount### has been initiated. This is for your appointment dated ###date### at ###time### and is being processed to your card/bank account. It will reflect in 7-10 working days."
	CounsellorAppointmentWorkClientEmailBody              = "Your appointment completed. your therapist send some attachment please check and do neccessary action according to give into attachment to improve your well being journey."
	RatingTitleForInternalReviewBody                      = "We are writing to inform you that ###therapistname### has received a rating of ###rating### for the session conducted on ###date### at ###time###. User Comment: ###content###"
	CounsellorApprovedContentBody                         = "Congratulations, your content ###content_name### has been uploaded on the clovemind app. We appreciate your contribution."
	ClientInPersonAppointmentBookTherapistBody            = "You have a new counselling session booked by ###clientName### for ###date### and ###time###. The appointment is booked for ###location###. Please ensure to Start and End the in-person meeting at the scheduled session time. For details, check 'In-Person Session' in the Clove app menu."
	ClientInPersonAppointmentBookClientBody               = "Your appointment for private & confidential talk with Clove therapist ###therapistName### is confirmed for  ###date### at ###time###. Your in-person counselling session will be at ###location###. For details, check 'In-Person Session' in the Clove app menu."
	ClientInPersonAppointmentRescheduleTherapistBody      = "Your client, ###clientName### has rescheduled the appointment to ###date### on ###time###. Please check the Upcoming Sessions in the in-person Section for further details."
	ClientInPersonAppointmentRescheduleClientBody         = "Your request to reschedule an in-person counselling session has been confirmed for ###time### on ###date### with Clove therapist ###therapistName###. Click in-person sessions on Clove app for more details."
	ClientInPersonAppointmentCancellationTherapistBody    = "Your upcoming session with ###clientName### for ###date### and ###time### at ###location### has been cancelled by the client."
	ClientInPersonAppointmentCancellationClientBody       = "Your scheduled in-person counselling session with Clove therapist ###therapistName### on ###date### and ###time### at ###location### has been cancelled successfully. You may reschedule your counselling session for the next available slot with  thearpist again at your convenience."
	TherapistInPersonAppointmentCancellationTherapistBody = "Due to your unavilability, all your scheduled in-person sessions for ###date### at ###location### have been cancelled successfully."
	TherapistInPersonAppointmentCancellationClientBody    = "We regret to inform you that your scheduled session with Clove therapist ###therapist### for ###date### and ###time### at ###location### has been cancelled due to therapist unavailbility. We request you to rebook a new appointment for the next available date and slot time."
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
	CounsellorRecordForClientTitle    = "Clove: New Client Record"
	CounsellorDocumentForClientTitle  = "Clove: Counsellor Notes"
	RatingTitleForInternalReviewTitle = "Rating for ###therapistName### : ###sessionDate###"
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
