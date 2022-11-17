package repository

import (
	"avito-tech/internal/models"
	"avito-tech/pkg/logger"
	"database/sql"
	"errors"
)


//go:generate mockgen -source=repository.go -destination=mock_repository/mock_repository.go

var (
	ErrNoRowsAffected = errors.New("now rows affected")
	ErrUserDoesntExist = errors.New("user doesnt exist")

	ErrNewAccNegativeBalance = errors.New("can not set a negative balance for a new account")

)

type Account interface {
	GetBalanceByID(id int) (acc *models.Account, err error)
	ChangeBalance(acc *models.AccountDebit) (err error)
	GetAll() (acc []models.Account, err error)
	MoneyTransaction(transaction *models.Transaction) (err error)
}

type TransactionHistory interface {
	GetByAccountID(req *models.TransactionHistoryReq) (tr []models.TransactionHistory, err error)
}

type Repository struct {
	Account
	TransactionHistory
}

func New(db *sql.DB, logger logger.Logger) (repository *Repository) {
	return &Repository{
		Account:            NewAccountRepository(db, logger),
		TransactionHistory: NewTransactionRepository(db, logger),
	}
}
