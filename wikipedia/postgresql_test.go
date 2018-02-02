package wikipedia

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var shaqClaimsJSON = []byte(`{
	"sex": [{"id": "Q6581097", "labels": {"en": {"value": "male", "language": "en"}}}]
}`)

var shaqClaimsPostgres = &Claims{
	Sex: []Wikidata{
		Wikidata{
			ID: "Q6581097",
			Labels: map[string]Text{
				"en": Text{Text: "male", Language: "en"},
			},
			Claims: &Claims{},
		},
	},
}

func TestPostgreSQL_Fetch(t *testing.T) {
	type args struct {
		query string
		lang  language.Tag
	}
	tests := []struct {
		name string
		args args
		want *Item
	}{
		{
			"shaq",
			args{"Shaquille O'Neal", language.MustParse("en")},
			&Item{
				Wikipedia: Wikipedia{
					Language: "en",
					Title:    "Shaquille O'Neal",
					Text:     "Shaquille O'Neal is a basketball player",
				},
				Wikidata: &Wikidata{
					ID:           "Q169452",
					Descriptions: shaqDescriptions,
					Aliases:      shaqAliases,
					Labels:       shaqLabels,
					Claims:       shaqClaimsPostgres,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			rows := sqlmock.NewRows(
				[]string{`w."id"`, `w."title"`, `w."text"`, `wd."labels"`, `wd."aliases"`, `wd."descriptions"`, `wd."claims"`},
			)
			rows = rows.AddRow(
				"Q169452", tt.args.query, "Shaquille O'Neal is a basketball player",
				[]byte(shaqRawLabels), []byte(shaqRawAliases), []byte(shaqRawDescriptions), shaqClaimsJSON,
			)
			mock.ExpectQuery("SELECT").WithArgs(tt.args.query).WillReturnRows(rows)

			p := &PostgreSQL{
				DB: db,
			}

			got, err := p.Fetch(tt.args.query, tt.args.lang)
			if err != nil {
				t.Fatal(err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestPostgreSQL_Dump(t *testing.T) {
	type args struct {
		wikidata bool
		lang     language.Tag
		rows     chan interface{}
		done     chan bool
	}
	tests := []struct {
		name string
		row  interface{}
		args args
	}{
		{
			"enwiki",
			shaqWikipedia,
			args{
				wikidata: false,
				lang:     language.MustParse("en"),
				rows:     make(chan interface{}),
				done:     make(chan bool),
			},
		},
		{
			"wikidata",
			shaqWikidata,
			args{
				wikidata: true,
				rows:     make(chan interface{}),
				done:     make(chan bool),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			// create table
			mock.ExpectExec("DROP TABLE IF EXISTS").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

			// insert data
			mock.ExpectBegin()
			mock.ExpectPrepare("COPY")
			mock.ExpectCommit()

			// create indices
			mock.ExpectBegin()
			mock.ExpectExec("CREATE INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("CREATE INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.name == "wikidata" {
				mock.ExpectExec("CREATE INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			// rename table
			mock.ExpectBegin()
			mock.ExpectExec("DROP TABLE").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("ALTER TABLE").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("ALTER INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("ALTER INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.name == "wikidata" {
				mock.ExpectExec("ALTER INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER INDEX").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			p := &PostgreSQL{
				DB: db,
			}

			go func() {
				tt.args.rows <- tt.row
			}()

			if err := p.Dump(tt.args.wikidata, tt.args.lang, tt.args.rows); err != nil {
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
		DB: db,
	}

	mock.ExpectExec("CREATE OR REPLACE FUNCTION").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE OR REPLACE FUNCTION").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE OR REPLACE FUNCTION").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE OR REPLACE FUNCTION").WillReturnResult(sqlmock.NewResult(1, 1))

	err = p.Setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
