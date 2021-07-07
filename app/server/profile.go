package server

import (
	"context"

	_ "github.com/dgrijalva/jwt-go"

	"github.com/shovanmaity/grpc-task/app/store"
	generated "github.com/shovanmaity/grpc-task/gen/go"
)

type ProfileServer struct {
	generated.UnimplementedProfileServiceServer
	Store *store.Store
}

func NewProfileServer(db *store.Store) *ProfileServer {
	return &ProfileServer{
		Store: db,
	}
}

func (ps *ProfileServer) GetProfile(ctx context.Context,
	profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	dbProfile, err := ps.Store.GetProfile(profile.Username)
	if err != nil {
		return nil, err
	}
	return dbProfile, nil
}

func (ps *ProfileServer) UpdateProfile(ctx context.Context,
	profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	dbProfile, err := ps.Store.UpdateProfile(profile.Username, profile)
	if err != nil {
		return nil, err
	}
	return dbProfile, nil
}
