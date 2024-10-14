package dto

import "time"

//////////////////// Error response ////////////////////

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Status    int       `json:"status" binding:"required,min=400,max=599"`
	Path      string    `json:"path" format:"url_path" binding:"required" example:"/api/v0/example"`
}

func NewErrorResponse(message string, status int, path string) ErrorResponse {
	return ErrorResponse{
		Timestamp: time.Now(),
		Message:   message,
		Status:    status,
		Path:      path,
	}
}

func NewErrorResponseFromError(error error, status int, path string) ErrorResponse {
	return ErrorResponse{
		Timestamp: time.Now(),
		Message:   error.Error(),
		Status:    status,
		Path:      path,
	}
}

//////////////////// I18N error ////////////////////

type I18nError struct {
	message string
	tag     string
	args    interface{}
}

func NewI18nError(message, tag string) I18nError {
	return I18nError{message: message, tag: tag, args: nil}
}

func NewI18nErrorWithArgs(message, tag string, args interface{}) I18nError {
	return I18nError{message: message, tag: tag, args: args}
}

func (e I18nError) Error() string {
	return e.message
}

func (e I18nError) Tag() string {
	return e.tag
}

func (e I18nError) Args() interface{} {
	return e.args
}
