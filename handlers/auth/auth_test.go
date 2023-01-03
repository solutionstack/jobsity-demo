package auth

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/mocks"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAuthHandler_Register(t *testing.T) {

	type TC struct {
		name                  string
		reqBody               string
		validationErr         bool
		serviceErr            error
		serviceReturn         uuid.UUID
		expectedHttpStatus    int
		expectedResponse      string
		expectedNumberOfCalls int
	}
	successId := uuid.New()
	testcases := []TC{
		{
			name:                  "tesRegisterUserOK",
			reqBody:               `{"first_name": "foo bar",    "email": "aa@gmail.com",  "password": "aaaaa"  }`,
			validationErr:         false,
			serviceErr:            nil,
			serviceReturn:         successId,
			expectedNumberOfCalls: 1,
			expectedResponse:      `{"id":"` + successId.String() + `"}`,
			expectedHttpStatus:    http.StatusOK,
		},
		{
			name:                  "tesBadRequest_InvalidEmail",
			reqBody:               `{"first_name": "foo bar",    "email": "@gmail.com",  "password": "aaaaa"  }`,
			validationErr:         true,
			serviceErr:            nil,
			serviceReturn:         uuid.Nil,
			expectedNumberOfCalls: 0,
			expectedResponse:      `{"message": "error in request: mail: no angle-addr"}`,
			expectedHttpStatus:    http.StatusBadRequest,
		},
		{
			name:                  "tesServiceError_DuplicateRegistration",
			reqBody:               `{"first_name": "foo bar",    "email": "aa@gmail.com",  "password": "aaaaa"  }`,
			validationErr:         false,
			serviceErr:            errors.New("user already exist"),
			serviceReturn:         uuid.Nil,
			expectedNumberOfCalls: 1,
			expectedResponse:      `{"message": "error in request: error while creating user: user already exist"}`,
			expectedHttpStatus:    http.StatusInternalServerError,
		},
	}
	for _, tc := range testcases {
		testFunc := func(t *testing.T) {
			authHandler, authServiceMock := initHandlerAndMock()
			handler := http.HandlerFunc(authHandler.Register)
			const serviceMethod = "CreateUser"
			const httpMethod = http.MethodPost
			const url = "/auth/signup"
			body := strings.NewReader(tc.reqBody)

			if !tc.validationErr {
				authServiceMock.On(serviceMethod, mock.Anything).Return(tc.serviceReturn, tc.serviceErr)
			}

			req, err := http.NewRequest(httpMethod, url, body)
			if err != nil {
				t.Fatal(err)
			}
			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, req)

			// assert
			assert.Equal(t, tc.expectedHttpStatus, responseRecorder.Code)
			if tc.expectedHttpStatus != http.StatusNoContent {
				assert.JSONEq(t, tc.expectedResponse, responseRecorder.Body.String())
			}
			authServiceMock.AssertNumberOfCalls(t, serviceMethod, tc.expectedNumberOfCalls)

			authServiceMock.AssertExpectations(t)
		}

		t.Run(tc.name, testFunc)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type TC struct {
		name                  string
		reqBody               string
		validationErr         bool
		serviceErr            error
		serviceReturn         *models.UserRecord
		expectedHttpStatus    int
		expectedResponse      string
		expectedNumberOfCalls int
	}

	OkRecord := &models.UserRecord{
		ID: uuid.MustParse("c0b996c9-7978-4a41-ada8-7e271be120ff"),
		Signup: models.Signup{
			FirstName: "foo bar",
			Email:     "aa@gmail.com",
			Password:  "aaaaa",
		},
	}
	testcases := []TC{
		{
			name:                  "tesLoginOK",
			reqBody:               `{   "email": "aa@gmail.com",  "password": "aaaaa"  }`,
			validationErr:         false,
			serviceErr:            nil,
			serviceReturn:         OkRecord,
			expectedNumberOfCalls: 1,
			expectedResponse:      `{"id": "c0b996c9-7978-4a41-ada8-7e271be120ff","first_name": "foo bar","email": "aa@gmail.com","password": "aaaaa"}`,
			expectedHttpStatus:    http.StatusOK,
		},

		{
			name:                  "tesLogin_NoUserRecordFound",
			reqBody:               `{   "email": "aa@gmail.com",  "password": "aaaaa"  }`,
			validationErr:         false,
			serviceErr:            nil,
			serviceReturn:         nil,
			expectedNumberOfCalls: 1,
			expectedResponse:      `{  "message": "error in request: no user record found"}`,
			expectedHttpStatus:    http.StatusNotFound,
		},
	}
	for _, tc := range testcases {
		testFunc := func(t *testing.T) {
			authHandler, authServiceMock := initHandlerAndMock()
			handler := http.HandlerFunc(authHandler.Login)
			const serviceMethod = "ValidateLogin"
			const httpMethod = http.MethodPost
			const url = "/auth/login"
			body := strings.NewReader(tc.reqBody)

			if !tc.validationErr {
				authServiceMock.On(serviceMethod, mock.Anything).Return(tc.serviceReturn, tc.serviceErr)
			}

			req, err := http.NewRequest(httpMethod, url, body)
			if err != nil {
				t.Fatal(err)
			}
			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, req)

			// assert
			assert.Equal(t, tc.expectedHttpStatus, responseRecorder.Code)
			if tc.expectedHttpStatus != http.StatusNoContent {
				assert.JSONEq(t, tc.expectedResponse, responseRecorder.Body.String())
			}
			authServiceMock.AssertNumberOfCalls(t, serviceMethod, tc.expectedNumberOfCalls)

			authServiceMock.AssertExpectations(t)
		}

		t.Run(tc.name, testFunc)
	}
}

func initHandlerAndMock() (AuthHandler, *mocks.Service) {
	authServiceMock := mocks.Service{}
	logger := zerolog.New(os.Stdout)

	handler := AuthHandler{
		logger: logger,
		svc:    &authServiceMock,
	}

	return handler, &authServiceMock
}
