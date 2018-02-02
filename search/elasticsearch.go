package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jivesearch/jivesearch/search/document"
	"github.com/jivesearch/jivesearch/search/vote"
	"github.com/olivere/elastic"
	"golang.org/x/text/language"
)

// ElasticSearch embeds our main Elasticsearch instance
type ElasticSearch struct {
	*document.ElasticSearch
}

// Fetch returns search results for a search query
// https://www.elastic.co/guide/en/elasticsearch/guide/current/one-lang-docs.html
// https://www.elastic.co/guide/en/elasticsearch/guide/current/_single_query_string.html#know-your-data
// The idea here is to first filter out docs that do not want to be indexed.
// We then search multiple fields for the search query, giving more weight to certain fields.
// We also are searching the standard analyzer and the language-specific analyzer.
// We weight the domain > path, path > title, title > description.
// We also give extra weight for bigram matches (need trigram????):
// https://www.elastic.co/guide/en/elasticsearch/guide/current/shingles.html
// Note: "It is not useful to mix not_analyzed fields with analyzed fields in multi_match queries."
// TODO: A better domain name method...we could use regex ('.*hendrix'), prefix query, etc.
func (e *ElasticSearch) Fetch(q string, lang language.Tag, region language.Region, number int, offset int, votes []vote.Result) (*Results, error) {
	res := &Results{}

	qu := elastic.NewBoolQuery().
		Filter(elastic.NewTermQuery("index", true)).
		Must(
			elastic.NewMultiMatchQuery(
				q,
				"domain^3", "path^2",
				"title^1.5", "title.lang^1.5",
				"description", "description.lang",
			).Type("cross_fields").MinimumShouldMatch("-25%"),
		).
		Should(
			elastic.NewMultiMatchQuery(
				q,
				"title.shingles",
				"description.shingles",
			).Type("cross_fields"),
		)

	// Boost results for regional queries (except for .me, .tv, etc. that are used for other purposes sometimes)
	// https://support.google.com/webmasters/answer/182192#1
	if t, err := region.TLD(); err == nil {
		tld := strings.ToLower(t.String())
		if tld != "us" && tld != "tv" && tld != "me" && tld != "co" && tld != "io" {
			qu = qu.Should(elastic.NewMatchQuery("tld", tld))
		}
	}

	a, err := e.Analyzer(lang)
	if err != nil {
		return res, err
	}

	idx := e.IndexName(a)

	o := e.Client.Search().Index(idx).Type(e.Type).Query(qu).From(offset).Size(number)

	// sort by votes
	if len(votes) > 0 {
		var script string

		for i, v := range votes {
			e := "if"
			if i > 0 {
				e = "else if"
			}

			script += fmt.Sprintf(`%v (doc['id'].value.equals('%v')) return %d;`, e, v.URL, v.Votes)
		}

		script += fmt.Sprintf("else return 0;")

		sort := elastic.NewScriptSort(
			elastic.NewScript(script).Lang("painless"), "number",
		).Order(false).Type("number")
		o = o.SortBy(sort)
	}

	out, err := o.Do(context.TODO())
	if err != nil {
		return res, err
	}

	res.Count = out.TotalHits()

	for _, u := range out.Hits.Hits {
		doc := &document.Document{}
		err := json.Unmarshal(*u.Source, doc)
		if err != nil {
			return res, err
		}

		// Rather than have the highlighting done here in elasticsearch
		// we should have a method on doc to highlight so we get consistent
		// highlighting regardless of the backend used???
		// For now, we have moved highlighting to a javascript function
		// if des, ok := u.Highlight["description"]; ok {
		//	for _, v := range des {
		//		doc.Description = v
		//	}
		//}

		doc.ID = u.Id
		res.Documents = append(res.Documents, doc)
	}

	return res, err
}
