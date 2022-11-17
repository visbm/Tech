package repository_test

import (
	"avito-tech/internal/models"
	"avito-tech/internal/repository"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_GetByAccountID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := repository.NewTransactionRepository(db, log)
	date := time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)

	testTable := []struct {
		name           string
		req            *models.TransactionHistoryReq
		mock           func(req *models.TransactionHistoryReq)
		expectedResult []models.TransactionHistory
		expectedError  bool
	}{
		{
			name: "OK",
			req: &models.TransactionHistoryReq{
				AccountID: 1,
				Limit:     0,
				Offset:    0,
				OrderBy:   "",
			},
			mock: func(req *models.TransactionHistoryReq) {
				rows := sqlmock.NewRows([]string{"transaction_id", "account_id", "amount", "date_time", "comment"}).
					AddRow(1, 2, 11.2, date, "salary")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT transaction_id ,account_id , amount , date_time ,comment FROM transactions WHERE account_id = $1")).WithArgs(req.AccountID).WillReturnRows(rows)
			},
			expectedResult: []models.TransactionHistory{
				{
					TransactionID: 1,
					AccountID:     2,
					Amount:        11.2,
					Date:          date,
					Comment:       "salary",
				},
			},
			expectedError: false,
		},
		{
			name: "no rows",
			req: &models.TransactionHistoryReq{
				AccountID: 1,
				Limit:     0,
				Offset:    0,
				OrderBy:   "amount",
			}, mock: func(req *models.TransactionHistoryReq) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT transaction_id ,account_id , amount , date_time ,comment FROM transactions WHERE account_id = $1")).WithArgs(req.AccountID).WillReturnError(errors.New("no rows"))
			},

			expectedResult: nil,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.req)
			result, err := r.GetByAccountID(tt.req)
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
