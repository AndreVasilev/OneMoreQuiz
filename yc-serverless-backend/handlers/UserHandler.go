package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"yc-serverless-backend/models"
	"yc-serverless-backend/repository"
)

func User(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		user, err := getUser(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if user == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		updateUser(user, req.URL.Query())
		position := repository.GetUsersUpperPositionCount(user.Score) + 1
		user.Position = &position

		rw.Header().Set("Content-Type", "application/json")
		json, err := json.Marshal(user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Write(json)
		rw.WriteHeader(200)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func updateUser(user *models.User, query url.Values) {
	hasChanges := false

	if query.Has("question_id") {
		id, err := strconv.ParseUint(query.Get("question_id"), 10, 64)
		if err == nil {
			user.LastQuestionId = id
			hasChanges = true
		}
	}
	successAnswer := false
	if query.Has("succeed") {
		successAnswer = query.Get("succeed") == "true"
		if successAnswer {
			user.SuccessAnswers += 1
			hasChanges = true
		}
	}
	if query.Has("score") {
		score, err := strconv.ParseUint(query.Get("score"), 10, 64)
		if err == nil && successAnswer {
			user.Score += score
			hasChanges = true
		}
	}

	if hasChanges {
		repository.UpdateUser(user)
	}
}
