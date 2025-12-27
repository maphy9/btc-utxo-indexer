package responses

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func NewLoginResponse(token string, refreshToken string) LoginResponse {
	return LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
}
