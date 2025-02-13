package cryptobot

type Balance struct {
	// Cryptocurrency type.
	CryptoAsset CryptoAsset `json:"currency_code"`

	// Total available amount.
	Available string `json:"available"`

	// Amount that is on hold and currenty unavailable.
	OnHold string `json:"onhold"`
}
