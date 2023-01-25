package models

type ErrorNotFound struct {
	msg string
}
func (e ErrorNotFound) Error() string {
	return e.msg
}