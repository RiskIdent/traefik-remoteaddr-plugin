package plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/RiskIdent/traefik-remoteaddr-plugin"
)

func TestInvalidConfig(t *testing.T) {
	cfg := plugin.CreateConfig()
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	_, err := plugin.New(context.Background(), next, cfg, "traefik-remoteaddr-plugin")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHeaderAddress(t *testing.T) {
	cfg := plugin.CreateConfig()
	cfg.Headers.Address = "X-Real-Address"
	req := testPlugin(t, cfg)
	assertHeader(t, req.Header, "X-Real-Address", "localhost:1234")
}

func TestHeaderIP(t *testing.T) {
	cfg := plugin.CreateConfig()
	cfg.Headers.IP = "X-Real-IP"
	req := testPlugin(t, cfg)
	assertHeader(t, req.Header, "X-Real-IP", "localhost")
}

func TestHeaderPort(t *testing.T) {
	cfg := plugin.CreateConfig()
	cfg.Headers.Port = "X-Real-Port"
	req := testPlugin(t, cfg)
	assertHeader(t, req.Header, "X-Real-Port", "1234")
}

func testPlugin(t *testing.T, cfg *plugin.Config) *http.Request {
	t.Helper()
	ctx := context.Background()
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	handler, err := plugin.New(ctx, next, cfg, "traefik-remoteaddr-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "localhost:1234"
	handler.ServeHTTP(recorder, req)

	t.Logf("request headers: %d", len(req.Header))
	for k, vals := range req.Header {
		for _, v := range vals {
			t.Logf("  %s=%q", k, v)
		}
	}

	return req
}

func assertHeader(t *testing.T, header http.Header, key, expected string) {
	t.Helper()

	if header.Get(key) != expected {
		t.Errorf("invalid header value\nwant: %s=%q\ngot:  %s=%q", key, expected, key, header.Get(key))
	}
}
