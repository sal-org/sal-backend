package constant

// format event-Target
// notification headings
const (
	// client
	ClientAppointmentReminderClientHeading                = "Session about to begin"
	ClientAppointmentFollowUpSessionReminderClientHeading = "Follow-up session about to begin"
	ClientEventReminderClientHeading                      = "Event Reminder"
	ClientAppointmentScheduleClientHeading                = "Booking confirmation"
	ClientAppointmentFollowUpSessionCounsellorHeading     = "Followup session booked!"
	ClientAppointmentFollowUpSessionClientHeading         = "Follow-up session confirmation"
	ClientAppointmentRescheduleClientHeading              = "Clove: Appointment has been rescheduled"
	ClientAppointmentCancelClientHeading                  = "Session cancelled!"
	ClientBulkAppointmentCancelClientHeading              = "All Session cancelled!"
	CounsellorAppointmentCancelClientHeading              = "Your booking cancellation"
	CounsellorAppointmentCancelCounsellorHeading          = "Appointment Cancellation"
	ClientPaymentSucessClientHeading                      = "Clove: Payment Confirmation"
	ClientAppointmentFeedbackHeading                      = "Your feedback"
	ClientAppointmentHasBeenStartedHeading                = "Your session has started!"
	ClientAppointmentGivenFeedbackHeading                 = "Thanks for your Feedback!"
	ClientEventPaymentSucessClientHeading                 = "Cafe booked successfully!"
	ClientSelectSadMoodHeading                            = "Want to speak to somebody now?"
	ClientCompletedProfileHeading                         = "Congratulations, you are registered."

	// counsellor
	ClientAppointmentReminderCounsellorHeading   = "Session about to begin"
	CounsellorAppointmentHasBeenStartedHeading   = "Your session has started!"
	CounsellorEventReminderCounsellorHeading     = "Event Reminder"
	CounsellorAccountSignupCounsellorHeading     = "Registration successful"
	ClientAppointmentScheduleCounsellorHeading   = "New session booked"
	ClientAppointmentRescheduleCounsellorHeading = "Clove: Appointment has been rescheduled"
	ClientAppointmentCancelCounsellorHeading     = "Session cancelled!"
	Client1AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client3AppointmentBookCounsellorHeading      = "Payment Confirmation"
	Client5AppointmentBookCounsellorHeading      = "Payment Confirmation"
)

// notification contents
const (
	// client
	ClientAppointmentHasBeenStartedContent            = "Hey ###clientname###, ###therapistname### has joined the session. Please join"
	ClientAppointmentRemiderClientContent             = "Hi ###user_name###, your session is starting soon at ###time###. Tap here to know more."
	ClientAppointmentFollowUpRemiderClientContent     = "Hi ###client_name###, your follow-up session with ###therapistname### is starting in 15 minutes. Please do find a quite place for yourself to take it."
	ClientAppointmentFollowUpRemiderCounsellorContent = "You have a scheduled follow-up session starting in 15 minutes with ###client_name### at ###time###. Check 'My sessions' within the app for details."
	ClientEventRemiderClientContent                   = "Event with ###counsellor_name### starts in 15 min. Please check your internet connectivity and get ready."
	ClientAppointmentScheduleClientContent            = "Your appointment for one to one chat is confirmed"
	ClientAppointmentFollowUpSessionCounsellorContent = "You have a followup counselling session booked by ###client_name### for ###date_time###. Check 'Booking section' for details."
	ClientAppointmentFollowUpSessionClientContent     = "Your follow-up appointment is confirmed with ###therpist_name###"
	ClientAppointmentRescheduleClientContent          = "Hey, you have rescheduled your session to ###date_time### with ###therapistname###, please view 'my sessions' for more details."
	ClientAppointmentCancelClientContent              = "You have cancelled your scheduled session on ###datetime### with ###therapistname###."
	ClientBulkAppointmentCancelClientContent          = "You have cancelled all your pending sessions with ###counsellor_name###. "
	CounsellorAppointmentCancelClientContent          = "Your session has been cancelled by ###therapist_name### for ###date_time### due to personal emergency. Please reschedule your appointment within the next 7 days or request for a session credit refund."
	CounsellorAppointmentCancelCounsellorContent      = "You have now cancelled your scheduled session on ###date_time### with ###client_name###."
	ClientPaymentSucessClientContent                  = "Payment of Rs. ###paid_amount### has been successful. You can manage your account anytime, anywhere from your phone. Click on My sessions to view details"
	ClientAppointmentFeedbackContent                  = "Your feedback is so important to us. Please rate your session with ###counsellor_name### to help us constantly improve our services"
	ClientAppointmentFeedbackGivenContent             = "Thank you for rating our Counselling Services and helping us in making a difference! Kudos!"
	ClientEventPaymentSucessClientContent             = "SAL Cafe ###cafe_name### has been booked successfully for Rs. ###paid_amount### on ###date### & ###time###."
	ClientSelectSadMoodContent                        = "We sensed you are ###mood###. Talk to our compassionate listeners or speak with our expert therapist. Book your appointment to speak to somebody now."
	ClientCompletedProfileContent                     = "Thank you for creating your profile with CLOVE"

	// counsellor
	CounsellorAppointmentHasBeenStartedContent   = "Hey ###therapistname###, ###clientname### has joined the session. Please join"
	ClientAppointmentReminderCounsellorContent   = "You have a scheduled counselling session starting in 15 minutes with ###clientname### at ###time###. Check 'Booking section' within the app for details."
	CounsellorEventReminderCounsellorContent     = "Your event will start in 15 min. Please check your internet connectivity and get ready."
	CounsellorAccountSignupCounsellorContent     = "Our offline team will contact you shortly."
	ClientAppointmentScheduleCounsellorContent   = "Your session is confirmed ###Date### and ###Time###. Check 'Booking section' for details. "
	ClientAppointmentRescheduleCounsellorContent = "Appointment has been rescheduled. Please click here and view the change in the 'Booking section' for details"
	ClientAppointmentCancelCounsellorContent     = "Your scheduled session with ###client_name### for ###date_time### has been cancelled by the client"
	Client1AppointmentBookCounsellorContent      = "Payment of Rs. ###paid_amount### has been successful. You can manage your account anytime, anywhere from your phone. Click on 'My sessions' to view details "
	Client3AppointmentBookCounsellorContent      = "Hi. ###client_name### has made a payment of Rs. ###paid_amount### for 3 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
	Client5AppointmentBookCounsellorContent      = "Hi. ###client_name### has made a payment of Rs. ###paid_amount### for 5 consultation sessions with you.  Manage your account anytime, anywhere from your phone. Click on 'My Sessions' to view Date and Time."
)
