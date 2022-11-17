package repository

import (
	"avito-tech/internal/models"
	"avito-tech/pkg/logger"
	"time"

	"database/sql"
)

type account struct {
	db     *sql.DB
	logger logger.Logger
}

func NewAccountRepository(db *sql.DB, logger logger.Logger) (repository Account) {
	return &account{
		db:     db,
		logger: logger,
	}
}

func (rep *account) GetBalanceByID(id int) (acc *models.Account, err error) {
	acc = &models.Account{}
	query := `SELECT account_id,
					balance
			FROM accounts
			WHERE account_id = $1`

	if err = rep.db.QueryRow(query, id).
		Scan(
			&acc.ID,
			&acc.Balance,
		); err != nil {
		rep.logger.Errorf("error occurred while getting message by id, err: %s", err)
		return nil, err
	}

	return acc, nil
}

func (rep *account) ChangeBalance(acc *models.AccountDebit) (err error) {
	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	query := `UPDATE accounts
			SET balance = balance + ($2)
			WHERE account_id = $1`

	result, err := tx.Exec(query,
		acc.AccountID,
		acc.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected != 1 {
		err = rep.insert(acc.AccountID, acc.Amount, acc.Comment)
		if err != nil {
			return err
		}
		return nil
	}

	th := &models.TransactionHistory{
		AccountID: acc.AccountID,
		Amount:    acc.Amount,
		Date:      time.Now(),
		Comment:   acc.Comment,
	}
	err = rep.addInTransactionsHistory(th)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rep *account) GetAll() (accounts []models.Account, err error) {
	query := `SELECT account_id, 
					balance
			FROM accounts`

	rows, err := rep.db.Query(query)
	if err != nil {
		rep.logger.Errorf("error occurred while getting all accounts. err: %s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		account := models.Account{}
		if err = rows.Scan(
			&account.ID,
			&account.Balance,
		); err != nil {
			rep.logger.Errorf("error occurred while getting all accounts. err: %s", err)
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (rep *account) MoneyTransaction(transaction *models.Transaction) (err error) {
	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	err = func() error {
		query := `UPDATE accounts 
 				SET balance = balance - $2 
  				WHERE account_id = $1;`

		result, err := tx.Exec(query,
			transaction.SenderID,
			transaction.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}

		if rowsAffected != 1 {
			tx.Rollback()
			return ErrUserDoesntExist
		}

		th := &models.TransactionHistory{
			AccountID: transaction.SenderID,
			Amount:    -(transaction.Amount),
			Date:      time.Now(),
			Comment:   transaction.Comment,
		}
		err = rep.addInTransactionsHistory(th)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	}()
	if err != nil {
		return err
	}

	err = func() error {
		query := `UPDATE accounts 
 				SET balance = balance + $2 
  				WHERE account_id = $1;`

		result, err := tx.Exec(query,
			transaction.ReceiverID,
			transaction.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}

		if rowsAffected != 1 {
			tx.Rollback()
			return ErrUserDoesntExist
		}

		th := &models.TransactionHistory{
			AccountID: transaction.ReceiverID,
			Amount:    transaction.Amount,
			Date:      time.Now(),
			Comment:   transaction.Comment,
		}
		err = rep.addInTransactionsHistory(th)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	}()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (rep *account) insert(id int, balance float64, comment string) (err error) {
	if balance <= 0 {
		return ErrNewAccNegativeBalance
	}
	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	query := `INSERT INTO accounts (account_id, balance)
			VALUES ($1, $2) ;`

	result, err := tx.Exec(query,
		id,
		balance)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return ErrUserDoesntExist
	}

	if rowsAffected != 1 {
		tx.Rollback()
		return ErrNoRowsAffected
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	transactionHistory := &models.TransactionHistory{
		AccountID: id,
		Amount:    balance,
		Date:      time.Now(),
		Comment:   comment,
	}
	err = rep.addInTransactionsHistory(transactionHistory)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rep *account) addInTransactionsHistory(tr *models.TransactionHistory) (err error) {

	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	query := `INSERT INTO transactions (account_id, amount, date_time, comment)
			VALUES ($1,$2, $3, $4)`

	result, err := tx.Exec(query,
		tr.AccountID,
		tr.Amount,
		tr.Date,
		tr.Comment)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected != 1 {
		tx.Rollback()
		return ErrNoRowsAffected
	}

	return tx.Commit()
}
