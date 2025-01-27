package tfe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
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

	// userAgentHeader is the name of the User-Agent HTTP header
	userAgentHeader = "User-Agent"

	// defaultContentType is the content-type we should request from Terraform Enterprise
	defaultContentType = "application/vnd.api+json"

	// defaultUserAgent is the default User-Agent for HTTP requests
	defaultUserAgent = "terraform-enterprise-go/0.1"

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
	HTTPClient *retryablehttp.Client

	// DefaultHeaders is a set of headers that will be added to every request.
	// This minimally includes the atlas user-agent string.
	DefaultHeader http.Header
}

type ClientOptions struct {
	BaseURL       string
	DefaultHeader http.Header
	NoVerifyTLS   bool
	CAPath        string
	CAFile        string
}

func (c *ClientOptions) SetToken(token string) {
	c.DefaultHeader.Set(authorizationHeader, fmt.Sprintf("Bearer %s", token))
}

func DefaultClientOptions() *ClientOptions {
	header := make(http.Header)
	header.Set(contentTypeHeader, defaultContentType)
	header.Set(userAgentHeader, defaultUserAgent)
	opts := &ClientOptions{
		BaseURL:       defaultBaseURL,
		DefaultHeader: header,
		NoVerifyTLS:   os.Getenv(atlasTLSNoVerifyEnvVar) != "",
		CAPath:        os.Getenv(atlasCAPathEnvVar),
		CAFile:        os.Getenv(atlasCAFileEnvVar),
	}
	opts.SetToken(os.Getenv(atlasTokenEnvVar))
	return opts
}

// retryPolicy implements retryablehttp.RetryPolicy. It copies the
// default behavior (retry on internal server errors and client
// network failures) and also retries on HTTP 429 too many requests.
func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	if resp.StatusCode == 429 {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func NewClient(opts *ClientOptions) (*Client, error) {
	u, err := url.Parse(opts.BaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		URL:           u,
		DefaultHeader: opts.DefaultHeader,
	}

	client.HTTPClient = retryablehttp.NewClient()
	client.HTTPClient.CheckRetry = retryPolicy
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
	client.HTTPClient.HTTPClient.Transport = t
	return client, nil
}
