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
	// defaultBaseURL is the default base URL for connecting
	// to the Terraform Enterprise API
	defaultBaseURL = "https://atlas.hashicorp.com/api/v2"

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

	// atlasTokenEnvVar is the name of the env var that contains
	// the user oauth token to talk to Terraform Enterprise
	atlasTokenEnvVar = "ATLAS_TOKEN"
)

type Client struct {
	// URL is the API endpoint
	URL *url.URL

	// HTTPClient is the underlying http client with which to make requests.
	HTTPClient *http.Client

	// DefaultHeaders is a set of headers that will be added to every request.
	// This minimally includes the atlas user-agent string.
	DefaultHeader http.Header
}

type ClientOptions struct {
	BaseUrl string
	DefaultHeader http.Header
	NoVerifyTLS bool
	CAPath string
	CAFile string
}

func DefaultClientOptions() *ClientOptions {
	header := make(http.Header)
	header.Set(contentTypeHeader, defaultContentType)
	header.Set(authorizationHeader, os.Getenv(atlasTokenEnvVar))
	return &ClientOptions{
		BaseUrl: defaultBaseURL,
		DefaultHeader: header,
		NoVerifyTLS: os.Getenv(atlasTLSNoVerifyEnvVar) != "",
		CAPath: os.Getenv(atlasCAPathEnvVar),
		CAFile: os.Getenv(atlasCAFileEnvVar),
	}
}

func NewClient(opts *ClientOptions) (*Client, error) {
	u, err := url.Parse(opts.BaseUrl)
	if err != nil {
		return nil, err
	}
	client := &Client{
		URL:   u,
		DefaultHeader: opts.DefaultHeader,
	}

	client.HTTPClient = cleanhttp.DefaultClient()
	tlsConfig := &tls.Config{}
	if opts.NoVerifyTLS {
		tlsConfig.InsecureSkipVerify = true
	}
	err = rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
		CAFile: opts.CAFile,
		CAPath: opts.CAPath,
	})
	if err != nil {
		return nil, err
	}
	t := cleanhttp.DefaultTransport()
	t.TLSClientConfig = tlsConfig
	client.HTTPClient.Transport = t
	return client, nil
}
