package models

type Organization struct {
	ID           uint32
	JwtSecret    string `json:"jwtSecret" gorm:"size:256"`
	JwtExpiresIn int64  `json:"jwtExpiresIn"`
}
