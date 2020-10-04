package contract

type GetStoryRequest struct {
	StoryID string `json:"story_id"`
}

type GetStoryResponse struct {
	Story Story `json:"story"`
}
