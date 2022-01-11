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
	CounsellorAccountSignupTextMessage             = "Hi ###counsellor_name###, thank you for successfully completing your profile. One of SAL's team members will contact you shortly to discuss the next steps."
	CounsellorOTPTextMessage                       = "###otp### is the OTP to register yourself on SAL Mobile App. You are only few steps away from taking your 1st counselling appointment."
	ClientOTPTextMessage                           = "Dear customer, ###otp### is the OTP to register yourself on SAL Mobile App. You can now access self-care content and book appointments for your emotional well-being!"
	ListenerOTPTextMessage                         = "###otp### is the OTP to register yourself on SAL Mobile App. You are only few steps away from your taking your 1st listening session."
	ClientAppointmentScheduleCounsellorTextMessage = "Hi ###counsellor_name###, you have a new counselling session booked by ###client_name### for ###date_time###. Please visit My sessions in the SAL Mobile app menu for more details."
)
