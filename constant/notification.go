package constant

// format event-Target
// notification headings
const (
	// client
	ClientAppointmentReminderClientHeading   = "Appointment Reminder"
	ClientEventReminderClientHeading         = "Event Reminder"
	ClientAppointmentScheduleClientHeading   = "Appointment Confirmation"
	ClientAppointmentRescheduleClientHeading = "Appointment Reschedule Success"
	ClientAppointmentCancelClientHeading     = "Booking Cancellation"
	ClientBulkAppointmentCancelClientHeading = "Order Cancellation"
	CounsellorAppointmentCancelClientHeading = "Appointment Cancellation"
	ClientPaymentSucessClientHeading         = "Payment Confirmation"
	ClientAppointmentFeedbackHeading         = "Feedback/Ratings"
	ClientEventPaymentSucessClientHeading    = "Cafe booked successfully!"

	// counsellor
	ClientAppointmentReminderCounsellorHeading   = "Appointment Reminder"
	CounsellorEventReminderCounsellorHeading     = "Event Reminder"
	CounsellorAccountSignupCounsellorHeading     = "Successfully Completed"
	ClientAppointmentScheduleCounsellorHeading   = "Appointment Booking"
	ClientAppointmentRescheduleCounsellorHeading = "Appointment Reschedule"
	ClientAppointmentCancelCounsellorHeading     = "Appointment Cancellation"
	Client1AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client3AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client5AppointmentBookCounsellorHeading      = "Payment Confirmation"
)

// notification contents
const (
	// client
	ClientAppointmentRemiderClientContent    = "Appointment for one to one chat with ###counsellor_name### starts in 15 min. Please check your internet connectivity and get ready for the video call"
	ClientEventRemiderClientContent          = "Event with ###counsellor_name### starts in 15 min. Please check your internet connectivity and get ready."
	ClientAppointmentScheduleClientContent   = "Appointment for one to one chat with ###counsellor_name### confirmed for ###date_time###"
	ClientAppointmentRescheduleClientContent = "Hello, your appointment ###date_time### has been rescheduled successfully with ###counsellor_name###. Please view 'My sessions' in Menu for more details"
	ClientAppointmentCancelClientContent     = "Hi ###client_name###. You have now cancelled your scheduled session on ###date_time### with ###counsellor_name###"
	ClientBulkAppointmentCancelClientContent = "Hi ###client_name###. You have cancelled your unused appointments with ###counsellor_name###"
	CounsellorAppointmentCancelClientContent = "Your scheduled session with ###counsellor_name### for ###date_time### has been cancelled by your counselor due to personal emergency. You may reschedule your appointment within the next 7 days or request for a session credit refund"
	ClientPaymentSucessClientContent         = "Hi ###client_name###. Payment of Rs. ###paid_amount### for your consultation booking has been Successful/Received. Manage your account anytime, anywhere from your phone. Click to view your Transaction ID #  and receipt"
	ClientAppointmentFeedbackContent         = "Please rate your session with ###counsellor_name### to help us constantly improve our services"
	ClientEventPaymentSucessClientContent    = "SAL Cafe ###cafe_name### has been booked successfully for Rs. ###paid_amount### on ###date### & ###time###."

	// counsellor
	ClientAppointmentReminderCounsellorContent   = "You have a counselling session to ###client_name### starts in 15 min. Please check your internet connectivity"
	CounsellorEventReminderCounsellorContent     = "Your event will start in 15 min. Please check your internet connectivity and get ready."
	CounsellorAccountSignupCounsellorContent     = "Thank you for creating your account and your profile. We'll let you know when your profile is approved"
	ClientAppointmentScheduleCounsellorContent   = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'My sessions' for details."
	ClientAppointmentRescheduleCounsellorContent = "Your client has rescheduled a current booking. Please click here and view the change under 'My Sessions'"
	ClientAppointmentCancelCounsellorContent     = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client"
	Client1AppointmentBookCounsellorContent      = "Hi. ###client_name### has made a payment of Rs. ###paid_amount### for a consultation session with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client3AppointmentBookCounsellorContent      = "Hi. ###client_name### has made a payment of Rs. ###paid_amount### for 3 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client5AppointmentBookCounsellorContent      = "Hi. ###client_name### has made a payment of Rs. ###paid_amount### for 5 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
)
