package cryptobot

import (
	"encoding/json"
	"time"
)

type AppStats struct {
	// Total volume of paid invoices in USD.
	Volume int64 `json:"volume"`

	// Conversion of all created invoices.
	Conversion int64 `json:"conversion"`

	// The unique number of users who have paid the invoice.
	UniqueUsers int64 `json:"unique_users_count"`

	// Total created invoice count.
	CreatedInvoices int64 `json:"created_invoice_count"`

	// Total paid invoice count.
	PaidInvoices int64 `json:"paid_invoice_count"`

	// The date on which the statistics calculation was started (ISO 8601 format).
	StartAt string `json:"start_at"`

	// The date on which the statistics calculation was ended (ISO 8601 format).
	EndAt string `json:"end_at"`
}

type AppStatsOptions struct {
	// Optional. Start date. Defaults last 24 hours.
	StartAt time.Time

	// Optional. End data. Defaults to current date.
	EndAt time.Time
}

func (aso AppStatsOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
	}{
		StartAt: aso.StartAt.Format(time.RFC3339),
		EndAt:   aso.EndAt.Format(time.RFC3339),
	})
}
