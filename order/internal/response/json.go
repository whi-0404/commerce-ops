package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func WriteJSON(
	w http.ResponseWriter,
	statusCode int,
	payload any,
) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(
	w http.ResponseWriter,
	statusCode int,
	data any,
) error {
	return WriteJSON(w, statusCode, SuccessResponse{
		Data: data,
	})
}

func WriteError(
	w http.ResponseWriter,
	statusCode int,
	code string,
	message string,
	requestID string,
) error {
	return WriteJSON(w, statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	})
}