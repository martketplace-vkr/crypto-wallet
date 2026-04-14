package app

import (
	"context"

	"github.com/martketplace-vkr/crypto-wallet/config"
	"github.com/martketplace-vkr/crypto-wallet/internal/app/cmp/server"
	repository "github.com/martketplace-vkr/crypto-wallet/internal/repository/pg"
	servicewallet "github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
	clientTransport "github.com/martketplace-vkr/crypto-wallet/internal/transport/grpc/v1/client"
	"github.com/martketplace-vkr/pkg/build"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
)

func Run(ctx context.Context, cfg *config.Config) error {
	pg := pgxsqlxcomponent.New(cfg.Postgres)

	repo := repository.New(pg.DB)
	service := servicewallet.New(repo, cfg.Wallet)
	clientHandler := clientTransport.New(service)
	grpcServer := server.New(cfg.Grpc, clientHandler)

	cmps := build.Components{
		pg,
		grpcServer,
	}

	app, err := build.NewApp(cmps)
	if err != nil {
		return err
	}

	return build.Run(ctx, app)
}
