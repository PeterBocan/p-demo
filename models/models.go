package models

import (
	"errors"
	"time"
)

type Account struct {
	ID             int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type OperationType int

// Description returns a string representation of the type of a operation.
func (t OperationType) Description() string {
	switch t {
	case NormalPurchase:
		return "Normal Purchase"
	case PurchaseWithInstallments:
		return "Purchase With Installments"
	case Withdrawal:
		return "Withdrawal"
	case CreditVoucher:
		return "Credit Voucher"

	default:
		return "Invalid Operation"
	}
}

// Validates the type of a variable.
func (t OperationType) Valid() error {
	switch t {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal, CreditVoucher:
		return nil
	default:
		return errors.New("invalid transaction type")
	}
}

const (
	InvalidOperation OperationType = iota
	NormalPurchase
	PurchaseWithInstallments
	Withdrawal
	CreditVoucher
)

type Transaction struct {
	ID        int           `json:"transaction_id"`
	AccountID int           `json:"account_id"`
	Type      OperationType `json:"operation_type_id"`
	Amount    float64       `json:"amount"`
	EventDate time.Time     `json:"event_date"`
}
