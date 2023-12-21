// a simplified version of [caarlos0/env](https://github.com/caarlos0/env/blob/main/env.go)

package structparse

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Config struct {
	TagKey       string
	LookupFn     func(key string) (string, bool)
	AllowMissing bool
}

func Parse(
	v any,
	config *Config,
) error {
	if reflect.ValueOf(v).Kind() != reflect.Pointer ||
		reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return errors.New("expect v to be a pointer to struct")
	}
	if config == nil {
		config = &Config{}
	}

	refType := reflect.TypeOf(v).Elem()
	refVal := reflect.ValueOf(v).Elem()

	var errs []error
	for i := 0; i < refType.NumField(); i++ {
		tagVal, ok := refType.Field(i).Tag.Lookup(config.TagKey)
		if !ok {
			continue
		}
		valParsed, ok := config.LookupFn(tagVal)
		if !ok {
			if !config.AllowMissing {
				errs = append(errs, fmt.Errorf("value not found for tag %s",
					tagVal,
				))
			}
			continue
		}
		if !refVal.Field(i).CanSet() {
			errs = append(errs, fmt.Errorf("cannot set value to struct field %s",
				refType.Field(i).Name,
			))
			continue
		}
		switch refType.Field(i).Type.Kind() {
		case reflect.String:
			refVal.Field(i).SetString(valParsed)
		case reflect.Bool:
			refVal.Field(i).SetBool(valParsed != "0")
		case reflect.Int:
			if val, err := strconv.ParseInt(valParsed, 10, 64); err != nil {
				errs = append(errs, fmt.Errorf("error when parsing tag %s: %w",
					tagVal,
					err,
				))
			} else {
				refVal.Field(i).SetInt(val)
			}
		case reflect.Float64:
			if val, err := strconv.ParseFloat(valParsed, 64); err != nil {
				errs = append(errs, fmt.Errorf("error when parsing tag %s: %w",
					tagVal,
					err,
				))
			} else {
				refVal.Field(i).SetFloat(val)
			}
		default:
			errs = append(errs, fmt.Errorf("Parse: unsupported type %s for tag %s",
				refType.Field(i).Type.Kind(),
				tagVal,
			))
		}
	}
	return errors.Join(errs...)
}
