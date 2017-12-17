package terraformenterprise

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

// RequestOptions is the list of options to pass to the request.
type RequestOptions struct {
	// Params is a map of key-value pairs that will be added to the Request.
	Params map[string]string

	// Headers is a map of key-value pairs that will be added to the Request.
	Headers map[string]string

	// Body is an io.Reader object that will be streamed or uploaded with the
	// Request. BodyLength is the final size of the Body.
	Body       io.Reader
	BodyLength int64
}

// Request creates a new HTTP request using the given verb and sub path.
func (c *Client) NewRequest(verb, spath string, ro *RequestOptions) (*http.Request, error) {
	// Ensure we have a RequestOptions struct (passing nil is an acceptable)
	if ro == nil {
		ro = new(RequestOptions)
	}

	// Create a new URL with the appended path
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	log.Printf("[INFO] request: %s %s", verb, u.Path)

	return c.rawRequest(verb, &u, ro)
}

// rawRequest accepts a verb, URL, and RequestOptions struct and returns the
// constructed http.Request and any errors that occurred
func (c *Client) rawRequest(verb string, u *url.URL, ro *RequestOptions) (*http.Request, error) {
	if verb == "" {
		return nil, fmt.Errorf("client: missing verb")
	}

	if u == nil {
		return nil, fmt.Errorf("client: missing URL.url")
	}

	if ro == nil {
		return nil, fmt.Errorf("client: missing RequestOptions")
	}

	// Add the token and other params
	var params = make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()

	// Create the request object
	request, err := http.NewRequest(verb, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}

	// set our default headers first
	for k, v := range c.DefaultHeader {
		request.Header[k] = v
	}

	// Add any request headers (auth will be here if set)
	for k, v := range ro.Headers {
		request.Header.Add(k, v)
	}

	// Add content-length if we have it
	if ro.BodyLength > 0 {
		request.ContentLength = ro.BodyLength
	}

	log.Printf("[DEBUG] raw request: %#v", request)

	return request, nil
}

// ErrAuth is the error returned if a 401 is returned by an API request.
var ErrAuth = fmt.Errorf("authentication failed")

// ErrNotFound is the error returned if a 404 is returned by an API request.
var ErrNotFound = fmt.Errorf("resource not found")

// CheckResp wraps http.Client.Do() and verifies that the request was
// successful. A non-200 request returns an error formatted to included any
// validation problems or otherwise.
func CheckResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that
	if err != nil {
		return resp, err
	}

	log.Printf("[INFO] response: %d (%s)", resp.StatusCode, resp.Status)

	switch resp.StatusCode {
	case 200:
		return resp, nil
	case 201:
		return resp, nil
	case 202:
		return resp, nil
	case 204:
		return resp, nil
	// case 400:
	// 	return nil, parseErr(resp)
	case 401:
		return nil, ErrAuth
	case 404:
		return nil, ErrNotFound
	// case 422:
	// 	return nil, parseErr(resp)
	default:
		return nil, fmt.Errorf("client: %s", resp.Status)
	}
}
