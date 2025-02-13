package cryptobot

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type CheckStatus string

const (
	CheckActive    CheckStatus = "active"
	CheckActivated             = "activated"
)

type Check struct {
	// Unique ID for this check.
	ID int64 `json:"check_id"`

	// Hash of the check.
	Hash string `json:"hash"`

	// Type of cryptocurrency.
	CryptoAsset CryptoAsset `json:"asset"`

	// Amount of the check.
	Amount string `json:"amount"`

	// URL for the user to activate the check.
	BotCheckURL string `json:"bot_check_url"`

	// Check status.
	Status CheckStatus `json:"status"`

	// Date the check was created (ISO 8601 format).
	CreatedAt string `json:"created_at"`

	// Date the check was activated (ISO 8601 format).
	ActivatedAt string `json:"activated_at"`
}

type NewCheck struct {
	// Type of cryptocurrency.
	CryptoAsset CryptoAsset `json:"asset"`

	// Amount of the check.
	Amount string `json:"amount"`

	// Optional. Telegram id of the user who will be able to activate the check.
	PinToUserID int64 `json:"pin_to_user_id,omitempty"`

	// Optional. Telegram user name who will be able to activate the check.
	PinToUsername string `json:"pin_to_username,omitempty"`
}

type CheckOptions struct {
	// Optional. Type of cryptocurrency to search by.
	CryptoAsset CryptoAsset `json:"asset,omitempty"`

	// Optional. Check ids to find.
	CheckIDs []int64 `json:"check_ids,omitempty"`

	// Optional. Status to search by.
	Status CheckStatus `json:"status,omitempty"`

	// Optional. Defaults to 0.
	Offset int64 `json:"offset,omitempty"`

	// Optional. Number of checks to be returned. Values between 1-1000 are accepted. Defaults to 100.
	Count int64 `json:"count,omitempty"`
}

type tempCheckOps struct {
	CryptoAsset string `json:"asset,omitempty"`
	CheckIDs    string `json:"check_ids,omitempty"`
	Status      string `json:"status,omitempty"`
	Offset      int64  `json:"offset,omitempty"`
	Count       int64  `json:"count,omitempty"`
}

func (co CheckOptions) MarshalJSON() ([]byte, error) {
	ids := make([]string, len(co.CheckIDs))

	for _, id := range co.CheckIDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}

	return json.Marshal(&tempCheckOps{
		CryptoAsset: string(co.CryptoAsset),
		CheckIDs:    strings.Join(ids, ","),
		Status:      string(co.Status),
		Offset:      co.Offset,
		Count:       co.Count,
	})
}

func validateNewCheck(nc NewCheck) error {
	var errs []error

	if len(nc.CryptoAsset) == 0 {
		errs = append(errs, errors.New("CryptoAsset cannot be empty"))
	}
	if len(nc.Amount) == 0 {
		errs = append(errs, errors.New("Amount cannot be empty"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func validateCheckOptions(ckops CheckOptions) error {
	var errs []error

	if ckops.Offset < 0 {
		errs = append(errs, errors.New("Offset cannot be less than 0"))
	}
	if ckops.Count != 0 && (ckops.Count < 1 || ckops.Count > 1000) {
		errs = append(errs, errors.New("Count needs to be within 1-1000 record range"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
