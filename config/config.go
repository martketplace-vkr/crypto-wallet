package config

import (
	"github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/server/grpc"
)

type Config struct {
	Grpc     grpc.Config             `validate:"required"`
	Postgres pgxsqlxcomponent.Config `validate:"required"`
	Wallet   wallet.Config           `validate:"required"`
}
