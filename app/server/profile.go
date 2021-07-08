package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return dbProfile, nil
}

func (ps *ProfileServer) UpdateProfile(ctx context.Context,
	profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	dbProfile, err := ps.Store.UpdateProfile(profile.Username, profile)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return dbProfile, nil
}
