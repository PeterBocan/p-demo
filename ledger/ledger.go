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

		l.txMux.RLock()
		accountTransactions := l.transactions[t.AccountID]
		l.txMux.RUnlock()

		creditAmount := t.Amount
		for idx := range accountTransactions {
			if accountTransactions[idx].Balance < 0 {
				amount := accountTransactions[idx].Balance

				if creditAmount+amount > 0.0 {
					accountTransactions[idx].Balance = 0.0
					creditAmount += amount
					continue
				}
				accountTransactions[idx].Balance = creditAmount + amount
				creditAmount = 0.0
				continue
			}
		}
		// the balance of the credit deposit is the remainder of the money left after settling the credit transactions
		t.Balance = creditAmount

	case models.NormalPurchase, models.PurchaseWithInstallments, models.Withdrawal:
		l.balanceMut.Lock()
		l.balances[t.AccountID] = balance + t.Amount
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

// GetTransactionHistory
func (l *Ledger) GetTransactionHistory(account int) ([]models.Transaction, error) {
	l.txMux.RLock()
	transactionHistory, ok := l.transactions[account]
	l.txMux.RUnlock()

	if !ok {
		return nil, ErrAccountDoesNotExistInLedger
	}

	return transactionHistory, nil
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
