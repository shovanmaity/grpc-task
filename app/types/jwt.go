package types

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	jwt.StandardClaims
	UserName  string
	SessionID string
}
