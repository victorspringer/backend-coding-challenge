package router

type createPayload struct {
	UserID  string  `json:"userId"`
	MovieID string  `json:"movieId"`
	Value   float32 `json:"value"`
}
