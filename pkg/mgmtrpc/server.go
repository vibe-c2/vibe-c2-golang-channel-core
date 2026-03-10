package mgmtrpc

import (
	"context"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/profile"
	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/runtime"
)

// Server is a placeholder for management RPC handlers.
type Server struct {
	Store runtime.ProfileStore
}

func NewServer(store runtime.ProfileStore) *Server {
	return &Server{Store: store}
}

func (s *Server) CreateProfile(ctx context.Context, channelID string, p profile.Profile) error {
	_ = ctx
	_ = channelID
	_ = p
	// TODO: Implement create profile RPC.
	return coreerrors.New(coreerrors.CodeNotImplemented, "CreateProfile not implemented")
}

func (s *Server) ReadProfile(ctx context.Context, channelID, profileID string) (profile.Profile, error) {
	_ = ctx
	_ = channelID
	_ = profileID
	// TODO: Implement read profile RPC.
	return profile.Profile{}, coreerrors.New(coreerrors.CodeNotImplemented, "ReadProfile not implemented")
}

func (s *Server) UpdateProfile(ctx context.Context, channelID string, p profile.Profile) error {
	_ = ctx
	_ = channelID
	_ = p
	// TODO: Implement update profile RPC.
	return coreerrors.New(coreerrors.CodeNotImplemented, "UpdateProfile not implemented")
}

func (s *Server) DeleteProfile(ctx context.Context, channelID, profileID string) error {
	_ = ctx
	_ = channelID
	_ = profileID
	// TODO: Implement delete profile RPC.
	return coreerrors.New(coreerrors.CodeNotImplemented, "DeleteProfile not implemented")
}

func (s *Server) ActivateProfile(ctx context.Context, channelID, profileID string) error {
	_ = ctx
	_ = channelID
	_ = profileID
	// TODO: Implement activate profile RPC.
	return coreerrors.New(coreerrors.CodeNotImplemented, "ActivateProfile not implemented")
}

func (s *Server) ValidateProfile(ctx context.Context, channelID string, p profile.Profile) error {
	_ = ctx
	_ = channelID
	_ = p
	// TODO: Implement validate profile RPC.
	return coreerrors.New(coreerrors.CodeNotImplemented, "ValidateProfile not implemented")
}
