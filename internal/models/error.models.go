package models

type ErrorResponseDetail struct {
	Error *ErrorResponse
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Status  int    `json:"status"`
}
