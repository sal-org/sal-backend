package constant

// miscellaneous constants
const (
	GSTPercent                               = 0.18  // 18%
	EventPrice                               = "100" // INR - amount to be paid by counsellor to create event
	MaximumAppointmentReschedule             = 2     // maximum number of times appointment can be rescheduled
	CounsellorCancellationCharges            = 0.10  // 10% of appointment amount is cancellation charges, if counsellor/therapist cancels appointment anytime
	ClientAppointmentCancellationCharges     = 0.15  // 15% of appointment amount is cancellation charges, if client cancels appointment before 4 hours
	ClientAppointmentBulkCancellationCharges = 0.15  // 15% of unused appointment amount is cancellation charges, if client cancels unused appointments
)

// urls
var URLs = map[string]string{
	"privacy": "https://sal-foundation.com/app_pp",
	"about":   "https://sal-foundation.com/about-sal",
	"faq":     "https://sal-foundation.com/app_faqs",
	"info":    "https://sal-foundation.com/about-sal",
	"tos":     "https://sal-foundation.com/apptos",
}
