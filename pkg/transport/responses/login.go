package responses

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Exp          string `json:"exp"`
	UserID       string `json:"id"`
}
