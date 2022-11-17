package handler_test

import (
	"avito-tech/internal/handler"
	"avito-tech/internal/models"
	"avito-tech/internal/repository/mock_repository"
	"avito-tech/pkg/logger"
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var log = logger.GetLogger()

func Test_MoneyTransaction(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAccount, tr *models.Transaction)

	testTable := []struct {
		name, inputBody     string
		mockBehavior        mockBehavior
		tr                  *models.Transaction
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: `{
				"receiver_id": 1,
				"sender_id": 2,
				"amount": 22.5,
				"comment": "You paid for me in a restaurant"
			}`,
			tr: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     22.5,
				Comment:    "You paid for me in a restaurant",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.Transaction) {
				s.EXPECT().MoneyTransaction(tr).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "\"Transaction was successful\"\n",
		},
		{
			name: "Invalid JSON request",
			inputBody: `{
				"receiver_id": 1,
				"sender_id": 2,
				amount": 20,
				"comment": "You paid for me in a restaurant"
			}`,
			tr: &models.Transaction{},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.Transaction) {
			},

			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while parsing json. err:invalid character 'a' looking for beginning of object key string \"\n",
		},
		{
			name: "Invalid receiver_id",
			inputBody: `{
				"receiver_id": -1,
				"sender_id": 2,
				"amount": 20,
				"comment": "You paid for me in a restaurant"
			}`,
			tr: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     22.5,
				Comment:    "owe you",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.Transaction) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while validating data. err:receiver_id: must be no less than 1. \"\n",
		},
		{
			name: "receiver_id: cannot be blank",
			inputBody: `{
				"sender_id": 2,
				"amount": 20,
				"comment": "You paid for me in a restaurant"
			}`,
			tr: &models.Transaction{
				ReceiverID: 1,
				SenderID:   2,
				Amount:     22.5,
				Comment:    "owe you",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.Transaction) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while validating data. err:receiver_id: cannot be blank. \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockAccount(c)
			testCase.mockBehavior(rep, testCase.tr)

			handler := handler.NewAccountHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("POST", "/transaction", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_changeBalance(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAccount, tr *models.AccountDebit)

	testTable := []struct {
		name, inputBody     string
		mockBehavior        mockBehavior
		tr                  *models.AccountDebit
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: `{
				"account_id": 1,
				"amount": 22.5,
				"comment": "Credit payment"
				}`,
			tr: &models.AccountDebit{
				AccountID: 1,
				Amount:    22.5,
				Comment:   "Credit payment",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.AccountDebit) {
				s.EXPECT().ChangeBalance(tr).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "\"Balance succsefully changed\"\n",
		},
		{
			name: "Invalid JSON request",
			inputBody: `{
				"account_id": 1,
				"amount": 22.5,
				comment": "Credit payment"
				}`,
			tr: &models.AccountDebit{},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.AccountDebit) {
			},

			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while parsing json. err:invalid character 'c' looking for beginning of object key string \"\n",
		},
		{
			name: "Invalid account_id",
			inputBody: `{
				"account_id": -1,
				"amount": 22.5,
				"comment": "Credit payment"
				}`,
			tr: &models.AccountDebit{
				AccountID: 1,
				Amount:    22.5,
				Comment:   "Credit payment",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.AccountDebit) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while validating data. err:account_id: must be no less than 1. \"\n",
		},
		{
			name: "account_id: cannot be blank",
			inputBody: `{
				"amount": 22.5,
				"comment": "Credit payment"
				}`,
			tr: &models.AccountDebit{
				AccountID: 1,
				Amount:    22.5,
				Comment:   "Credit payment",
			},
			mockBehavior: func(s *mock_repository.MockAccount, tr *models.AccountDebit) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while validating data. err:account_id: cannot be blank. \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockAccount(c)
			testCase.mockBehavior(rep, testCase.tr)

			handler := handler.NewAccountHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("POST", "/changeBalance", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_getBalanceByID(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAccount, id int)

	testTable := []struct {
		name, url           string
		id                  int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",
			id:   1,
			url:  "/get/balance/1",
			mockBehavior: func(s *mock_repository.MockAccount, id int) {
				s.EXPECT().GetBalanceByID(id).Return(&models.Account{ID: 1, Balance: 22}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"account_id\":1,\"balance\":22}\n",
		},
		{
			name: "Invalid link",
			id:   1,
			url:  "/get/balance/",
			mockBehavior: func(s *mock_repository.MockAccount, id int) {
			},

			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
		{
			name: "no such acc",
			id:   100,
			url:  "/get/balance/100",
			mockBehavior: func(s *mock_repository.MockAccount, id int) {
				s.EXPECT().GetBalanceByID(id).Return(nil, errors.New("sql: no rows in result set"))

			},
			expectedStatusCode:  500,
			expectedRequestBody: "\"error occurred while getting all account. err:sql: no rows in result set \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockAccount(c)
			testCase.mockBehavior(rep, testCase.id)

			handler := handler.NewAccountHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", testCase.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_currencyBalance(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAccount, id int)

	testTable := []struct {
		name, url           string
		id                  int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{

		{
			name: "Invalid link",
			id:   1,
			url:  "/get/balance/USD/",
			mockBehavior: func(s *mock_repository.MockAccount, id int) {
			},

			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
		{
			name: "no such acc",
			id:   100,
			url:  "/get/balance/USD/100",
			mockBehavior: func(s *mock_repository.MockAccount, id int) {
				s.EXPECT().GetBalanceByID(id).Return(nil, errors.New("sql: no rows in result set"))

			},
			expectedStatusCode:  500,
			expectedRequestBody: "\"error occurred while getting all account. err:sql: no rows in result set \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockAccount(c)
			testCase.mockBehavior(rep, testCase.id)

			handler := handler.NewAccountHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", testCase.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_getAll(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAccount)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",

			mockBehavior: func(s *mock_repository.MockAccount) {
				s.EXPECT().GetAll().Return([]models.Account{{ID: 1, Balance: 22}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"account_id\":1,\"balance\":22}]\n",
		},
		{
			name: "no accounts",
		
			mockBehavior: func(s *mock_repository.MockAccount) {
				s.EXPECT().GetAll().Return(nil, errors.New("sql: no rows in result set"))

			},
			expectedStatusCode:  500,
			expectedRequestBody: "\"error occurred while getting all accounts. err:sql: no rows in result set \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockAccount(c)
			testCase.mockBehavior(rep)

			handler := handler.NewAccountHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", "/get/all", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
