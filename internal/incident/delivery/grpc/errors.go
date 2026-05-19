package grpc

import (
	"github.com/tbikbulatov/go-pulseops/internal/platform/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch apperror.CodeOf(err) {
	case apperror.CodeNotFound:
		return status.Error(codes.NotFound, err.Error())
	case apperror.CodeInvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	case apperror.CodeAborted:
		return status.Error(codes.Aborted, err.Error())
	case apperror.CodeFailedPrecondition:
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
