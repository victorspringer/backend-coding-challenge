package router

type createPayload struct {
	Username    string `json:"username"`
	MD5Password string `json:"md5Password"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
}

type credentialsPayload struct {
	Username    string `json:"username"`
	MD5Password string `json:"md5Password"`
}
