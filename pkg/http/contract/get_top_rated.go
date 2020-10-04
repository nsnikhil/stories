package contract

type TopRatedStoriesRequest struct {
	OffSet int `json:"off_set"`
	Limit  int `json:"limit"`
}

type TopRatedStoriesResponse struct {
	Stories []Story `json:"stories"`
}
