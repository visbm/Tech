package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Account struct {
	ID      int     `json:"account_id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ReceiverID int     `json:"receiver_id"`
	SenderID   int     `json:"sender_id"`
	Amount     float64 `json:"amount"`
	Comment    string  `json:"comment"`
}

type AccountDebit struct {
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
	Comment   string  `json:"comment"`
}

type ExchangeResponse struct {
	Result float64 `json:"result"`
}

type TransactionHistory struct {
	TransactionID int       `json:"transaction_ID"`
	AccountID     int       `json:"account_id"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Comment       string    `json:"comment"`
}

type TransactionHistoryReq struct {
	AccountID int    `json:"account_id"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	OrderBy   string `json:"order_by"`
}

func (t *Transaction) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.ReceiverID, validation.Required, validation.Min(1)),
		validation.Field(&t.SenderID, validation.Required, validation.Min(1), validation.NotIn(t.ReceiverID)),
		validation.Field(&t.Amount, validation.Required, validation.Min(float64(1))),
		validation.Field(&t.Comment, validation.Required, validation.Length(5, 50)),
	)
}

func (a *AccountDebit) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.AccountID, validation.Required, validation.Min(1)),
		validation.Field(&a.Amount, validation.Required, validation.NotIn(float64(0))),
		validation.Field(&a.Comment, validation.Required, validation.Length(5, 50)),
	)
}

func (tr *TransactionHistoryReq) Validate() error {
	return validation.ValidateStruct(
		tr,
		validation.Field(&tr.AccountID, validation.Required, validation.Min(1)),
		validation.Field(&tr.OrderBy, validation.In("date_time", "amount",
			"date_time ASC", "amount ASC",
			"date_time DESC", "amount DESC",
			"date_time asc", "amount asc",
			"date_time desc", "amount desc")),
	)
}
