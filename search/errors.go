package search

import "errors"

// ErrInvalidSearchResultCodeReturned - error raised when a search request to the Transfer360 API server returns a non 200 result code
var ErrInvalidSearchResultCodeReturned = errors.New("unexpected search result code returned")

// ErrInvalidSearchResultBody - error raised when a JSON body is expected on sending or return
var ErrInvalidSearchResultBody = errors.New("missing or invalid search result body returned")

var ErrTimeOutStatusCode = errors.New("time out code (504)")
var ErrUnableToHandleStatusCode = errors.New("time out code (503)")
