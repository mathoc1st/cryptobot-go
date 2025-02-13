package cryptobot

type updateType string

const (
	updateInvoicePaid = "invoice_paid"
)

type Update struct {
	// Non-unique update ID.
	ID int64 `json:"update_id"`

	// Webhook update type.
	Type updateType `json:"update_type"`

	// Date the request was sent (ISO 8601 format).
	RequestDate string  `json:"request_date"`
	Payload     Invoice `json:"payload"`
}
