package responses

type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenResponse(token string, refreshToken string) TokenResponse {
	return TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
}
