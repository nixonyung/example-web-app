package env

import (
	"fmt"
	"os"

	"source.local/common/internal/structparse"
)

func Parse(v any) error {
	if err := structparse.Parse(
		v,
		&structparse.Config{
			TagKey:       "env",
			LookupFn:     os.LookupEnv,
			AllowMissing: false,
		},
	); err != nil {
		return fmt.Errorf("structparse.Parse: %w", err)
	}
	return nil
}
