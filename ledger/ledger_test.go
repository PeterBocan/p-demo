package ledger

import (
	"testing"

	"github.com/PeterBocan/p-demo/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLedgerSettlesFirstTransaction(t *testing.T) {
	l := New()

	accountID := 1245
	assert.NoError(t, l.AddAccount(accountID))
	wTx := models.Transaction{ID: 1, AccountID: accountID, Amount: -50.0, Balance: -50.0, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(wTx))

	depositTx := models.Transaction{ID: 2, AccountID: accountID, Amount: 40.0, Type: models.CreditVoucher}
	assert.NoError(t, l.Transact(depositTx))

	balance, err := l.Balance(accountID)
	assert.NoError(t, err)
	assert.Equal(t, -10.0, balance)

	transactions, err := l.GetTransactionHistory(accountID)
	assert.NoError(t, err)

	require.Len(t, transactions, 2)

	withdrawalTx := transactions[0]
	assert.Equal(t, -10.0, withdrawalTx.Balance)
}

func TestLedgerSettleBalancesTwoTransactions(t *testing.T) {
	l := New()

	accountID := 1245
	assert.NoError(t, l.AddAccount(accountID))
	wTx := models.Transaction{ID: 1, AccountID: accountID, Amount: -50.0, Balance: -50.0, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(wTx))

	w2Tx := models.Transaction{ID: 2, AccountID: accountID, Amount: -20.0, Balance: -20.0, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(w2Tx))

	depositTx := models.Transaction{ID: 3, AccountID: accountID, Amount: 80.0, Type: models.CreditVoucher}
	assert.NoError(t, l.Transact(depositTx))

	balance, err := l.Balance(accountID)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, balance)

	transactions, err := l.GetTransactionHistory(accountID)
	assert.NoError(t, err)

	require.Len(t, transactions, 3)

	withdrawalTx := transactions[0]
	assert.Equal(t, 0.0, withdrawalTx.Balance)

	withdrawal2Tx := transactions[1]
	assert.Equal(t, 0.0, withdrawal2Tx.Balance)

	depositHistoryTx := transactions[2]
	assert.Equal(t, 10.0, depositHistoryTx.Balance)
}

func TestLedgerSettleBalancesFourTransactions(t *testing.T) {
	l := New()

	accountID := 1245
	assert.NoError(t, l.AddAccount(accountID))
	wTx := models.Transaction{ID: 1, AccountID: accountID, Amount: -50.0, Balance: -50.0, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(wTx))

	w2Tx := models.Transaction{ID: 2, AccountID: accountID, Amount: -23.5, Balance: -23.5, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(w2Tx))

	w3Tx := models.Transaction{ID: 3, AccountID: accountID, Amount: -18.7, Balance: -18.7, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(w3Tx))
	w4Tx := models.Transaction{ID: 4, AccountID: accountID, Amount: 60.0, Balance: 60.0, Type: models.CreditVoucher}
	assert.NoError(t, l.Transact(w4Tx))
	w5Tx := models.Transaction{ID: 5, AccountID: accountID, Amount: 100.0, Balance: 100.0, Type: models.CreditVoucher}
	assert.NoError(t, l.Transact(w5Tx))

	depositTx := models.Transaction{ID: 6, AccountID: accountID, Amount: -50.0, Balance: -50.0, Type: models.Withdrawal}
	assert.NoError(t, l.Transact(depositTx))

	// balance, err := l.Balance(accountID)
	// assert.NoError(t, err)
	// assert.Equal(t, 17.0, balance)

	transactions, err := l.GetTransactionHistory(accountID)
	assert.NoError(t, err)

	require.Len(t, transactions, 6)

	withdrawalTx := transactions[0]
	assert.Equal(t, 0.0, withdrawalTx.Balance)

	withdrawal2Tx := transactions[1]
	assert.Equal(t, 0.0, withdrawal2Tx.Balance)

	withdrawal3Tx := transactions[2]
	assert.Equal(t, 0.0, withdrawal3Tx.Balance)

	deposit1Tx := transactions[3]
	assert.Equal(t, 0.0, deposit1Tx.Balance)

	deposit2Tx := transactions[4]
	assert.Equal(t, 67.8, deposit2Tx.Balance)

	withdrawal4Tx := transactions[5]
	assert.Equal(t, -50.0, withdrawal4Tx.Balance)
}
