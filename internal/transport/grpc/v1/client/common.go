package client

import (
	"errors"

	servicewallet "github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateID(field string, value int64) error {
	if value <= 0 {
		return status.Errorf(codes.InvalidArgument, "%s must be greater than zero", field)
	}

	return nil
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, servicewallet.ErrInvalidArgument),
		errors.Is(err, servicewallet.ErrUnsupported):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, servicewallet.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return err
	}
}
