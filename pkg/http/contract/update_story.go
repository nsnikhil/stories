package contract

type UpdateStoryRequest struct {
	Story Story `json:"story"`
}

type UpdateStoryResponse struct {
	Success bool `json:"success"`
}
