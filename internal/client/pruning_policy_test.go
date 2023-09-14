package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
)

type testPruningPolicyStruct struct {
	server           *httptest.Server
	expectedResponse client.ResponsePruningPolicy
	expectedErr      error
}

func TestCreateValidPruningPolicy(t *testing.T) {
	testResPolicy := client.ResponsePruningPolicy{
		ID:      "fake-test-id",
		Enabled: true,
		Rules: []client.PruningPolicyRuleAPI{
			{
				Field:    "tag",
				Operator: "eq",
				Values:   []string{"test"},
			},
		},
	}
	mAccount, err := json.Marshal(testResPolicy)
	if err != nil {
		t.Error(err)
	}
	tc := testPruningPolicyStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(mAccount); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: testResPolicy,
		expectedErr:      nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass", true)
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.CreatePruningPolicy(ctx, "fake", "fake", client.CreatePruningPolicy{
		Enabled: true,
		Rules: []client.PruningPolicyRuleAPI{
			{
				Field:    "tag",
				Operator: "eq",
				Values:   []string{"test"},
			},
		},
	})
	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestCreateInvalidPruningPolicy(t *testing.T) {
	testResPolicy := client.ResponsePruningPolicy{}
	mAccount, err := json.Marshal(testResPolicy)
	if err != nil {
		t.Fatal(err)
	}
	tc := testPruningPolicyStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(mAccount); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: testResPolicy,
		expectedErr:      client.ErrEmptyResError,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass", true)
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()

	resp, err := testClient.CreatePruningPolicy(ctx, "fake", "fake", client.CreatePruningPolicy{
		Enabled: true,
		Rules: []client.PruningPolicyRuleAPI{
			{
				Field:    "tag",
				Operator: "eq",
				Values:   []string{"test"},
			},
		},
	})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestCreatePruningPolicyErrUmarshaling(t *testing.T) {
	tc := testPruningPolicyStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponsePruningPolicy{},
		expectedErr:      client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass", true)
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()

	resp, err := testClient.CreatePruningPolicy(ctx, "fake", "fake", client.CreatePruningPolicy{})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestUpdatePruningPolicyFailed(t *testing.T) {
	tc := testPruningPolicyStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponsePruningPolicy{},
		expectedErr:      client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass", true)
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.UpdatePruningPolicy(ctx, "fakeid", "fakeid", client.CreatePruningPolicy{Enabled: true}, "fakeid")

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}
