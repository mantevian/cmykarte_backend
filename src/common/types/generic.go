package types

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func GenericResponseSuccess(message string) GenericResponse {
	return GenericResponse{
		Success: true,
		Message: message,
	}
}

func GenericResponseError(message string) GenericResponse {
	return GenericResponse{
		Success: false,
		Message: message,
	}
}
