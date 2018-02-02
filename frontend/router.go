package frontend

import (
	"jivesearch/config"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"willnorris.com/go/imageproxy"
)

// Router sets up the routes & handlers
func (f *Frontend) Router(cfg config.Provider) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.NewRoute().Name("search").Methods("GET").Path("/").Handler(
		f.middleware(appHandler(f.searchHandler)),
	)
	router.NewRoute().Name("autocomplete").Methods("GET").Path("/autocomplete").Handler(
		f.middleware(appHandler(f.autocompleteHandler)),
	)
	router.NewRoute().Name("vote").Methods("POST").Path("/vote").Handler(
		f.middleware(appHandler(f.voteHandler)),
	)
	router.NewRoute().Name("favicon").Methods("GET").Path("/favicon.ico").Handler(
		http.FileServer(http.Dir("static")),
	)

	// How do we exclude viewing the entire static directory of /static path?
	router.NewRoute().Name("static").Methods("GET").PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// make hmac key available to our templates
	key := cfg.GetString("hmac.secret")
	os.Setenv("hmac_secret", key)

	p := imageproxy.NewProxy(nil, nil)
	p.Verbose = false // otherwise logs the image fetched
	//p.UserAgent = cfg.GetString("useragent") // not implemented yet: https://github.com/willnorris/imageproxy/pull/83
	p.SignatureKey = []byte(key)
	p.Timeout = 2 * time.Second
	router.NewRoute().Name("image").Methods("GET").PathPrefix("/image/").Handler(http.StripPrefix("/image", p))

	/* To generate new HMAC secret...
	// DON'T RUN IN PLAYGROUND! Will get same secret each time ;)
	b := make([]byte, 32) // s/b at least 32
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	fmt.Println(base64.URLEncoding.EncodeToString(b))
	*/

	return router
}
