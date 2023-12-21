package httpbase

import (
	"net/url"
)

type URLConfig struct {
	Scheme string
	Host   string
	Path   string
	Values url.Values
	User   *url.Userinfo
}

func URL(config *URLConfig) string {
	return (&url.URL{
		Scheme:   config.Scheme,
		Host:     config.Host,
		Path:     config.Path,
		RawQuery: config.Values.Encode(),
		User:     config.User,
	}).String()
}
