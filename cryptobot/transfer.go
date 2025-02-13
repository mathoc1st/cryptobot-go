package cryptobot

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type TransferStatus string

const (
	TransferCompleted TransferStatus = "completed"
)

type Transfer struct {
	// Unique ID for this transfer.
	ID int64 `json:"transfer_id"`

	// Unique UTF-8 string.
	SpendID string `json:"spend_id"`

	// Telegram user id the transfer was sent to.
	UserID int64 `json:"user_id"`

	// Type of cryptocurrency.
	CryptoAsset CryptoAsset `json:"asset"`

	// Amount of the transfer.
	Amount string `json:"amount"`

	// Transfer status.
	Status TransferStatus `json:"status"`

	// Date the transfer was completed (ISO 8601 format).
	CompletedAt string `json:"completed_at"`

	// Optional. Comment for this transfer.
	Comment string `json:"comment,omitempty"`
}

type NewTransfer struct {
	// Telegram user id the transfer will be sent to.
	UserID int64 `json:"user_id"`

	// Type of cryptocurrency.
	CryptoAsset CryptoAsset `json:"asset"`

	// Amount of the transfer. The minimum and maximum limits for each supported cryptocurrency are roughly $1–$25,000 USD.
	Amount string `json:"amount"`

	// Random UTF-8 string. Shoud be unique for every transfer for idempotent requests. 64 characters max.
	SpendID string `json:"spend_id"`

	// Optional. Comment for the transfer. Users will see this comment in the notification about the transfer. 1024 characters max.
	Comment string `json:"comment,omitempty"`

	// Optional. Weither or not to notify the user about the transfer. Defaults to false.
	DisableSendNotification bool `json:"disable_send_notification"`
}

type TransferOptions struct {
	// Optiona. Type of cryptocurrency to search by.
	CryptoAsset CryptoAsset

	// Optional. Transfer ids to find.
	TransferIDs []int64

	// Optional. Unique UTF-8 transfer string to search by.
	SpendID string

	// Optional. Defaults to 0.
	Offset int64

	// Optional. Number of transfers to be returned. Values between 1-1000 are accepted. Defaults to 100.
	Count int64
}

type tempTrOps struct {
	CryptoAsset string `json:"asset,omitempty"`
	TransferIDs string `json:"transfer_ids,omitempty"`
	SpendID     string `json:"spend_id,omitempty"`
	Offset      int64  `json:"offset,omitempty"`
	Count       int64  `json:"count,omitempty"`
}

func (to TransferOptions) MarshalJSON() ([]byte, error) {
	ids := make([]string, len(to.TransferIDs))

	for _, id := range to.TransferIDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}

	return json.Marshal(&tempTrOps{
		CryptoAsset: string(to.CryptoAsset),
		TransferIDs: strings.Join(ids, ","),
		SpendID:     to.SpendID,
		Offset:      to.Offset,
		Count:       to.Count,
	})
}

func validateNewTransfer(nt NewTransfer) error {
	var errs []error

	if len(nt.CryptoAsset) == 0 {
		errs = append(errs, errors.New("CryptoAsset cannot be empty"))
	}
	if len(nt.SpendID) == 0 {
		errs = append(errs, errors.New("SpendID cannot be empty"))
	}
	if len(nt.SpendID) > 64 {
		errs = append(errs, errors.New("SpendID cannot exceed 64 characters"))
	}
	if len(nt.Comment) > 1024 {
		errs = append(errs, errors.New("Comment cannot exceed 1024 characters"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func validateTransferOptions(trops TransferOptions) error {
	var errs []error

	if len(trops.SpendID) > 64 {
		errs = append(errs, errors.New("SpendID cannot exceed 64 characters"))
	}
	if trops.Offset < 0 {
		errs = append(errs, errors.New("Offset cannot be less than 0"))
	}
	if trops.Count != 0 && (trops.Count < 1 || trops.Count > 1000) {
		errs = append(errs, errors.New("Count needs to be within 1-1000 record range"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
