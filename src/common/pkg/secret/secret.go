package secret

import (
	"fmt"
	"os"

	"source.local/common/internal/structparse"
)

func Parse(v any) error {
	return structparse.Parse(
		v,
		"secret",
		func(key string) (string, bool) {
			if b, err := os.ReadFile(fmt.Sprintf("/run/secrets/%s", key)); err != nil {
				return "", false
			} else {
				return string(b), true
			}
		},
	)
}
