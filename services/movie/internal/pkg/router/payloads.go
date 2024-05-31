package router

type createPayload struct {
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originalTitle"`
	Poster        string   `json:"poster"`
	Genres        []string `json:"genres"`
}
