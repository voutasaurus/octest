package config

import (
	"os"
	"strconv"
)

// Envvar stores details for an environment variable
type Envvar struct {
	Key   string
	Value string
	Set   bool
}

// Env reads an environment variable from the OS
func Env(key string) Envvar {
	value, set := os.LookupEnv(key)
	return Envvar{Key: key, Value: value, Set: set}
}

// WithDefault returns the value of the environment variable if it is set.
// Otherwise it returns the provided default value.
func (e Envvar) WithDefault(value string) string {
	if e.Set {
		return e.Value
	}
	return value
}

// Required returns the value of the environment variable if it is set.
// Otheriwse it will call provided error func and return an empty string.
//
// Example:
//  config.Env("SOME_ENVIRONMENT_VARIABLE").Required(func(key string) { log.Fatalf("%q must be set", key) })
//
func (e Envvar) Required(errlog func(key string)) string {
	if !e.Set {
		errlog(e.Key)
	}
	return e.Value
}

func (e Envvar) WithDefaultInt(value int, errlog func(key string, parseErr error)) int {
	if !e.Set {
		return value
	}
	v, err := strconv.Atoi(e.Value)
	if err != nil {
		errlog(e.Key, err)
	}
	return v
}
