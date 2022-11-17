package repository_test

import (
	"avito-tech/internal/models"
	"avito-tech/internal/repository"
	"avito-tech/pkg/logger"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var log = logger.GetLogger()

func Test_GetBalanceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := repository.NewAccountRepository(db, log)

	testTable := []struct {
		name           string
		id             int
		mock           func(id int)
		expectedResult *models.Account
		expectedError  bool
	}{
		{
			name: "OK",
			id:   1,
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"account_id", "balance"}).
					AddRow(1, 11.623413)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT account_id, balance FROM accounts WHERE account_id = $1")).WithArgs(id).WillReturnRows(rows)
			},
			expectedResult: &models.Account{
				ID:      1,
				Balance: 11.623413,
			},
			expectedError: false,
		},
		{
			name: "no rows",
			id:   1,
			mock: func(id int) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT account_id, balance FROM accounts WHERE account_id = $1")).WithArgs(id).WillReturnError(errors.New("no rows"))
			},

			expectedResult: nil,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.id)
			result, err := r.GetBalanceByID(tt.id)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := repository.NewAccountRepository(db, log)

	testTable := []struct {
		name           string
		mock           func()
		expectedResult []models.Account
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"account_id", "balance"}).
					AddRow(1, 11.623413).AddRow(2, 22)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT account_id, balance FROM accounts")).WillReturnRows(rows)
			},
			expectedResult: []models.Account{
				{
					ID:      1,
					Balance: 11.623413,
				},
				{
					ID:      2,
					Balance: 22,
				},
			},

			expectedError: false,
		},
		{
			name: "no rows",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT account_id, balance FROM accounts")).WillReturnError(errors.New("no rows"))
			},

			expectedResult: nil,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := r.GetAll()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_ChangeBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := repository.NewAccountRepository(db, log)

	testTable := []struct {
		name          string
		acc           *models.Transaction
		mock          func(acc *models.Transaction)
		expectedError bool
	}{
		{
			name: "OK",
			acc: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     2,
				Comment:    "2",
			},
			mock: func(acc *models.Transaction) {
				mock.ExpectBegin()
				query := `UPDATE accounts
						SET balance = balance - $2
						WHERE account_id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(acc.SenderID, acc.Amount).WillReturnError(errors.New("transaction has already been committed or rolled back"))
				mock.ExpectRollback()
			},

			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.acc)
			err := r.MoneyTransaction(tt.acc)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
