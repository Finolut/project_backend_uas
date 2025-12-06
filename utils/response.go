package utils

import (
	"encoding/json"
	"net/http"
)

// ApiResponse adalah response standard untuk semua endpoint
type ApiResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse mengirim response success
func SuccessResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim response error
func ErrorResponse(w http.ResponseWriter, code int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   err,
	})
}

// ExtractUserFromContext mengambil user claims dari JWT
func ExtractUserFromContext(r *http.Request) map[string]interface{} {
	ctx := r.Context()
	user, ok := ctx.Value("user").(map[string]interface{})
	if !ok {
		return nil
	}
	return user
}

// ExtractRoleFromContext mengambil role dari JWT
func ExtractRoleFromContext(r *http.Request) string {
	user := ExtractUserFromContext(r)
	if user == nil {
		return ""
	}
	role, ok := user["role"].(string)
	if !ok {
		return ""
	}
	return role
}
