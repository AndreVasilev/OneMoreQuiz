package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"yc-serverless-backend/models"
	"yc-serverless-backend/repository"
)

func Question(rw http.ResponseWriter, req *http.Request) {

	// Check request

	var fromIndex uint64 = 0
	var user *models.User

	if req.Method == "POST" {
		u, err := getUser(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		user = u
		if user != nil {
			fromIndex = user.LastQuestionId
		}
	} else if req.Method != "GET" {
		rw.WriteHeader(http.StatusNotFound)
	}

	// Get questions

	questions := repository.GetQuestions(fromIndex)
	// Reset questions counter if user has answered all
	if len(questions) == 0 && user != nil {
		user.LastQuestionId = 0
		repository.UpdateUser(user)
		questions = repository.GetQuestions(0)
	}

	// Write response data

	if questions != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(questions)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		} else {
			_, err := rw.Write(jsonResp)
			if err != nil {
				log.Println(err.Error())
			}
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
