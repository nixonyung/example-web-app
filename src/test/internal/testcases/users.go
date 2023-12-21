package testcases

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"source.local/api-gateway/pkg/apigateway"
	"source.local/common/pkg/httpbase"
)

var defaultHeaders = map[string]string{
	"X-Test": "1",
}

func getUser(userID uuid.UUID) error {
	var resBody []struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Role      string    `json:"role"`
		Balance   float64   `json:"balance"`
		CreatedAt time.Time `json:"created_at"`
	}

	if err := httpbase.MakeRequest(
		http.MethodGet,
		apigateway.URL("/users"),
		&httpbase.MakeRequestConfig{
			ReqParams: map[string]string{
				"id": userID.String(),
			},
			ReqHeaders: defaultHeaders,
			ResBody:    &resBody,
		},
	); err != nil {
		return err
	}
	return nil
}

func TestUsers() error {
	var postUserRes struct {
		UserID uuid.UUID `json:"id"`
	}

	if err := httpbase.MakeRequest(
		http.MethodPost,
		apigateway.URL("/users"),
		&httpbase.MakeRequestConfig{
			ReqHeaders: defaultHeaders,
			ReqBody: map[string]any{
				"username": "newuser1",
				"password": "password1",
			},
			ResBody: &postUserRes,
		},
	); err != nil {
		return err
	}
	if err := getUser(postUserRes.UserID); err != nil {
		return err
	}

	// negative test: duplicated username
	if err := httpbase.MakeRequest(
		http.MethodPost,
		apigateway.URL("/users"),
		&httpbase.MakeRequestConfig{
			ReqHeaders: defaultHeaders,
			ReqBody: map[string]any{
				"username": "newuser1",
				"password": "password1",
			},
			ResBody: &postUserRes,
		},
	); err != nil {
		log.Println(err)
	} else {
		return errors.New("POST /users should return error on duplicated username")
	}

	if err := httpbase.MakeRequest(
		http.MethodPost,
		apigateway.URL("/users/add-balance"),
		&httpbase.MakeRequestConfig{
			ReqHeaders: defaultHeaders,
			ReqBody: map[string]any{
				"user_id": postUserRes.UserID,
				"amount":  100.0,
			},
		},
	); err != nil {
		return err
	}
	if err := getUser(postUserRes.UserID); err != nil {
		return err
	}

	if err := httpbase.MakeRequest(
		http.MethodDelete,
		apigateway.URL("/users"),
		&httpbase.MakeRequestConfig{
			ReqParams: map[string]string{
				"id": postUserRes.UserID.String(),
			},
			ReqHeaders: defaultHeaders,
		},
	); err != nil {
		return err
	}
	if err := getUser(postUserRes.UserID); err != nil {
		return err
	}

	return nil
}
