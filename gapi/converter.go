package gapi

import (
	db "simplebank/db/sqlc"
	"simplebank/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		UserName:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: timestamppb.New(user.PasswordCahngedAt),
		CreatedAt:        timestamppb.New(user.CreatedAt),
	}
}
