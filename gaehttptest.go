package gaehttptest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"appengine"
)

type ResponseWriter struct {
	StatusCode int
	header     http.Header
	Body       *bytes.Buffer
}

func newResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		StatusCode: 200,
		header:     make(http.Header),
		Body:       &bytes.Buffer{},
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func Get(c appengine.Context, urlStr string, urlValues url.Values) (*ResponseWriter, error) {
	if enc := urlValues.Encode(); enc != "" {
		urlStr = fmt.Sprintf("%v?%v", urlStr, enc)
	}
	return do(c, "GET", urlStr, "", nil)
}

func Post(c appengine.Context, urlStr string, bodyType string, body io.Reader) (*ResponseWriter, error) {
	return do(c, "POST", urlStr, bodyType, body)
}

func Put(c appengine.Context, urlStr string, bodyType string, body io.Reader) (*ResponseWriter, error) {
	return do(c, "PUT", urlStr, bodyType, body)
}

func Delete(c appengine.Context, urlStr string, urlValues url.Values) (*ResponseWriter, error) {
	if enc := urlValues.Encode(); enc != "" {
		urlStr = fmt.Sprintf("%v?%v", urlStr, enc)
	}
	return do(c, "DELETE", urlStr, "", nil)
}

func do(c appengine.Context, method, urlStr string, bodyType string, body io.Reader) (*ResponseWriter, error) {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		r.Header.Set("Content-Type", bodyType)
	}
	return Do(c, r)
}

func Do(c appengine.Context, r *http.Request) (*ResponseWriter, error) {
	r.Header.Set("X-AppEngine-Inbound-AppId", "dev~"+appengine.AppID(c))
	w := newResponseWriter()
	http.DefaultServeMux.ServeHTTP(w, r)
	return w, nil
}
