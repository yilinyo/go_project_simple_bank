package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func inValidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	br := &errdetails.BadRequest{FieldViolations: violations}
	st := status.New(codes.InvalidArgument, "invalid user request")
	st, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}
	return st.Err()
}
