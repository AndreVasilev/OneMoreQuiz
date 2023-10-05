package handlers

import (
	"encoding/json"
	"io"
	"yc-serverless-backend/models"
	"yc-serverless-backend/repository"
	"yc-serverless-backend/tg"
)

type Body struct {
	InitData string `json:"initData"`
}

func getValidInitData(reqBody io.ReadCloser) (*models.WebAppInitData, *string, error) {
	data, err := io.ReadAll(reqBody)
	if err != nil {
		return nil, nil, err
	}
	var body Body
	json.Unmarshal(data, &body)
	return tg.Validate(body.InitData), &body.InitData, nil
}

func getUser(reqBody io.ReadCloser) (*models.User, error) {
	initData, initDataString, err := getValidInitData(reqBody)
	if err != nil || initData == nil || initData.User == nil {
		return nil, err
	}
	user := repository.GetUser(initData.User.Id)
	if user == nil {
		user = repository.AddUser(initData.User.Id, *initDataString)
	}
	return user, nil
}
