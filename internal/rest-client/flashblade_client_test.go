package client


import (
	"testing"
        "regexp"
        "strings"
	"net/http"
	"net/http/httptest"
	"os"
)

func TestNewRestClient(t *testing.T) {

	vers, _ := os.ReadFile("../../test/data/versions.json")
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	        valid := regexp.MustCompile(`^/api/login$`)
		if r.URL.Path == "/api/api_version" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(vers))
		} else if valid.MatchString(r.URL.Path) {
			w.Header().Set("x-auth-token", "faketoken")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items":[{"username":"fakeuser"}]}`))
		}
	   }))
        endp := strings.Split(server.URL, "/")
        e := endp[len(endp)-1]
        t.Run("login", func(t *testing.T) {
            defer server.Close()
            c := NewRestClient(e, "fake-api-token", "latest", false)
            if c.EndPoint != e || c.ApiToken != "fake-api-token" || c.XAuthToken != "faketoken" {
                t.Errorf("expected (%s, fake-api-token, faketoken), got (%s %s %s)", e, c.EndPoint, c.ApiToken, c.XAuthToken)
            }
        })
}
