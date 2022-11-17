package handler_test

import (
	"avito-tech/internal/handler"
	"avito-tech/internal/models"
	"avito-tech/internal/repository/mock_repository"
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_getTransactionsByAccountID(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockTransactionHistory, req *models.TransactionHistoryReq)

	testTable := []struct {
		name, inputBody     string
		tr                  *models.TransactionHistoryReq
		req                 *models.TransactionHistoryReq
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: `{ 
				"account_id": 1, 
				"limit":1,
				"offset": 1, 
				"order_by" : "date_time DESC"  
			}`,
			req: &models.TransactionHistoryReq{
				AccountID: 1,
				Limit:     1,
				Offset:    1,
				OrderBy:   "date_time DESC",
			},
			mockBehavior: func(s *mock_repository.MockTransactionHistory, req *models.TransactionHistoryReq) {
				s.EXPECT().GetByAccountID(req).Return([]models.TransactionHistory{{TransactionID: 1, AccountID: 2, Amount: 22, Comment: "comment"}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"transaction_ID\":1,\"account_id\":2,\"amount\":22,\"date\":\"0001-01-01T00:00:00Z\",\"comment\":\"comment\"}]\n",
		},
		{
			name: "Invalid JSON request",
			inputBody: `{ 
				"account_id": 1, 
				"limit":1,
				offset": 1, 
				"order_by" : "date_time DESC"  
			}`,
			mockBehavior: func(s *mock_repository.MockTransactionHistory, req *models.TransactionHistoryReq) {
			},

			expectedStatusCode:  400,
			expectedRequestBody: "\"error occurred while parsing json. err:invalid character 'o' looking for beginning of object key string \"\n",
		},
		{
			name: "no such acc",
			req: &models.TransactionHistoryReq{
				AccountID: 1,
				Limit:     1,
				Offset:    1,
				OrderBy:   "date_time DESC",
			},
			inputBody: `{ 
				"account_id": 1, 
				"limit":1,
				"offset": 1, 
				"order_by" : "date_time DESC"  
			}`,
			mockBehavior: func(s *mock_repository.MockTransactionHistory, req *models.TransactionHistoryReq) {
				s.EXPECT().GetByAccountID(req).Return(nil, errors.New("sql: no rows in result set"))

			},
			expectedStatusCode:  500,
			expectedRequestBody: "\"error occurred while getting all transactions. err:sql: no rows in result set \"\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rep := mock_repository.NewMockTransactionHistory(c)
			testCase.mockBehavior(rep, testCase.req)

			handler := handler.NewTransactionHandler(log, rep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", "/get/transactions", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
