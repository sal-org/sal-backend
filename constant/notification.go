package constant

// format event-Target
// notification headings
const (
	// client
	ClientAppointmentScheduleClientHeading   = "Appointment Confirmation"
	ClientAppointmentRescheduleClientHeading = "Appointment Reschedule Success"
	ClientAppointmentCancelClientHeading     = "Booking Cancellation"
	CounsellorAppointmentCancelClientHeading = "Appointment Cancellation"
	ClientPaymentSucessClientHeading         = "Payment Confirmation"

	// counsellor
	CounsellorAccountSignupCounsellorHeading     = "Successfully Completed"
	ClientAppointmentScheduleCounsellorHeading   = "Appointment Booking"
	ClientAppointmentRescheduleCounsellorHeading = "Appointment Reschedule"
	ClientAppointmentCancelCounsellorHeading     = "Appointment Cancellation"
	Client1AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client3AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client5AppointmentBookCounsellorHeading      = "Payment Confirmation"

	// listener
	ListenerAccountSignupListenerHeading       = "Successfully Completed"
	ClientAppointmentScheduleListenerHeading   = "Appointment Details"
	ClientAppointmentRescheduleListenerHeading = "Appointment Reschedule"
	ClientAppointmentCancelListenerHeading     = "Appointment Cancellation"

	// therapist
	TherapistAccountSignupTherapistHeading      = "Successfully Completed"
	ClientAppointmentScheduleTherapistHeading   = "Appointment Booking"
	ClientAppointmentRescheduleTherapistHeading = "Appointment Reschedule"
	ClientAppointmentCancelTherapistHeading     = "Appointment Cancellation"
	Client1AppointmentBookTherapistHeading      = "Payment Confirmation"
	Client3AppointmentBookTherapistHeading      = "Payment Confirmation"
	Client5AppointmentBookTherapistHeading      = "Payment Confirmation"
)

// notification contents
const (
	// client
	ClientAppointmentScheduleClientContent   = "Appointment for one to one chat with ###counsellor_name### confirmed for ###date_time###"
	ClientAppointmentRescheduleClientContent = "Hello, your appointment ###date_time### has been rescheduled successfully with ###counsellor_name###, please view 'My Sessions' for more details"
	ClientAppointmentCancelClientContent     = "Hi ###client_name###. You have now cancelled your scheduled session on ###date_time### with ###counsellor_name###. Refund if applicable will be processed shortly as per our Cancellation & Refund Policy"
	CounsellorAppointmentCancelClientContent = "Your scheduled session with ###counsellor_name### for ###date_time### has been cancelled by your counselor due to personal emergency. You may reschedule your appointment within the next 7 days or request for a session credit refund"
	ClientPaymentSucessClientContent         = "Hi ###client_name###. Payment of Rs. ###paid_amount### for your consultation booking has been successfully received. Manage your account anytime, anywhere from your phone."

	// counsellor
	CounsellorAccountSignupCounsellorContent     = "Thank you for creating your account and your profile. We'll let you know when your profile is approved"
	ClientAppointmentScheduleCounsellorContent   = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'My Sessions' for details."
	ClientAppointmentRescheduleCounsellorContent = "Your Client ###client_name### has rescheduled the appointment. Please go to app and view the change under 'My Sessions'"
	ClientAppointmentCancelCounsellorContent     = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client"
	Client1AppointmentBookCounsellorContent      = "###client_name### has made a payment of Rs. ###paid_amount### for a consultation session with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client3AppointmentBookCounsellorContent      = "###client_name### has made a payment of Rs. ###paid_amount### for 3 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client5AppointmentBookCounsellorContent      = "###client_name### has made a payment of Rs. ###paid_amount### for 5 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."

	// listener
	ListenerAccountSignupListenerContent       = "Thank you for creating your account and your profile. We'll let you know when your profile is approved"
	ClientAppointmentScheduleListenerContent   = "###client_name### has scheduled a call with you on ###date_time###. Please click on 'My sessions' for the meeting info"
	ClientAppointmentRescheduleListenerContent = "Your Client ###client_name### has rescheduled the appointment. Please go to app and view the change under 'My Sessions'"
	ClientAppointmentCancelListenerContent     = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client"

	// therapist
	TherapistAccountSignupTherapistContent      = "Thank you for creating your account and your profile. We'll let you know when your profile is approved"
	ClientAppointmentScheduleTherapistContent   = "You have a new counselling session booked by ###client_name### for ###date_time###. Check 'My Sessions' for details."
	ClientAppointmentRescheduleTherapistContent = "Your Client ###client_name### has rescheduled the appointment. Please go to app and view the change under 'My Sessions'"
	ClientAppointmentCancelTherapistContent     = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client"
	Client1AppointmentBookTherapistContent      = "###client_name### has made a payment of Rs. ###paid_amount### for a consultation session with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client3AppointmentBookTherapistContent      = "###client_name### has made a payment of Rs. ###paid_amount### for 3 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client5AppointmentBookTherapistContent      = "###client_name### has made a payment of Rs. ###paid_amount### for 5 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
)
