package contract

type AddStoryRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type AddStoryResponse struct {
	Success bool `json:"success"`
}
