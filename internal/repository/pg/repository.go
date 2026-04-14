package pg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/martketplace-vkr/crypto-wallet/internal/service/wallet"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDepositAddress(ctx context.Context, userID int64, network, asset string) (*wallet.DepositAddress, error) {
	query := `
		select
			id,
			user_id,
			network,
			asset,
			address,
			derivation_index,
			derivation_path,
			created_at,
			updated_at
		from crypto_wallet.deposit_address
		where user_id = $1
			and network = $2
			and asset = $3
	`

	var result wallet.DepositAddress
	if err := r.db.GetContext(ctx, &result, query, userID, network, asset); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) NextDerivationIndex(ctx context.Context) (uint32, error) {
	query := `select nextval('crypto_wallet.deposit_address_derivation_index_seq')`

	var idx int64
	if err := r.db.GetContext(ctx, &idx, query); err != nil {
		return 0, err
	}

	return uint32(idx - 1), nil
}

func (r *Repository) CreateDepositAddress(ctx context.Context, address wallet.DepositAddress) (*wallet.DepositAddress, error) {
	query := `
		insert into crypto_wallet.deposit_address (
			user_id,
			network,
			asset,
			address,
			derivation_index,
			derivation_path
		) values (
			:user_id,
			:network,
			:asset,
			:address,
			:derivation_index,
			:derivation_path
		)
		returning
			id,
			user_id,
			network,
			asset,
			address,
			derivation_index,
			derivation_path,
			created_at,
			updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var created wallet.DepositAddress
	if err := rows.StructScan(&created); err != nil {
		return nil, err
	}

	return &created, nil
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
