package server

import "errors"

var NoAuthTokenError = errors.New("no auth token found")
