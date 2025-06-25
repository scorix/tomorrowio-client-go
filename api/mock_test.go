package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockServer struct {
	server *httptest.Server
	t      *testing.T
}

func NewMockServer(t *testing.T) *MockServer {
	m := &MockServer{t: t}
	m.server = httptest.NewServer(m.Handler())
	t.Cleanup(m.server.Close)

	return m
}

func (m *MockServer) URL() string {
	return m.server.URL
}

func (m *MockServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v4/weather/forecast", mockWeatherForecast(m.t))
	return mux
}

func mockWeatherForecast(t *testing.T) http.HandlerFunc {
	t.Helper()

	jsonData, err := os.ReadFile("testdata/forecast-response.json")
	require.NoError(t, err)

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, bytes.NewReader(jsonData))
		require.NoError(t, err)
	}
}
