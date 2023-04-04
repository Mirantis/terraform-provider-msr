package client

import "errors"

var (
	ErrEmptyClientArgs = errors.New("MSR client did not receive host, username and/or password")
	ErrRequestCreation = errors.New("creating request failed in MSR client")
	ErrMarshaling      = errors.New("marshalling struct failed in MSR client")
	ErrUnmarshaling    = errors.New("unmarshalling struc failed in MSR client")
	ErrEmptyResError   = errors.New("request returned empty ResponseError struct in MSR client")
	ErrResponseError   = errors.New("request returned ResponseError in MSR client")
	ErrUnauthorizedReq = errors.New("unauthorized request in MSR client")
	ErrEmptyStruct     = errors.New("empty struct passed in MSR client")
	ErrInvalidFilter   = errors.New("passing invalid account retrieval filter in MSR client")
	ErrIDHasNoRepoName = errors.New("ID doesn't contain repository name in MSR client")
)
