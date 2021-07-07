package server

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/shovanmaity/grpc-task/app/store"
	"github.com/shovanmaity/grpc-task/app/types"
	generated "github.com/shovanmaity/grpc-task/gen/go"
)

type LoginServer struct {
	generated.UnimplementedLoginServiceServer
	Store *store.Store
}

func NewLoginServer(db *store.Store) *LoginServer {
	return &LoginServer{
		Store: db,
	}
}
func (lc *LoginServer) Register(ctx context.Context, cred *generated.CredentialMessage) (*generated.TokenMessage, error) {
	if err := lc.Store.InsertToken(cred); err != nil {
		return nil, err
	}
	if _, err := lc.Store.InsertProfile(&generated.ProfileMessage{Username: cred.Username}); err != nil {
		return nil, err
	}
	session := lc.Store.GetSession(cred.Username)
	claim := types.JWTClaim{
		UserName:  cred.Username,
		SessionID: session,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	jwtString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(types.JWTSigningSecret()))
	if err != nil {
		return nil, err
	}
	return &generated.TokenMessage{Jwt: jwtString}, nil
}
func (lc *LoginServer) Login(ctx context.Context, cred *generated.CredentialMessage) (*generated.TokenMessage, error) {
	token := string(md5.New().Sum([]byte(cred.Username + ":" + cred.Password)))
	tokenInDB, err := lc.Store.GetToken(cred.Username)
	if err != nil {
		return nil, err
	}
	if token != tokenInDB {
		return nil, fmt.Errorf("username or password mismatch")
	}
	session := lc.Store.GetSession(cred.Username)
	claim := types.JWTClaim{
		UserName:  cred.Username,
		SessionID: session,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	jwtString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(types.JWTSigningSecret()))
	if err != nil {
		return nil, err
	}
	return &generated.TokenMessage{Jwt: jwtString}, nil
}

func (lc *LoginServer) Logout(ctx context.Context, msg *generated.EmptyMessage) (*generated.EmptyMessage, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
	}
	auths, ok := md[types.AuthorizationHeaderKey]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
	}
	jwtToken := ""
	for _, auth := range auths {
		parts := strings.Split(auth, " ")
		if len(parts) == 2 {
			if parts[0] == "Token" || parts[0] == "token" {
				jwtToken = parts[1]
			}
		}
	}
	if jwtToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
	}
	claim, err := jwt.ParseWithClaims(jwtToken, &types.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if ok := token.Method.Alg() == jwt.SigningMethodHS256.Alg(); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(types.JWTSigningSecret()), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	claimTyped, ok := claim.Claims.(*types.JWTClaim)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token claim")
	}
	lc.Store.RemoveSession(claimTyped.UserName)
	return &generated.EmptyMessage{}, nil
}
