package repository

import "yc-serverless-backend/models"

type Interface interface {
	GetQuestions(from uint64) []models.Question
	GetUser(userId int64) *models.User
	AddUser(id int64, initData string) *models.User
	UpdateUser(user *models.User)
	GetUsersUpperPositionCount(score uint64) int
}
