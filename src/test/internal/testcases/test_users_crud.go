package testcases

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"source.local/api-gateway/pkg/apigateway"
	"source.local/common/pkg/floats"
	"source.local/common/pkg/httpbase"
	"source.local/common/pkg/logger"
	"source.local/test/internal/testbase"
)

func TestUsersCRUD() {
	var (
		newUsername = "test_users_crud_username"
		newPassword = "test_users_crud_password"

		postUserRes struct {
			UserID uuid.UUID `json:"id"`
		}
	)

	if b, err := httpbase.MakeRequest(&httpbase.MakeRequestConfig{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   apigateway.Host,
			Path:   "users",
		},
		ReqHeader: testbase.Header,
		ReqBody: map[string]any{
			"username": newUsername,
			"password": newPassword,
		},
	}); err != nil {
		logger.Default.Fatal(err)
	} else if err := json.Unmarshal(b, &postUserRes); err != nil {
		logger.Default.Fatal(err)
	} else if users, err := testbase.GetUser(postUserRes.UserID); err != nil {
		logger.Default.Fatal(err)
	} else if len(users) != 1 {
		logger.Default.Fatal(errors.New("a unique user should exist after POST /users"))
	}

	// negative test: create user with duplicated username
	if _, err := httpbase.MakeRequest(&httpbase.MakeRequestConfig{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   apigateway.Host,
			Path:   "users",
		},
		ReqHeader: testbase.Header,
		ReqBody: map[string]any{
			"username": newUsername,
			"password": newPassword,
		},
	}); err != nil {
		logger.Default.Printf("%v", err)
	} else {
		logger.Default.Fatal(errors.New("POST /users should return error on duplicated username"))
	}

	if _, err := httpbase.MakeRequest(&httpbase.MakeRequestConfig{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   apigateway.Host,
			Path:   "users/add-balance",
		},
		ReqHeader: testbase.Header,
		ReqBody: map[string]any{
			"user_id": postUserRes.UserID,
			"amount":  100.0,
		},
	}); err != nil {
		logger.Default.Fatal(err)
	} else if users, err := testbase.GetUser(postUserRes.UserID); err != nil {
		logger.Default.Fatal(err)
	} else if len(users) != 1 {
		logger.Default.Fatal(errors.New("a unique user should exist after POST /users/add-balance"))
	} else if !floats.AlmostEqual(users[0].Balance, 100.0) {
		logger.Default.Fatal(errors.New("wrong balance after POST /users/add-balance"))
	}

	if _, err := httpbase.MakeRequest(&httpbase.MakeRequestConfig{
		Method: http.MethodDelete,
		URL: &url.URL{
			Scheme: "http",
			Host:   apigateway.Host,
			Path:   "users",
			RawQuery: url.Values{
				"id": []string{postUserRes.UserID.String()},
			}.Encode(),
		},
		ReqHeader: testbase.Header,
	}); err != nil {
		logger.Default.Fatal(err)
	} else if users, err := testbase.GetUser(postUserRes.UserID); err != nil {
		logger.Default.Fatal(err)
	} else if len(users) != 0 {
		logger.Default.Fatal(errors.New("user should not exist after DELETE /users"))
	}
}
