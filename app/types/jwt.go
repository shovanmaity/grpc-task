package types

import (
	"github.com/golang-jwt/jwt"
)

// JWT token details
type JWTClaim struct {
	jwt.StandardClaims
	UserName  string
	SessionID string
}
