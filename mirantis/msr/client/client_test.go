package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/Mirantis/terraform-provider-msr/mirantis/msr/client"
)

type testClientStruct struct {
	server           *httptest.Server
	expectedResponse client.HealthResponse
	expectedErr      error
}

func TestMSRClientHealthy(t *testing.T) {
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"error": "", "healthy":true}`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Error:   "",
			Healthy: true,
		},
		expectedErr: nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(tc.expectedResponse.Healthy, resp) {
		t.Errorf("expected (%v), got (%v)", tc.expectedResponse.Healthy, resp)
	}
	if err != tc.expectedErr {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientUnhealthy(t *testing.T) {
	unhealthyRes := client.HealthResponse{Healthy: false}
	bodyRes, _ := json.Marshal(unhealthyRes)
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: unhealthyRes,
		expectedErr:      nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create client new client")
	}
	ctx := context.Background()
	isHealthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(unhealthyRes.Healthy, isHealthy) {
		t.Errorf("expected (%v),\n got (%v)", unhealthyRes.Healthy, isHealthy)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientBadrequest(t *testing.T) {
	resError := client.ResponseError{
		Errors: []client.Errors{
			{
				Code:    strconv.Itoa(http.StatusBadRequest),
				Message: "Bad request",
			},
		},
	}
	bodyRes, err := json.Marshal(resError)
	if err != nil {
		t.Errorf("couldn't marshal struct %+v", resError)
		return
	}
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Healthy: false,
		},
		expectedErr: client.ErrResponseError,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Fatal("Couldn't create Client")
	}
	ctx := context.Background()
	healthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(healthy, tc.expectedResponse.Healthy) {
		t.Errorf("expected (%v),\n got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientUnauthorized(t *testing.T) {
	resError := client.ResponseError{
		Errors: []client.Errors{
			{
				Code:    strconv.Itoa(http.StatusUnauthorized),
				Message: "Bad creds",
			},
		},
	}
	bodyRes, err := json.Marshal(resError)
	if err != nil {
		t.Errorf("couldn't marshal struct %+v", resError)
		return
	}
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Healthy: false,
		},
		expectedErr: client.ErrUnauthorizedReq,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Fatalf("Couldn't create client")
	}
	ctx := context.Background()
	healthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(healthy, tc.expectedResponse.Healthy) {
		t.Errorf("expected (%v),\n got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestDoRequestWrongErrorStruct(t *testing.T) {
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Healthy: false,
		},
		expectedErr: client.ErrUnmarshaling,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Fatalf("Couldn't create client")
	}
	ctx := context.Background()
	healthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(healthy, tc.expectedResponse.Healthy) {
		t.Errorf("expected (%v),\n got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestDoRequestWrongErrorStructField(t *testing.T) {
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{"errors":[{"code":true, "message":"lol"}]}`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      client.ErrUnmarshaling,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Fatalf("Couldn't create client")
	}
	ctx := context.Background()
	healthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(healthy, tc.expectedResponse.Healthy) {
		t.Errorf("expected (%v),\n got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestEmptyUsernameField(t *testing.T) {
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{ "msg": "ok"`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      client.ErrEmptyClientArgs,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestEmptyPasswordField(t *testing.T) {
	tc := testClientStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{ "msg": "ok"`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      client.ErrEmptyClientArgs,
	}

	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}
