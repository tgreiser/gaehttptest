package gaehttptest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"appengine"
)

type responseWriter struct {
	StatusCode int
	header     http.Header
	Body       *bytes.Buffer
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		StatusCode: 200,
		header:     make(http.Header),
		Body:       &bytes.Buffer{},
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func Get(c appengine.Context, urlStr string, urlValues url.Values) (*responseWriter, error) {
	if enc := urlValues.Encode(); enc != "" {
		urlStr = fmt.Sprintf("%v?%v", urlStr, enc)
	}
	return do(c, "GET", urlStr, "", nil)
}

func Post(c appengine.Context, urlStr string, bodyType string, body io.Reader) (*responseWriter, error) {
	return do(c, "POST", urlStr, bodyType, body)
}

func Put(c appengine.Context, urlStr string, bodyType string, body io.Reader) (*responseWriter, error) {
	return do(c, "PUT", urlStr, bodyType, body)
}

func Delete(c appengine.Context, urlStr string, urlValues url.Values) (*responseWriter, error) {
	if enc := urlValues.Encode(); enc != "" {
		urlStr = fmt.Sprintf("%v?%v", urlStr, enc)
	}
	return do(c, "DELETE", urlStr, "", nil)
}

func do(c appengine.Context, method, urlStr string, bodyType string, body io.Reader) (*responseWriter, error) {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		r.Header.Set("Content-Type", bodyType)
	}
	return RoundTrip(c, r)
}

func RoundTrip(c appengine.Context, r *http.Request) (*responseWriter, error) {
	r.Header.Set("X-AppEngine-Inbound-AppId", "dev~"+appengine.AppID(c))
	w := newResponseWriter()
	http.DefaultServeMux.ServeHTTP(w, r)
	return w, nil
}
