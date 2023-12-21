package httpbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"source.local/common/pkg/logger"
)

type MakeRequestConfig struct {
	Method    string
	URL       *url.URL
	ReqHeader map[string]string
	ReqBody   map[string]any
}

func MakeRequest(config *MakeRequestConfig) ([]byte, error) {
	// helper function
	wrappedError := func(message string, err error) error {
		return fmt.Errorf("MakeRequest method=%s url=%s req_header=%+v req_body=%+v: %s: %w",
			config.Method,
			config.URL.String(),
			config.ReqHeader,
			config.ReqBody,
			message,
			err,
		)
	}

	reqBodyBytes, err := json.Marshal(config.ReqBody)
	if err != nil {
		return nil, wrappedError("marshal reqBody error", err)
	}

	req, err := http.NewRequest(
		config.Method,
		config.URL.String(),
		bytes.NewBuffer(reqBodyBytes),
	)
	if err != nil {
		return nil, wrappedError("http.NewRequest error", err)
	}
	{
		// set req header
		req.Header.Set("Content-Type", "application/json")
		for k, v := range config.ReqHeader {
			req.Header.Set(k, v)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, wrappedError("http.DefaultClient.Do error", err)
	}
	if res.StatusCode >= 400 {
		return nil, wrappedError("", fmt.Errorf("got status %s", res.Status))
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, wrappedError("reading res.Body", err)
	} else {
		res.Body.Close()
	}

	(&logger.Logger{LocOffset: 1}).Printf("MakeRequest [%s] %s \"%s\" reqBody=%s resBody=%s",
		res.Status,
		config.Method,
		config.URL.String(),
		reqBodyBytes,
		resBodyBytes,
	)
	return resBodyBytes, nil
}
