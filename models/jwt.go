package models

type RefreshJwtToken struct {
	ID        uint
	SessionID string
	Role      string
	OrgSlug   string
	Version   string
}

type AccessJwtToken struct {
	ID        uint
	SessionID string
	Role      string
	OrgSlug   string
	Version   string
}

type Token struct {
	ID           string
	UserID       uint
	AccessToken  string
	RefreshToken string
}

type JwtTokensData struct {
	RefreshToken RefreshJwtToken
	AccessToken  AccessJwtToken
}

type JwtTokens struct {
	RefreshToken string
	AccessToken  string
}
