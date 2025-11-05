package gapi

import (
	"context"
	"database/sql"

	db "github.com/emrecolak-23/go-bank/db/sqlc"
	"github.com/emrecolak-23/go-bank/pb"
	"github.com/emrecolak-23/go-bank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	err = utils.CheckPasswordIsValid(user.HashedPassword, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password: %v", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}

	// Get user agent and client IP from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	userAgent := ""
	clientIP := ""
	if ok {
		if userAgents := md.Get("user-agent"); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
		if clientIPs := md.Get("x-forwarded-for"); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiresAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiresAt),
	}

	return rsp, nil
}
