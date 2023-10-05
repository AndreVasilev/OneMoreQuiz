package models

type User struct {
	Id             int64
	LastQuestionId uint64
	Score          uint64
	SuccessAnswers uint64
	TgData         *string
	Position       *int
}
