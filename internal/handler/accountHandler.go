package handler

import (
	"avito-tech/internal/models"
	"avito-tech/internal/repository"
	"avito-tech/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	getBalanceByID  = "/get/balance/{id:[0-9]+}"
	currencyBalance = "/get/balance/{currency}/{id:[0-9]+}"

	getAllAccounts = "/get/all"
	transaction    = "/transaction"
	changeBalance  = "/changeBalance"

	apiLink = "http://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%g"
)

var apiKey = os.Getenv("CURRENCY_API_KEY")

type accountHandler struct {
	logger  logger.Logger
	accRepo repository.Account
}

func NewAccountHandler(logger logger.Logger, accRepo repository.Account) *accountHandler {
	return &accountHandler{
		logger:  logger,
		accRepo: accRepo,
	}
}

func (fh *accountHandler) Register(router *mux.Router) {
	router.HandleFunc(getBalanceByID, fh.getBalanceByID).Methods("GET")
	router.HandleFunc(getAllAccounts, fh.getAll).Methods("GET")
	router.HandleFunc(transaction, fh.moneyTransaction).Methods("POST")
	router.HandleFunc(changeBalance, fh.changeBalance).Methods("POST")
	router.HandleFunc(currencyBalance, fh.currencyBalance).Methods("GET")
}

func (fh *accountHandler) getBalanceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while getting id. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while  getting id. err:%s ", err))
		return
	}
	account, err := fh.accRepo.GetBalanceByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while getting account . err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting all account. err:%s ", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)

}

func (fh *accountHandler) getAll(w http.ResponseWriter, r *http.Request) {
	users, err := fh.accRepo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while getting all users . err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting all accounts. err:%s ", err))
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (fh *accountHandler) moneyTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Transaction{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while parsing json. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while parsing json. err:%s ", err))
		return
	}

	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while validating data. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while validating data. err:%s ", err))
		return
	}

	err = fh.accRepo.MoneyTransaction(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while money transaction between users. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while money transaction between users. err:%s ", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Transaction was successful")
}

func (fh *accountHandler) changeBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.AccountDebit{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while parsing json. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while parsing json. err:%s ", err))
		return
	}
	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while validating data. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while validating data. err:%s ", err))
		return
	}

	err = fh.accRepo.ChangeBalance(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while changing balance. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while changing balance. err:%s ", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balance succsefully changed")
}

func (fh *accountHandler) currencyBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while getting id. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting id. err:%s ", err))
		return
	}

	balance, err := fh.accRepo.GetBalanceByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while getting account. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting all account. err:%s ", err))
		return
	}

	fh.convert(w, vars["currency"], balance.Balance)
}

func (fh *accountHandler) convert(w http.ResponseWriter, currency string, balance float64) {

	url := fmt.Sprintf(apiLink, currency, balance)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while making request.err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while  making request.err:%s ", err))
		return
	}
	request.Header.Set("apikey", apiKey)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.logger.Errorf("error occurred while getting exchange rates.err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting exchange rates.err:%s ", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fh.logger.Errorf("error occurred while getting exchange rates.err:%s ", err)
			json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting exchange rates.err:%s ", err))
			return
		}
		fh.logger.Errorf("error occurred while getting exchange rates.err:%s ", string(b))
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting exchange rates.err:%s ", string(b)))
		return
	}

	exchange := &models.ExchangeResponse{}

	err = json.NewDecoder(resp.Body).Decode(exchange)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.logger.Errorf("error occurred while parsing json. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while parsing json. err:%s ", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exchange)

}
