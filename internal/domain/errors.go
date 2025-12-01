package domain

type AppError struct {
	Status int
	Code   string
	Msg    string
}

func (e *AppError) Error() string {
	return e.Msg
}

func NotFound(msg string) *AppError { return &AppError{Status: 404, Code: "not_found", Msg: msg} }
func AlreadyExists(msg string) *AppError {
	return &AppError{Status: 409, Code: "already_exists", Msg: msg}
}
func Failed(msg string) *AppError      { return &AppError{Status: 500, Code: "failed", Msg: msg} }
func InvalidData(msg string) *AppError { return &AppError{Status: 400, Code: "invalid_data", Msg: msg} }
func Expired(msg string) *AppError     { return &AppError{Status: 400, Code: "expired", Msg: msg} }
func NotExpired(msg string) *AppError  { return &AppError{Status: 400, Code: "not_expired", Msg: msg} }
func Unauthorized(msg string) *AppError {
	return &AppError{Status: 401, Code: "unauthorized", Msg: msg}
}
