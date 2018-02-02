// Dumper downloads and dumps wikipedia/wikidata/wikiquotes data to a postgresql database.
package main

import (
	"database/sql"
	"fmt"
	"jivesearch/config"
	"jivesearch/log"
	"jivesearch/wikipedia"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func setup(v *viper.Viper) (*wikipedia.PostgreSQL, error) {
	v.SetEnvPrefix("jivesearch")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetDefaults(v)

	var err error

	if v.GetBool("debug") {
		log.Debug.SetOutput(os.Stdout)
	}

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

	if err != nil {
		panic(err)
	}

	p.DB.SetMaxIdleConns(0)

	return p, err
}

func files(v *viper.Viper, supported []language.Tag) ([]*wikipedia.File, error) {
	files := []*wikipedia.File{}

	if v.GetBool("wikipedia.data") {
		files = append(files, wikipedia.NewFile(wikipedia.WikiDataURL, language.English))
	}

	if v.GetBool("wikipedia.text") {
		f, err := wikipedia.CirrusLinks(supported)
		if err != nil {
			return nil, err
		}

		files = append(files, f...)
	}

	return files, nil
}

func languages(cfg config.Provider) ([]language.Tag, []language.Tag) {
	sup := cfg.GetStringSlice("languages")
	supported := []language.Tag{}

	for _, l := range sup {
		supported = append(supported, language.MustParse(l))
	}

	return wikipedia.Languages(supported)

}

func main() {
	v := viper.New()

	p, err := setup(v)
	if err != nil {
		panic(err)
	}

	defer p.DB.Close()

	supported, unsupported := languages(v)
	for _, lang := range unsupported {
		log.Info.Printf("wikipedia does not support langugage %q\n", lang)
	}

	files, err := files(v, supported)
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		log.Info.Fatalln("what files do you want to parse?")
	}

	download := make(chan *wikipedia.File)
	parse := make(chan *wikipedia.File)

	// Wikipedia seems to allow 3 concurrent downloads.
	var dwg sync.WaitGroup
	for i := 0; i < 3; i++ {
		dwg.Add(1)
		go func() {
			defer dwg.Done()
			for f := range download {
				if _, err := os.Stat(f.ABS); os.IsNotExist(err) {
					log.Info.Printf("downloading %v\n", f.URL.String())
					if err := f.Download(); err != nil {
						panic(err)
					}
				}
				parse <- f
			}
		}()
	}

	// parsing
	var wg sync.WaitGroup
	workers := v.GetInt("wikipedia.workers")

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range parse {
				if err := f.Parse(v.GetInt("wikipedia.truncate")); err != nil {
					panic(err)
				}
			}
		}()
	}

	dir := v.GetString("wikipedia.dir")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	for _, f := range files {
		f.SetDumper(p).SetABS(dir)
		download <- f
	}

	close(download)
	dwg.Wait()

	close(parse)
	wg.Wait()
}
