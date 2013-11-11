// # Api Proxy
//
// ## Purpose:
// Provides a fast, intermediate server to conceal your api keys when proxying
// to an api
//
// ## Inputs:
// 1. PORT: which port to run on ( e.g. 8080 )
// 2. URL_ROOT: the destination root url ( e.g. http://foo.com/api/ )
// 3. URL_SUFFIX: the query params you wish to transparently append to the each
//	  request. ( e.g. a suffix of "language=en&key=XXXXXX" would result in the
//    url http://foo.com/test?name=bob being requested transparently as
//    http://foo.com/test?name=bob&language=en&key=XXXXXX )

package main

import (
	"github.com/semanticart/simpleapiproxy/apiproxy"
	"net/http"
	"os"
)

func main() {
	apiProxy := apiproxy.Proxy(os.Getenv("URL_ROOT"), os.Getenv("URL_SUFFIX"))
	http.ListenAndServe(":"+os.Getenv("PORT"), apiProxy)
}
