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
