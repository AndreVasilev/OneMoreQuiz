package tg

import (
	"encoding/json"
	"net/url"
	"yc-serverless-backend/models"
)

func NewWebAppInitData(query url.Values) *models.WebAppInitData {
	data := models.WebAppInitData{}
	for k, v := range query {
		if k == "user" {
			var user *models.WebAppUser
			json.Unmarshal([]byte(v[0]), &user)
			data.User = user
		}
	}
	return &data
}
