package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// MSRAPIVERSION the
	MSRAPIVERSION = "api/v0"
	ENZIENDPOINT  = "enzi/v0"

	// MsrURL - Default MSR URL
	DEFAULTMSRURL = "http://localhost:80"
)

// Client MSR client
type Client struct {
	MsrURL     string
	HTTPClient *http.Client
	Creds      AuthStruct
}

// AuthStruct basicauth struct
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Errors struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseError structure from MSR
type ResponseError struct {
	Errors []Errors `json:"errors"`
}

// NewDefaultClient creates a new MSR SSL safe Client
func NewDefaultClient(host, username, password string) (Client, error) {
	if username == "" || password == "" || host == "" {
		return Client{}, ErrEmptyClientArgs
	}

	return NewClient(username, password, host, &http.Client{})

}

// NewUnsafeSSLClient creates a new unsafe MSR HTTP Client
func NewUnsafeSSLClient(host, username, password string) (Client, error) {
	if username == "" || password == "" || host == "" {
		return Client{}, ErrEmptyClientArgs
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return NewClient(username, password, host, &http.Client{Transport: tr})
}

// NewClient creates a new MSR API Client from raw components
func NewClient(username, password, MsrURL string, HTTPClient *http.Client) (Client, error) {
	creds := AuthStruct{
		Username: username,
		Password: password,
	}
	return Client{
		HTTPClient: HTTPClient,
		MsrURL:     MsrURL,
		Creds:      creds,
	}, nil
}

// doRequest - performing the actual HTTP request
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.Creds.Username, c.Creds.Password)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("%w: Status code: %d", ErrUnauthorizedReq, res.StatusCode)
		}
		errStruct := &ResponseError{}

		if err := json.Unmarshal(body, errStruct); err != nil {
			return nil, fmt.Errorf("%w: Status code: %d", ErrUnmarshaling, res.StatusCode)
		}

		if len(errStruct.Errors) <= 0 {
			return nil, fmt.Errorf("%w: Status code: %d", ErrEmptyResError, res.StatusCode)
		}

		errMsg := errors.New(errStruct.Errors[0].Message)

		return nil, fmt.Errorf("%w: Status code: %d. ErrMsg: %s", ErrResponseError, res.StatusCode, errMsg)
	}

	return body, err
}

func (c *Client) createMsrUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", c.MsrURL, MSRAPIVERSION, endpoint)
}

func (c *Client) createEnziUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", c.MsrURL, ENZIENDPOINT, endpoint)
}
