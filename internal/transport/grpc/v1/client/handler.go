package client

import (
	"context"

	clientpb "github.com/martketplace-vkr/crypto-wallet/pkg/api/grpc/v1/client"
)

type Handler struct {
	service service
	clientpb.UnimplementedCryptoWalletClientServiceServer
}

func New(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetOrCreateDepositAddress(ctx context.Context, req *clientpb.GetOrCreateDepositAddressRequest) (*clientpb.GetOrCreateDepositAddressResponse, error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}

	address, err := h.service.GetOrCreateDepositAddress(ctx, req)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &clientpb.GetOrCreateDepositAddressResponse{
		DepositAddress: toProto(address),
	}, nil
}

func (h *Handler) GetDepositAddress(ctx context.Context, req *clientpb.GetDepositAddressRequest) (*clientpb.GetDepositAddressResponse, error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}

	address, err := h.service.GetDepositAddress(ctx, req.GetUserId(), req.GetNetwork(), req.GetAsset())
	if err != nil {
		return nil, toStatusError(err)
	}

	return &clientpb.GetDepositAddressResponse{
		DepositAddress: toProto(address),
	}, nil
}
