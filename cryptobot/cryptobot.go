package cryptobot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	Mainnet = "https://pay.crypt.bot/api"
	Testnet = "https://testnet-pay.crypt.bot/api"
)

type resultConstraint interface {
	json.RawMessage | Invoice | Check | Transfer | AppStats | []Balance | []ExchangeRate | bool | struct {
		Items []Invoice `json:"items"`
	} | struct {
		Items []Check `json:"items"`
	} | struct {
		Items []Transfer `json:"items"`
	}
}

type response[T resultConstraint] struct {
	Ok     bool            `json:"ok"`
	Error  json.RawMessage `json:"error"`
	Result T               `json:"result"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	// Cryptobot API token
	Token string
	// Mainnet or Testnet
	Endpoint string
	Client   HTTPClient
}

type cryptobot struct {
	token    string
	client   HTTPClient
	endpoint string
}

// New creates a new crypto bot instance. There are two endpoints: Testnet and Mainnet.
// Testnet is used for testing and Mainnet for production. You need a different token for each of the networks.
// It uses the default http client if none is provided.
func New(cf Config) (*cryptobot, error) {
	if len(cf.Token) == 0 {
		return nil, errors.New("no token was provided for crypto bot")
	}
	if len(cf.Endpoint) == 0 {
		return nil, errors.New("no endpoint was provided for crypto bot")
	}
	if cf.Client == nil {
		cf.Client = http.DefaultClient
	}

	return &cryptobot{token: cf.Token, endpoint: cf.Endpoint, client: cf.Client}, nil
}

func (cb cryptobot) makeRequest(method, url string, r io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Crypto-Pay-API-Token", cb.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := cb.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HandleUpdate is meant for proccessing webhook update messages. It includes verification of update message integrety.
// You are free to implement your own handler. This is just a minimal implementation.
func (cb cryptobot) HandleUpdate(r *http.Request) (Update, error) {
	sig := r.Header.Get("crypto-pay-api-signature")
	if len(sig) == 0 {
		return Update{}, errors.New("crypto-pay-api-signature header was not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Update{}, fmt.Errorf("failed to read the update body: %w", err)
	}

	hkey := sha256.Sum256([]byte(cb.token))

	h := hmac.New(sha256.New, hkey[:])
	if _, err := h.Write(body); err != nil {
		return Update{}, fmt.Errorf("failed to create a new sha256 hmac: %w", err)
	}

	if sig != fmt.Sprintf("%x", h.Sum(nil)) {
		return Update{}, errors.New("failed to verify the update")
	}

	var u Update

	if err := json.Unmarshal(body, &u); err != nil {
		return Update{}, fmt.Errorf("failed to unmarshal the update: %w", err)
	}

	return u, nil
}

// GetMe returns basic application information. The return of the getMe API method is not documented.
// To mitigate any potential issues GetMe returns raw json.
func (cb cryptobot) GetMe() (json.RawMessage, error) {
	murl, err := url.JoinPath(cb.endpoint, "/getMe")
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("GET", murl, nil)
	if err != nil {
		return nil, err
	}

	var res response[json.RawMessage]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// CreateInvoice takes in a new invoice and returns the invoice on success.
func (cb cryptobot) CreateInvoice(in NewInvoice) (Invoice, error) {
	if err := validateNewInvoice(in); err != nil {
		return Invoice{}, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/createInvoice")
	if err != nil {
		return Invoice{}, err
	}

	data, err := json.Marshal(in)
	if err != nil {
		return Invoice{}, err
	}

	body, err := cb.makeRequest("GET", murl, bytes.NewReader(data))
	if err != nil {
		return Invoice{}, err
	}

	var res response[Invoice]

	if err := json.Unmarshal(body, &res); err != nil {
		return Invoice{}, err
	}

	if !res.Ok {
		return Invoice{}, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// DeleteInvoice takes in the id of the invoice you want to delete. The bool indicates whether the deletion was successful.
func (cb cryptobot) DeleteInvoice(id int64) (bool, error) {
	murl, err := url.JoinPath(cb.endpoint, "/deleteInvoice")
	if err != nil {
		return false, err
	}

	data, err := json.Marshal(struct {
		InvoiceID int64 `json:"invoice_id"`
	}{InvoiceID: id})

	if err != nil {
		return false, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	var res response[bool]

	if err := json.Unmarshal(body, &res); err != nil {
		return false, err
	}

	if !res.Ok {
		return false, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// GetInvoices takes in invoice search options and returns found invoices on success.
func (cb cryptobot) GetInvoices(inop InvoiceOptions) ([]Invoice, error) {
	if err := validateInvoiceOptions(inop); err != nil {
		return nil, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/getInvoices")
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(inop)
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var res response[struct {
		Items []Invoice `json:"items"`
	}]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result.Items, nil
}

// CreateCheck takes in a new check and returns the check on success.
func (cb cryptobot) CreateCheck(nc NewCheck) (Check, error) {
	if err := validateNewCheck(nc); err != nil {
		return Check{}, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/createCheck")
	if err != nil {
		return Check{}, err
	}

	data, err := json.Marshal(nc)
	if err != nil {
		return Check{}, err
	}

	body, err := cb.makeRequest("GET", murl, bytes.NewReader(data))
	if err != nil {
		return Check{}, err
	}

	var res response[Check]

	if err := json.Unmarshal(body, &res); err != nil {
		return Check{}, err
	}

	if !res.Ok {
		return Check{}, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// DeleteCheck takes in the id of the check you want to delete. The bool indicates whether the deletion was successful.
func (cb cryptobot) DeleteCheck(id int64) (bool, error) {
	murl, err := url.JoinPath(cb.endpoint, "/deleteCheck")
	if err != nil {
		return false, err
	}

	data, err := json.Marshal(struct {
		CheckID int64 `json:"check_id"`
	}{CheckID: id})

	if err != nil {
		return false, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	var res response[bool]

	if err := json.Unmarshal(body, &res); err != nil {
		return false, err
	}

	if !res.Ok {
		return false, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// GetChecks takes in check search options and returns found checks on success.
func (cb cryptobot) GetChecks(ckops CheckOptions) ([]Check, error) {
	if err := validateCheckOptions(ckops); err != nil {
		return nil, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/getChecks")
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(ckops)
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var res response[struct {
		Items []Check `json:"items"`
	}]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result.Items, nil
}

// CreateTransfer takes in a new transfer and returns the transfer on success.
func (cb cryptobot) CreateTransfer(nt NewTransfer) (Transfer, error) {
	if err := validateNewTransfer(nt); err != nil {
		return Transfer{}, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/transfer")
	if err != nil {
		return Transfer{}, err
	}

	data, err := json.Marshal(nt)
	if err != nil {
		return Transfer{}, err
	}

	body, err := cb.makeRequest("GET", murl, bytes.NewReader(data))
	if err != nil {
		return Transfer{}, err
	}

	var res response[Transfer]

	if err := json.Unmarshal(body, &res); err != nil {
		return Transfer{}, err
	}

	if !res.Ok {
		return Transfer{}, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// GetTransfers takes in transfer search options and returns found transfers on success.
func (cb cryptobot) GetTransfers(trops TransferOptions) ([]Transfer, error) {
	if err := validateTransferOptions(trops); err != nil {
		return nil, err
	}

	murl, err := url.JoinPath(cb.endpoint, "/getTransfers")
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(trops)
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var res response[struct {
		Items []Transfer `json:"items"`
	}]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result.Items, nil
}

// GetBalance return the current application balance.
func (cb cryptobot) GetBalance() ([]Balance, error) {
	murl, err := url.JoinPath(cb.endpoint, "/getBalance")
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("GET", murl, nil)
	if err != nil {
		return nil, err
	}

	var res response[[]Balance]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// GetExchangeRates return exchange rates of supported currencies.
func (cb cryptobot) GetExchangeRates() ([]ExchangeRate, error) {
	murl, err := url.JoinPath(cb.endpoint, "/getExchangeRates")
	if err != nil {
		return nil, err
	}

	body, err := cb.makeRequest("GET", murl, nil)
	if err != nil {
		return nil, err
	}

	var res response[[]ExchangeRate]

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if !res.Ok {
		return nil, errors.New(string(res.Error))
	}

	return res.Result, nil
}

// GetAppStats takes in application statistics search options and return found application statistics on success.
func (cb cryptobot) GetAppStats(asops AppStatsOptions) (AppStats, error) {
	murl, err := url.JoinPath(cb.endpoint, "/getStats")
	if err != nil {
		return AppStats{}, err
	}

	data, err := json.Marshal(asops)
	if err != nil {
		return AppStats{}, err
	}

	body, err := cb.makeRequest("POST", murl, bytes.NewReader(data))
	if err != nil {
		return AppStats{}, err
	}

	var res response[AppStats]

	if err := json.Unmarshal(body, &res); err != nil {
		return AppStats{}, err
	}

	if !res.Ok {
		return AppStats{}, errors.New(string(res.Error))
	}

	return res.Result, nil
}
