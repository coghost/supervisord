package supervisord

import (
	"fmt"
	"net/http"
)

// basicAuthTransport is an http.RoundTripper that wraps another http.RoundTripper
// and injects basic auth credentials into each request.
type basicAuthTransport struct {
	rt       http.RoundTripper
	username string
	password string
}

func newBasicAuth(username, password string) *basicAuthTransport {
	return &basicAuthTransport{
		rt:       http.DefaultTransport,
		username: username,
		password: password,
	}
}

func (b basicAuthTransport) String() string {
	return fmt.Sprintf("%s, %d", b.username, len(b.password))
}

func (b basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if b.username != "" && b.password != "" {
		req.SetBasicAuth(b.username, b.password)
	}

	return b.rt.RoundTrip(req)
}
