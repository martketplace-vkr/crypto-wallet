package wallet

type Config struct {
	Mnemonic string `validate:"required"`
}
