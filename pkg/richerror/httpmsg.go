package richerror

import (
	"errors"
	"net/http"
)

type MessageResponse struct {
	Message   string         `json:"message"`
	ErrorCode string         `json:"error_code"`
	Meta      map[string]any `json:"meta,omitempty"`
}

func ToHTTP(err error) (MessageResponse, int) {
	resp := MessageResponse{
		Message:   "bad request",
		ErrorCode: "ERR_BAD_REQUEST",
	}

	var re *RichError
	ok := errors.As(err, &re)
	if !ok {
		return resp, http.StatusBadRequest
	}

	resp.Message = re.Error()
	resp.Meta = re.Meta()
	status := mapKindToStatus(re.Kind())

	if status >= 500 {
		resp.Message = "something went wrong"
	}

	resp.ErrorCode = statusToCode(status)

	return resp, status
}

func mapKindToStatus(kind Kind) int {
	switch kind {
	case KindInvalid:
		return http.StatusUnprocessableEntity
	case KindNotFound:
		return http.StatusNotFound
	case KindForbidden:
		return http.StatusForbidden
	case KindUnexpected:
		return http.StatusInternalServerError
	case KindRateLimit:
		return http.StatusTooManyRequests
	case KindBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}

func statusToCode(status int) string {
	return http.StatusText(status)
}
