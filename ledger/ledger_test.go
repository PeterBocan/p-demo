package ledger

import (
	"testing"

	"github.com/PeterBocan/p-demo/models"
	"github.com/stretchr/testify/assert"
)

func TestLedgerAddAccount(t *testing.T) {
	l := New()

	assert.NoError(t, l.AddAccount(1234))
}

func TestLedgerTransactionFailsWithNoAccount(t *testing.T) {
	l := New()

	tx := models.Transaction{ID: 1234, AccountID: 123445678, Type: models.NormalPurchase}
	err := l.Transact(tx)
	assert.ErrorIs(t, err, ErrAccountDoesNotExistInLedger)
}

func TestLedgerTransactionFailsNotEnoughFunds(t *testing.T) {
	l := New()

	accountID := 1245

	depositTx := models.Transaction{ID: 72434286, AccountID: accountID, Amount: 100.0, Type: models.CreditVoucher}

	assert.NoError(t, l.AddAccount(accountID))
	assert.NoError(t, l.Transact(depositTx))

	balance, err := l.Balance(accountID)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)

	tx := models.Transaction{ID: 1234, AccountID: accountID, Amount: 562342.0, Type: models.NormalPurchase}
	assert.ErrorIs(t, l.Transact(tx), ErrNotEnoughFunds)
}
