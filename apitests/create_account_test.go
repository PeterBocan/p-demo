package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PeterBocan/p-demo/api"
	"github.com/PeterBocan/p-demo/ledger"
	"github.com/PeterBocan/p-demo/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAccountSucceeds(t *testing.T) {
	l := ledger.New()
	server := api.NewServer(l)

	body, err := json.Marshal(models.Account{DocumentNumber: "1234"})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))

	server.ServeHTTP(w, request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatingAccountAndFetchingAccountInformationSucceeds(t *testing.T) {
	l := ledger.New()
	server := api.NewServer(l)

	body, err := json.Marshal(models.Account{DocumentNumber: "1234"})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))

	server.ServeHTTP(w, request)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := httptest.NewRecorder()
	request2 := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)

	server.ServeHTTP(w2, request2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var resp models.Account
	err = json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 1, resp.ID)
}

func TestCreatingAccountAndDepositMoney(t *testing.T) {
	l := ledger.New()
	server := api.NewServer(l)

	body, err := json.Marshal(models.Account{DocumentNumber: "1234"})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))

	server.ServeHTTP(w, request)
	assert.Equal(t, http.StatusOK, w.Code)

	tx, err := json.Marshal(models.Transaction{ID: 12533, AccountID: 1, Type: models.CreditVoucher, Amount: 100.0, EventDate: time.Now().UTC()})
	assert.NoError(t, err)
	w2 := httptest.NewRecorder()
	request2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(tx))

	server.ServeHTTP(w2, request2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var resp models.Transaction
	err = json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 12533, resp.ID)
	assert.Equal(t, 1, resp.AccountID)
}
