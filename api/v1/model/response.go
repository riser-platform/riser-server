package model

type APIResponse struct {
	Message string `json:"message"`
}

type ValidationResponse struct {
	APIResponse
	ValidationErrors error `json:"validationErrors"`
}
