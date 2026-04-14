package client

import (
	"context"

	servicewallet "github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
	clientpb "github.com/martketplace-vkr/crypto-wallet/pkg/api/grpc/v1/client"
)

type service interface {
	GetOrCreateDepositAddress(ctx context.Context, req *clientpb.GetOrCreateDepositAddressRequest) (*servicewallet.DepositAddress, error)
	GetDepositAddress(ctx context.Context, userID int64, network, asset string) (*servicewallet.DepositAddress, error)
}
