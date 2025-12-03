package gapi

import (
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/pb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}
