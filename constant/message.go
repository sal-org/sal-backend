package constant

// config
const (
	CorefactorsSendSMSEndpoint = "https://teleduce.corefactors.in/sendsms/"
	CorefactorsAPIKey          = "1d0ed788-d667-4e71-ba94-83b18e53a41e"
	TextMessageFrom            = "SALUBR"
)

// routes
const (
	TransactionalRouteTextMessage = "0"
)

const (
	InstantSendTextMessage = true
	LaterSendTextMessage   = false
)

// text messages
const (
	CounsellorAccountSignupTextMessage                       = "Hey ###counsellor_name###, the onboarding process has been completed successfully. We will contact you shortly to discuss the next steps. - Team Clove"
	CounsellorOTPTextMessage                                 = "Your Clove Mobile App registration OTP is ###otp###. You are now only few steps away from your getting your 1st appointment. - Team Clove"
	ClientOTPTextMessage                                     = "###otp### is the OTP to register yourself on the Clove Mobile App. Please do not share it with anyone.  - Team Clove"
	ClientProfileTitleMessage                                = "Hey ###client_name###, welcome to Clove! Access self-care assessments and content, daily mood diaries and book counselling sessions with Listeners & Therapists.  - Team Clove"
	ListenerOTPTextMessage                                   = "###otp### is the OTP to register yourself on SAL Mobile App. You are only few steps away from your taking your 1st listening session."
	ClientAppointmentScheduleCounsellorTextMessage           = "Hi ###counsellor_name###, you have a new counselling session booked by ###client_name### for ###date_time###. Please visit My sessions in the SAL Mobile app menu for more details."
	ClientAppointmentConfirmationTextMessage                 = "Hi ###client_name###, your online appointment with ###counsellor_name### is confirmed for ###date### at ###time###. - SAL Team"
	ClientPaymentConfirmationTextMeassge                     = "Hi ###client_name###, payment of INR ###amount### for ###bought### session(s) has been received. Please click ###Aaplink### to access the receipt. - SAL Team"
	ClientAppointmentRescheduleClientTextMeassge             = "Hello ###client_name##, your appointment has been re-scheduled successfully with ###counsellor_name###. Please view 'My Sessions' in the Menu section for more details. - SAL Team"
	ClientAppointmentCancellationToCounsellorTextMessage     = "Update - Your scheduled session has been cancelled by ###client_name###. - SAL Team"
	ClientAppointmentCancellationTextMessage                 = "Hi ###client_name###, your scheduled session on ###slot_bought### with ###counsellor_name### has been cancelled successfully. - SAL Team"
	CounsellorAppointmentCancellationToClientTextMessage     = "Hi ###client_name###, your scheduled session has been cancelled by your therapist. You may reschedule the appointment within the next 7 days or cancel it. - SAL Team"
	ClientAppointmentRescheduleClientToCounsellorTextMeassge = "Update - Your client, ###client_name### has rescheduled the appointment. Please click here ###navigatelink### to navigate into the app and view the change in 'Upcoming Sessions'. - SAL Team"
)
