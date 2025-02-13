# CryptoBot Go

A minimal implementaion of the [@CryptoBot](https://t.me/CryptoBot) API.

## Installation

```sh
go get github.com/mathoc1st/cryptobot-go
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mathoc1st/cryptobot-go/cryptobot"
)

func main() {
	// Initialize the Cryptobot client with your configuration
	cb, err := cryptobot.New(cryptobot.Config{
		Token:    "<API-TOKEN>", // Replace with your actual token
		Endpoint: cryptobot.Mainnet,
	})
	if err != nil {
		panic(err) // Handle error gracefully in production
	}

	// Fetch basic information
	info, err := cb.GetMe()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(info)
}
```

## Examples

### Creating a new invoice

```go
	in, err := cb.CreateInvoice(cryptobot.NewInvoice{
		CurrencyType:         cryptobot.Fiat,
		Fiat:                 cryptobot.USD,
		AcceptedCryptoAssets: []cryptobot.CryptoAsset{cryptobot.USDC, cryptobot.TON},
		Amount:               "50",
		Description:          "Hello world",
		HiddenMessage:        "Hello",
		PaidBtnName:          cryptobot.ViewItem,
		PaidBtnUrl:           "https://example.com",
		Payload:              "World",
		AllowComments:        true,
		AllowAnonymous:       false,
		ExpiresIn:            600,
	})
```

### Creating a new check

```go
	ch, err := cb.CreateCheck(cryptobot.NewCheck{
		CryptoAsset: cryptobot.USDT,
		Amount: "50",
	})
```

### Creating a new transfer 

```go
	tr, err := cb.CreateTransfer(cryptobot.NewTransfer{
		UserID: 123123123,
		CryptoAsset: cryptobot.USDT,
		Amount: "50",
		SpendID: "randomid",
		Comment: "Hello World",
	})
```

### Retrieving invoices

```go
	ins, err := cb.GetInvoices(cryptobot.InvoiceOptions{
		CryptoAsset: cryptobot.USDT,
		Status:      cryptobot.InvoiceActive,
	})
```

### Retrieving checks

```go
	chs, err := cb.GetChecks(cryptobot.CheckOptions{
		CryptoAsset: cryptobot.USDT,
		Status:      cryptobot.CheckActivated,
	})
```

### Retrieving transfers

```go
	trs, err := cb.GetTransfers(cryptobot.TransferOptions{
		CryptoAsset: cryptobot.USDT,
	})
```


### Retrieving application balance

```go
	b, err := cb.GetBalance()
```

### Retrieving exchange rates

```go
	rate, err := cb.GetExchangeRates()
```

### Retrieving application statistics

```go
	st, err := cb.GetAppStats(cryptobot.AppStatsOptions{
		StartAt: time.Now().Add(-12 * time.Hour),
		EndAt:   time.Now(),
	})
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
