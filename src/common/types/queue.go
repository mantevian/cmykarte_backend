package types

type QueueRequest struct {
}

type QueueResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func QueueError(message string) QueueResponse {
	return QueueResponse{
		Success: false,
		Message: message,
	}
}
