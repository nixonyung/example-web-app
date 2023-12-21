package env

import (
	"os"

	"source.local/common/internal/structparse"
)

func Parse(v any) error {
	return structparse.Parse(
		v,
		"env",
		func(key string) (string, bool) {
			return os.LookupEnv(key)
		},
	)
}
