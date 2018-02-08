package tfe

import (
	"net/http"
	"testing"

	"github.com/hashicorp/atlas-go/v2/testutils"
	"github.com/stretchr/testify/assert"
)

func testingClientServer(t *testing.T) (*Client, *testutils.Server) {
	server := testutils.NewTestServer(t)
	header := make(http.Header)
	header.Set(authorizationHeader, "Bearer abcd1234")
	header.Set(contentTypeHeader, defaultContentType)
	opts := &ClientOptions{
		BaseURL:       server.BaseURL.String(),
		DefaultHeader: header,
		NoVerifyTLS:   true,
		CAPath:        "",
		CAFile:        "",
	}
	client, err := NewClient(opts)
	assert.NoError(t, err)
	return client, server
}
