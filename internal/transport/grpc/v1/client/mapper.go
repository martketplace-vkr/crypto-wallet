package client

import (
	servicewallet "github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
	domainpb "github.com/martketplace-vkr/crypto-wallet/pkg/api/grpc/v1/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProto(address *servicewallet.DepositAddress) *domainpb.DepositAddress {
	if address == nil {
		return nil
	}

	return &domainpb.DepositAddress{
		Id:              address.ID,
		UserId:          address.UserID,
		Network:         address.Network,
		Asset:           address.Asset,
		Address:         address.Address,
		DerivationIndex: address.DerivationIndex,
		DerivationPath:  address.DerivationPath,
		CreatedAt:       timestamppb.New(address.CreatedAt),
		UpdatedAt:       timestamppb.New(address.UpdatedAt),
	}
}
