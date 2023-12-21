package testbase

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"source.local/api-gateway/pkg/apigateway"
	"source.local/common/pkg/httpbase"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUser(userID uuid.UUID) ([]User, error) {
	var users []User
	if b, err := httpbase.MakeRequest(&httpbase.MakeRequestConfig{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   apigateway.Host,
			Path:   "users",
			RawQuery: url.Values{
				"id": []string{userID.String()},
			}.Encode(),
		},
		ReqHeader: Header,
	}); err != nil {
		return nil, err
	} else if err := json.Unmarshal(b, &users); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}
