package types

import (
	"os"
)

// JWTSigningSecret used to sign the claim and get the claim from token
var JWTSigningSecret = func() string {
	if os.Getenv("JWT_SIGNING_SECRET") == "" {
		return "secret"
	}
	return os.Getenv("JWT_SIGNING_SECRET")
}

// ServerPort used to run gRPC server. Default port 8080
var ServerPort = func() string {
	if os.Getenv("SERVER_PORT") == "" {
		return "8080"
	}
	return os.Getenv("SERVER_PORT")
}

// GatewayPort used to run gRPC gateway. Default port 8090
var GatewayPort = func() string {
	if os.Getenv("GATEWAY_PORT") == "" {
		return "8090"
	}
	return os.Getenv("GATEWAY_PORT")
}
