package cli

import "errors"

var ErrNoStdin = errors.New(`stdin is missing. Please pass the result of "terraform version -json" to stdin`)
