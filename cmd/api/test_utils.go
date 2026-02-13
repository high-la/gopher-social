package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/high-la/gopher-social/internal/auth"
	"github.com/high-la/gopher-social/internal/store"
	"github.com/high-la/gopher-social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {

	t.Helper()

	// logger := zap.NewNop().Sugar() // for less detailed logs
	logger := zap.Must(zap.NewProduction()).Sugar() // for detail(including stack trace) logs
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {

	reqRecorder := httptest.NewRecorder()
	mux.ServeHTTP(reqRecorder, req)

	return reqRecorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}

}
