package utils

import (
	"encoding/json"
	"go-fullstack-starter/schema"
	"net/http"
)

func SendResponse(w http.ResponseWriter, r *http.Request, response schema.Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	if err != nil {
		w.WriteHeader(statusCode)
		response := schema.Response{
			Status:    "ERROR",
			Message:   message + err.Error(),
			RequestID: r.Context().Value(schema.RequestIDKey{}).(string),
		}
		SendResponse(w, r, response)
	}
}

func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		HandleError(w, r, http.StatusMethodNotAllowed, "Method "+method+" not allowed", nil)
	}
}

func FeatureNotImplementedHandler(w http.ResponseWriter, r *http.Request, feature string) {
	HandleError(w, r, http.StatusNotImplemented, "Feature not yet implemented: "+feature, nil)
}
