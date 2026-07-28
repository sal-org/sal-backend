package model

// CounsellorOrderCreateRequest .
type CounsellorOrderCreateRequest struct {
	ClientID     string `json:"client_id"`
	CounsellorID string `json:"counsellor_id"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	CouponCode   string `json:"coupon_code"`
	NoSessions   string `json:"no_session"`
}

// ListenerOrderCreateRequest .
type ListenerOrderCreateRequest struct {
	ClientID   string `json:"client_id"`
	ListenerID string `json:"listener_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
}

// TherapistOrderCreateRequest .
type TherapistOrderCreateRequest struct {
	ClientID    string `json:"client_id"`
	TherapistID string `json:"therapist_id"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	CouponCode  string `json:"coupon_code"`
	NoSessions  string `json:"no_session"`
}

// CounsellorOrderPaymentCompleteRequest .
type CounsellorOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// TherapistOrderPaymentCompleteRequest .
type TherapistOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// ListenerOrderPaymentCompleteRequest .
type ListenerOrderPaymentCompleteRequest struct {
	OrderID string `json:"order_id"`
}

// AppointmentRescheduleRequest .
type AppointmentRescheduleRequest struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

// AppointmentBookRequest .
type AppointmentBookRequest struct {
	AppointmentSlotID string `json:"appointment_slot_id"`
	Date              string `json:"date"`
	Time              string `json:"time"`
}

// EventOrderCreateRequest .
type EventOrderCreateRequest struct {
	UserID       string `json:"user_id"`
	EventOrderID string `json:"event_order_id"`
	CouponCode   string `json:"coupon_code"`
}

// EventBlockOrderCreateRequest .
type EventBlockOrderCreateRequest struct {
	CounsellorID string `json:"counsellor_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	TopicID      string `json:"topic_id"`
	Date         string `json:"date"`
	Photo        string `json:"photo"`
	Time         string `json:"time"`
	Price        string `json:"price"`
}

// EventOrderPaymentCompleteRequest .
type EventOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// EventBlockOrderPaymentCompleteRequest .
type EventBlockOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// CounsellorProfileAddRequest .
type CounsellorProfileAddRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Gender            string `json:"gender"`
	Phone             string `json:"phone"`
	Photo             string `json:"photo"`
	Email             string `json:"email"`
	Price             string `json:"price"`
	Multiple_Sessions string `json:"multiple_sessions"`
	Price3            string `json:"price_3"`
	Price5            string `json:"price_5"`
	Education         string `json:"education"`
	Experience        string `json:"experience"`
	About             string `json:"about"`
	Timezone          string `json:"timezone"`
	TopicIDs          string `json:"topic_ids"`
	LanguageIDs       string `json:"language_ids"`
	Resume            string `json:"resume"`
	Certificate       string `json:"certificate"`
	Aadhar            string `json:"aadhar"`
	Linkedin          string `json:"linkedin"`
	DeviceID          string `json:"device_id"`
	PayoutPercentage  string `json:"payout_percentage"`
	PayeeName         string `json:"payee_name"`
	BankAccountNumber string `json:"bank_account_no"`
	IFSC              string `json:"ifsc"`
	BranchName        string `json:"branch_name"`
	BankName          string `json:"bank_name"`
	BankAccountType   string `json:"bank_account_type"`
	PAN               string `json:"pan"`
}

// ListenerProfileAddRequest .
type ListenerProfileAddRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	AgeGroup    string `json:"age_group"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
	Email       string `json:"email"`
	Occupation  string `json:"occupation"`
	Aadhar      string `json:"aadhar"`
	About       string `json:"about"`
	Timezone    string `json:"timezone"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	DeviceID    string `json:"device_id"`
}

// TherapistProfileAddRequest .
type TherapistProfileAddRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Gender            string `json:"gender"`
	Phone             string `json:"phone"`
	Photo             string `json:"photo"`
	Email             string `json:"email"`
	Price             string `json:"price"`
	Multiple_Sessions string `json:"multiple_sessions"`
	Price3            string `json:"price_3"`
	Price5            string `json:"price_5"`
	Education         string `json:"education"`
	Experience        string `json:"experience"`
	About             string `json:"about"`
	Timezone          string `json:"timezone"`
	TopicIDs          string `json:"topic_ids"`
	LanguageIDs       string `json:"language_ids"`
	Resume            string `json:"resume"`
	Certificate       string `json:"certificate"`
	Aadhar            string `json:"aadhar"`
	Linkedin          string `json:"linkedin"`
	DeviceID          string `json:"device_id"`
	PayoutPercentage  string `json:"payout_percentage"`
	PayeeName         string `json:"payee_name"`
	BankAccountNumber string `json:"bank_account_no"`
	IFSC              string `json:"ifsc"`
	BranchName        string `json:"branch_name"`
	BankName          string `json:"bank_name"`
	BankAccountType   string `json:"bank_account_type"`
	PAN               string `json:"pan"`
}

// // ClientProfileUpdateRequest .
// type ClientProfileUpdateRequest struct {
// 	FirstName   string `json:"first_name"`
// 	LastName    string `json:"last_name"`
// 	Location    string `json:"location"`
// 	Timezone    string `json:"timezone"`
// 	DeviceID    string `json:"device_id"`
// 	DateOfBirth string `json:"date_of_birth"`
// 	Photo       string `json:"photo"`
// 	TopicIDs    string `json:"topic_ids"`
// 	Gender      string `json:"gender"`
// }

// CounsellorProfileUpdateRequest .
type CounsellorProfileUpdateRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Gender            string `json:"gender"`
	Photo             string `json:"photo"`
	Price             string `json:"price"`
	Multiple_Sessions string `json:"multiple_sessions"`
	Price3            string `json:"price_3"`
	Price5            string `json:"price_5"`
	Education         string `json:"education"`
	Experience        string `json:"experience"`
	About             string `json:"about"`
	Timezone          string `json:"timezone"`
	TopicIDs          string `json:"topic_ids"`
	LanguageIDs       string `json:"language_ids"`
	Resume            string `json:"resume"`
	Certificate       string `json:"certificate"`
	Aadhar            string `json:"aadhar"`
	Linkedin          string `json:"linkedin"`
	DeviceID          string `json:"device_id"`
	PayoutPercentage  string `json:"payout_percentage"`
	PayeeName         string `json:"payee_name"`
	BankAccountNumber string `json:"bank_account_no"`
	IFSC              string `json:"ifsc"`
	BranchName        string `json:"branch_name"`
	BankName          string `json:"bank_name"`
	BankAccountType   string `json:"bank_account_type"`
	PAN               string `json:"pan"`
}

// ListenerProfileUpdateRequest .
type ListenerProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	AgeGroup    string `json:"age_group"`
	Photo       string `json:"photo"`
	Occupation  string `json:"occupation"`
	Aadhar      string `json:"aadhar"`
	About       string `json:"about"`
	Timezone    string `json:"timezone"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	DeviceID    string `json:"device_id"`
}

// TherapistProfileUpdateRequest .
type TherapistProfileUpdateRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Gender            string `json:"gender"`
	Photo             string `json:"photo"`
	Price             string `json:"price"`
	Multiple_Sessions string `json:"multiple_sessions"`
	Price3            string `json:"price_3"`
	Price5            string `json:"price_5"`
	Education         string `json:"education"`
	Experience        string `json:"experience"`
	About             string `json:"about"`
	Timezone          string `json:"timezone"`
	TopicIDs          string `json:"topic_ids"`
	LanguageIDs       string `json:"language_ids"`
	Resume            string `json:"resume"`
	Certificate       string `json:"certificate"`
	Aadhar            string `json:"aadhar"`
	Linkedin          string `json:"linkedin"`
	DeviceID          string `json:"device_id"`
	PayoutPercentage  string `json:"payout_percentage"`
	PayeeName         string `json:"payee_name"`
	BankAccountNumber string `json:"bank_account_no"`
	IFSC              string `json:"ifsc"`
	BranchName        string `json:"branch_name"`
	BankName          string `json:"bank_name"`
	BankAccountType   string `json:"bank_account_type"`
	PAN               string `json:"pan"`
}

// AssessmentAddRequest .
type AssessmentAddRequest struct {
	UserID       string `json:"user_id" validate:"required,min=5,max=45"`
	Name         string `json:"name" validate:"required"`
	Age          string `json:"age" validate:"required,numeric"`
	Gender       string `json:"gender" validate:"required,oneof=Male Female Other"`
	Phone        string `json:"phone" validate:"required,len=12,numeric"`
	AssessmentID string `json:"assessment_id" validate:"required,min=5,max=45"`
	Feedback     string `json:"feedback" validate:"omitempty,min=2"`
	Details      []struct {
		AssessmentQuestionID       string `json:"assessment_question_id" validate:"required,min=5,max=45"`
		AssessmentQuestionOptionID string `json:"assessment_question_option_id" validate:"required,min=5,max=45"`
		Score                      string `json:"score" validate:"required"`
	} `json:"details" validate:"required,min=1,dive"`
}

type ClientUpdateProfileRequestInAdminPanel struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Phone       string `json:"phone" validate:"required,len=12,numeric"`
	Email       string `json:"email" validate:"required,email,max=100"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Gender      string `json:"gender" validate:"required,oneof=Male Female Other"`
	Status      string `json:"status" validate:"required,oneof=0 1"`
	ModifiedBy  string `json:"modified_by"`
}

type ContentAddRequestInAdminPanel struct {
	CounsellorID    string `json:"counsellor_id"`
	Title           string `json:"title" validate:"required,min=3,max=200"`
	Description     string `json:"description" validate:"required,max=1000"`
	Photo           string `json:"photo" validate:"required"`
	BackgroundPhoto string `json:"background_photo" validate:"required"`
	ShareContent    string `json:"share_content" validate:"required"`
	Content         string `json:"content" validate:"required"`
	Type            string `json:"type" validate:"required,oneof=1 2 3"`
	Redirection     string `json:"redirection" validate:"omitempty, oneof=1 2 3"`
	CategoryID      string `json:"category_id" validate:"required"`
	Training        string `json:"training" validate:"required,oneof=0 1"`
	MoodID          string `json:"mood_id" validate:"required"`
	Duration        string `json:"duration" validate:"required"`
}

type ContentInWebAddRequestInAdminPanel struct {
	CounsellorID    string `json:"counsellor_id"`
	Title           string `json:"title" validate:"required,min=3,max=200"`
	SubTitle        string `json:"subtitle" validate:"required,min=3,max=1000"`
	Description     string `json:"description" validate:"required,max=1000"`
	Photo           string `json:"photo" validate:"required"`
	BackgroundPhoto string `json:"background_photo" validate:"required"`
	ShareContent    string `json:"share_content" validate:"required"`
	Content         string `json:"content" validate:"required"`
	Type            string `json:"type" validate:"required,oneof=1 2 3"`
	Redirection     string `json:"redirection" validate:"omitempty,oneof=1 2 3 4"`
	CategoryID      string `json:"category_id" validate:"required"`
	ResourceID      string `json:"resource_id" validate:"required,oneof=1 2 3 4"`
	ContentMood     string `json:"content_mode" validate:"required,oneof=1 2 3 4"`
	MoodID          string `json:"mood_id" validate:"required"`
	Duration        string `json:"duration" validate:"required"`
}

type ContentInWebUpdateRequestInAdminPanel struct {
	CounsellorID    string `json:"counsellor_id"`
	Title           string `json:"title" validate:"required,min=3,max=200"`
	SubTitle        string `json:"subtitle" validate:"required,min=3,max=1000"`
	Description     string `json:"description" validate:"required,max=1000"`
	Photo           string `json:"photo" validate:"required"`
	BackgroundPhoto string `json:"background_photo" validate:"required"`
	ShareContent    string `json:"share_content" validate:"required"`
	Content         string `json:"content" validate:"required"`
	Type            string `json:"type" validate:"required,oneof=1 2 3"`
	Redirection     string `json:"redirection" validate:"omitempty,oneof=1 2 3 4"`
	CategoryID      string `json:"category_id" validate:"required"`
	ResourceID      string `json:"resource_id" validate:"required,oneof=1 2 3 4"`
	ContentMood     string `json:"content_mode" validate:"required,oneof=1 2 3 4"`
	MoodID          string `json:"mood_id" validate:"required"`
	Duration        string `json:"duration" validate:"required"`
	Status          string `json:"status" validate:"required"`
}

type ContentUpdateRequestInAdminPanel struct {
	CounsellorID    string `json:"counsellor_id"`
	Title           string `json:"title" validate:"required,min=3,max=200"`
	Description     string `json:"description" validate:"required,max=1000"`
	Photo           string `json:"photo" validate:"required"`
	BackgroundPhoto string `json:"background_photo" validate:"required"`
	ShareContent    string `json:"share_content" validate:"required"`
	Content         string `json:"content" validate:"required"`
	Type            string `json:"type" validate:"required,oneof=1 2 3"`
	Redirection     string `json:"redirection" validate:"omitempty, oneof=1 2 3"`
	CategoryID      string `json:"category_id" validate:"required"`
	Training        string `json:"training" validate:"required,oneof=0 1"`
	MoodID          string `json:"mood_id" validate:"required"`
	Duration        string `json:"duration" validate:"required"`
	Status          string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type CorporatePartnerAddRequestInAdminPanel struct {
	PartnerName string `json:"partnerName" validate:"required,min=2,max=100"`
	Domain      string `json:"domain" validate:"required,min=4,max=25"`
	AccessCode  string `json:"accessCode" validate:"required,min=4,max=20"`
}

type CorporatePartnerUpdateRequestInAdminPanel struct {
	PartnerName string `json:"partnerName" validate:"required,min=2,max=100"`
	Domain      string `json:"domain" validate:"required,min=4,max=25"`
	AccessCode  string `json:"accessCode" validate:"required,min=4,max=20"`
	Status      string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type CorporatePartnerAddressAddRequestInAdminPanel struct {
	PartnerName string `json:"partnerName" validate:"required,min=2,max=100"`
	Domain      string `json:"domain" validate:"required,min=4,max=25"`
	Address     string `json:"address" validate:"required,min=4,max=300"`
}

type CorporatePartnerAddressUpdateRequestInAdminPanel struct {
	PartnerName string `json:"partnerName" validate:"required,min=2,max=100"`
	Domain      string `json:"domain" validate:"required,min=4,max=25"`
	Address     string `json:"address" validate:"required,min=4,max=300"`
	Status      string `json:"status" validate:"required,oneof= 0 1 2"`
}

type UpdateCounsellorProfileRequestInAdminPanel struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"required,len=12,numeric"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Gender    string `json:"gender" validate:"required,oneof=Male Female Other"`

	Price          string `json:"price" validate:"required,numeric"`
	Price3         string `json:"price_3" validate:"required,numeric"`
	Price5         string `json:"price_5" validate:"required,numeric"`
	CorporatePrice string `json:"corporate_price" validate:"required,numeric"`

	Education  string `json:"education" validate:"required,max=500"`
	Experience string `json:"experience" validate:"required,max=500"`
	About      string `json:"about" validate:"required,max=5000"`

	PayoutPercentage string `json:"payout_percentage" validate:"required,numeric"`

	PayeeName       string `json:"payee_name" validate:"required,max=100"`
	BankAccountNo   string `json:"bank_account_no" validate:"required,numeric,min=9,max=18"`
	IFSC            string `json:"ifsc" validate:"required,len=11,alphanum"`
	BranchName      string `json:"branch_name" validate:"required,max=100"`
	BankName        string `json:"bank_name" validate:"required,max=100"`
	BankAccountType string `json:"bank_account_type" validate:"required,oneof=Savings Current"`

	PAN string `json:"pan" validate:"required,len=10,alphanum"`

	CorporateTherapist string `json:"corporate_therpist" validate:"required,oneof=0 1 2 3"`
	Status             string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type AddMotivationalQuotaInAdminPanel struct {
	Quota  string `json:"quote" validate:"required,min=2,max=400"`
	MoodID string `json:"mood_id" validate:"required,numeric,max=2"`
}

type UpdateTherapistProfileRequestInAdminPanel struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"required,len=12,numeric"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Gender    string `json:"gender" validate:"required,oneof=Male Female Other"`

	Price          string `json:"price" validate:"required,numeric"`
	Price3         string `json:"price_3" validate:"required,numeric"`
	Price5         string `json:"price_5" validate:"required,numeric"`
	CorporatePrice string `json:"corporate_price" validate:"required,numeric"`

	Education  string `json:"education" validate:"required,max=500"`
	Experience string `json:"experience" validate:"required,max=500"`
	About      string `json:"about" validate:"required,max=5000"`

	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	GapYears  string `json:"gap_years" validate:"required,numeric"`
	GapMonths string `json:"gap_months" validate:"required,numeric"`
	Location  string `json:"location" validate:"required,max=200"`
	Video     string `json:"video" validate:"required"`

	PayoutPercentage string `json:"payout_percentage" validate:"required,numeric"`

	PayeeName       string `json:"payee_name" validate:"required,max=100"`
	BankAccountNo   string `json:"bank_account_no" validate:"required,numeric,min=9,max=18"`
	IFSC            string `json:"ifsc" validate:"required,len=11,alphanum"`
	BranchName      string `json:"branch_name" validate:"required,max=100"`
	BankName        string `json:"bank_name" validate:"required,max=100"`
	BankAccountType string `json:"bank_account_type" validate:"required,oneof=Savings Current"`

	PAN string `json:"pan" validate:"required,len=10,alphanum"`

	CorporateTherapist string `json:"corporate_therpist" validate:"required,oneof=0 1 2 3"`
	Status             string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type AddTherapistProfileRequest struct {
	FirstName           string `json:"first_name" validate:"required,min=2,max=50"`
	LastName            string `json:"last_name" validate:"required,min=2,max=50"`
	Pronoun             string `json:"pronoun" validate:"required,min=2,max=50"`
	Gender              string `json:"gender" validate:"required,oneof=Male Female Other"`
	Location            string `json:"location" validate:"required,max=200"`
	Phone               string `json:"phone" validate:"required,len=12,numeric"`
	Photo               string `json:"photo" validate:"required"`
	Email               string `json:"email" validate:"required,email,max=100"`
	Price               string `json:"price" validate:"required,numeric"`
	MultipleSessions    string `json:"multiple_sessions" validate:"required,numeric"`
	Price3              string `json:"price_3" validate:"omitempty,numeric"`
	Price5              string `json:"price_5" validate:"omitempty,numeric"`
	Education           string `json:"education" validate:"required,max=500"`
	Experience          string `json:"experience" validate:"required,max=500"`
	About               string `json:"about" validate:"required,max=5000"`
	TherapeuticApproach string `json:"therapeutic_approach" validate:"required,max=5000"`
	StartDate           string `json:"start_date" validate:"required,datetime=2006-01-02"`
	GapYears            string `json:"gap_years" validate:"required,numeric"`
	GapMonths           string `json:"gap_months" validate:"required,numeric"`
	PayoutPercentage    string `json:"payout_percentage" validate:"required,numeric"`
	PayeeName           string `json:"payee_name" validate:"required,max=100"`
	BankAccountNo       string `json:"bank_account_no" validate:"required,numeric,min=9,max=18"`
	IFSC                string `json:"ifsc" validate:"required,len=11,alphanum"`
	BranchName          string `json:"branch_name" validate:"required,max=100"`
	BankName            string `json:"bank_name" validate:"required,max=100"`
	BankAccountType     string `json:"bank_account_type" validate:"required,oneof=Savings Current"`
	PAN                 string `json:"pan" validate:"required,len=10,alphanum"`
	Resume              string `json:"resume"`
	Aadhar              string `json:"aadhar"`
	Linkedin            string `json:"linkedin"`
	DeviceID            string `json:"device_id"`
	TopicIDs            string `json:"topic_ids" validate:"required"`
	LanguageIDs         string `json:"language_ids" validate:"required"`
	Certificate         string `json:"certificate" validate:"required"`
	CorporateTherapist  string `json:"corporate_therpist" validate:"required,oneof=0 1 2 3"`
	Timezone            string `json:"timezone" validate:"required"`
	Status              string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type UpdateTherapistProfileRequest struct {
	FirstName           string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName            string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Pronoun             string `json:"pronoun" validate:"omitempty,min=2,max=50"`
	Gender              string `json:"gender" validate:"omitempty,oneof=Male Female Other"`
	Location            string `json:"location" validate:"omitempty,max=200"`
	Phone               string `json:"phone" validate:"omitempty,len=12,numeric"`
	Photo               string `json:"photo" validate:"omitempty"`
	Email               string `json:"email" validate:"omitempty,email,max=100"`
	Price               string `json:"price" validate:"omitempty,numeric"`
	MultipleSessions    string `json:"multiple_sessions" validate:"omitempty,numeric"`
	Price3              string `json:"price_3" validate:"omitempty,numeric"`
	Price5              string `json:"price_5" validate:"omitempty,numeric"`
	Education           string `json:"education" validate:"omitempty,max=500"`
	Experience          string `json:"experience" validate:"omitempty,max=500"`
	About               string `json:"about" validate:"omitempty,max=5000"`
	TherapeuticApproach string `json:"therapeutic_approach" validate:"omitempty,max=5000"`
	StartDate           string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	GapYears            string `json:"gap_years" validate:"omitempty,numeric"`
	GapMonths           string `json:"gap_months" validate:"omitempty,numeric"`
	PayoutPercentage    string `json:"payout_percentage" validate:"omitempty,numeric"`
	PayeeName           string `json:"payee_name" validate:"omitempty,max=100"`
	BankAccountNo       string `json:"bank_account_no" validate:"omitempty,numeric,min=9,max=18"`
	IFSC                string `json:"ifsc" validate:"omitempty,len=11,alphanum"`
	BranchName          string `json:"branch_name" validate:"omitempty,max=100"`
	BankName            string `json:"bank_name" validate:"omitempty,max=100"`
	BankAccountType     string `json:"bank_account_type" validate:"omitempty,oneof=Savings Current"`
	PAN                 string `json:"pan" validate:"omitempty,len=10,alphanum"`
	Resume              string `json:"resume"`
	Aadhar              string `json:"aadhar"`
	Linkedin            string `json:"linkedin"`
	DeviceID            string `json:"device_id"`
	TopicIDs            string `json:"topic_ids" validate:"omitempty"`
	LanguageIDs         string `json:"language_ids" validate:"omitempty"`
	Certificate         string `json:"certificate" validate:"omitempty"`
	CorporateTherapist  string `json:"corporate_therpist" validate:"omitempty,oneof=0 1 2 3"`
	Timezone            string `json:"timezone" validate:"omitempty"`
	Status              string `json:"status" validate:"omitempty,oneof=0 1 2 3"`
}

type AddWebinarSessionInAdminPanel struct {
	CounsellorName          string `json:"counsellor_name" validate:"required,min=2,max=200"`
	Title                   string `json:"title" validate:"required,min=2,max=max=100"`
	CounsellorPhoto         string `json:"counsellor_photo" validate:"required"`
	CounsellorQualification string `json:"counsellor_qualification" validate:"required"`
	CounsellorAbout         string `json:"counsellor_about" validate:"required,min=10"`
	CounsellorExperience    string `json:"counsellor_experience" validate:"required"`
	CounsellorRating        string `json:"counsellor_rating" validate:"required"`
	About                   string `json:"about" validate:"required,max=2000"`
	PartnerName             string `json:"partner_name" validate:"required,min=2,max=100"`
	WhyAttend               string `json:"why_attend" validate:"required,min=100,max=2000"`
	Address                 string `json:"address" validate:"required,min=100,max=2000"`
	Photo                   string `json:"photo" validate:"required"`
	BackgroundPhoto         string `json:"background_photo" validate:"required"`
	Date                    string `json:"date" validate:"required,datetime=2006-01-02"`
	Time                    string `json:"time" validate:"required,numeric"`
	Duration                string `json:"duration" validate:"required,numeric"`
	Mode                    string `json:"mode" validate:"required"`
	Status                  string `json:"status" validate:"required,oneof=0 1 2 3 4 5"`
}

type BookAgainRemainingAppointmentRequestInB2C struct {
	AppointmentSlotID string `json:"appointment_slot_id" validate:"required,min=3,max=30"`
	Date              string `json:"date" validate:"required,datetime=2006-01-02"`
	Time              string `json:"time" validate:"required,numeric"`
}

type RescheduleAppointmentRequest struct {
	Date string `json:"date" validate:"required,datetime=2006-01-02"`
	Time string `json:"time" validate:"required,numeric"`
}

type AddAppointmentRatingRequest struct {
	ClientID      string `json:"client_id" validate:"required,min=3,max=20"`
	AppointmentID string `json:"appointment_id" validate:"required,min=3,max=22"`
	CounsellorID  string `json:"counsellor_id" validate:"required,min=3,max=22"`
	Rating        string `json:"rating" validate:"required,oneof=1 2 3 4 5"`
	RatingTypes   string `json:"rating_types" validate:"required,min=3,max=50"`
	RatingComment string `json:"rating_comment" validate:"required,min=1"`
}

type CancellationReasonAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" validate:"required,min=3"`
}

type AddNoSlotAppointmentRequest struct {
	ClientID     string `json:"client_id" validate:"required,min=3,max=20"`
	CounsellorID string `json:"counsellor_id" validate:"required,min=3,max=22"`
	Type         string `json:"type" validate:"required,oneof=1 2 3 4 5"`
}

type AddNoSlotInPersonAppointmentRequest struct {
	ClientID        string `json:"client_id" validate:"required,min=3,max=20"`
	CounsellorID    string `json:"counsellor_id" validate:"required,min=3,max=22"`
	CompanyName     string `json:"companyName" validate:"required,min=2,max=46"`
	CompanyLocation string `json:"companyLocation" validate:"required,min=2"`
	Type            string `json:"type" validate:"required,oneof=1 2 3 4 5"`
}

type OrderAppointmentForB2BRequest struct {
	ClientID     string `json:"client_id" validate:"required,min=3,max=20"`
	CounsellorID string `json:"listener_id" validate:"required,min=3,max=22"`
	Date         string `json:"date" validate:"required,datetime=2006-01-02"`
	Time         string `json:"time" validate:"required,numeric"`
}

type OrderConfirmAppointmentForB2BRequest struct {
	OrderID string `json:"order_id" validate:"required,min=3,max=20"`
}

type CounsellorOrderAppointmentForB2CRequest struct {
	ClientID     string `json:"client_id" validate:"required,min=3,max=45"`
	CounsellorID string `json:"counsellor_id" validate:"required,min=3,max=45"`
	Date         string `json:"date" validate:"required,datetime=2006-01-02"`
	Time         string `json:"time" validate:"required,numeric"`
	NoSession    string `json:"no_session" validate:"required,oneof=1 3 5"`
	CouponCode   string `json:"coupon_code"`
}

type CounsellorOrderConfirmAppointmentForB2CRequest struct {
	OrderID       string `json:"order_id" validate:"required,min=3,max=45"`
	PaymentMethod string `json:"payment_method" validate:"required,min=2,max=45"`
	PaymentID     string `json:"payment_id" validate:"required,min=3"`
}

type CancelInPersonEventForB2BRequest struct {
	UserID             string `json:"user_id" validate:"required,min=3,max=45"`
	OrderID            string `json:"order_id" validate:"required,min=3,max=45"`
	CancellationReason string `json:"cancellation_reason" validate:"required,min=3"`
}

type RequestInPersonEventForB2BRequest struct {
	ClientID string `json:"client_id" validate:"required,min=3,max=45"`
	OrderID  string `json:"order_id" validate:"required,min=3,max=45"`
}

type RateInPersonEventForB2BRequest struct {
	UserID    string `json:"user_id" validate:"required,min=3,max=45"`
	OrderID   string `json:"order_id" validate:"required,min=3,max=45"`
	Question1 string `json:"question1" validate:"required,min=3"`
	Question2 string `json:"question2" validate:"required,min=3"`
	Question3 string `json:"question3" validate:"required,min=3"`
	Question4 string `json:"question4" validate:"required,min=3"`
	Question5 string `json:"question5" validate:"required,min=3"`
}

type OrderEventForB2CRequest struct {
	UserID       string `json:"user_id" validate:"required,min=3,max=45"`
	EventOrderID string `json:"event_order_id" validate:"required,min=3,max=45"`
	CouponCode   string `json:"coupon_code"`
}

type OrderConfirmationEventForB2CRequest struct {
	OrderID       string `json:"order_id" validate:"required,min=3,max=45"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	PaymentID     string `json:"payment_id" validate:"required"`
}

type OrderInPersonEventForB2BRequest struct {
	UserID       string `json:"user_id" validate:"required,min=3,max=45"`
	EventOrderID string `json:"event_order_id" validate:"required,min=3,max=45"`
}

type OrderWebniarForB2BRequest struct {
	ClientID  string `json:"client_id" validate:"required,min=3,max=45"`
	WebinarID string `json:"webinar_id" validate:"required,min=3,max=45"`
}

type AddMoodInClientRequest struct {
	ClientID string `json:"client_id" validate:"required,min=3,max=45"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Age      string `json:"age" validate:"required,numeric"`
	Gender   string `json:"gender" validate:"required,oneof=Male Female Other"`
	Phone    string `json:"phone" validate:"required,len=12,numeric"`
	MoodID   string `json:"mood_id" validate:"required,numeric"`
	Notes    string `json:"notes" validate:"required,min=2,max=1000"`
	Date     string `json:"date" validate:"required,datetime=2006-01-02"`
}

type ClientProfileAddRequest struct {
	FirstName          string `json:"first_name" validate:"required,min=2,max=100"`
	LastName           string `json:"last_name" validate:"required,min=2,max=100"`
	Phone              string `json:"phone" validate:"required,len=12,numeric"`
	Email              string `json:"email" validate:"required,email"`
	DateOfBirth        string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Photo              string `json:"photo"`
	TopicIDs           string `json:"topic_ids" validate:"required"`
	Gender             string `json:"gender" validate:"required,oneof=Male Female Other"`
	Location           string `json:"location" validate:"required"`
	Timezone           string `json:"timezone" validate:"required"`
	DeviceID           string `json:"device_id" validate:"required"`
	Platform           string `json:"platform" validate:"required,oneof=ios android web"`
	Version            string `json:"version" validate:"required"`
	NotificationStatus string `json:"notification_status" validate:"required,oneof=0 1 2"`
}

type ClientB2BProfileAddRequest struct {
	EmpID              string `json:"emp_id" validate:"required"`
	FirstName          string `json:"first_name" validate:"required,min=2,max=100"`
	LastName           string `json:"last_name" validate:"required,min=2,max=100"`
	Phone              string `json:"phone" validate:"required,len=12,numeric"`
	Email              string `json:"email" validate:"required,email"`
	OTP                string `json:"otp" validate:"required,min=4,max=6"`
	DateOfBirth        string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Photo              string `json:"photo"`
	TopicIDs           string `json:"topic_ids" validate:"required"`
	Gender             string `json:"gender" validate:"required,oneof=Male Female Other"`
	Location           string `json:"location" validate:"required"`
	CorDarpartment     string `json:"cor_darpartment" validate:"required"`
	Timezone           string `json:"timezone" validate:"required"`
	DeviceID           string `json:"device_id" validate:"required"`
	Platform           string `json:"platform" validate:"required,oneof=ios android web"`
	Version            string `json:"version" validate:"required"`
	NotificationStatus string `json:"notification_status" validate:"required,oneof=0 1 2"`
}

type ClientB2BFamilyProfileAddRequest struct {
	ClientID           string `json:"client_id" validate:"required,min=4,max=45"`
	EmpID              string `json:"emp_id" validate:"required"`
	Relation           string `json:"relation" validate:"required"`
	FirstName          string `json:"first_name" validate:"required,min=2,max=100"`
	LastName           string `json:"last_name" validate:"required,min=2,max=100"`
	Phone              string `json:"phone" validate:"required,len=12,numeric"`
	Email              string `json:"email" validate:"required,email"`
	DateOfBirth        string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Photo              string `json:"photo"`
	TopicIDs           string `json:"topic_ids" validate:"required"`
	Gender             string `json:"gender" validate:"required,oneof=Male Female Other"`
	Location           string `json:"location" validate:"required"`
	CorDarpartment     string `json:"cor_darpartment" validate:"required"`
	Timezone           string `json:"timezone" validate:"required"`
	DeviceID           string `json:"device_id" validate:"required"`
	Platform           string `json:"platform" validate:"required,oneof=ios android web"`
	Version            string `json:"version" validate:"required"`
	NotificationStatus string `json:"notification_status" validate:"required,oneof=0 1 2"`
}

type ClientProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Photo       string `json:"photo"`
	TopicIDs    string `json:"topic_ids"`
	Gender      string `json:"gender" validate:"omitempty,oneof=Male Female Other"`
	DeviceID    string `json:"device_id"`
}

type TherapistOrderAppointmentForB2CRequest struct {
	ClientID    string `json:"client_id" validate:"required,min=3,max=45"`
	TherapistID string `json:"therapist_id" validate:"required,min=3,max=45"`
	Date        string `json:"date" validate:"required,datetime=2006-01-02"`
	Time        string `json:"time" validate:"required,numeric"`
	NoSession   string `json:"no_session" validate:"required,oneof=1 3 5"`
	CouponCode  string `json:"coupon_code"`
}

type TherapistOrderConfirmAppointmentForB2CRequest struct {
	OrderID       string `json:"order_id" validate:"required,min=3,max=45"`
	PaymentMethod string `json:"payment_method" validate:"required,min=2,max=45"`
	PaymentID     string `json:"payment_id" validate:"required,min=3"`
}

type SlotUpdateAvailabilityModel struct {
	Date  string
	Key   string
	Value string
}

type NewVersionOfCounsellorRecordFormRequest struct {
	RecordID                              string `json:"record_id" validate:"omitempty,min=4,max=45"`
	CounsellorID                          string `json:"counsellor_id" validate:"omitempty,min=4,max=45"`
	ClientID                              string `json:"client_id" validate:"omitempty,min=4,max=45"`
	AppointmentID                         string `json:"appointment_id" validate:"omitempty,min=4,max=45"`
	SessionFor                            string `json:"session_for" validate:"omitempty,oneof=Self Couple Family"`
	SessionType                           string `json:"session_type" validate:"required,oneof=Self Couple"`
	SessionMood                           string `json:"session_mood" validate:"omitempty,oneof=In-Person Virtual"`
	SessionDate                           string `json:"session_date" validate:"omitempty,datetime=2006-01-02"`
	FamilyRelation                        string `json:"family_relation" validate:"required"`
	InTime                                string `json:"in_time" validate:"omitempty,datetime=15:04"`
	OutTime                               string `json:"out_time" validate:"required,datetime=15:04"`
	NoShow                                string `json:"noshow" validate:"required,oneof=0 1"`
	IncompleteSession                     string `json:"incomplete_session" validate:"omitempty,oneof=0 1"`
	PresentingConcerns                    string `json:"presenting_concers" validate:"omitempty"`
	MentalHealth                          string `json:"mental_health" validate:"omitempty,numeric"`
	MentalHealthCheck                     string `json:"mental_health_check" validate:"omitempty"`
	DowngradingHighRiskCase               string `json:"downgrading_high_risk_case" validate:"omitempty"`
	IsClinicalPsychologistRequiredReason  string `json:"is_clinical_psychologist_required_reason" validate:"omitempty"`
	PsychiatricInterventionRequiredReason string `json:"psychiatric_intervention_required_reason" validate:"omitempty"`
	Category                              string `json:"category" validate:"omitempty"`
	SubCategory                           string `json:"sub_category" validate:"omitempty"`
	EmotionalState                        string `json:"emotional_state" validate:"omitempty"`
	TherapyNotes                          string `json:"therapy_notes" validate:"omitempty"`
	GoalsAchieved                         string `json:"goals_achieved" validate:"omitempty,oneof=No Yes"`
	GoalsAchievedReason                   string `json:"goals_achieved_reason" validate:"omitempty"`
	TotalSessionNeeded                    string `json:"total_session_needed" validate:"omitempty,numeric"`
	TakenSessions                         string `json:"taken_sessions" validate:"omitempty,numeric"`
	NextSessionPlan                       string `json:"next_session_plan" validate:"omitempty"`
	NextFollowUpDate                      string `json:"next_follow_up_date" validate:"omitempty,datetime=2006-01-02"`
	ClientNotes                           string `json:"client_notes" validate:"omitempty"`
	SelfWorkMaterial                      string `json:"self_work_material" validate:"omitempty"`
	Assessment                            string `json:"assessment" validate:"omitempty"`
	Links                                 string `json:"links" validate:"omitempty"`
}

type AppFeedbackRequest struct {
	UserID      string `json:"user_id" validate:"required,min=3,max=45"`
	Type        string `json:"type" validate:"required"`
	ConcernArea string `json:"concern_area" validate:"omitempty"`
	Details     string `json:"details" validate:"required"`
	Attach1     string `json:"attach_1" validate:"omitempty"`
	Attach2     string `json:"attach_2" validate:"omitempty"`
	Attach3     string `json:"attach_3" validate:"omitempty"`
}

type UpdateListenerProfileRequestInAdminPanel struct {
	FirstName  string `json:"first_name" validate:"required,min=2,max=50"`
	LastName   string `json:"last_name" validate:"required,min=2,max=50"`
	Phone      string `json:"phone" validate:"required,len=12,numeric"`
	Email      string `json:"email" validate:"required,email,max=100"`
	Gender     string `json:"gender" validate:"required,oneof=Male Female Other"`
	Occupation string `json:"occupation" validate:"required,min=2,max=100"`
	AgeGroup   string `json:"age_group" validate:"required,min=2,max=50"`
	About      string `json:"about" validate:"required,min=10,max=5000"`
	Status     string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type AddNotificationRequestInAdminPanel struct {
	Title            string `json:"title" validate:"required,min=2,max=50"`
	Body             string `json:"body" validate:"required,min=2,max=200"`
	UserIds          string `json:"user_ids" validate:"required,min=3,max=20"`
	Type             string `json:"type" validate:"required,oneof= 1 2 3 4 5 6"`
	NotificationType string `json:"notification_type" validate:"required,oneof=1 2 3 4 5 6"`
	UserType         string `json:"user_type" validate:"required,oneof=1 2 3 4 5"`
	PartnerName      string `json:"partner_name" validate:"omitempty,min=6,max=90"`
	PartnerLocation  string `json:"partner_location" validate:"omitempty,min=6,max=120"`
}

type AddUserProfileRequestInAdminPanel struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=2,max=200"`
	Type     string `json:"type" validate:"required,oneof= 1 2 3 4"`
}

type UpdateUserProfileRequestInAdminPanel struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=2,max=200"`
	Type     string `json:"type" validate:"required,oneof= 1 2 3 4"`
	Status   string `json:"status" validate:"required,oneof= 0 1 2 3 4"`
}

type AttachPermissionAddRequestInAdminPanel struct {
	RoleID      string `json:"role_id" validate:"required,min=2,max=25"`
	Username    string `json:"username" validate:"required,min=2,max=50"`
	Password    string `json:"password" validate:"required,min=2,max=200"`
	ProfileName string `json:"profile_name" validate:"required,min=3,max=50"`
}

type AttachPermissionUpdateRequestInAdminPanel struct {
	RoleID      string `json:"role_id" validate:"required,min=2,max=25"`
	Username    string `json:"username" validate:"required,min=2,max=50"`
	Password    string `json:"password" validate:"required,min=2,max=200"`
	ProfileName string `json:"profile_name" validate:"required,min=3,max=50"`
	Status      string `json:"status" validate:"required,oneof= 0 1 2 3 4"`
}

type CreateProfileForRoleRequestInAdminPanel struct {
	ProfileName string `json:"profile_name" validate:"required,min=2,max=100"`

	PCAdd  string `json:"pc_add" validate:"omitempty,oneof=check"`
	PCEdit string `json:"pc_edit" validate:"omitempty,oneof=check"`
	PCView string `json:"pc_view" validate:"omitempty,oneof=check"`

	AssessmentAdd  string `json:"assessment_add" validate:"omitempty,oneof=check"`
	AssessmentEdit string `json:"assessment_edit" validate:"omitempty,oneof=check"`
	AssessmentView string `json:"assessment_view" validate:"omitempty,oneof=check"`

	HomeAdd  string `json:"home_add" validate:"omitempty,oneof=check"`
	HomeEdit string `json:"home_edit" validate:"omitempty,oneof=check"`
	HomeView string `json:"home_view" validate:"omitempty,oneof=check"`

	SlotAdd  string `json:"slot_add" validate:"omitempty,oneof=check"`
	SlotEdit string `json:"slot_edit" validate:"omitempty,oneof=check"`
	SlotView string `json:"slot_view" validate:"omitempty,oneof=check"`

	InpersonCafeAdd  string `json:"inperson_cafe_add" validate:"omitempty,oneof=check"`
	InpersonCafeEdit string `json:"inperson_cafe_edit" validate:"omitempty,oneof=check"`
	InpersonCafeView string `json:"inperson_cafe_view" validate:"omitempty,oneof=check"`

	LinkAdd  string `json:"link_add" validate:"omitempty,oneof=check"`
	LinkEdit string `json:"link_edit" validate:"omitempty,oneof=check"`
	LinkView string `json:"link_view" validate:"omitempty,oneof=check"`

	NotiAdd  string `json:"noti_add" validate:"omitempty,oneof=check"`
	NotiEdit string `json:"noti_edit" validate:"omitempty,oneof=check"`
	NotiView string `json:"noti_view" validate:"omitempty,oneof=check"`

	ContAdd  string `json:"cont_add" validate:"omitempty,oneof=check"`
	ContEdit string `json:"cont_edit" validate:"omitempty,oneof=check"`
	ContView string `json:"cont_view" validate:"omitempty,oneof=check"`

	MQAdd  string `json:"mq_add" validate:"omitempty,oneof=check"`
	MQEdit string `json:"mq_edit" validate:"omitempty,oneof=check"`
	MQView string `json:"mq_view" validate:"omitempty,oneof=check"`

	CentAdd  string `json:"cent_add" validate:"omitempty,oneof=check"`
	CentEdit string `json:"cent_edit" validate:"omitempty,oneof=check"`
	CentView string `json:"cent_view" validate:"omitempty,oneof=check"`

	CounAdd  string `json:"coun_add" validate:"omitempty,oneof=check"`
	CounEdit string `json:"coun_edit" validate:"omitempty,oneof=check"`
	CounView string `json:"coun_view" validate:"omitempty,oneof=check"`

	PartAdd  string `json:"part_add" validate:"omitempty,oneof=check"`
	PartEdit string `json:"part_edit" validate:"omitempty,oneof=check"`
	PartView string `json:"part_view" validate:"omitempty,oneof=check"`

	PartLocAdd  string `json:"part_loc_add" validate:"omitempty,oneof=check"`
	PartLocEdit string `json:"part_loc_edit" validate:"omitempty,oneof=check"`
	PartLocView string `json:"part_loc_view" validate:"omitempty,oneof=check"`

	ListAdd  string `json:"list_add" validate:"omitempty,oneof=check"`
	ListEdit string `json:"list_edit" validate:"omitempty,oneof=check"`
	ListView string `json:"list_view" validate:"omitempty,oneof=check"`

	TherAdd  string `json:"ther_add" validate:"omitempty,oneof=check"`
	TherEdit string `json:"ther_edit" validate:"omitempty,oneof=check"`
	TherView string `json:"ther_view" validate:"omitempty,oneof=check"`

	AppointAdd  string `json:"appoint_add" validate:"omitempty,oneof=check"`
	AppointEdit string `json:"appoint_edit" validate:"omitempty,oneof=check"`
	AppointView string `json:"appoint_view" validate:"omitempty,oneof=check"`

	CafeAdd  string `json:"cafe_add" validate:"omitempty,oneof=check"`
	CafeEdit string `json:"cafe_edit" validate:"omitempty,oneof=check"`
	CafeView string `json:"cafe_view" validate:"omitempty,oneof=check"`

	ReptAdd  string `json:"rept_add" validate:"omitempty,oneof=check"`
	ReptEdit string `json:"rept_edit" validate:"omitempty,oneof=check"`
	ReptView string `json:"rept_view" validate:"omitempty,oneof=check"`
}

type CouponAddRequestInAdminPanel struct {
	CouponCode  string `json:"coupon_code" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,max=500"`

	ClientID     string `json:"client_id"`
	CounsellorID string `json:"counsellor_id"`
	TherapistID  string `json:"therapist_id"`

	Discount             string `json:"discount" validate:"required,numeric"`
	MinimumOrderValue    string `json:"minimum_order_value" validate:"required,numeric"`
	MaximumDiscountValue string `json:"maximum_discount_value" validate:"required,numeric"`
	ValidForOrder        string `json:"valid_for_order" validate:"required,numeric"`

	Type      string `json:"type" validate:"required,oneof=1 2"`
	OrderType string `json:"order_type" validate:"required,oneof=0 1 2"`

	StartBy string `json:"start_by" validate:"required,datetime=2006-01-02 15:04:05"`
	EndBy   string `json:"end_by" validate:"required,datetime=2006-01-02 15:04:05"`
}

type CouponUpdateRequestInAdminPanel struct {
	CouponCode  string `json:"coupon_code" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,max=500"`

	ClientID     string `json:"client_id"`
	CounsellorID string `json:"counsellor_id"`
	TherapistID  string `json:"therapist_id"`

	Discount             string `json:"discount" validate:"required,numeric"`
	MinimumOrderValue    string `json:"minimum_order_value" validate:"required,numeric"`
	MaximumDiscountValue string `json:"maximum_discount_value" validate:"required,numeric"`
	ValidForOrder        string `json:"valid_for_order" validate:"required,numeric"`

	Type      string `json:"type" validate:"required,oneof=1 2"`
	OrderType string `json:"order_type" validate:"required,oneof=0 1 2"`

	StartBy string `json:"start_by" validate:"required,datetime=2006-01-02 15:04:05"`
	EndBy   string `json:"end_by" validate:"required,datetime=2006-01-02 15:04:05"`
	Status  string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type EventUpdateRequestInAdminPanel struct {
	CounsellorID string `json:"counsellor_id" validate:"required,min=4,max=25"`
	Title        string `json:"title" validate:"required,min=3,max=100"`
	Description  string `json:"description" validate:"required,min=10,max=1000"`
	Photo        string `json:"photo" validate:"required"`
	Date         string `json:"date" validate:"required,datetime=2006-01-02"`
	Time         string `json:"time" validate:"required,numeric"`
	Duration     string `json:"duration" validate:"required"` // in minutes
	Price        string `json:"price" validate:"required,numeric"`
	Status       string `json:"status" validate:"required,oneof=0 1 2 3 4"`
}

type InPersonEventAddRequestInAdminPanel struct {
	CounsellorID    string `json:"counsellor_id" validate:"required"`
	Title           string `json:"title" validate:"required,min=3,max=100"`
	Description     string `json:"description" validate:"required,min=10,max=1000"`
	TotalSeat       string `json:"total_seat" validate:"required,numeric"`
	Document        string `json:"document"`
	CarryThings     string `json:"carry_things" validate:"required,max=500"`
	CompanyName     string `json:"company_name" validate:"required,min=2,max=100"`
	CompanyLocation string `json:"company_location" validate:"required,min=2,max=255"`
	Address         string `json:"address" validate:"required,min=5,max=500"`
	Photo           string `json:"photo" validate:"required"`
	BackgroundPhoto string `json:"background_photo" validate:"required"`
	Date            string `json:"date" validate:"required,datetime=2006-01-02"`
	Time            string `json:"time" validate:"required,required"`
	Duration        string `json:"duration" validate:"required,numeric"`
	Status          string `json:"status" validate:"required,oneof=0 1 2 3 4"`
}

type AvailabilityUpdateRequestInAdminPanel struct {
	ID              string `json:"id"`
	CounsellorID    string `json:"counsellor_id" validate:"required,min=3,max=16"`
	Date            string `json:"date" validate:"required,datetime=2006-01-02"`
	FromTime        string `json:"fromTime" validate:"required,datetime=15:04"`
	ToTime          string `json:"toTime" validate:"required,datetime=15:04"`
	CompanyName     string `json:"companyName" validate:"required,min=3,max=100"`
	CompanyLocation string `json:"companyLocation" validate:"required,min=3,max=100"`
	RoomNo          string `json:"roomNo" validate:"required,min=1,max=80"`
	Address         string `json:"address" validate:"required,min=3,max=200"`
	Status          string `json:"status" validate:"required,oneof=0 1 2"`
	Zero            string `json:"0" validate:"required,oneof=0 1"`
	One             string `json:"1" validate:"required,oneof=0 1"`
	Two             string `json:"2" validate:"required,oneof=0 1"`
	Three           string `json:"3" validate:"required,oneof=0 1"`
	Four            string `json:"4" validate:"required,oneof=0 1"`
	Five            string `json:"5" validate:"required,oneof=0 1"`
	Six             string `json:"6" validate:"required,oneof=0 1"`
	Seven           string `json:"7" validate:"required,oneof=0 1"`
	Eight           string `json:"8" validate:"required,oneof=0 1"`
	Nine            string `json:"9" validate:"required,oneof=0 1"`
	Ten             string `json:"10" validate:"required,oneof=0 1"`
	Eleven          string `json:"11" validate:"required,oneof=0 1"`
	Twelve          string `json:"12" validate:"required,oneof=0 1"`
	Thirteen        string `json:"13" validate:"required,oneof=0 1"`
	Fourteen        string `json:"14" validate:"required,oneof=0 1"`
	Fifteen         string `json:"15" validate:"required,oneof=0 1"`
	Sixteen         string `json:"16" validate:"required,oneof=0 1"`
	Seventeen       string `json:"17" validate:"required,oneof=0 1"`
	Eighteen        string `json:"18" validate:"required,oneof=0 1"`
	Nineteen        string `json:"19" validate:"required,oneof=0 1"`
	Twenty          string `json:"20" validate:"required,oneof=0 1"`
	TwentyOne       string `json:"21" validate:"required,oneof=0 1"`
	TwentyTwo       string `json:"22" validate:"required,oneof=0 1"`
	TwentyThree     string `json:"23" validate:"required,oneof=0 1"`
	TwentyFour      string `json:"24" validate:"required,oneof=0 1"`
	TwentyFive      string `json:"25" validate:"required,oneof=0 1"`
	TwentySix       string `json:"26" validate:"required,oneof=0 1"`
	TwentySeven     string `json:"27" validate:"required,oneof=0 1"`
	TwentyEight     string `json:"28" validate:"required,oneof=0 1"`
	TwentyNine      string `json:"29" validate:"required,oneof=0 1"`
	Thirty          string `json:"30" validate:"required,oneof=0 1"`
	ThirtyOne       string `json:"31" validate:"required,oneof=0 1"`
	ThirtyTwo       string `json:"32" validate:"required,oneof=0 1"`
	ThirtyThree     string `json:"33" validate:"required,oneof=0 1"`
	ThirtyFour      string `json:"34" validate:"required,oneof=0 1"`
	ThirtyFive      string `json:"35" validate:"required,oneof=0 1"`
	ThirtySix       string `json:"36" validate:"required,oneof=0 1"`
	ThirtySeven     string `json:"37" validate:"required,oneof=0 1"`
	ThirtyEight     string `json:"38" validate:"required,oneof=0 1"`
	ThirtyNine      string `json:"39" validate:"required,oneof=0 1"`
	Forty           string `json:"40" validate:"required,oneof=0 1"`
	FortyOne        string `json:"41" validate:"required,oneof=0 1"`
	FortyTwo        string `json:"42" validate:"required,oneof=0 1"`
	FortyThree      string `json:"43" validate:"required,oneof=0 1"`
	FortyFour       string `json:"44" validate:"required,oneof=0 1"`
	FortyFive       string `json:"45" validate:"required,oneof=0 1"`
	FortySix        string `json:"46" validate:"required,oneof=0 1"`
	FortySeven      string `json:"47" validate:"required,oneof=0 1"`
}

type AvailabilityUpdateTherapistRequest struct {
	ID                 string `json:"id"`
	CounsellorID       string `json:"counsellor_id" validate:"required,min=3,max=16"`
	WeekDay            string `json:"weekday" validate:"omitempty,oneof=0 1 2 3 4 5 6 7"`
	Format             string `json:"format" validate:"required"`
	Dates              string `json:"dates" validate:"omitempty"`
	AvailabilityStatus string `json:"availability_status" validate:"required,oneof=0 1 2"`
	Break              string `json:"break" validate:"required,oneof=0 1 2"`
	Status             string `json:"status" validate:"required,oneof=0 1 2"`
	Zero               string `json:"0" validate:"required,oneof=0 1"`
	One                string `json:"1" validate:"required,oneof=0 1"`
	Two                string `json:"2" validate:"required,oneof=0 1"`
	Three              string `json:"3" validate:"required,oneof=0 1"`
	Four               string `json:"4" validate:"required,oneof=0 1"`
	Five               string `json:"5" validate:"required,oneof=0 1"`
	Six                string `json:"6" validate:"required,oneof=0 1"`
	Seven              string `json:"7" validate:"required,oneof=0 1"`
	Eight              string `json:"8" validate:"required,oneof=0 1"`
	Nine               string `json:"9" validate:"required,oneof=0 1"`
	Ten                string `json:"10" validate:"required,oneof=0 1"`
	Eleven             string `json:"11" validate:"required,oneof=0 1"`
	Twelve             string `json:"12" validate:"required,oneof=0 1"`
	Thirteen           string `json:"13" validate:"required,oneof=0 1"`
	Fourteen           string `json:"14" validate:"required,oneof=0 1"`
	Fifteen            string `json:"15" validate:"required,oneof=0 1"`
	Sixteen            string `json:"16" validate:"required,oneof=0 1"`
	Seventeen          string `json:"17" validate:"required,oneof=0 1"`
	Eighteen           string `json:"18" validate:"required,oneof=0 1"`
	Nineteen           string `json:"19" validate:"required,oneof=0 1"`
	Twenty             string `json:"20" validate:"required,oneof=0 1"`
	TwentyOne          string `json:"21" validate:"required,oneof=0 1"`
	TwentyTwo          string `json:"22" validate:"required,oneof=0 1"`
	TwentyThree        string `json:"23" validate:"required,oneof=0 1"`
	TwentyFour         string `json:"24" validate:"required,oneof=0 1"`
	TwentyFive         string `json:"25" validate:"required,oneof=0 1"`
	TwentySix          string `json:"26" validate:"required,oneof=0 1"`
	TwentySeven        string `json:"27" validate:"required,oneof=0 1"`
	TwentyEight        string `json:"28" validate:"required,oneof=0 1"`
	TwentyNine         string `json:"29" validate:"required,oneof=0 1"`
	Thirty             string `json:"30" validate:"required,oneof=0 1"`
	ThirtyOne          string `json:"31" validate:"required,oneof=0 1"`
	ThirtyTwo          string `json:"32" validate:"required,oneof=0 1"`
	ThirtyThree        string `json:"33" validate:"required,oneof=0 1"`
	ThirtyFour         string `json:"34" validate:"required,oneof=0 1"`
	ThirtyFive         string `json:"35" validate:"required,oneof=0 1"`
	ThirtySix          string `json:"36" validate:"required,oneof=0 1"`
	ThirtySeven        string `json:"37" validate:"required,oneof=0 1"`
	ThirtyEight        string `json:"38" validate:"required,oneof=0 1"`
	ThirtyNine         string `json:"39" validate:"required,oneof=0 1"`
	Forty              string `json:"40" validate:"required,oneof=0 1"`
	FortyOne           string `json:"41" validate:"required,oneof=0 1"`
	FortyTwo           string `json:"42" validate:"required,oneof=0 1"`
	FortyThree         string `json:"43" validate:"required,oneof=0 1"`
	FortyFour          string `json:"44" validate:"required,oneof=0 1"`
	FortyFive          string `json:"45" validate:"required,oneof=0 1"`
	FortySix           string `json:"46" validate:"required,oneof=0 1"`
	FortySeven         string `json:"47" validate:"required,oneof=0 1"`
}

type CreateOrderEvent struct {
	UserID       string `json:"user_id" validate:"required"`
	EventOrderID string `json:"event_order_id" validate:"required"`
	CouponCode   string `json:"coupon_code"`
}

type OrderPaymentCompleteEvent struct {
	OrderID       string `json:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	PaymentID     string `json:"payment_id" validate:"required"`
}

type AssessmentOption struct {
	Option string `json:"option" validate:"required"`
	Score  string `json:"score" validate:"required,numeric"`
	Order  string `json:"order" validate:"required,numeric"`
	Status string `json:"status" validate:"required,oneof=0 1"`
}

type AssessmentQuestion struct {
	Question string             `json:"question" validate:"required"`
	Order    string             `json:"order" validate:"required,numeric"`
	Status   string             `json:"status" validate:"required,oneof=0 1"`
	Options  []AssessmentOption `json:"options" validate:"required, min=1,dive"`
}

type AssessmentScore struct {
	MinScore string `json:"min" validate:"required,numeric"`
	MaxScore string `json:"max" validate:"required,numeric"`
	Result   string `json:"result" validate:"required"`
}

// AssessmentAddRequest .
type AssessmentAddRequestInAdminPanel struct {
	Title       string               `json:"title" validate:"required,min=3,max=100"`
	SubTitles   string               `json:"subtitles" validate:"required,min=20,max=500"`
	Photo       string               `json:"photo" validate:"required"`
	Duration    string               `json:"duration" validate:"required,min=1,max=15"`
	Type        string               `json:"type" validate:"required,oneof=1 2"`
	Instruction string               `json:"instruction" validate:"required,min=10"`
	Source      string               `json:"source" validate:"required"`
	Reference   string               `json:"reference" validate:"required"`
	Order       string               `json:"order" validate:"required,numeric"`
	Status      string               `json:"status" validate:"required,oneof=0 1"`
	Questions   []AssessmentQuestion `json:"questions" validate:"required,min=1,dive"`
	Scores      []AssessmentScore    `json:"scores" validate:"required,min=1,dive"`
}

type AssessmentUpdateOption struct {
	AssessmentQuestionOptionID string `json:"assessment_question_option_id" validate:"required" min:"5" max:"25"`
	Option                     string `json:"option" validate:"required" min:"3"`
	Score                      string `json:"score" validate:"required,numeric"`
	Order                      string `json:"order" validate:"required,numeric"`
	Status                     string `json:"status" validate:"required,oneof=0 1"`
}

type AssessmentUpdateQuestion struct {
	AssessmentQuestionID string                   `json:"assessment_question_id" validate:"required" min:"5" max:"25"`
	Question             string                   `json:"question" validate:"required" min:"3"`
	Order                string                   `json:"order" validate:"required,numeric"`
	Status               string                   `json:"status" validate:"required,oneof=0 1"`
	Options              []AssessmentUpdateOption `json:"options" validate:"required,min=1,dive"`
}

type AssessmentUpdateScore struct {
	ScoreID  string `json:"id" validate:"required" min:"5" max:"25"`
	MinScore string `json:"min" validate:"required,numeric"`
	MaxScore string `json:"max" validate:"required,numeric"`
	Result   string `json:"result" validate:"required"`
}

// AssessmentUpdateRequest .
type AssessmentUpdateRequestInAdminPanel struct {
	AssessmentID string                     `json:"assessment_id" validate:"required,min=5,max=18"`
	Title        string                     `json:"title" validate:"required,min=3,max=100"`
	SubTitles    string                     `json:"subtitles" validate:"required,min=20,max=500"`
	Photo        string                     `json:"photo" validate:"required"`
	Duration     string                     `json:"duration" validate:"required,min=1,max=15"`
	Type         string                     `json:"type" validate:"required,oneof=1 2"`
	Instruction  string                     `json:"instruction" validate:"required,min=10"`
	Source       string                     `json:"source" validate:"required"`
	Reference    string                     `json:"reference" validate:"required"`
	Order        string                     `json:"order" validate:"required,numeric"`
	Status       string                     `json:"status" validate:"required,oneof=0 1"`
	Questions    []AssessmentUpdateQuestion `json:"questions" validate:"required,min=1,dive"`
	Scores       []AssessmentUpdateScore    `json:"scores" validate:"required,min=1,dive"`
}

type InPersonCounsellorConnectWithCorporateAddRequest struct {
	CounsellorID    string `json:"counsellor_id" validate:"required,min=5,max=16"`
	PartnerName     string `json:"partner_name" validate:"required"`
	PartnerLocation string `json:"partner_location" validate:"required"`
}

type InPersonCounsellorConnectWithCorporateUpdateRequest struct {
	CounsellorID    string `json:"counsellor_id" validate:"required,min=5,max=16"`
	PartnerName     string `json:"partner_name" validate:"required"`
	PartnerLocation string `json:"partner_location" validate:"required"`
	Status          string `json:"status" validate:"required,oneof=0 1 2 3"`
}

type CafeAttendedAddRequest struct {
	OrderID   string `json:"order_id" validate:"required,min=9,max=25"`
	ClientIDs []struct {
		UserID string `json:"user_id" validate:"required"`
	} `json:"clientids" validate:"required,dive"`
}

// MoodAddRequest .
type MoodAddRequest struct {
	ClientID string `json:"client_id"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	MoodID   string `json:"mood_id"`
	Date     string `json:"date"`
	Notes    string `json:"notes"`
}

type EmailDataForEvent struct {
	First_Name  string
	Last_Name   string
	Type        string
	Title       string
	Description string
	Photo       string
	Topic_Name  string
	Date        string
	Time        string
	Duration    string
	Price       string
}

type EmailDataForCounsellorProfile struct {
	Media_URL            string
	First_Name           string
	Last_Name            string
	Pronoun              string
	Gender               string
	Location             string
	Type                 string
	Phone                string
	Email                string
	Photo                string
	Education            string
	CounsellingStartDate string
	CounsellingGap       string
	Experience           string
	TherapeuticApproach  string
	About                string
	Resume               string
	Certificate          string
	Aadhar               string
	Linkedin             string
	Status               string
}

type EmailDataForWebClientB2CProfile struct {
	Name            string
	Phone           string
	Email           string
	CompanyName     string
	CompanyLocation string
	CompanySize     string
	Message         string
}

type EmailDataForCounsellorCancellation struct {
	First_Name            string
	Last_Name             string
	Previous_Date         string
	Previous_Client_Name  string
	Previous_Client_Email string
	Latest_Date           string
	Lastest_Client_Name   string
	Lastest_Client_Email  string
}

type SlotUpdateModelInTherapistAvailability struct {
	Date  string
	Key   string
	Value string
}

type AppSummaryReport struct {
	AppointmentTotal              string `json:"appointment_total"`
	ClientTotal                   string `json:"client_total"`
	AppointmentsInPersonTotal     string `json:"appointments_inperson_total"`
	EmeCaseVirtualTotal           string `json:"emecase_virtual_total"`
	EmeCaseInPersonTotal          string `json:"emecase_inperson_total"`
	ContentsTotal                 string `json:"contents_total"`
	MoodsTotal                    string `json:"moods_total"`
	AssessmentsTotal              string `json:"assessments_total"`
	TotalRatingTotal              string `json:"total_rating_total"`
	AvgRatingTotal                string `json:"avgrating_total"`
	AppointmentsCancellationTotal string `json:"appointments_cancellation_total"`
	AppointmentsNoShowTotal       string `json:"appointments_noshow_total"`
}

type PaymentRequest struct {
	MerchantKey string `json:"merchant_key"`
	Amount      string `json:"amount"`
	OrderId     string `json:"order_id"`
	ProductInfo string `json:"product_info"`
	FirstName   string `json:"first_name"`
	Email       string `json:"email"`
}

type PaymentVerify struct {
	Status             int                          `json:"status"`
	Msg                string                       `json:"msg"`
	TransactionDetails map[string]map[string]string `json:"transaction_details"`
}

type EmailDataForCounsellorRecord struct {
	TherapistName                 string
	SessionFor                    string
	First_Name                    string
	Last_Name                     string
	Gender                        string
	Age                           string
	NoShow                        string
	PresentingConcerns            string
	PsychiatricIntervention       string
	PsychiatricInterventionReason string
	TherapyNotes                  string
	SubCategory                   string
	EmotionalState                string
	NextFollowDate                string
	SessionMode                   string
	SessionDate                   string
	InTime                        string
	OutTime                       string
	MentalHealth                  string
	TherapeuticGoal               string
	TherapyPlan                   string
	AssessmentTool                string
	ClientNotes                   string
	ClientAttach                  string
	SendingStatus                 string
}

type EmailDataForCounsellorRecordForLastestVersion struct {
	TherapistName                         string
	Client_First_Name                     string
	Client_Last_Name                      string
	Client_Gender                         string
	Client_Age                            string
	SessionFor                            string
	SessionType                           string
	FamilyRelation                        string
	SessionMode                           string
	SessionDate                           string
	InTime                                string
	OutTime                               string
	NoShow                                string
	PresentingConcerns                    string
	MentalHealthScale                     string
	MentalHealthCheck                     string
	DowngradingHighRiskCase               string
	IsClinicalPsychologistRequiredReason  string
	PsychiatricInterventionRequiredReason string
	Category                              string
	SubCategory                           string
	EmotionalState                        string
	TotalSessionNeeded                    string
	TakenSessions                         string
	TherapyNotes                          string
	NextSessionPlan                       string
	GoalsAchieved                         string
	GoalsAchievedReason                   string
	NextFollowDate                        string
	ClientNotes                           string
	Assessment                            string
	SelfWorkMaterial                      string
}

type EmailDataForFeedback struct {
	ClientFirstName string
	ClientLastName  string
	ClientEmail     string
	ConcernsType    string
	ConcernArea     string
	Details         string
	Attach1         string
	Attach2         string
	Attach3         string
	MediaURL        string
}

type EmailDataForCounsellorVisit struct {
	Client_Name     string
	Client_Location string
	InTime          string
	OutTime         string
}

type EmailDataForPaymentReceipt struct {
	Date        string
	ReceiptNo   string
	ReferenceNo string
	SPrice      string
	Qty         string
	Total       string
	TPrice      string
	CouponC     string
	Discount    string
	TotalP      string
}

type CancellationUpdateRequest struct {
	CancellationReason string `json:"cancellation_reason"`
}

type CounsellorCommentRequest struct {
	CommentForClient string `json:"commentforclient"`
	Attachment       string `json:"attachments"`
}

type AssessmentDownloadAIS struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Score    string `json:"score"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Answer4  string `json:"answer4"`
	Answer5  string `json:"answer5"`
	Answer6  string `json:"answer6"`
	Answer7  string `json:"answer7"`
	Answer8  string `json:"answer8"`
	Answer9  string `json:"answer9"`
	Answer10 string `json:"answer10"`
}

type AssessmentDownloadGAD7Model struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Age     string `json:"age"`
	Gender  string `json:"gender"`
	Score   string `json:"score"`
	Answer1 string `json:"answer1"`
	Answer2 string `json:"answer2"`
	Answer3 string `json:"answer3"`
	Answer4 string `json:"answer4"`
	Answer5 string `json:"answer5"`
	Answer6 string `json:"answer6"`
	Answer7 string `json:"answer7"`
	Answer8 string `json:"answer8"`
}

type AssessmentDownloadGWBModel struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Age     string `json:"age"`
	Gender  string `json:"gender"`
	Score   string `json:"score"`
	Answer1 string `json:"answer1"`
	Answer2 string `json:"answer2"`
	Answer3 string `json:"answer3"`
	Answer4 string `json:"answer4"`
	Answer5 string `json:"answer5"`
	Answer6 string `json:"answer6"`
}

type AssessmentDownloadSRSModel struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Score    string `json:"score"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Answer4  string `json:"answer4"`
	Answer5  string `json:"answer5"`
	Answer6  string `json:"answer6"`
	Answer7  string `json:"answer7"`
	Answer8  string `json:"answer8"`
	Answer9  string `json:"answer9"`
	Answer10 string `json:"answer10"`
	Answer11 string `json:"answer11"`
	Answer12 string `json:"answer12"`
	Answer13 string `json:"answer13"`
	Answer14 string `json:"answer14"`
	Answer15 string `json:"answer15"`
	Answer16 string `json:"answer16"`
	Answer17 string `json:"answer17"`
	Answer18 string `json:"answer18"`
	Answer19 string `json:"answer19"`
	Answer20 string `json:"answer20"`
	Answer21 string `json:"answer21"`
}

type AssessmentDownloadBurnOutModel struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Score    string `json:"score"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Answer4  string `json:"answer4"`
	Answer5  string `json:"answer5"`
	Answer6  string `json:"answer6"`
	Answer7  string `json:"answer7"`
	Answer8  string `json:"answer8"`
	Answer9  string `json:"answer9"`
	Answer10 string `json:"answer10"`
	Answer11 string `json:"answer11"`
	Answer12 string `json:"answer12"`
}

type AssessmentDownloadSelfEsteemModel struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Score    string `json:"score"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Answer4  string `json:"answer4"`
	Answer5  string `json:"answer5"`
	Answer6  string `json:"answer6"`
	Answer7  string `json:"answer7"`
	Answer8  string `json:"answer8"`
	Answer9  string `json:"answer9"`
	Answer10 string `json:"answer10"`
}

type AssessmentDownloadPSYCHOLOGICALWELLBEINGModel struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Score    string `json:"score"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Answer4  string `json:"answer4"`
	Answer5  string `json:"answer5"`
	Answer6  string `json:"answer6"`
	Answer7  string `json:"answer7"`
	Answer8  string `json:"answer8"`
	Answer9  string `json:"answer9"`
	Answer10 string `json:"answer10"`
	Answer11 string `json:"answer11"`
	Answer12 string `json:"answer12"`
	Answer13 string `json:"answer13"`
	Answer14 string `json:"answer14"`
	Answer15 string `json:"answer15"`
	Answer16 string `json:"answer16"`
	Answer17 string `json:"answer17"`
	Answer18 string `json:"answer18"`
}

type AssessmentDownloadBDIModel struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Age        string `json:"age"`
	Gender     string `json:"gender"`
	Score      string `json:"score"`
	Answer1    string `json:"answer1"`
	Answer2    string `json:"answer2"`
	Answer3    string `json:"answer3"`
	Answer4    string `json:"answer4"`
	Answer5    string `json:"answer5"`
	Answer6    string `json:"answer6"`
	Answer7    string `json:"answer7"`
	Answer8    string `json:"answer8"`
	Answer9    string `json:"answer9"`
	Answer10   string `json:"answer10"`
	Answer11   string `json:"answer11"`
	Answer12   string `json:"answer12"`
	Answer13   string `json:"answer13"`
	Answer14   string `json:"answer14"`
	Answer15   string `json:"answer15"`
	Answer16   string `json:"answer16"`
	Answer17   string `json:"answer17"`
	Answer18   string `json:"answer18"`
	Answer19   string `json:"answer19"`
	Answer20   string `json:"answer20"`
	Answer21   string `json:"answer21"`
	Response1  string `json:"response1"`
	Response2  string `json:"response2"`
	Response3  string `json:"response3"`
	Response4  string `json:"response4"`
	Response5  string `json:"response5"`
	Response6  string `json:"response6"`
	Response7  string `json:"response7"`
	Response8  string `json:"response8"`
	Response9  string `json:"response9"`
	Response10 string `json:"response10"`
	Response11 string `json:"response11"`
	Response12 string `json:"response12"`
	Response13 string `json:"response13"`
	Response14 string `json:"response14"`
	Response15 string `json:"response15"`
	Response16 string `json:"response16"`
	Response17 string `json:"response17"`
	Response18 string `json:"response18"`
	Response19 string `json:"response19"`
	Response20 string `json:"response20"`
	Response21 string `json:"response21"`
}

type TokenRequest struct {
	UserUUID string `json:"user_uuid,omitempty"` // Optional for App Token
	Expire   uint32 `json:"expire"`
}

type TokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

type ClientAppointmentConfirmation struct {
	First_Name      string `json:"first_name"`
	Counsellor_Name string `json:"counsellor_name"`
	Date_Time       string `json:"date_time"`
}

type ClientRequestS struct {
	Region              string `json:"region"`
	ResourceExpiredHour int    `json:"resourceExpiredHour"`
	Scene               int    `json:"scene"`
}

type PostRequestForAgora struct {
	CName         string         `json:"cname"`
	Uid           string         `json:"uid"`
	ClientRequest ClientRequestS `json:"clientRequest"`
}

type PostRequestForHTMLToPDF struct {
	HtmlContent string `json:"html_content"`
}

type RecordingConfigModel struct {
	MaxIdleTime int `json:"maxIdleTime"`
	// StreamMode         string `json:"streamMode"`
	StreamTypes int `json:"streamTypes"`
	// AudioProfile       int `json:"audioProfile"`
	ChannelType        int `json:"channelType"`
	VideoStreamType    int `json:"videoStreamType"`
	TranscodingConfigs TranscodingConfig
	// SubscribeVideoUids []string `json:"subscribeVideoUids"`
	// SubscribeAudioUids []string `json:"subscribeAudioUids"`
	// SubscribeUidGroup  int      `json:"subscribeUidGroup"`
}

type TranscodingConfig struct {
	Height           int `json:"height"`
	Width            int `json:"width"`
	Bitrate          int `json:"bitrate"`
	Fps              int `json:"fps"`
	MixedVideoLayout int `json:"mixedVideoLayout"`
}

type Tags struct {
	Security string `json:"security"`
}

type ExtensionParamsModel struct {
	Tag string `json:"tag"`
}

type StorageConfigModel struct {
	AccessKey       string               `json:"accessKey"`
	Bucket          string               `json:"bucket"`
	SecretKey       string               `json:"secretKey"`
	Vendor          int                  `json:"vendor"`
	Region          int                  `json:"region"`
	FileNamePrefix  []string             `json:"fileNamePrefix"`
	ExtensionParams ExtensionParamsModel `json:"extensionParams"`
}

type RecordingFileConfigModel struct {
	AvFileType []string `json:"avFileType"`
}

type ClientRequestForStartCall struct {
	Token               string                   `json:"token"`
	RecordingConfig     RecordingConfigModel     `json:"recordingConfig"`
	RecordingFileConfig RecordingFileConfigModel `json:"recordingFileConfig"`
	StorageConfig       StorageConfigModel       `json:"storageConfig"`
}

type ClientRequestForUpdateStartCall struct {
	UID           string        `json:"uid"`
	Cname         string        `json:"cname"`
	ClientRequest ClientRequest `json:"clientRequest"`
}
type AudioUIDList struct {
	SubscribeAudioUids []string `json:"subscribeAudioUids"`
}
type VideoUIDList struct {
	SubscribeVideoUids []string `json:"subscribeVideoUids"`
}
type StreamSubscribe struct {
	AudioUIDList AudioUIDList `json:"audioUidList"`
	VideoUIDList VideoUIDList `json:"videoUidList"`
}
type ClientRequest struct {
	StreamSubscribe StreamSubscribe `json:"streamSubscribe"`
}

type AgoraCallStartModel struct {
	Uid           string                    `json:"uid"`
	CName         string                    `json:"cname"`
	ClientRequest ClientRequestForStartCall `json:"clientRequest"`
}

type ClientRequestForStopCall struct {
	Async_stop bool `json:"async_stop"`
}

type AgoraCallStopModel struct {
	CName         string                   `json:"cname"`
	Uid           string                   `json:"uid"`
	ClientRequest ClientRequestForStopCall `json:"clientRequest"`
}

type AgoraCallStartResponse struct {
	Sid        string `json:"sid"`
	ResourceID string `json:"resourceId"`
}

type AgoraCallStatus struct {
	ResourceID     string `json:"resourceId"`
	Sid            string `json:"sid"`
	ServerResponse struct {
		FileListMode string `json:"fileListMode"`
		FileList     []struct {
			FileName       string `json:"fileName"`
			TrackType      string `json:"trackType"`
			UID            string `json:"uid"`
			MixedAllUser   bool   `json:"mixedAllUser"`
			IsPlayable     bool   `json:"isPlayable"`
			SliceStartTime int64  `json:"sliceStartTime"`
		} `json:"fileList"`
		Status         int   `json:"status"`
		Slicestarttime int64 `json:"sliceStartTime"`
	} `json:"serverResponse"`
}

// type AgoraCallStopResponseModel struct {
// 	ResourceID     string `json:"resourceId"`
// 	Sid            string `json:"sid"`
// 	ServerResponse struct {
// 		FileListMode string `json:"fileListMode"`
// 		FileList     []struct {
// 			FileName       string `json:"fileName"`
// 			TrackType      string `json:"trackType"`
// 			UID            string `json:"uid"`
// 			MixedAllUser   bool   `json:"mixedAllUser"`
// 			IsPlayable     bool   `json:"isPlayable"`
// 			SliceStartTime int64  `json:"sliceStartTime"`
// 		} `json:"fileList"`
// 		UploadingStatus string `json:"uploadingStatus"`
// 	} `json:"serverResponse"`
// }

type AgoraCallStopResponseModel struct {
	Code int `json:"Code"`
	Body struct {
		ResourceID     string `json:"resourceId"`
		Sid            string `json:"sid"`
		ServerResponse struct {
			FileListMode string `json:"fileListMode"`
			FileList     []struct {
				FileName       string `json:"fileName"`
				TrackType      string `json:"trackType"`
				UID            string `json:"uid"`
				MixedAllUser   bool   `json:"mixedAllUser"`
				IsPlayable     bool   `json:"isPlayable"`
				SliceStartTime int64  `json:"sliceStartTime"`
			} `json:"fileList"`
			UploadingStatus string `json:"uploadingStatus"`
		} `json:"serverResponse"`
	} `json:"Body"`
}

type EmailBodyMessageWithNameModel struct {
	Name string
}

type EmailBodyMessageModel struct {
	Name    string
	Message string
}

type EmailBodyWithAccessCodeMessageModel struct {
	Name       string
	Message    string
	AccessCode string
}

type EmailBodyMessageModelWithDocu struct {
	Name     string
	Message  string
	Message1 string
	Message2 string
	Message3 string
	Message4 string
}

type EmailRecipientModel struct {
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
}

type NotificationAllowSettingModel struct {
	UserType string `json:"userType"`
	Status   string `json:"status"`
}

type DocumentList struct {
	DocumentName string
	DocumentLink string
}

type OneSignalNotificationBulkData struct {
	AppID            string            `json:"app_id"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
	IncludedSegments []string          `json:"included_segments"`
	Data             map[string]string `json:"data"`
}
