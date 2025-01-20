package types

import "time"

type JWTToken struct {
	Token string
	Exp   time.Time
}
type LoginResponse struct {
	AuthToken    JWTToken
	RefreshToken JWTToken
}
