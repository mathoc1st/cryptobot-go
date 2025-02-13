package cryptobot

type ExchangeRate struct {
	// Whether or not the received rate is up-to-date.
	IsValid bool `json:"is_valid"`

	// True if Source is crypto.
	IsCrypto bool `json:"is_crypto"`

	// True if Source is fiat.
	IsFiat bool `json:"is_fiat"`

	// Type of cryptocurrency.
	Source CryptoAsset `json:"source"`

	// Type of fiat currency.
	Target CurrencyCode `json:"target"`

	// The current rate of the source asset valued in the target currency.
	Rate string `json:"rate"`
}
