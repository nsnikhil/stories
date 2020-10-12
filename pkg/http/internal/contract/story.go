package contract

type Story struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	ViewCount int64  `json:"view_count"`
	UpVotes   int64  `json:"up_votes"`
	DownVotes int64  `json:"down_votes"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
