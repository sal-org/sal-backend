package constant

// miscellaneous constants
const (
	GSTPercent                               = 18   // 18% GST Write as 18 Only
	EventPrice                               = "0"  // INR - amount to be paid by counsellor to create event
	MaximumAppointmentReschedule             = 2    // maximum number of times appointment can be rescheduled
	CounsellorCancellationCharges            = 0.20 // 20% of appointment amount is cancellation charges, if counsellor/therapist cancels appointment anytime
	ClientAppointmentCancellationCharges     = 0.20 // 20% of appointment amount is cancellation charges, if client cancels appointment before 4 hours
	ClientAppointmentBulkCancellationCharges = 0.20 // 20% of unused appointment amount is cancellation charges, if client cancels unused appointments
	CounsellorPayoutPercentage               = 50   // 50% Pay the Counsellor for 1 counselling session
	CounsellorPayoutPercentageColumns        = "50" // these only col to insert in database
)

// lastest version
const (
	ClientIOSVersion        = "2.15"
	ClientAndroidVersion    = "1.16.0"
	TherapistIOSVersion     = "2.11"
	TherapistAndroidVersion = "1.15.0"
)

// urls
var URLs = map[string]string{
	"privacy": "https://salapp.sal-foundation.com/app_pp/",
	"about":   "https://sal-foundation.com/about-sal",
	"faq":     "https://salapp.sal-foundation.com/app_faqs/",
	"info":    "https://salapp.sal-foundation.com/info/",
	"tos":     "https://salapp.sal-foundation.com/app_tos/",
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

const (
	EventDuration = "60" // Event duration 60 mins
)

// default photo for client and listerner

const (
	DefaultPhotoForClientAndListerner = "miscellaneous/kmPlQ123NSP123.jpg"
)
