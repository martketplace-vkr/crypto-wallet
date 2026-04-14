create schema if not exists crypto_wallet;

create sequence if not exists crypto_wallet.deposit_address_derivation_index_seq
    start with 1
    increment by 1
    minvalue 1;

create table if not exists crypto_wallet.deposit_address (
    id bigserial primary key,
    user_id bigint not null,
    network text not null,
    asset text not null,
    address text not null,
    derivation_index integer not null,
    derivation_path text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (user_id, network, asset),
    unique (network, address),
    unique (network, asset, derivation_index)
);

create index if not exists idx_crypto_wallet_deposit_address_user
    on crypto_wallet.deposit_address(user_id, network, asset);
