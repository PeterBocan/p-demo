package ledger

import (
	"errors"
	"sync"

	"github.com/PeterBocan/p-demo/models"
)

// ErrAccountDoesNotExistInLedger returns from the transaction processing, if the account has not been
// created within the ledger
var ErrAccountDoesNotExistInLedger = errors.New("the account does not exist in the ledger")

// ErrAccountAlreadyExistInLedger returns that the accouunt with the specified ID already exists.
var ErrAccountAlreadyExistInLedger = errors.New("the account does exist in the ledger already")

// ErrNotEnoughFunds returns from the transaction processing if the current balance is not enough.
var ErrNotEnoughFunds = errors.New("the account does not have enough funds")

// Ledger keeps the current balance associated with the account (the key of the map)
// and the associated amount of money with it. The ledger currently doesn't hold any history.
type Ledger struct {
	balances   map[int]float64
	balanceMut sync.RWMutex

	transactions map[int][]models.Transaction
	txMux        sync.RWMutex
}

// New returns a new instance of an empty ledger
func New() *Ledger {
	return &Ledger{
		balances:     map[int]float64{},
		transactions: map[int][]models.Transaction{},
	}
}

// Transact adds or subtracts the amount to the account
func (l *Ledger) Transact(t models.Transaction) error {
	l.balanceMut.RLock()
	balance, ok := l.balances[t.AccountID]
	l.balanceMut.RUnlock()

	if !ok {
		return ErrAccountDoesNotExistInLedger
	}

	switch t.Type {
	case models.CreditVoucher: // execute deposit
		l.balanceMut.Lock()
		l.balances[t.AccountID] = balance + t.Amount
		l.balanceMut.Unlock()

	case models.NormalPurchase, models.PurchaseWithInstallments, models.Withdrawal:
		if balance < t.Amount {
			return ErrNotEnoughFunds
		}

		l.balanceMut.Lock()
		l.balances[t.AccountID] = balance - t.Amount
		l.balanceMut.Unlock()

	default:
		return errors.New("invalid transaction received")
	}

	l.txMux.Lock()
	l.transactions[t.AccountID] = append(l.transactions[t.AccountID], t)
	l.txMux.Unlock()

	return nil
}

// AddAccount adds an account to the ledger. This method must be called before you can transact
// with this account.
func (l *Ledger) AddAccount(account int) error {
	l.balanceMut.RLock()
	_, ok := l.balances[account]
	l.balanceMut.RUnlock()

	if ok {
		return ErrAccountAlreadyExistInLedger
	}

	l.balanceMut.Lock()
	l.balances[account] = 0
	l.balanceMut.Unlock()
	l.txMux.Lock()
	l.transactions[account] = []models.Transaction{}
	l.txMux.Unlock()

	return nil
}

// Balance returns current balance associated with the account
func (l *Ledger) Balance(account int) (float64, error) {
	l.balanceMut.RLock()
	balance, ok := l.balances[account]
	l.balanceMut.RUnlock()

	if !ok {
		return 0, ErrAccountDoesNotExistInLedger
	}

	return balance, nil
}
