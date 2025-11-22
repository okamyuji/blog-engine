package presenter

import (
	"encoding/json"
	"net/http"
)

// JSONResponse JSONレスポンスを返す
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		}
	}
}

// ErrorResponse エラーレスポンス構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// JSONError JSONエラーレスポンスを返す
func JSONError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	JSONResponse(w, statusCode, response)
}

// SuccessResponse 成功レスポンス構造体
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSONSuccess 成功レスポンスを返す
func JSONSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}

	JSONResponse(w, http.StatusOK, response)
}
