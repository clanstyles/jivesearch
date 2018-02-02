// Command frontend demonstrates how to run the web app
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/jivesearch/jivesearch/bangs"
	"github.com/jivesearch/jivesearch/config"
	"github.com/jivesearch/jivesearch/frontend"
	"github.com/jivesearch/jivesearch/log"
	"github.com/jivesearch/jivesearch/search"
	"github.com/jivesearch/jivesearch/search/document"
	"github.com/jivesearch/jivesearch/search/vote"
	"github.com/jivesearch/jivesearch/suggest"
	"github.com/jivesearch/jivesearch/wikipedia"
	"github.com/lib/pq"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

var (
	f          *frontend.Frontend
	httpClient = &http.Client{
		Timeout: 2 * time.Second,
	}
)

func setup(v *viper.Viper) *http.Server {
	v.SetEnvPrefix("jivesearch")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetDefaults(v)

	if v.GetBool("debug") {
		log.Debug.SetOutput(os.Stdout)
	}

	frontend.ParseTemplates()
	f = &frontend.Frontend{}

	f.Bangs = bangs.New()

	router := f.Router(v)

	return &http.Server{
		Addr:    ":" + strconv.Itoa(v.GetInt("frontend.port")),
		Handler: http.TimeoutHandler(router, 5*time.Second, "Sorry, we took too long to get back to you"),
	}
}

func main() {
	v := viper.New()
	s := setup(v)

	// Set the backend for our core search results
	client, err := elastic.NewClient(
		elastic.SetURL(v.GetString("elasticsearch.url")),
		elastic.SetSniff(false),
	)

	if err != nil {
		panic(err)
	}

	f.Search = &search.ElasticSearch{
		ElasticSearch: &document.ElasticSearch{
			Client: client,
			Index:  v.GetString("elasticsearch.search.index"),
			Type:   v.GetString("elasticsearch.search.type"),
		},
	}

	// Set the backend for our autocomplete & phrase suggestor
	f.Suggest = &suggest.ElasticSearch{
		Client: client,
		Index:  v.GetString("elasticsearch.query.index"),
		Type:   v.GetString("elasticsearch.query.type"),
	}

	exists, err := f.Suggest.IndexExists()
	if err != nil {
		panic(err)
	}

	if !exists {
		if err := f.Suggest.Setup(); err != nil {
			panic(err)
		}
	}

	// Setup the voting backend. Tables will be setup automatically.
	// The database needs to be setup beforehand.
	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s database=%s sslmode=require",
			v.GetString("postgresql.user"),
			v.GetString("postgresql.password"),
			v.GetString("postgresql.host"),
			v.GetString("postgresql.database"),
		),
	)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetMaxIdleConns(0)

	f.Vote = &vote.PostgreSQL{
		DB:    db,
		Table: v.GetString("postgresql.votes.table"),
	}

	if err := f.Vote.Setup(); err != nil {
		switch err.(type) {
		case *pq.Error:
			if err.(*pq.Error).Error() != vote.ErrScoreFnExists.Error() {
				panic(err)
			}
		default:
			panic(err)
		}
	}

	// supported languages
	supported, unsupported := languages(v)
	for _, lang := range unsupported {
		log.Info.Printf("wikipedia does not support langugage %q\n", lang)
	}

	f.Wikipedia.Matcher = language.NewMatcher(supported)
	f.Wikipedia.Fetcher = &wikipedia.PostgreSQL{
		DB: db,
	}

	if err := f.Wikipedia.Setup(); err != nil {
		panic(err)
	}

	// see notes on customizing languages in search/document/document.go
	f.Document.Languages = document.Languages(supported)
	f.Document.Matcher = language.NewMatcher(f.Document.Languages)

	log.Info.Printf("Listening at http://127.0.0.1%v", s.Addr)
	log.Info.Fatal(s.ListenAndServe())
}

func languages(cfg config.Provider) ([]language.Tag, []language.Tag) {
	supported := []language.Tag{}

	for _, l := range cfg.GetStringSlice("languages") {
		supported = append(supported, language.MustParse(l))
	}

	return wikipedia.Languages(supported)
}
