package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/shovanmaity/grpc-task/app/store"
	"github.com/shovanmaity/grpc-task/app/types"
)

var unauthenticatedFullMethods map[string]int = make(map[string]int)

func init() {
	unauthenticatedFullMethods["/LoginService/Register"] = 0
	unauthenticatedFullMethods["/LoginService/Login"] = 0
	unauthenticatedFullMethods["/LoginService/Logout"] = 0

}

type Interceptor struct {
	DB *store.Store
}

func New(db *store.Store) *Interceptor {
	return &Interceptor{
		DB: db,
	}
}

func (inc *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := unauthenticatedFullMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}
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
		if inc.DB.IsValidSession(claimTyped.UserName, claimTyped.SessionID) {
			return handler(ctx, req)
		}
		return nil, status.Errorf(codes.Unauthenticated, "invalid session")
	}
}
