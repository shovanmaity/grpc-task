package types

import (
	"os"
)

var JWTSigningSecret = func() string {
	if os.Getenv("JWT_SIGNING_SECRET") == "" {
		return "secret"
	}
	return os.Getenv("JWT_SIGNING_SECRET")
}
