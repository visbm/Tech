package models_test

import (
	"avito-tech/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {

	testCases := []struct {
		name    string
		r       *models.Transaction
		isValid bool
	}{
		{
			name: "pass",
			r: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     236.5,
				Comment:    "First transaction",
			},
			isValid: true,
		},
		{
			name: "pass",
			r: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     236.5,
				Comment:    "First transaction",
			},
			isValid: true,
		},
		{
			name: "invalid sender id",
			r: &models.Transaction{
				ReceiverID: 1,
				SenderID:   -2,
				Amount:     236.5,
				Comment:    "First transaction",
			},
			isValid: false,
		},
		{
			name: "invalid receiver id",
			r: &models.Transaction{
				ReceiverID: -10,
				SenderID:   2,
				Amount:     236.5,
				Comment:    "First transaction",
			},
			isValid: false,
		},
		{
			name: "invalid amount",
			r: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     -236.5,
				Comment:    "First transaction",
			},
			isValid: false,
		},
		{
			name: "invalid reciever = sender",
			r: &models.Transaction{
				ReceiverID: 2,
				SenderID:   2,
				Amount:     236.5,
				Comment:    "First transaction",
			},

			isValid: false,
		},
		{
			name: "no comm",
			r: &models.Transaction{
				ReceiverID: 2,
				SenderID:   2,
				Amount:     236.5,
			},

			isValid: false,
		},
		{
			name: "short comm",
			r: &models.Transaction{
				ReceiverID: 2,
				SenderID:   2,
				Amount:     236.5,
				Comment:    "1",
			},

			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.r.Validate())
			} else {
				assert.Error(t, tc.r.Validate())
			}
		})
	}
}

func TestAccountDebit_Validate(t *testing.T) {

	testCases := []struct {
		name    string
		r       *models.AccountDebit
		isValid bool
	}{
		{
			name: "pass",
			r: &models.AccountDebit{
				AccountID: 2,
				Amount:    5,
				Comment:    "First transaction",

			},
			isValid: true,
		},
		{
			name: "pass",
			r: &models.AccountDebit{
				AccountID: 25,
				Amount:    551.5,
				Comment:    "First transaction",

			},
			isValid: true,
		},
		{
			name: "invalid account id",
			r: &models.AccountDebit{
				AccountID: -2,
				Amount:    5,
				Comment:    "First transaction",

			},
			isValid: false,
		},
		{
			name: "invalid amount is blank",
			r: &models.AccountDebit{
				AccountID: 2,
				Amount:    0,
				Comment:    "First transaction",

			},
			isValid: false,
		},
		{
			name: "no comm",
			r: &models.AccountDebit{
				AccountID: 2,
				Amount:    -5,
				Comment:    "",

			},
			isValid: false,
		},
		{
			name: "short comment",
			r: &models.AccountDebit{
				AccountID: 2,
				Amount:    -5,
				Comment:    "3",

			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.r.Validate())
			} else {
				assert.Error(t, tc.r.Validate())
			}
		})
	}
}

func TestTransactionHistoryReq_Validate(t *testing.T) {

	testCases := []struct {
		name    string
		r       *models.TransactionHistoryReq
		isValid bool
	}{
		{
			name: "pass",
			r: &models.TransactionHistoryReq{
				AccountID: 1,
				OrderBy: "date_time",
			},
			isValid: true,
		},
		{
			name: "pass",
			r: &models.TransactionHistoryReq{
				AccountID: 1,
				OrderBy: "amount",
			},
			isValid: true,
		},
		{
			name: "pass",
			r: &models.TransactionHistoryReq{
				AccountID: 1,
				OrderBy: "date_time asc",
			},
			isValid: true,
		},
		{
			name: "pass",
			r: &models.TransactionHistoryReq{
				AccountID: 1,
				OrderBy: "date_time DESC",
			},
			isValid: true,
		},
		{
			name: "wrong order by",
			r: &models.TransactionHistoryReq{
				AccountID: 1,
				OrderBy: "id",
			},
			isValid: false,
		},{
			name: "wrong id",
			r: &models.TransactionHistoryReq{
				AccountID: -1,
				OrderBy: "date_time",
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.r.Validate())
			} else {
				assert.Error(t, tc.r.Validate())
			}
		})
	}
}