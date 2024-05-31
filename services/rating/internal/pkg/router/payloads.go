package router

type createPayload struct {
	UserID  string `json:"userId"`
	MovieID string `json:"movieId"`
	Value   int    `json:"value"`
}
