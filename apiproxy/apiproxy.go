package apiproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Set the proxied request's host to the destination host (instead of the
// source host).  e.g. http://foo.com proxying to http://bar.com will ensure
// that the proxied requests appear to be coming from http://bar.com
//
// For both this function and queryCombiner (below), we'll be wrapping a
// Handler with our own HandlerFunc so that we can do some intermediate work
func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

// Append additional query params to the original URL query.
func queryCombiner(handler http.Handler, addon string) http.Handler {
	// first parse the provided string to pull out the keys and values
	values, err := url.ParseQuery(addon)
	if err != nil {
		log.Fatal("addon failed to parse")
	}

	// now we apply our addon params to the existing query
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		for k, _ := range values {
			query.Add(k, values.Get(k))
		}

		r.URL.RawQuery = query.Encode()
		handler.ServeHTTP(w, r)
	})
}

// Allow cross origin resource sharing
func addCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		handler.ServeHTTP(w, r)
	})
}

// Combine the two functions above with http.NewSingleHostReverseProxy
func Proxy(remoteUrl string, queryAddon string) http.Handler {
	// pull the root url we're proxying to from an environment variable.
	serverUrl, err := url.Parse(remoteUrl)
	if err != nil {
		log.Fatal("URL failed to parse")
	}

	// initialize our reverse proxy
	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)
	// wrap that proxy with our sameHost function
	singleHosted := sameHost(reverseProxy)
	// wrap that with our query param combiner
	combined := queryCombiner(singleHosted, queryAddon)
	// and finally allow CORS
	return addCORS(combined)
}
