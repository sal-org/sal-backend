package model

// type OneSignalNotificationData struct {
// 	AppID            string              `json:"app_id"`
// 	Headings         map[string]string   `json:"headings"`
// 	Contents         map[string]string   `json:"contents"`
// 	IncludePlayerIDs []string            `json:"include_player_ids"`
// 	Data             map[string]string   `json:"data"`
// 	BigPicture       string              `json:"big_picture"`
// 	IosAttachments   IosAttachmentsModel `json:"ios_attachments"`
// 	// SendAfter        string            `json:"send_after"`
// }

type IosAttachmentsModel struct {
	ID1 string `json:"id1"`
}

type OneSignalNotificationWithImage struct {
	AppID          string              `json:"app_id"`
	Headings       map[string]string   `json:"headings"`
	Contents       map[string]string   `json:"contents"`
	IncludeAliases IncludeAliase       `json:"include_aliases"`
	Channels       []string            `json:"target_channel"`
	Data           map[string]string   `json:"data"`
	BigPicture     string              `json:"big_picture"`
	IosAttachments IosAttachmentsModel `json:"ios_attachments"`
}

type OneSignalNotification struct {
	AppID          string            `json:"app_id"`
	Headings       map[string]string `json:"headings"`
	Contents       map[string]string `json:"contents"`
	IncludeAliases IncludeAliase     `json:"include_aliases"`
	Channels       []string          `json:"target_channel"`
	Data           map[string]string `json:"data"`
	// BigPicture       string            `json:"big_picture"`
	// URL              string            `json:"url"`
}

type IncludeAliase struct {
	ExternalID []string `json:"external_id"`
}
