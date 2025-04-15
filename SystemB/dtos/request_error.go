package dtos

import "strconv"

type RequestError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *RequestError) Error() string {
	return e.Message + " - StatusCode: " + strconv.Itoa(e.StatusCode)
}
