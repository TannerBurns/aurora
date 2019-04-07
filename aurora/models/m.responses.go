package models

type Exception struct {
	Message string `json:"message"`
}

type JwtToken struct {
	Token         string `json:"token"`
	UserProfileID int    `json:"userprofile_id"`
}

type RespError struct {
	Error string `json:"error"`
}
