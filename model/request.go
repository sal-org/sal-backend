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

// ClientProfileAddRequest .
type ClientProfileAddRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	Photo       string `json:"photo"`
	TopicIDs    string `json:"topic_ids"`
	Gender      string `json:"gender"`
	Location    string `json:"location"`
	Timezone    string `json:"timezone"`
	DeviceID    string `json:"device_id"`
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

// ClientProfileUpdateRequest .
type ClientProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Location    string `json:"location"`
	Timezone    string `json:"timezone"`
	DeviceID    string `json:"device_id"`
	DateOfBirth string `json:"date_of_birth"`
	Photo       string `json:"photo"`
	TopicIDs    string `json:"topic_ids"`
	Gender      string `json:"gender"`
}

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

// AppointmentRatingAdd .
type AppointmentRatingAdd struct {
	AppointmentID string `json:"appointment_id"`
	Rating        string `json:"rating"`
	RatingTypes   string `json:"rating_types"`
	RatingComment string `json:"rating_comment"`
	ClientID      string `json:"client_id"`
	CounsellorID  string `json:"counsellor_id"`
}

// AssessmentAddRequest .
type AssessmentAddRequest struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Age          string `json:"age"`
	Gender       string `json:"gender"`
	Phone        string `json:"phone"`
	AssessmentID string `json:"assessment_id"`
	Feedback     string `json:"feedback"`
	Details      []struct {
		AssessmentQuestionID       string `json:"assessment_question_id"`
		AssessmentQuestionOptionID string `json:"assessment_question_option_id"`
		Score                      string `json:"score"`
	} `json:"details"`
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
	Media_URL   string
	First_Name  string
	Last_Name   string
	Gender      string
	Type        string
	Phone       string
	Email       string
	Photo       string
	Education   string
	Experience  string
	About       string
	Resume      string
	Certificate string
	Aadhar      string
	Linkedin    string
	Status      string
}

type EmailDataForCounsellorRecord struct {
	First_Name      string
	Last_Name       string
	Gender          string
	Age             string
	Department      string
	Location        string
	SessionMode     string
	SessionNo       string
	SessionDate     string
	InTime          string
	OutTime         string
	TherapeuticGoal string
	TherapyPlan     string
	AssessmentTool  string
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

type AgoraCallStopResponseModel struct {
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
}

type EmailBodyMessageModel struct {
	Name    string
	Message string
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
