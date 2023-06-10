package utils

import (
	"github.com/pkg/errors"
)

var (
	ErrSecretTypeNotServiceAccountToken = errors.New("the secret is not of type service-account-token")
)
