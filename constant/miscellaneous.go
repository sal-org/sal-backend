package constant

// miscellaneous constants
const (
	GSTPercent                               = 18   // 18% GST Write as 18 Only
	EventPrice                               = "0"  // INR - amount to be paid by counsellor to create event
	MaximumAppointmentReschedule             = 2    // maximum number of times appointment can be rescheduled
	CounsellorCancellationCharges            = 0.10 // 10% of appointment amount is cancellation charges, if counsellor/therapist cancels appointment anytime
	ClientAppointmentCancellationCharges     = 0.20 // 15% of appointment amount is cancellation charges, if client cancels appointment before 4 hours
	ClientAppointmentBulkCancellationCharges = 0.20 // 15% of unused appointment amount is cancellation charges, if client cancels unused appointments
)

// urls
var URLs = map[string]string{
	"privacy": "http://salapp.sal-foundation.com/app_pp",
	"about":   "https://sal-foundation.com/about-sal",
	"faq":     "http://salapp.sal-foundation.com/app_faqs/",
	"info":    "https://sal-foundation.com/about-sal",
	"tos":     "http://salapp.sal-foundation.com/app_tos/",
}

// default images used for events
var EventImages = []string{
	"event/event_image_1.png",
	"event/event_image_2.png",
	"event/event_image_3.png",
	"event/event_image_4.png",
	"event/event_image_5.png",
	"event/event_image_6.png",
}
