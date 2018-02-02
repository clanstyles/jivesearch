Jive Search is a completely open source search engine that respects your privacy. 

[Documentation](https://godoc.org/github.com/jivesearch/jivesearch)

go get -u github.com/jivesearch/jivesearch

###Crawler
cd $GOPATH/src/github.com/jivesearch/jivesearch/search/crawler && go run ./cmd/crawler.go --workers=75 --time=5m --debug=true

###Frontend
cd $GOPATH/src/github.com/jivesearch/jivesearch/frontend && go run ./cmd/frontend.go --debug=true

###Wikipedia Dump File
cd $GOPATH/src/github.com/jivesearch/jivesearch/wikipedia/cmd/dumper && go run dumper.go --workers=3 --dir=/path/to/wiki/files --text=true --data=true --truncate=400