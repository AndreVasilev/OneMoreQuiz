package models

type Question struct {
	Id       uint64
	Question *string
	A        *string
	B        *string
	C        *string
	D        *string
	Answer   *string
}
