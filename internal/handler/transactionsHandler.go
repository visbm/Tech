package handler

import (
	"avito-tech/internal/models"
	"avito-tech/internal/repository"
	"avito-tech/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	getTransactionsByAccountID = "/get/transactions"
)

type transactionHandler struct {
	logger logger.Logger
	trRepo repository.TransactionHistory
}

func NewTransactionHandler(logger logger.Logger, trRepo repository.TransactionHistory) *transactionHandler {
	return &transactionHandler{
		logger: logger,
		trRepo: trRepo,
	}
}

func (trh *transactionHandler) Register(router *mux.Router) {
	router.HandleFunc(getTransactionsByAccountID, trh.getTransactionsByAccountID).Methods("GET")
}

func (trh *transactionHandler) getTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.TransactionHistoryReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		trh.logger.Errorf("error occurred while parsing json. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while parsing json. err:%s ", err))
		return
	}
	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		trh.logger.Errorf("error occurred while validating data. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while validating data. err:%s ", err))
		return
	}
	transactions, err := trh.trRepo.GetByAccountID(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		trh.logger.Errorf("error occurred while getting transactions. err:%s ", err)
		json.NewEncoder(w).Encode(fmt.Sprintf("error occurred while getting all transactions. err:%s ", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)

}
