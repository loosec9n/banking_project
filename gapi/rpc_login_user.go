package gapi

import (
	"context"
	"database/sql"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/utils"
	"simplebank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	validations := validateLoginUserRequest(req)
	if validations != nil {
		return nil, invalidArgumentError(validations)
	}
	user, err := server.store.GetUser(ctx, req.GetUserName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user form DB %s", err)
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password found %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenLife,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no token was created %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no refresh token was created %s", err)
	}

	mtdt := server.extractMetadata(ctx)

	session, err := server.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session %s", err)
	}

	resp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return resp, nil

}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUserName()); err != nil {
		violations = append(violations, fieldViolations("user_name", err))
	}
	// if err := val.ValidateFullName(req.GetFullName()); err != nil {
	// 	violations = append(violations, fieldViolations("full_name", err))
	// }
	// if err := val.ValidateEmail(req.GetEmail()); err != nil {
	// 	violations = append(violations, fieldViolations("email", err))
	// }
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolations("password", err))
	}
	return violations
}
