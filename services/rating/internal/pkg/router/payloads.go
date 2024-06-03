package router

type upsertPayload struct {
	UserID  string  `json:"userId"`
	MovieID string  `json:"movieId"`
	Value   float32 `json:"value"`
}
