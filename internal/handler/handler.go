package handler

import (
	"avito-tech/internal/repository"
	"avito-tech/pkg/logger"

	"github.com/gorilla/mux"
)

type Handler struct {
	logger             logger.Logger
	repository         *repository.Repository
	accountHandler     *accountHandler
	transactionHandler *transactionHandler
}

func NewHandler(logger logger.Logger, repository *repository.Repository) *Handler {
	return &Handler{
		logger:             logger,
		repository:         repository,
		accountHandler:     NewAccountHandler(logger, repository.Account),
		transactionHandler: NewTransactionHandler(logger, repository.TransactionHistory),
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	h.accountHandler.Register(router)
	h.transactionHandler.Register(router)
	return router
}
