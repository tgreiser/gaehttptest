HTTP testing for Google Appengine Go
====================================

Example:

    import (
        "net/http"
        "net/url"

        httptest "github.com/vmihailenco/gaehttptest"
    )

    func TestGetPosts(c appengine.context) {
        w, err := httptest.Get(c, "/posts", url.Values{})
        if err != nil {
            panic(err)
        }
        if code := w.StatusCode; code != http.StatusOK {
            panic(fmt.Sprintf("got %v, expected 200 OK", code))
        }
        if ct := w.Header.Get("Content-Type"); ct != "text/html" {
            panic("got %v, expected text/html", ct)
        }
        if body := w.Body.String(); body != "<html></html>" {
            panic(fmt.Sprintf("unexpected response: %v", body))
        }
    }
