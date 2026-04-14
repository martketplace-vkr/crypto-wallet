package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientpb "github.com/martketplace-vkr/crypto-wallet/pkg/api/grpc/v1/client"
)

const (
	networkTRON = "TRON"
	assetUSDT   = "USDT"
)

type DepositAddress struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"user_id"`
	Network         string    `db:"network"`
	Asset           string    `db:"asset"`
	Address         string    `db:"address"`
	DerivationIndex uint32    `db:"derivation_index"`
	DerivationPath  string    `db:"derivation_path"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type repository interface {
	GetDepositAddress(ctx context.Context, userID int64, network, asset string) (*DepositAddress, error)
	NextDerivationIndex(ctx context.Context) (uint32, error)
	CreateDepositAddress(ctx context.Context, address DepositAddress) (*DepositAddress, error)
}

type Service struct {
	repository repository
	cfg        Config
}

func New(repository repository, cfg Config) *Service {
	return &Service{
		repository: repository,
		cfg:        cfg,
	}
}

func (s *Service) GetOrCreateDepositAddress(ctx context.Context, req *clientpb.GetOrCreateDepositAddressRequest) (*DepositAddress, error) {
	if err := validateRequest(req.GetUserId(), req.GetNetwork(), req.GetAsset()); err != nil {
		return nil, err
	}

	existing, err := s.repository.GetDepositAddress(ctx, req.GetUserId(), req.GetNetwork(), req.GetAsset())
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) && !isNotFound(err) {
		return nil, err
	}

	index, err := s.repository.NextDerivationIndex(ctx)
	if err != nil {
		return nil, err
	}

	address, path, err := DeriveTRONAddress(s.cfg.Mnemonic, index)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.CreateDepositAddress(ctx, DepositAddress{
		UserID:          req.GetUserId(),
		Network:         req.GetNetwork(),
		Asset:           req.GetAsset(),
		Address:         address,
		DerivationIndex: index,
		DerivationPath:  path,
	})
	if err != nil {
		if isUniqueRace(err) {
			return s.repository.GetDepositAddress(ctx, req.GetUserId(), req.GetNetwork(), req.GetAsset())
		}
		return nil, err
	}

	return created, nil
}

func (s *Service) GetDepositAddress(ctx context.Context, userID int64, network, asset string) (*DepositAddress, error) {
	if err := validateRequest(userID, network, asset); err != nil {
		return nil, err
	}

	address, err := s.repository.GetDepositAddress(ctx, userID, network, asset)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return address, nil
}

func validateRequest(userID int64, network, asset string) error {
	if userID <= 0 {
		return fmt.Errorf("%w: user_id must be greater than zero", ErrInvalidArgument)
	}
	if network != networkTRON {
		return fmt.Errorf("%w: only TRON network is supported", ErrUnsupported)
	}
	if asset != assetUSDT {
		return fmt.Errorf("%w: only USDT asset is supported", ErrUnsupported)
	}

	return nil
}
