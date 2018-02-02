package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jivesearch/jivesearch/config"
	"github.com/jivesearch/jivesearch/log"
	"github.com/jivesearch/jivesearch/wikipedia"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type fetcher struct {
	wikipedia.Fetcher
}

func (f *fetcher) handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", f.wikiHandler)
	r.HandleFunc("/favicon.ico", favHandler)
	return r
}

func (f *fetcher) wikiHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.FormValue("q"))

	// Language parsing is outside the scope of
	// this package. See frontend/search.go for an example.
	l := language.MustParse("en")

	item, err := f.Fetch(q, l)
	if err != nil {
		log.Info.Println(err)
	}

	j, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func favHandler(w http.ResponseWriter, r *http.Request) {}

func setup() (*wikipedia.PostgreSQL, error) {
	v := viper.New()
	v.SetEnvPrefix("jivesearch")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetDefaults(v)

	if v.GetBool("debug") {
		log.Debug.SetOutput(os.Stdout)
	}

	var err error

	p := &wikipedia.PostgreSQL{}

	p.DB, err = sql.Open("postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s database=%s sslmode=require",
			v.GetString("postgresql.user"),
			v.GetString("postgresql.password"),
			v.GetString("postgresql.host"),
			v.GetString("postgresql.database"),
		),
	)

	p.DB.SetMaxIdleConns(0)

	return p, err
}

func main() {
	p, err := setup()
	if err != nil {
		panic(err)
	}

	defer p.DB.Close()

	f := &fetcher{p}

	if err := f.Setup(); err != nil {
		panic(err)
	}

	port := 8000
	log.Info.Printf("Listening at http://localhost:%d", port)
	log.Info.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), f.handler()))
}
