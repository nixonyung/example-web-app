package httpbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type MakeRequestConfig struct {
	ReqParams  map[string]string
	ReqHeaders map[string]string
	ReqBody    map[string]any
	ResBody    any
}

func MakeRequest(
	method string,
	url string,
	config *MakeRequestConfig,
) error {
	newError := func(message string, err error) error {
		return fmt.Errorf("MakeRequest method=%s url=%s config=%+v : %s : %w",
			method,
			url,
			config,
			message,
			err,
		)
	}

	//
	// construct req
	//

	if len(config.ReqParams) != 0 {
		urlBuilder := strings.Builder{}
		urlBuilder.WriteString(url)
		urlBuilder.WriteString("?")
		isFirstParam := true
		for k, v := range config.ReqParams {
			if !isFirstParam {
				urlBuilder.WriteString("&")
			}
			urlBuilder.WriteString(k)
			urlBuilder.WriteString("=")
			urlBuilder.WriteString(v)
			isFirstParam = false
		}
		url = urlBuilder.String()
	}

	reqBodyBytes, err := json.Marshal(config.ReqBody)
	if err != nil {
		return newError("marshal reqBody error", err)
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqBodyBytes),
	)
	if err != nil {
		return newError("http.NewRequest error", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range config.ReqHeaders {
		req.Header.Set(k, v)
	}

	// send req
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return newError("http.DefaultClient.Do error", err)
	}
	defer res.Body.Close()

	// parse resBody
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return newError("reading res.Body", err)
	}
	if res.StatusCode >= 400 {
		return newError(fmt.Sprintf("got %s", res.Status), err)
	}
	if config.ResBody != nil {
		if err := json.Unmarshal(resBodyBytes, config.ResBody); err != nil {
			return newError("unmarshal resBody error", err)
		}
	}
	log.Printf("MakeRequest [%s] %s \"%s\" reqBody=%s resBody=%s",
		res.Status,
		method,
		url,
		reqBodyBytes,
		resBodyBytes,
	)
	return nil
}
