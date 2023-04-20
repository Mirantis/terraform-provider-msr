package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Mirantis/terraform-provider-msr/mirantis/msr/client"
)

type testRepoStruct struct {
	server           *httptest.Server
	expectedResponse client.ResponseRepo
	expectedErr      error
}

func TestCreateValidRepo(t *testing.T) {
	testRepo := client.ResponseRepo{
		Name: "test-repo",
	}
	mRepo, err := json.Marshal(testRepo)
	if err != nil {
		t.Error(err)
	}
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write(mRepo); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: testRepo,
		expectedErr:      nil,
	}
	defer tc.server.Close()

	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.CreateRepo(ctx, "fakeid", client.CreateRepo{Name: testRepo.Name})
	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestCreateInvalidRepo(t *testing.T) {
	testRepo := client.ResponseRepo{}

	mRepo, err := json.Marshal(testRepo)
	if err != nil {
		t.Fatal(err)
	}
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(mRepo); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: testRepo,
		expectedErr:      client.ErrEmptyResError,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()

	resp, err := testClient.CreateRepo(ctx, "fakeid", client.CreateRepo{Name: "fake"})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestCreateEmptyRepo(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponseRepo{},
		expectedErr:      client.ErrEmptyStruct,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()

	resp, err := testClient.CreateRepo(ctx, "fakeid", client.CreateRepo{})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestCreateRepoUnmarshalErr(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponseRepo{},
		expectedErr:      client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()

	resp, err := testClient.CreateRepo(ctx, "fakeid", client.CreateRepo{Name: "fake"})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestDeleteRepoSuccess(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedErr: nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	err = testClient.DeleteRepo(ctx, "fakename")

	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestDeleteRepoFailed(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedErr: client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	err = testClient.DeleteAccount(ctx, "fakename")

	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestReadRepoSuccess(t *testing.T) {
	resRepo := client.ResponseRepo{
		Name: "fake repo",
	}
	mResRepo, err := json.Marshal(resRepo)
	if err != nil {
		t.Fatal(err)
	}
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(mResRepo); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: resRepo,
		expectedErr:      nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.ReadRepo(ctx, "fakeorg", "fakename")

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestReadRepoFailed(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponseRepo{},
		expectedErr:      client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.ReadRepo(ctx, "fakeorg", "fakename")

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestUpdateRepoSuccess(t *testing.T) {
	uRepo := client.ResponseRepo{
		Name:       "fakename",
		ScanOnPush: true,
	}
	mURepo, err := json.Marshal(uRepo)
	if err != nil {
		t.Fatal(err)
	}
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(mURepo); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: uRepo,
		expectedErr:      nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.UpdateRepo(ctx, "fakeorg", "fakeid", client.UpdateRepo{ScanOnPush: true})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestUpdateRepoEmpty(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponseRepo{},
		expectedErr:      client.ErrEmptyStruct,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.UpdateRepo(ctx, "fakeorg", "fakename", client.UpdateRepo{})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}

func TestUpdateRepoFailed(t *testing.T) {
	tc := testRepoStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.ResponseRepo{},
		expectedErr:      client.ErrUnmarshaling,
	}
	defer tc.server.Close()
	testClient, err := client.NewDefaultClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.UpdateRepo(ctx, "fakeorg", "fakeid", client.UpdateRepo{ScanOnPush: true})

	if !reflect.DeepEqual(tc.expectedResponse, resp) {
		t.Errorf("expected resp: (%+v),\n got (%+v)", tc.expectedResponse, resp)
	}
	if !errors.Is(err, tc.expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", tc.expectedErr, err)
	}
}
