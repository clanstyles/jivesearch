package vote

import (
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGet(t *testing.T) {
	type u struct {
		url   string
		votes int
	}

	for _, c := range []struct {
		name string
		q    string
		l    int
		o    int
		urls []u
		want []Result
	}{
		{
			name: "basic",
			q:    "a search term",
			l:    1,
			o:    1,
			urls: []u{
				{"https://www.example.com/a-path-to-nowhere", 2},
			},
			want: []Result{
				Result{
					URL:   "https://www.example.com/a-path-to-nowhere",
					Votes: 2,
				},
			},
		},
		{
			name: "multiple urls",
			q:    "another search term",
			l:    1,
			o:    1,
			urls: []u{
				{"http://example.com/?a=query", 150},
				{"https://www.example.com/a-path-to-somewhere", -419},
			},
			want: []Result{
				Result{
					URL:   "http://example.com/?a=query",
					Votes: 150,
				},
				Result{
					URL:   "https://www.example.com/a-path-to-somewhere",
					Votes: -419,
				},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			urls := []string{}
			rows := sqlmock.NewRows([]string{"url", "vote"})
			for _, u := range c.urls {
				urls = append(urls, u.url)
				rows = rows.AddRow(u.url, u.votes)
			}

			mock.ExpectQuery("SELECT").
				WithArgs(c.q, c.l).
				WillReturnRows(rows)

			p := &PostgreSQL{
				DB:    db,
				Table: "votes",
			}

			got, err := p.Get(c.q, c.l)
			if err != nil {
				t.Fatal(err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want: %+v", got, c.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	for _, c := range []struct {
		name string
		q    string
		url  string
		vote int
	}{
		{
			name: "upvote",
			q:    "a search term",
			url:  "https://www.example.com/a-path-to-nowhere",
			vote: 1,
		},
		{
			name: "downvote",
			q:    "cats",
			url:  "https://www.cat.com/?something=nothing",
			vote: -1,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			v := &Vote{
				Query: c.q,
				URL:   c.url,
				Vote:  c.vote,
				Date:  now(),
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			mock.ExpectExec("INSERT INTO votes").
				WithArgs(v.Query, v.URL, v.Domain, v.Vote, v.Date).
				WillReturnResult(sqlmock.NewResult(1, 1))

			p := &PostgreSQL{
				DB:    db,
				Table: "votes",
			}

			if err := p.Insert(v); err != nil {
				t.Fatal(err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSetup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	p := &PostgreSQL{
		DB:    db,
		Table: "votes",
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE INDEX IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("CREATE FUNCTION").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = p.Setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
