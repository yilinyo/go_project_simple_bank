package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/pb"
	"github.com/yilinyo/project_bank/token"
	"github.com/yilinyo/project_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to login user: %s", err)
	}

	//check password
	if isValid := util.CheckPassword(req.GetPassword(), user.HashedPassword); !isValid {

		return nil, status.Errorf(codes.Internal, "Wrong username or password: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, "admin", server.config.AccessTokenDuration, token.TokenTypeAccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token: %s", err)
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, "admin", server.config.RefreshTokenDuration, token.TokenTypeRefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token: %s", err)
	}
	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %s", err)
	}

	resp := &pb.UserLoginResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return resp, nil
}
