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
	CounsellorAccountSignupTextMessage   = "Hey ###counsellor_name###, the onboarding process has been completed successfully. We will contact you shortly to discuss the next steps. - Team Clove"
	CounsellorOTPTextMessage             = "###otp### is the OTP to access the Mobile App. Please do not share it with anyone. - Team SAL"
	ClientOTPTextMessage                 = "###otp### is the OTP to register yourself on the Clove Mobile App. Please do not share it with anyone.  - Team Clove"
	ClientProfileTitleMessage            = "Hey ###client_name###, welcome to Clove! Access self-care assessments and content, daily mood diaries and book counselling sessions with Listeners & Therapists.  - Team Clove"
	ListenerOTPTextMessage               = "###otp### is the OTP to register yourself on SAL Mobile App. You are only few steps away from your taking your 1st listening session."
	ClientAppointmentReminderTextMessage = "Hey ###user_name###, just a quick reminder for your upcoming session starting soon at ###time### with ###userName###. - Team Salubrium" //"Hey ###user_name###, just a quick appointment reminder for your upcoming session with ###userName### in the Clove mobile app. It is starting soon at ###time###. - Team Salubrium."
	// ClientAppointmentScheduleCounsellorTextMessage           = "Hi ###counsellor_name###, you have a new counselling session booked by ###client_name### for ###date_time###. Please visit My sessions in the SAL Mobile app menu for more details."
	ClientAppointmentConfirmationTextMessage = "Hi ###userName###, your appointment with ###user_Name### is confirmed in the Clove mobile app. Date: ###date### Time: ###time### - Team Salubrium"
	ClientRefundAmonutTextMessage            = "Refund initiated: Rs. ###amount### for your appointment dated ###date### at ###time### in the Clove mobile app is being processed to your card/bank account and will reflect in 7-10 working days. - Team Salubrium"
	// ClientPaymentConfirmationTextMeassge                     = "Hi ###client_name###, payment of INR ###amount### for ###bought### session(s) has been received. Please click ###Aaplink### to access the receipt. - SAL Team"
	ClientAppointmentRescheduleClientTextMeassge         = "Hi ###clientName###, your upcoming session with ###therapistName### has been rescheduled to ###date### at ###time###. For more information, view ‘My Sessions’ in the Clove mobile app - Team Salubrium"
	ClientAppointmentCancellationToCounsellorTextMessage = "Hi ###userName###, your upcoming session with ###user_Name### for ###date### at ###time### has been cancelled in the Clove mobile app. - Team Salubrium"
	// ClientAppointmentCancellationTextMessage                 = "Hi ###client_name###, your scheduled session on ###slot_bought### with ###counsellor_name### has been cancelled successfully. - SAL Team"
	// CounsellorAppointmentCancellationToClientTextMessage     = "Hi ###client_name###, your scheduled session has been cancelled by your therapist. You may reschedule the appointment within the next 7 days or cancel it. - SAL Team"
	ClientAppointmentRescheduleClientToCounsellorTextMeassge         = "Hi ###counsellorName###, your upcoming session with ###clientName### has been rescheduled to ###date### at ###time###. For more information, go to the ‘Bookings’ section in the Clove mobile app - Team Salubrium"
	ClientAppointmentRequestSMSToCounsellorTextMessage               = "Hi ###counsellorName###, your client ###clientName###, wants a session with you on the Clove Mind app. Update your 15 days availability within 24 hrs and we will let client know. Team Salubrium"
	ClientAppointmentRequestSMSToClientTextMessage                   = "Hi ###clientName###, thank you for requesting an appointment with ###counsellorName### on the Clove Mind app. It has been duly noted. Team Salubrium"
	ClientAppointmentRequestSMSToNotAcceptedClientTextMessage        = "Hello, your therapist ###counsellorName###, has still not updated the availability in Clove Mind app. You may choose to wait longer or select any other therapist. Team Salubrium"
	CounsellorAppointmentRequestCreatedAvailabilityClientTextMessage = "Hello ###clientName###, your therapist ###counsellorName### has updated availability on the Clove Mind app. Login to book your session and continue your wellbeing journey. Team Salubrium"
)
