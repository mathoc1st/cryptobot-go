package cryptobot

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"
)

const (
	testToken = "API_TOKEN"
)

var cbot Client

func TestMain(m *testing.M) {
	fmt.Println("Initializing a new crypto bot instance...")

	cb, err := NewClient(Config{
		Token:    testToken,
		Endpoint: Testnet,
	})

	if err != nil {
		fmt.Println("Initialization failed: ", err)
		os.Exit(1)
	}

	cbot = cb

	code := m.Run()

	os.Exit(code)
}

func TestGetMe(t *testing.T) {
	_, err := cbot.GetMe()
	if err != nil {
		t.Error(err)
	}
}

func TestInvoice(t *testing.T) {
	tdata := []struct {
		name  string
		input NewInvoice
	}{
		{
			name: "crypto",
			input: NewInvoice{
				CurrencyType:   Crypto,
				CryptoAsset:    USDT,
				Amount:         "5",
				Description:    "Test",
				HiddenMessage:  "Test",
				PaidBtnName:    ViewItem,
				PaidBtnUrl:     "https://google.com",
				Payload:        "Test",
				AllowComments:  true,
				AllowAnonymous: false,
				ExpiresIn:      42069,
			},
		},
		{
			name: "fiat",
			input: NewInvoice{
				CurrencyType:         Fiat,
				Fiat:                 EUR,
				AcceptedCryptoAssets: []CryptoAsset{TON},
				Amount:               "4",
				Description:          "Test",
				HiddenMessage:        "Test",
				PaidBtnName:          OpenChannel,
				PaidBtnUrl:           "https://google.com",
				Payload:              "Test",
				AllowComments:        false,
				AllowAnonymous:       true,
				ExpiresIn:            42069,
			},
		},
	}

	for _, test := range tdata {
		var in Invoice
		t.Run(fmt.Sprintf("creating new %s invoice", test.name), func(t *testing.T) {
			got, err := cbot.CreateInvoice(test.input)
			if err != nil {
				t.Fatal(err)
			}
			assertInvoices(t, test.input, got)
			in = got
		})

		t.Run(fmt.Sprintf("getting %s invoice", test.name), func(t *testing.T) {
			_, err := cbot.GetInvoices(InvoiceOptions{InvoiceIDs: []int64{in.ID}})
			if err != nil {
				t.Error(err)
			}
		})

		t.Run(fmt.Sprintf("deleting %s invoice", test.name), func(t *testing.T) {
			_, err := cbot.DeleteInvoice(in.ID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tdata := []struct {
		input NewCheck
	}{
		{
			input: NewCheck{
				CryptoAsset: TON,
				Amount:      "0.01",
				PinToUserID: 123123,
			},
		},
		{
			input: NewCheck{
				CryptoAsset:   TON,
				Amount:        "0.01",
				PinToUsername: "user",
			},
		},
	}

	for _, test := range tdata {
		var ch Check
		t.Run("creating a new check", func(t *testing.T) {
			got, err := cbot.CreateCheck(test.input)
			if err != nil {
				t.Fatal(err)
			}
			assertChecks(t, test.input, got)
			ch = got
		})

		t.Run("getting the check", func(t *testing.T) {
			_, err := cbot.GetChecks(CheckOptions{CheckIDs: []int64{ch.ID}})
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("deleting the check", func(t *testing.T) {
			_, err := cbot.DeleteCheck(ch.ID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	id, err := rand64CharHex()
	if err != nil {
		t.Error("failed to generate a random SpendID: ", err)
	}

	tdata := []struct {
		input NewTransfer
	}{
		{
			input: NewTransfer{
				UserID:      1844235715,
				SpendID:     id,
				CryptoAsset: TON,
				Amount:      "0.35",
			},
		},
	}

	for _, test := range tdata {
		var tr Transfer
		t.Run("creating a new transfer", func(t *testing.T) {
			got, err := cbot.CreateTransfer(test.input)
			if err != nil {
				t.Fatal(err)
			}
			assertTransfers(t, test.input, got)
			tr = got
		})

		t.Run("getting the transfer", func(t *testing.T) {
			_, err = cbot.GetTransfers(TransferOptions{TransferIDs: []int64{tr.ID}})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBalance(t *testing.T) {
	_, err := cbot.GetBalance()
	if err != nil {
		t.Error("failed to get the balance: ", err)
	}
}

func TestRates(t *testing.T) {
	_, err := cbot.GetExchangeRates()
	if err != nil {
		t.Error("failed to get the exchange rates: ", err)
	}
}

func TestStats(t *testing.T) {
	statops := AppStatsOptions{
		StartAt: time.Now().Add(-12 * time.Hour),
		EndAt:   time.Now(),
	}

	_, err := cbot.GetAppStats(statops)
	if err != nil {
		t.Error("failed to get the application statistics: ", err)
	}
}

func assertTransfers(t *testing.T, want NewTransfer, got Transfer) {
	t.Helper()

	var errs []error

	if want.UserID != got.UserID {
		errs = append(errs, fmt.Errorf("got user id %d, want %d", got.UserID, want.UserID))
	}
	if want.SpendID != got.SpendID {
		errs = append(errs, fmt.Errorf("got spend id %s, want %s", got.SpendID, want.SpendID))
	}
	if want.CryptoAsset != got.CryptoAsset {
		errs = append(errs, fmt.Errorf("got asset %s, want %s", got.CryptoAsset, want.CryptoAsset))
	}
	if want.Comment != got.Comment {
		errs = append(errs, fmt.Errorf("got comment %s, want %s", got.Comment, want.Comment))
	}

	if len(errs) == 0 {
		return
	}

	t.Error(errors.Join(errs...))
}

func assertInvoices(t *testing.T, want NewInvoice, got Invoice) {
	t.Helper()

	var errs []error
	if want.CurrencyType != got.CurrencyType {
		errs = append(errs, fmt.Errorf("got currency type %s, want %s", got.CurrencyType, want.CurrencyType))
	}
	if want.CryptoAsset != got.CryptoAsset {
		errs = append(errs, fmt.Errorf("got asset %s, want %s", got.CryptoAsset, want.CryptoAsset))
	}
	if want.Amount != got.Amount {
		errs = append(errs, fmt.Errorf("got amount %s, want %s", got.Amount, want.Amount))
	}
	if want.Fiat != got.Fiat {
		errs = append(errs, fmt.Errorf("got fiat %s, want %s", got.Fiat, want.Fiat))
	}
	if want.AcceptedCryptoAssets != nil && !slices.Equal(want.AcceptedCryptoAssets, got.AcceptedCryptoAssets) {
		errs = append(errs, fmt.Errorf("got AcceptedAssets %s, want %s", got.Fiat, want.Fiat))
	}
	if want.Description != got.Description {
		errs = append(errs, fmt.Errorf("got description %s, want %s", got.Description, want.Description))
	}
	if want.HiddenMessage != got.HiddenMessage {
		errs = append(errs, fmt.Errorf("got hidden message %s, want %s", got.HiddenMessage, want.HiddenMessage))
	}
	if want.PaidBtnName != got.PaidBtnName {
		errs = append(errs, fmt.Errorf("got paid button name %s, want %s", got.PaidBtnName, want.PaidBtnName))
	}
	if want.PaidBtnUrl != got.PaidBtnUrl {
		errs = append(errs, fmt.Errorf("got paid button url %s, want %s", got.PaidBtnUrl, want.PaidBtnUrl))
	}
	if want.Payload != got.Payload {
		errs = append(errs, fmt.Errorf("got payload %s, want %s", got.Payload, want.Payload))
	}
	if want.AllowAnonymous != got.AllowAnonymous {
		errs = append(errs, fmt.Errorf("got allow anonymous %v, want %v", got.AllowAnonymous, want.AllowAnonymous))
	}
	if want.AllowComments != got.AllowComments {
		errs = append(errs, fmt.Errorf("got allow comments %v, want %v", got.AllowComments, want.AllowComments))
	}

	if len(errs) == 0 {
		return
	}

	t.Error(errors.Join(errs...))
}

func assertChecks(t *testing.T, want NewCheck, got Check) {
	t.Helper()

	var errs []error

	if want.CryptoAsset != got.CryptoAsset {
		errs = append(errs, fmt.Errorf("got asset %s, want %s", got.CryptoAsset, want.CryptoAsset))
	}
	if want.Amount != got.Amount {
		errs = append(errs, fmt.Errorf("got amount %s, want %s", got.CryptoAsset, want.CryptoAsset))
	}

	if len(errs) == 0 {
		return
	}

	t.Error(errors.Join(errs...))
}

func rand64CharHex() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
