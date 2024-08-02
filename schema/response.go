package schema

type Response struct {
	Data      interface{} `json:"data,omitempty"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
}
