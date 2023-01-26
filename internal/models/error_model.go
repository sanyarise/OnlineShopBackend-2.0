package models

type ErrorNotFound struct {
	Message string
}

func (e ErrorNotFound) Error() string {
	return e.Message
}
