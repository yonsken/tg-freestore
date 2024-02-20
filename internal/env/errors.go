package env

import "fmt"

var _ error = (*RequiredEnvVarError)(nil)

type RequiredEnvVarError struct {
	envVar string
}

func NewRequiredEnvVarError(envVar string) *RequiredEnvVarError {
	return &RequiredEnvVarError{envVar: envVar}
}

func (r *RequiredEnvVarError) Error() string {
	return fmt.Sprintf("required environment variable %q is not set", r.envVar)
}
