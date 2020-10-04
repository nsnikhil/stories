package contract

type MostViewedStoriesRequest struct {
	OffSet int `json:"off_set"`
	Limit  int `json:"limit"`
}

type MostViewedStoriesResponse struct {
	Stories []Story `json:"stories"`
}
