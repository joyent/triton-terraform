package helpers

import (
	local "github.com/joyent/gosdc/localservices/cloudapi"
	"github.com/joyent/gosign/auth"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// Server creates a local test double for testing API responses
type Server struct {
	Server *httptest.Server
	API    *local.CloudAPI
	Creds  *auth.Credentials
}

// NewServer returns a Server
func NewServer() (*Server, error) {
	s := new(Server)

	mux := http.NewServeMux()
	s.Server = httptest.NewServer(mux)

	key, err := ioutil.ReadFile(TestKeyFile)
	if err != nil {
		return nil, err
	}

	authentication, err := auth.NewAuth(TestAccount, string(key), "rsa-sha256")

	s.Creds = &auth.Credentials{
		UserAuthentication: authentication,
		SdcKeyId:           TestKeyID,
		SdcEndpoint:        auth.Endpoint{URL: s.Server.URL},
	}

	s.API = local.New(s.Server.URL, TestAccount)
	s.API.SetupHTTP(mux)

	return s, nil
}

// URL returns the URL of the server
func (s *Server) URL() string {
	return s.Server.URL
}

// Stop stops the server for teardown
func (s *Server) Stop() {
	s.Server.Close()
}