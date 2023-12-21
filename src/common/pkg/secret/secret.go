package secret

import (
	"fmt"
	"os"

	"source.local/common/internal/structparse"
)

func Parse(v any) error {
	if err := structparse.Parse(
		v,
		&structparse.Config{
			TagKey: "secret",
			LookupFn: func(key string) (string, bool) {
				if b, err := os.ReadFile(fmt.Sprintf("/run/secrets/%s", key)); err != nil {
					return "", false
				} else {
					return string(b), true
				}
			},
			AllowMissing: false,
		},
	); err != nil {
		return fmt.Errorf("structparse.Parse: %w", err)
	}
	return nil
}
