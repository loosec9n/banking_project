package gapi

import (
	"context"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/utils"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hashing failed - %s", err)
	}

	arg := db.CreateUserParams{
		Username:     req.GetUserName(),
		HashPassword: hashPassword,
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
	}

	//user details from database
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists - %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create username - %s", err)
	}

	//converting user from database to user in api
	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}
