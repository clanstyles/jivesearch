package document

import (
	"context"
	"fmt"

	"github.com/jivesearch/jivesearch/log"
	"github.com/olivere/elastic"
	"golang.org/x/text/language"
)

// ElasticSearch hold connection and index settings
type ElasticSearch struct {
	Client *elastic.Client
	Index  string
	Type   string
}

var langAnalyzer = make(map[language.Tag]string)

// IndexName returns the language-specific index
// e.g. "search-english", "search-french"
func (e *ElasticSearch) IndexName(a string) string {
	return e.Index + "-" + a
}

// Analyzer returns the appropriate analyzer for a given language.
func (e *ElasticSearch) Analyzer(lang language.Tag) (string, error) {
	var analyzer string
	var ok bool

	// Get the parent tag until an analyzer is found
	tag, _, _ := Matcher.Match(lang)
	for ; tag != language.Und; tag = tag.Parent() {
		if analyzer, ok = langAnalyzer[tag]; ok {
			return analyzer, nil
		}
	}

	return "", fmt.Errorf("analyzer should not be blank. Lang: %s", lang)
}

// Setup will create our main search index
// and language-specific indices for the content
func (e *ElasticSearch) Setup() error {
	// We create one index per analyzer: search-english, search-spanish, etc...
	// This is a list of all elasticsearch analyzers
	analyzers := []string{"arabic", "armenian", "basque", "brazilian",
		"bulgarian", "catalan", "cjk", "czech", "danish", "dutch",
		"english", "finnish", "french", "galician", "german", "greek",
		"hindi", "hungarian", "indonesian", "irish", "italian", "latvian",
		"lithuanian", "norwegian", "persian", "portuguese", "romanian",
		"russian", "sorani", "spanish", "swedish", "turkish", "thai",
	}

	for _, a := range analyzers {
		idx := e.IndexName(a)
		exists, err := e.Client.IndexExists(idx).Do(context.TODO())
		if err != nil {
			return err
		}

		if !exists {
			log.Info.Println("Creating index:", idx)
			if _, err = e.Client.CreateIndex(idx).Body(e.mapping(a)).Do(context.TODO()); err != nil {
				return err
			}
		}
	}

	return nil
}

// mapping is the mapping of our main search Index.
// https://www.elastic.co/guide/en/elasticsearch/guide/current/one-lang-docs.html
func (e *ElasticSearch) mapping(a string) string {
	// Notes: Does the domain_name_analyzer and path_analyzer deal with rtl text (arabic, hebrew, etc)???
	// Also, needs better analyzer for domain...
	// e.g. search for "jimi hendrix" should return jimihendrix.com
	m := fmt.Sprintf(`{
		"settings": {
			"analysis": {
				"filter": {
					"my_shingle_filter": {
						"type":             "shingle",
						"min_shingle_size": 2, 
						"max_shingle_size": 2, 
						"output_unigrams":  false   
					}
				},
				"analyzer": {
					"my_shingle_analyzer": {
						"type":             "custom",
						"tokenizer":        "standard",
						"filter": [
							"lowercase",
							"my_shingle_filter" 
						]
					},
					"domain_name_analyzer": {
						"tokenizer": "domain_name_tokenizer"
					},
					"path_analyzer": {
						"tokenizer": "path_tokenizer"
					}
				},
				"tokenizer": {
					"domain_name_tokenizer": {
						"type": "path_hierarchy",
						"delimiter": "."
					},
					"path_tokenizer": {
						"type": "path_hierarchy",
						"delimiter": "/",
						"replacement": " "
					}
				}
			}
		},
		"mappings": {
			"document": {
				"_all": {
					"enabled": false
				},
				"dynamic": "strict",
				"properties": {
					"title": {
						"type": "text",
						"fields": {
							"lang": {
								"type":     "text",
								"analyzer": "%v" 
							},
							"shingles": {
								"type": 	"text",
								"analyzer": "my_shingle_analyzer"
							}
						}
					},
					"description": {
						"type": "text",
						"fields": {
							"lang": {
								"type":     "text",
								"analyzer": "%v" 
							},
							"shingles": {
								"type": 	"text",
								"analyzer": "my_shingle_analyzer"
							}
						}
					},
					"id": {
						"type": "keyword"
					},
					"keywords": {
						"type": "text"
					},
					"scheme": {
						"type": "keyword"
					},
					"domain": {
						"type": "text",
						"analyzer": "domain_name_analyzer"
					},
					"tld": {
						"type": "keyword"
					},						
					"host": {
						"type": "keyword"
					},						
					"path_parts": {
						"type":  "text",
						"analyzer": "path_analyzer"
					},
					"index": {
						"type": "boolean"
					},
					"crawled": {
						"type": "date",
						"format": "basic_date"
					},
					"date": {
						"type": "date",
						"format": "strict_date_optional_time"
					},
					"status": {
						"type": "short"
					},
					"canonical": {
						"type": "text",
						"index": "false"
					},
					"mime": {
						"type": "keyword"
					}
				}
			}
		}
	}`, a, a)

	return m
}

func init() {
	// These are the most commonly used languages mapped to an elasticsearch analyzer
	// TODO: fill in the rest of this map. Also, we haven't mapped the Basque, Galician,
	// Irish, and Sorani analyzers to any languages yet.

	//langAnalyzer[language.Afrikaans] = ""                   // af
	//langAnalyzer[language.Amharic] = ""                      // am
	langAnalyzer[language.Arabic] = "arabic" // ar
	//langAnalyzer[language.ModernStandardArabic] = ""         // ar-001
	//langAnalyzer[language.Azerbaijani] = ""                // az
	langAnalyzer[language.Bulgarian] = "bulgarian" // bg
	//langAnalyzer[language.Bengali] = ""                    // bn
	langAnalyzer[language.Catalan] = "catalan"              // ca
	langAnalyzer[language.Czech] = "czech"                  //  cs
	langAnalyzer[language.Danish] = "danish"                //  da
	langAnalyzer[language.German] = "german"                //  de
	langAnalyzer[language.Greek] = "greek"                  //  el
	langAnalyzer[language.English] = "english"              //  en
	langAnalyzer[language.AmericanEnglish] = "english"      //  en-US
	langAnalyzer[language.BritishEnglish] = "english"       //  en-GB
	langAnalyzer[language.Spanish] = "spanish"              //  es
	langAnalyzer[language.EuropeanSpanish] = "spanish"      //  es-ES
	langAnalyzer[language.LatinAmericanSpanish] = "spanish" //  es-419
	//langAnalyzer[language.Estonian] = ""                     //  et
	langAnalyzer[language.Persian] = "persian" //  fa
	langAnalyzer[language.Finnish] = "finnish" //  fi
	//langAnalyzer[language.Filipino] = ""                     //  fil
	langAnalyzer[language.French] = "french"         //  fr
	langAnalyzer[language.CanadianFrench] = "french" //  fr-CA
	//langAnalyzer[language.Gujarati] = ""                     //  gu
	//langAnalyzer[language.Hebrew] = ""                       //  he
	langAnalyzer[language.Hindi] = "hindi" //  hi
	//langAnalyzer[language.Croatian] = ""                 //  hr
	langAnalyzer[language.Hungarian] = "hungarian"   //  hu
	langAnalyzer[language.Armenian] = "armenian"     //  hy
	langAnalyzer[language.Indonesian] = "indonesian" //  id
	//langAnalyzer[language.Icelandic] = ""                   //  is
	langAnalyzer[language.Italian] = "italian" //  it
	langAnalyzer[language.Japanese] = "cjk"    //  ja
	//langAnalyzer[language.Georgian] = ""                    //  ka
	//langAnalyzer[language.Kazakh] = ""                    //  kk
	//langAnalyzer[language.Khmer] = ""                    //  km
	//langAnalyzer[language.Kannada] = ""                      //  kn
	langAnalyzer[language.Korean] = "cjk" //  ko
	//langAnalyzer[language.Kirghiz] = ""                      //  ky
	//langAnalyzer[language.Lao] = ""                          //  lo
	langAnalyzer[language.Lithuanian] = "lithuanian" //  lt
	langAnalyzer[language.Latvian] = "latvian"       //  lv
	//langAnalyzer[language.Macedonian] = ""                //  mk
	//langAnalyzer[language.Malayalam] = ""                    //  ml
	//langAnalyzer[language.Mongolian] = ""                 //  mn
	//langAnalyzer[language.Marathi] = ""                    //  mr
	//langAnalyzer[language.Malay] = ""                      //  ms
	//langAnalyzer[language.Burmese] = ""                      //  my
	//langAnalyzer[language.Nepali] = ""                      //  ne
	langAnalyzer[language.Dutch] = "dutch"         //  nl
	langAnalyzer[language.Norwegian] = "norwegian" //  no
	//langAnalyzer[language.Punjabi] = ""                //  pa
	//langAnalyzer[language.Polish] = ""                     //  pl
	langAnalyzer[language.Portuguese] = "portuguese"         //  pt
	langAnalyzer[language.BrazilianPortuguese] = "brazilian" //  pt-BR
	langAnalyzer[language.EuropeanPortuguese] = "portuguese" //  pt-PT
	langAnalyzer[language.Romanian] = "romanian"             //  ro
	langAnalyzer[language.Russian] = "russian"               //  ru
	//langAnalyzer[language.Sinhala] = ""                      //  si
	//langAnalyzer[language.Slovak] = ""                       //  sk
	//langAnalyzer[language.Slovenian] = ""                    //  sl
	//langAnalyzer[language.Albanian] = ""                  //  sq
	//langAnalyzer[language.Serbian] = ""                      //  sr
	//langAnalyzer[language.SerbianLatin] = ""                 //  sr-Latn
	langAnalyzer[language.Swedish] = "swedish" //  sv
	//langAnalyzer[language.Swahili] = ""                    //  sw
	//langAnalyzer[language.Tamil] = ""                        //  ta
	//langAnalyzer[language.Telugu] = ""                       //  te
	langAnalyzer[language.Thai] = "thai"       //  th
	langAnalyzer[language.Turkish] = "turkish" //  tr
	//langAnalyzer[language.Ukrainian] = ""                   //  uk
	//langAnalyzer[language.Urdu] = ""                        //  ur
	//langAnalyzer[language.Uzbek] = ""                      //  uz
	langAnalyzer[language.Vietnamese] = "cjk"         //  vi
	langAnalyzer[language.Chinese] = "cjk"            //  zh
	langAnalyzer[language.SimplifiedChinese] = "cjk"  //  zh-Hans
	langAnalyzer[language.TraditionalChinese] = "cjk" //  zh-Hant
	//langAnalyzer[language.Zulu] = ""                       //  zu
}
