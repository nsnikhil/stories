package contract

type DeleteStoryRequest struct {
	StoryID string `json:"story_id"`
}

type DeleteStoryResponse struct {
	Success bool `json:"success"`
}
