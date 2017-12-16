package terraformenterprise

import (
	"net/url"
	"net/http"
	"os"
	"crypto/tls"
	
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-rootcerts"
)

const (
	// apiDefaultEndpoint is the default base URL for connecting
	// to the Terraform Enterprise API
	apiDefaultEndpoint = "https://atlas.hashicorp.com/api/v2"

	// authorizationHeader is the name of the Authorization HTTP header
	authorizationHeader = "Authorization"

	// contentTypeHeader is the name of the Content-Type HTTP header
	contentTypeHeader = "Content-Type"

	// defaultContentType is the content-type we should request from Terraform Enterprise
	defaultContentType = "application/vnd.api+json"
	
	// atlasCAFileEnvVar is the environment variable that causes the client to
	// load trusted certs from a file
	atlasCAFileEnvVar = "ATLAS_CAFILE"

	// atlasCAPathEnvVar is the environment variable that causes the client to
	// load trusted certs from a directory
	atlasCAPathEnvVar = "ATLAS_CAPATH"

	// atlasTLSNoVerifyEnvVar disables TLS verification, similar to curl -k
	// This defaults to false (verify) and will change to true (skip
	// verification) with any non-empty value
	atlasTLSNoVerifyEnvVar = "ATLAS_TLS_NOVERIFY"
)

type Client struct {
	// URL is the API endpoint
	URL *url.URL

	// token is the Authorization token
	token string

	// HTTPClient is the underlying http client with which to make requests.
	HTTPClient *http.Client

	// DefaultHeaders is a set of headers that will be added to every request.
	// This minimally includes the atlas user-agent string.
	DefaultHeader http.Header
}

func NewClient(token string) (*Client, error) {
	u, err := url.Parse(apiDefaultEndpoint)
	if err != nil {
		return nil, err
	}
	c := &Client{
		URL:   u,
		token: token,
		DefaultHeader: make(http.Header),
	}
	c.init()
	c.DefaultHeader.Set(contentTypeHeader, defaultContentType)
	return c, nil
}

// init() sets defaults on the client.
func (c *Client) init() error {
	c.HTTPClient = cleanhttp.DefaultClient()
	tlsConfig := &tls.Config{}
	if os.Getenv(atlasTLSNoVerifyEnvVar) != "" {
		tlsConfig.InsecureSkipVerify = true
	}
	err := rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
		CAFile: os.Getenv(atlasCAFileEnvVar),
		CAPath: os.Getenv(atlasCAPathEnvVar),
	})
	if err != nil {
		return err
	}
	t := cleanhttp.DefaultTransport()
	t.TLSClientConfig = tlsConfig
	c.HTTPClient.Transport = t
	return nil
}
