package model

type OneSignalNotificationData struct {
	AppID            string            `json:"app_id"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
	IncludePlayerIDs []string          `json:"include_player_ids"`
	Data             map[string]string `json:"data"`
}
