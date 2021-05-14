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

// CounsellorOrderPaymentCompleteRequest .
type CounsellorOrderPaymentCompleteRequest struct {
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

// ClientEventOrderCreateRequest .
type ClientEventOrderCreateRequest struct {
	ClientID     string `json:"client_id"`
	EventOrderID string `json:"event_order_id"`
	CouponCode   string `json:"coupon_code"`
}

// CounsellorEventOrderCreateRequest .
type CounsellorEventOrderCreateRequest struct {
	CounsellorID string `json:"counsellor_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	TopicID      string `json:"topic_id"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Duration     string `json:"duration"`
	Price        string `json:"price"`
}

// ClientEventOrderPaymentCompleteRequest .
type ClientEventOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// CounsellorEventOrderPaymentCompleteRequest .
type CounsellorEventOrderPaymentCompleteRequest struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentID     string `json:"payment_id"`
}

// ClientProfileAddRequest .
type ClientProfileAddRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Location  string `json:"location"`
	DeviceID  string `json:"device_id"`
}

// CounsellorProfileAddRequest .
type CounsellorProfileAddRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
	Email       string `json:"email"`
	Price       string `json:"price"`
	Price3      string `json:"price_3"`
	Price5      string `json:"price_5"`
	Education   string `json:"education"`
	Experience  string `json:"experience"`
	About       string `json:"about"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	Resume      string `json:"resume"`
	Certificate string `json:"certificate"`
	Aadhar      string `json:"aadhar"`
	Linkedin    string `json:"linkedin"`
	DeviceID    string `json:"device_id"`
}

// ListenerProfileAddRequest .
type ListenerProfileAddRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
	Email       string `json:"email"`
	Occupation  string `json:"occupation"`
	Experience  string `json:"experience"`
	About       string `json:"about"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	DeviceID    string `json:"device_id"`
}

// ClientProfileUpdateRequest .
type ClientProfileUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location  string `json:"location"`
	DeviceID  string `json:"device_id"`
}

// CounsellorProfileUpdateRequest .
type CounsellorProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Photo       string `json:"photo"`
	Price       string `json:"price"`
	Price3      string `json:"price_3"`
	Price5      string `json:"price_5"`
	Education   string `json:"education"`
	Experience  string `json:"experience"`
	About       string `json:"about"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	Resume      string `json:"resume"`
	Certificate string `json:"certificate"`
	Aadhar      string `json:"aadhar"`
	Linkedin    string `json:"linkedin"`
	DeviceID    string `json:"device_id"`
}

// ListenerProfileUpdateRequest .
type ListenerProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Photo       string `json:"photo"`
	Occupation  string `json:"occupation"`
	Experience  string `json:"experience"`
	About       string `json:"about"`
	TopicIDs    string `json:"topic_ids"`
	LanguageIDs string `json:"language_ids"`
	DeviceID    string `json:"device_id"`
}
