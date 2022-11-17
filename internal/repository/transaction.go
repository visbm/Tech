package repository

import (
	"avito-tech/internal/models"
	"avito-tech/pkg/logger"
	"fmt"

	"database/sql"
)

type transactionHistory struct {
	db     *sql.DB
	logger logger.Logger
}

func NewTransactionRepository(db *sql.DB, logger logger.Logger) (repository TransactionHistory) {
	return &transactionHistory{
		db:     db,
		logger: logger,
	}
}

func (rep *transactionHistory) GetByAccountID(req *models.TransactionHistoryReq) (transactions []models.TransactionHistory, err error) {

	var rows *sql.Rows

	query := `SELECT transaction_id,
					account_id,
					amount,
					date_time,
					COMMENT
			FROM transactions
			WHERE account_id = $1`

	if req.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", req.OrderBy)
	}

	if req.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, req.Offset)
	}

	rows, err = rep.db.Query(query, req.AccountID)
	if err != nil {
		rep.logger.Errorf("error occurred while getting transaction history. err: %s", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		tr := models.TransactionHistory{}
		if err = rows.Scan(
			&tr.TransactionID,
			&tr.AccountID,
			&tr.Amount,
			&tr.Date,
			&tr.Comment,
		); err != nil {
			rep.logger.Errorf("error occurred while getting transaction history. err: %s", err)
			return nil, err
		}

		transactions = append(transactions, tr)
	}

	return transactions, nil

}
