package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal()
	}

	// test 1 UnAuthenticated user
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {

		// check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// .
		reqRecorder := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, reqRecorder.Code)
	})

	// test 2 Authenticated user
	t.Run("should allow authenticated requests", func(t *testing.T) {

		// check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		reqRecorder := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, reqRecorder.Code)
	})
}
