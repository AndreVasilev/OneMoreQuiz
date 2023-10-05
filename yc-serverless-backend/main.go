package main

import (
	"log"
	"net/http"
	"yc-serverless-backend/handlers"
	"yc-serverless-backend/repository"
)

func main() {
	questions := repository.GetQuestions(0)
	log.Printf("Questions total: %d", len(questions))
}

func QuestionHandler(rw http.ResponseWriter, req *http.Request) {
	handlers.Question(rw, req)
}

func UserHandler(rw http.ResponseWriter, req *http.Request) {
	handlers.User(rw, req)
}
