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
	tests := []struct {
		name       string
		cfg        plugin.ConfigHeaders
		wantHeader string
		wantValue  string
	}{
		{
			name:       "full address",
			cfg:        plugin.ConfigHeaders{Address: "X-Real-Address"},
			wantHeader: "X-Real-Address",
			wantValue:  "127.0.0.1:1234",
		},
		{
			name:       "only ip",
			cfg:        plugin.ConfigHeaders{IP: "X-Real-IP"},
			wantHeader: "X-Real-IP",
			wantValue:  "127.0.0.1",
		},
		{
			name:       "only port",
			cfg:        plugin.ConfigHeaders{Port: "X-Real-Port"},
			wantHeader: "X-Real-Port",
			wantValue:  "1234",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := plugin.CreateConfig()
			cfg.Headers = tc.cfg
			req := testPlugin(t, cfg)
			assertHeader(t, req.Header, tc.wantHeader, tc.wantValue)
		})
	}
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

	req.RemoteAddr = "127.0.0.1:1234"
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
