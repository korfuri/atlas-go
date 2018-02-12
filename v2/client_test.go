package tfe

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	client, server := TestingClientServer(t)
	defer server.Stop()

	server.Mux.HandleFunc("/somepath", func(rw http.ResponseWriter, r *http.Request) {})

	request, err := client.NewRequest("GET", "/somepath", nil)
	assert.NoError(t, err)
	response, err := client.HTTPClient.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, response.Status, "200 OK")
}
