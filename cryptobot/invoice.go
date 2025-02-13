package cryptobot

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type CurrencyType string

// Types of available currencies.
const (
	Crypto CurrencyType = "crypto"
	Fiat   CurrencyType = "fiat"
)

// Identifies different cryptocurrency assets (e.g., ETH, BTC etc.).
type CryptoAsset string

// All the available cryptocurrency types.
const (
	USDT CryptoAsset = "USDT"
	TON              = "TON"
	BTC              = "BTC"
	ETH              = "ETH"
	LTC              = "LTC"
	BNB              = "BNB"
	TRX              = "TRX"
	USDC             = "USDC"
)

type CurrencyCode string

// Types of available fiat currency codes.
const (
	// US Dollar
	USD CurrencyCode = "USD"
	// Euro
	EUR = "EUR"
	// Russian Ruble
	RUB = "RUB"
	// Belarusian Ruble
	BYN = "BYN"
	// Ukrainian Hryvnia
	UAH = "UAH"
	// British Pound Sterling
	GBP = "GBP"
	// Chinese Yuan
	CNY = "CNY"
	// Kazakhstani Tenge
	KZT = "KZT"
	// Uzbekistani Som
	UZS = "UZS"
	// Georgian Lari
	GEL = "GEL"
	// Turkish Lira
	TRY = "TRY"
	// Armenian Dram
	AMD = "AMD"
	// Thai Baht
	THB = "THB"
	// Indian Rupee
	INR = "INR"
	// Brazilian Real
	BRL = "BRL"
	// Indonesian Rupiah
	IDR = "IDR"
	// Azerbaijani Manat
	AZN = "AZN"
	// United Arab Emirates Dirham
	AED = "AED"
	// Polish Zloty
	PLN = "PLN"
	// Israeli New Shekel
	ILS = "ILS"
)

type InvoiceStatus string

const (
	InvoicePaid    InvoiceStatus = "paid"
	InvoiceActive                = "active"
	InvoiceExpired               = "expired"
)

type ButtonName string

const (
	ViewItem    ButtonName = "viewItem"
	OpenChannel            = "openChannel"
	OpenBot                = "openBot"
	Callback               = "callback"
)

type Invoice struct {
	// Unique ID for this invoice.
	ID int64 `json:"invoice_id"`

	// Hash of the invoice.
	Hash string `json:"hash"`

	// Type of currency.
	CurrencyType CurrencyType `json:"currency_type"`

	// Available only if the CurrencyType is set to crypto. Cryptocurrency type to pay the invoice with.
	CryptoAsset CryptoAsset `json:"asset,omitempty"`

	// Available only if the CurrencyType is set to fiat. Cryptocurrency types that can be used to pay the invoice.
	AcceptedCryptoAssets []CryptoAsset `json:"accepted_assets,omitempty"`

	// Available only if the CurrencyType is set to fiat. Fiat currency type.
	Fiat CurrencyCode `json:"fiat,omitempty"`

	// Amount of the invoice.
	Amount string `json:"amount"`

	// Available only if CurrencyType is fiat and Status is invoicePaid. Cryptocurrency that was used to pay the invoice.
	PaidAsset CryptoAsset `json:"paid_asset,omitempty"`

	// Available only if CurrencyType is fiat and Status is invoicePaid. Amount of the invoice for which the invoice was paid.
	PaidAmount string `json:"paid_amount,omitempty"`

	// Available only if CurrencyType is fiat and Status is invoicePaid. The rate of the PaidAsset value in the fiat currency.
	PaidFiatRate string `json:"paid_fiat_rate,omitempty"`

	// Available only if Status is invoicePaid. Cryptocurrency that was used to pay the invoice fee.
	FeeAsset string `json:"fee_asset,omitempty"`

	// Available only if Status is invoicePaid. Fee amount that was charged for the invoice.
	FeeAmount int64 `json:"fee_amount,omitempty"`

	// URL for the user to pay the invoice using Crypto Bot.
	BotInvoiceURL string `json:"bot_invoice_url"`

	// URL for the user to pay the invoice using Mini App version of Crypto Bot.
	MiniAppInvoiceURL string `json:"mini_app_invoice_url"`

	// URL for the user to pay the invoice using WebApp version of Crypto Bot.
	WebAppInvoiceURL string `json:"web_app_invoice_url"`

	// Optional. Description for this invoice.
	Description string `json:"description,omitempty"`

	// Invoice status.
	Status InvoiceStatus `json:"status"`

	// Date the invoice was created (ISO 8601 format).
	CreatedAt string `json:"created_at"`

	// Available only if Status is invoicePaid. Rate of the cryptocurrency that was used during payment in USD.
	PaidUSDRate string `json:"paid_usd_rate,omitempty"`

	// Whether or not the user can add a comment to the invoice.
	AllowComments bool `json:"allow_comments"`

	// Whether or not the user can pay anonymously.
	AllowAnonymous bool `json:"allow_anonymous"`

	// Available only if the expiration date was set when the invoice was created. Expiration date of the invoice.
	ExpirationDate string `json:"expiration_date,omitempty"`

	// Available only if Status is invoicePaid. Date the invoice was paid (ISO 8601 format).
	PaidAt string `json:"paid_at,omitempty"`

	// Whether or not the invoice was paid anonymously.
	PaidAnonymously bool `json:"paid_anonymously"`

	// Optional. User's comment.
	Comment string `json:"comment,omitempty"`

	// Optional. Hidden message that is shown to the user when the invoice is paid.
	HiddenMessage string `json:"hidden_message,omitempty"`

	// Optional. Payload that is attached to the invoice.
	Payload string `json:"payload,omitempty"`

	// Optional. Type of the button that is shown to the user when the invoice is paid.
	PaidBtnName ButtonName `json:"paid_btn_name,omitempty"`

	// Available only if PaidBtnName was set. URL attached to the button.
	PaidBtnUrl string `json:"paid_btn_url,omitempty"`
}

type NewInvoice struct {
	// Type of currency that should be used to pay the invoice.
	CurrencyType CurrencyType

	// Should be set if CurrencyType is crypto. Type of cryptocurrency to pay the invoice with.
	CryptoAsset CryptoAsset

	// Should be set if CurrencyType is fiat. Type of fiat currency to pay the invoice with.
	Fiat CurrencyCode

	// Should be set if CurrencyType is fiat. Cryptocurrency types that can be used to pay the invoice with.
	AcceptedCryptoAssets []CryptoAsset

	// Amount the user will have to pay.
	Amount string

	// Optional. Description for the invoice. 1024 characters max.
	Description string

	// Optional. Hidden message that will be shown to the user once the invoice is paid. 2048 characters max.
	HiddenMessage string

	// Optional. Type of the button that will be shown to the user once the invoice is paid.
	PaidBtnName ButtonName

	// Should be set if PaidBtnName is set. URL that will be attached to the button.
	PaidBtnUrl string

	// Optional. Payload to attach to the invoice. 4096 characters max.
	Payload string

	// Whether or not a user can add comments to the payment.
	AllowComments bool

	// Whether or not a user can pay the invoice anonymously.
	AllowAnonymous bool

	// Optional. Expiration time of the invoice in seconds. Values between 1-2678400 are accepted.
	ExpiresIn int64
}

type tempNewInvoice struct {
	CurrencyType         CurrencyType `json:"currency_type"`
	CryptoAsset          CryptoAsset  `json:"asset,omitempty"`
	Fiat                 CurrencyCode `json:"fiat,omitempty"`
	AcceptedCryptoAssets string       `json:"accepted_assets,omitempty"`
	Amount               string       `json:"amount"`
	Description          string       `json:"description,omitempty"`
	HiddenMessage        string       `json:"hidden_message,omitempty"`
	PaidBtnName          ButtonName   `json:"paid_btn_name,omitempty"`
	PaidBtnUrl           string       `json:"paid_btn_url,omitempty"`
	Payload              string       `json:"payload,omitempty"`
	AllowComments        bool         `json:"allow_comments"`
	AllowAnonymous       bool         `json:"allow_anonymous"`
	ExpiresIn            int64        `json:"expires_in,omitempty"`
}

func (in NewInvoice) MarshalJSON() ([]byte, error) {
	var as []string

	for _, a := range in.AcceptedCryptoAssets {
		as = append(as, string(a))
	}

	return json.Marshal(tempNewInvoice{
		CurrencyType:         in.CurrencyType,
		CryptoAsset:          in.CryptoAsset,
		Fiat:                 in.Fiat,
		AcceptedCryptoAssets: strings.Join(as, ","),
		Amount:               in.Amount,
		Description:          in.Description,
		HiddenMessage:        in.HiddenMessage,
		PaidBtnName:          in.PaidBtnName,
		PaidBtnUrl:           in.PaidBtnUrl,
		Payload:              in.Payload,
		AllowComments:        in.AllowComments,
		AllowAnonymous:       in.AllowAnonymous,
		ExpiresIn:            in.ExpiresIn,
	})
}

type InvoiceOptions struct {
	// Optional. Type of cryptocurrency to search by.
	CryptoAsset CryptoAsset `json:"asset,omitempty"`

	// Optional. Type of fiat currency to search by.
	Fiat CurrencyCode `json:"fiat,omitempty"`

	// Optional. Invoice ids to find.
	InvoiceIDs []int64 `json:"invoice_ids,omitempty"`

	// Optional. Status to search by.
	Status InvoiceStatus `json:"status,omitempty"`

	// Optional. Defaults to 0.
	Offset int64 `json:"offset,omitempty"`

	// Optional. Number of invoices to be returned. Values between 1-1000 are accepted. Defaults to 100.
	Count int64 `json:"count,omitempty"`
}

type tempInOps struct {
	CryptoAsset  string `json:"asset,omitempty"`
	FiatCurrency string `json:"fiat,omitempty"`
	InvoiceIDs   string `json:"invoice_ids,omitempty"`
	Status       string `json:"status,omitempty"`
	Offset       int64  `json:"offset,omitempty"`
	Count        int64  `json:"count,omitempty"`
}

func (no InvoiceOptions) MarshalJSON() ([]byte, error) {
	ids := make([]string, len(no.InvoiceIDs))

	for _, id := range no.InvoiceIDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}

	return json.Marshal(&tempInOps{
		CryptoAsset:  string(no.CryptoAsset),
		FiatCurrency: string(no.Fiat),
		InvoiceIDs:   strings.Join(ids, ","),
		Status:       string(no.Status),
		Offset:       no.Offset,
		Count:        no.Count,
	})
}

func validateNewInvoice(in NewInvoice) error {
	var errs []error
	if len(in.CurrencyType) == 0 {
		errs = append(errs, errors.New("CurrencyType cannot be empty"))
	}
	if in.CurrencyType == Crypto && len(in.CryptoAsset) == 0 {
		errs = append(errs, errors.New("CryptoAsset cannot be empty"))
	}
	if in.CurrencyType == Fiat && len(in.AcceptedCryptoAssets) == 0 {
		errs = append(errs, errors.New("AcceptedCryptoAssets cannot be empty"))
	}
	if in.CurrencyType == Fiat && len(in.Fiat) == 0 {
		errs = append(errs, errors.New("FiatCurrency cannot be empty"))
	}
	if len(in.Amount) == 0 {
		errs = append(errs, errors.New("Amount cannot be empty"))
	}
	if len(in.PaidBtnName) != 0 && len(in.PaidBtnUrl) == 0 {
		errs = append(errs, errors.New("PaidBtnUrl cannot be empty"))
	}
	if len(in.Payload) > 4096 {
		errs = append(errs, errors.New("Payload should not exceed 4096 characters"))
	}
	if in.ExpiresIn != 0 && (in.ExpiresIn < 1 || in.ExpiresIn > 2678400) {
		errs = append(errs, errors.New("expiration time should be within 1-2678400 second range"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func validateInvoiceOptions(inop InvoiceOptions) error {
	var errs []error
	if inop.Offset < 0 {
		errs = append(errs, errors.New("Offset cannot be less than 0"))
	}
	if inop.Count != 0 && (inop.Count < 1 || inop.Count > 1000) {
		errs = append(errs, errors.New("Count needs to be within 1-1000 record range"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
