package controller

import "errors"

var (
	ErrConfigNotFound            = errors.New("configuration file .tfprovidercheck.yaml isn't found")
	ErrProviderNameIsRequired    = errors.New("providers[].name is required")
	ErrDisallowedProvider        = errors.New("this Terraform Provider is disallowed")
	ErrDisallowedProviderVersion = errors.New("this Terraform Provider version is disallowed")
)
