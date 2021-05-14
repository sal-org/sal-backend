package constant

// notification headings
const (
	// client
	ClientListenerAppointmentRescheduleHeading   = "Appointment Reschedule Success"
	ClientCounsellorAppointmentRescheduleHeading = "Appointment Reschedule Success"
	ClientAppointmentCancelHeading               = "Booking Cancellation"

	// counsellor
	CounsellorAccountSignupHeading = "Successfully Completed"

	// listener
)

// notification contents
const (
	// client
	ClientListenerAppointmentRescheduleContent   = "Hello, your appointment ###date_time### has been rescheduled successfully with ###listener_name###, please view my sessions for more details"
	ClientCounsellorAppointmentRescheduleContent = "Hello, your appointment ###date_time### has been rescheduled successfully with ###counsellor_name###, please view my sessions for more details"
	ClientAppointmentCancelContent               = "Hi ###client_name### You have now cancelled your scheduled session on ###date_time### with ###counsellor_name###. Refund if applicable will be processed shortly as per our Cancellation & Refund Policy"

	// counsellor
	CounsellorAccountSignupContent = "Thank you for creating your account and your profile. We'll let you know when your profile is approved"

	// listener
)
