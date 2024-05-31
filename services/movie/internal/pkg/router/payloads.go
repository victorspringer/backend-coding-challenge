package router

type createPayload struct {
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originalTitle"`
	Overview      string   `json:"overview"`
	Poster        string   `json:"poster"`
	Genres        []string `json:"genres"`
	Keywords      []string `json:"keywords"`
}
