package gapi

import (
	"context"
	"errors"
	"strings"

	"github.com/yilinyo/project_bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authorizations := md.Get(authorizationHeader)
	if len(authorizations) == 0 {
		return nil, errors.New("missing authorization header")
	}

	accessHeader := authorizations[0]
	//Bearer xxxxxxx
	fields := strings.Fields(accessHeader)
	if len(fields) != 2 {
		return nil, errors.New("invalid authorization header format")
	}
	authType := fields[0]
	if strings.ToLower(authType) != authorizationType {
		return nil, errors.New("invalid authorization type")
	}
	payload, err := server.tokenMaker.VerifyToken(fields[1], token.TokenTypeAccessToken)
	if err != nil {
		return nil, err
	}
	return payload, nil

}
