package gapi

import (
	"context"

	db "github.com/emrecolak-23/go-bank/db/sqlc"
	"github.com/emrecolak-23/go-bank/pb"
	"github.com/emrecolak-23/go-bank/utils"
	"github.com/emrecolak-23/go-bank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateUsersFullname(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmailAddress(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email_address", err))
	}

	return violations
}
