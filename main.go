package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func main() {
	claims := jwt.MapClaims{}
	claims["sessionid"] = "session"
	claims["exp"] = time.Now().Unix()
	jwtSigningSecret := os.Getenv("JWT_SIGNING_SECRET")
	jwtString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSigningSecret))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jwtString)
}
