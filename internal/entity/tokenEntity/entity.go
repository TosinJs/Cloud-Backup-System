package tokenEntity

type JWTPayload struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}
