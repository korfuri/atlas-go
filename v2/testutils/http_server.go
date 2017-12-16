package testutils

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Server struct {
	BaseURL *url.URL

	t *testing.T

	listener net.Listener

	server *http.Server

	Mux *http.ServeMux

	workspaces []string
}

func NewTestServer(t *testing.T) *Server {
	srv := &Server{t: t}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv.listener = listener

	srv.BaseURL = &url.URL{
		Scheme: "http",
		Host:   listener.Addr().String(),
	}

	srv.Mux = http.NewServeMux()
	server := &http.Server{}
	server.Handler = srv.Mux
	srv.server = server
	go server.Serve(listener)

	return srv
}

func (srv *Server) Stop() {
	srv.listener.Close()
}

// CheckHeaders asserts that the Authorization and Content-Type
// type Server struct headers are correct in an incoming request.
func (srv *Server) CheckHeaders(req *http.Request) {
	assert.Contains(srv.t, req.Header, "Authorization")
	assert.Contains(srv.t, req.Header, "Content-Type")
	assert.Equal(srv.t, []string{"Bearer abcd1234"}, req.Header["Authorization"])
	assert.Equal(srv.t, []string{"application/vnd.api+json"}, req.Header["Content-Type"])
}
