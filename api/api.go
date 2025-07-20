package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/PeterBocan/p-demo/ledger"
	"github.com/PeterBocan/p-demo/models"
	"github.com/gorilla/mux"
)

type Server struct {
	handler http.Handler
	ledger  *ledger.Ledger

	accounts map[int]models.Account
	accMux   sync.RWMutex
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.handler.ServeHTTP(w, req)
}

// NewServer returns a new instance of Server, which handles HTTP requests.
func NewServer(ledger *ledger.Ledger) *Server {
	r := mux.NewRouter()
	server := &Server{handler: r, ledger: ledger, accounts: map[int]models.Account{}}
	r.HandleFunc("/accounts", server.CreateAccount).Methods(http.MethodPost)
	r.HandleFunc("/accounts/{id}", server.GetAccountInformation).Methods(http.MethodGet)
	r.HandleFunc("/transactions", server.CreateTransaction).Methods(http.MethodPost)

	return server
}

func (s *Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("malformed request received - connection closed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var account models.Account
	err = json.Unmarshal(bytes, &account)
	if err != nil {
		log.Print("malformed request received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	s.accMux.RLock()
	id := len(s.accounts)
	s.accMux.RUnlock()
	account.ID = id + 1

	err = s.ledger.AddAccount(account.ID)
	if errors.Is(err, ledger.ErrAccountAlreadyExistInLedger) {
		log.Print("account with this ID already exists")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`account already exists`))
		return
	}

	s.addAccount(account)

	bytes, err = json.Marshal(account)
	if err != nil {
		log.Print("internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
		return
	}
	w.Write(bytes)
}

func (s *Server) GetAccountInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]
	if accountID == "" {
		log.Print("malformed request received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	accountNumber, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		log.Print("malformed request received - invalid account number")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	s.accMux.RLock()
	account, ok := s.accounts[int(accountNumber)]
	s.accMux.RUnlock()

	if !ok {
		log.Print("account does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	bytes, err := json.Marshal(account)
	if err != nil {
		log.Print("could not marshal account into bytes")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
		return
	}

	w.Write(bytes)
}

func (s *Server) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("malformed request received - connection closed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tx models.Transaction
	err = json.Unmarshal(bytes, &tx)
	if err != nil {
		log.Print("malformed request received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	if tx.AccountID == 0 {
		log.Print("malformed request received - not valid account ID")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	if err := tx.Type.Valid(); err != nil {
		log.Print("malformed request received - not valid operation type")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}

	err = s.ledger.Transact(tx)
	if err != nil {
		log.Printf("could not process the transaction: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid payload provided`))
		return
	}
	w.Write(bytes)
}

func (s *Server) addAccount(account models.Account) {
	s.accMux.Lock()
	defer s.accMux.Unlock()
	s.accounts[account.ID] = account
}
