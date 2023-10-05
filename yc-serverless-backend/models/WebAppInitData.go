package models

type WebAppInitData struct {
	User *WebAppUser
}

type WebAppUser struct {
	Id              int64   `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        *string `json:"last_name,omitempty"`
	Username        *string `json:"username,omitempty"`
	LanguageCode    *string `json:"language_code,omitempty"`
	IsPremium       *bool   `json:"is_premium,omitempty"`
	AllowsWriteToPm *bool   `json:"allows_write_to_pm,omitempty"`
}
