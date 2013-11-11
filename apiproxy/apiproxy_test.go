package apiproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestSameHost(t *testing.T) {
	destinationServer := httptest.NewServer(http.HandlerFunc(hostDumper))
	destinationHost := hostFor(destinationServer)

	reverseProxy := reverseProxyFor(destinationServer)
	defaultReverseProxyHost := hostFor(httptest.NewServer(reverseProxy))

	if defaultReverseProxyHost == destinationHost {
		t.Errorf("expected reverseProxy host %s to not equal destinationServer host %s", defaultReverseProxyHost, destinationHost)
	}

	correctedHost := hostFor(httptest.NewServer(sameHost(reverseProxy)))
	if correctedHost != destinationHost {
		t.Errorf("expected correctedHost %s to equal destinationServer host %s", correctedHost, destinationHost)
	}
}

func TestQueryCombiner(t *testing.T) {
	initialQuery := "q=some-test&x=y"
	toAdd := "key=1234&name=bob"
	expectedModified := toAdd + "&" + initialQuery

	destinationServer := httptest.NewServer(http.HandlerFunc(queryDumper))
	if queryFor(destinationServer, initialQuery) != initialQuery {
		t.Errorf("expected initial query %s to be unmodified but was %s", initialQuery, queryFor(destinationServer, initialQuery))
	}

	modifiedServer := httptest.NewServer(queryCombiner(http.HandlerFunc(queryDumper), toAdd))
	modifiedQuery := queryFor(modifiedServer, initialQuery)
	if modifiedQuery != expectedModified {
		t.Errorf("expected %s after modification but was %s", expectedModified, modifiedQuery)
	}
}

func hostFor(s *httptest.Server) string {
	resp, err := http.Get(s.URL)
	return readBody(resp, err)
}

func readBody(resp *http.Response, err error) string {
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	return string(content[:])
}

func queryFor(s *httptest.Server, baseQuery string) string {
	resp, err := http.Get(s.URL + "?" + baseQuery)
	return readBody(resp, err)
}

func reverseProxyFor(server *httptest.Server) *httputil.ReverseProxy {
	serverUrl, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal("URL failed to parse")
	}
	return httputil.NewSingleHostReverseProxy(serverUrl)
}

func hostDumper(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.Host)
}

func queryDumper(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.RawQuery)
}
