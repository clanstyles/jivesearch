package wikipedia

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jivesearch/jivesearch/log"
	"github.com/lib/pq"
	"golang.org/x/text/language"
)

// PostgreSQL contains our client and database info
type PostgreSQL struct {
	*sql.DB
}

type tableType = string

const wikidataTable tableType = "wikidata"
const wikiTable tableType = "wiki"

type table struct {
	Type      tableType
	name      string
	temporary string
	columns   []column
	rows      chan interface{}
}

type column struct {
	name  string
	t     string
	index bool
}

// Scan unmarshals jsonb data
func (l *Labels) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), l)
}

// Scan unmarshals jsonb data
func (a *Aliases) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), a)
}

// Scan unmarshals jsonb data
func (d *Descriptions) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), d)
}

// Scan unmarshals jsonb data
// http://www.booneputney.com/development/gorm-golang-jsonb-value-copy/
func (c *Claims) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), c)
}

// Fetch retrieves an Item from PostgreSQL
// https://www.wikidata.org/w/api.php
func (p *PostgreSQL) Fetch(query string, lang language.Tag) (*Item, error) {
	item := &Item{
		Wikipedia: Wikipedia{
			Language: lang.String(),
		},
		Wikidata: &Wikidata{
			Claims: &Claims{},
		},
	}

	// iterate through Data
	tags, objects, stmts := []string{}, []string{}, []string{}

	data := reflect.Indirect(reflect.ValueOf(item.Claims))
	for i := 0; i < data.NumField(); i++ {
		tag := strings.Split(data.Type().Field(i).Tag.Get("json"), ",")[0]

		tags = append(tags, fmt.Sprintf(`"%v"`, tag))
		objects = append(objects, fmt.Sprintf(`jsonb_build_object('%v', "%v".item)`, tag, tag)) // `'influences', influences."item"`

		switch data.Field(i).Interface().(type) {
		case []string, []Text:
			switch tag {
			case "image", "flag":
				/* two urls to get images:
				1. https://commons.wikimedia.org/w/thumb.php?width=500&f=Junior-Jaguar-Belize-Zoo.jpg
					a. Will resize the image. Requesting a larger size than the original will result in an error
				2. https://upload.wikimedia.org/wikipedia/commons/2/21/Junior-Jaguar-Belize-Zoo.jpg
					a. 2 & 21 represent the first and first two characters of the image md5 hash
				*/
				stmts = append(stmts, fmt.Sprintf(`
					"%v" AS (
						SELECT build_image(claims->'%v') item
						FROM item
					)`, tag, tag),
				)
			default:
				stmts = append(stmts, fmt.Sprintf(`
					"%v" AS (
						SELECT claims->'%v' item
						FROM item
					)`, tag, tag),
				)
			}

		case []DateTime:
			stmts = append(stmts, fmt.Sprintf(`
				"%v" AS (
					SELECT build_datetime(claims->'%v') item
					FROM item
				)`, tag, tag),
			)
		case []Quantity:
			stmts = append(stmts, fmt.Sprintf(`
				"%v" AS (
					SELECT build_quantity(claims->'%v') item
					FROM item
				)`, tag, tag),
			)
		case []Wikidata:
			stmts = append(stmts, fmt.Sprintf(`
				"%v" AS (
					SELECT build_item(claims->'%v') item
					FROM item
				)`, tag, tag),
			)
		case []Coordinate:

		default: // e.g. has qualifiers
			var elem reflect.Value
			field := reflect.Indirect(reflect.ValueOf(item.Claims)).Field(i)

			typ := field.Type().Elem()
			if typ.Kind() == reflect.Ptr {
				elem = reflect.New(typ.Elem())
			}
			if typ.Kind() == reflect.Struct {
				elem = reflect.New(typ).Elem()
			}

			var inner []string

			for j := 0; j < reflect.Indirect(elem).NumField(); j++ {
				t := strings.Split(elem.Type().Field(j).Tag.Get("json"), ",")[0]

				switch elem.Field(j).Interface().(type) {
				case []string:
					inner = append(inner, fmt.Sprintf("'%v', x.d->'%v'", t, t))
				case []Wikidata:
					inner = append(inner, fmt.Sprintf("'%v', build_item(x.d->'%v')", t, t))
				case []Quantity:
					inner = append(inner, fmt.Sprintf("'%v', build_quantity(x.d->'%v')", t, t))
				case []DateTime:
					inner = append(inner, fmt.Sprintf("'%v', build_datetime(x.d->'%v')", t, t))
				default:
					log.Info.Printf(" unsupported field for %v\n", t)

				}
			}

			stmts = append(stmts, fmt.Sprintf(`
				"%v" AS (
					SELECT jsonb_agg(
						jsonb_build_object(
							%v
						)
					) item
					FROM item
					JOIN LATERAL (
						SELECT * FROM jsonb_array_elements(item.claims->'%v')
					) as x(d) on true
				)`, tag, strings.Join(inner, ", "), tag,
			))
		}
	}

	// Note: We cannot build 1 large jsonb_build_object as PostgreSQL has a 100 item limit.
	sql := fmt.Sprintf(`
		WITH item AS (
			SELECT w."id", w."title", w."text", wd."labels", wd."aliases", wd."descriptions", wd."claims" FROM 
			%vwiki w 
			LEFT JOIN wikidata wd
			ON w.id = wd.id
			WHERE LOWER(title) = LOWER($1)
			AND wd.id IS NOT NULL
		),
		%v
		
		SELECT
			item."id", item."title", item."text", item."labels", item."aliases", item."descriptions", %v "claims"
		FROM item, %v
	`, item.Language, strings.Join(stmts, ", "), strings.Join(objects, " || "), strings.Join(tags, ", "))

	err := p.DB.QueryRow(sql, query).Scan(&item.Wikidata.ID, &item.Title, &item.Text, &item.Labels, &item.Aliases, &item.Descriptions, &item.Claims)

	return item, err
}

type transaction = func(tx *sql.Tx) error

func (p *PostgreSQL) executeTransaction(t transaction) (err error) {
	tx, err := p.DB.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	err = t(tx)
	return
}

// Dump creates a temporary table and dumps rows via our transaction
func (p *PostgreSQL) Dump(wikidata bool, lang language.Tag, rows chan interface{}) error {
	t := &table{
		rows: rows,
	}

	switch wikidata {
	case true:
		t.Type = wikidataTable
		t.name = wikidataTable
	default:
		t.Type = wikiTable

		n := strings.Replace(lang.String(), "-", "_", -1)
		t.name = fmt.Sprintf("%v%v", strings.ToLower(n), wikiTable) // enwiki, cebwiki, etc...
	}

	t.temporary = t.name + "_tmp"

	if err := t.setColumns(); err != nil {
		return err
	}

	if _, err := p.DB.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %v`, t.temporary)); err != nil {
		return err
	}

	if _, err := p.DB.Exec(t.createTable()); err != nil {
		return err
	}

	if err := p.executeTransaction(t.insert); err != nil {
		return err
	}

	if err := p.executeTransaction(t.addIndices); err != nil {
		return err
	}

	return p.executeTransaction(t.rename)
}

func (t *table) setColumns() error {
	var err error

	switch t.Type {
	case wikiTable:
		t.columns = []column{
			column{"id", "text", true},
			column{"title", "text", true},
			column{"text", "text", false},
		}
	case wikidataTable:
		t.columns = []column{
			column{"id", "text", true},
			column{"labels", "jsonb", true},
			column{"aliases", "jsonb", true},
			column{"descriptions", "jsonb", false},
			column{"claims", "jsonb", true},
		}
	default:
		err = fmt.Errorf("unknown postgresql table type %v", t.Type)
	}

	return err
}

func (t *table) createTable() string {
	c := fmt.Sprintf(`CREATE TABLE %v (pk serial PRIMARY KEY,`, t.temporary)

	cols := []string{}
	for _, col := range t.columns {
		cols = append(cols, fmt.Sprintf("%v %v NOT NULL", col.name, col.t))
	}

	c += strings.Join(cols, ", ") + ")"
	return c
}

func (t *table) insert(tx *sql.Tx) (err error) {
	cols := []string{}
	for _, col := range t.columns {
		cols = append(cols, col.name)
	}

	stmt, err := tx.Prepare(pq.CopyIn(t.temporary, cols...))
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
	}()

	// dump the rows
	for row := range t.rows {
		r := []interface{}{}

		switch row.(type) {
		case *Wikipedia:
			w := row.(*Wikipedia)
			r = []interface{}{w.ID, w.Title, w.Text}
		case *Wikidata:
			w := row.(*Wikidata)

			r = []interface{}{w.ID}

			// The following are all jsonb columns.
			val := reflect.ValueOf(w).Elem()
			for i := 1; i < val.NumField(); i++ {
				j, err := json.Marshal(val.Field(i).Interface())
				if err != nil {
					return err
				}

				r = append(r, string(j))
			}
		default:
			return fmt.Errorf("unknown datatype for %+v", r)
		}

		if _, err = stmt.Exec(r...); err != nil {
			return err
		}
	}

	return nil
}

func (t *table) indexName(tbl, col string) string {
	return fmt.Sprintf("%v_%v", tbl, col)
}

// addIndices adds indexes to our temporary table
func (t *table) addIndices(tx *sql.Tx) (err error) {
	for _, c := range t.columns {
		if !c.index {
			continue
		}

		col := c.name
		if c.name == "title" {
			col = fmt.Sprintf("LOWER(%v)", col)
		}

		var using string
		if c.t == "jsonb" {
			using = "USING gin"
		}

		idx := fmt.Sprintf(`CREATE INDEX %v ON %v %v (%v)`, t.indexName(t.temporary, c.name), t.temporary, using, col)
		if _, err = tx.Exec(idx); err != nil {
			return err
		}
	}

	return err
}

// rename renames the t.temporary table to t.name
func (t *table) rename(tx *sql.Tx) (err error) {
	_, err = tx.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %v`, t.name))
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`ALTER TABLE %v RENAME to %v`, t.temporary, t.name))
	if err != nil {
		return err
	}

	for _, c := range t.columns {
		if !c.index {
			continue
		}

		_, err = tx.Exec(fmt.Sprintf(`ALTER INDEX %v RENAME to %v`,
			t.indexName(t.temporary, c.name), t.indexName(t.name, c.name)),
		)
		if err != nil {
			return err
		}
	}

	return err
}

// Setup creates our functions
func (p *PostgreSQL) Setup() error {
	buildItem := `
	CREATE OR REPLACE FUNCTION build_item(jsonb) 
	RETURNS jsonb immutable strict language sql as $$
	   SELECT jsonb_agg(                                               
			jsonb_build_object(                                     
				 'id', x.d->'id',                                
				 'labels', wikidata.labels                       
			)                                                       
		  )                                                               
		FROM jsonb_array_elements($1) AS x(d)                           
	   LEFT JOIN wikidata ON (x.d->>'id') = wikidata.id     
	$$;
	`

	buildDateTime := `
	CREATE OR REPLACE FUNCTION build_datetime(jsonb) 
	RETURNS jsonb immutable strict language sql as $$
		SELECT jsonb_agg(                                               
			jsonb_build_object(                                     
				'value', x.d->'value',                          
				'calendar', jsonb_build_object(                 
					'id', x.d->'calendar'->>'id',           
					'labels', wikidata.labels               
				)                                               
			)                                                       
		)                                                               
		FROM jsonb_array_elements($1) AS x(d)                           
		LEFT JOIN wikidata on (x.d->'calendar'->>'id') = wikidata.id 
	$$;
	`

	buildQuantity := `
	CREATE OR REPLACE FUNCTION build_quantity(jsonb) 
	RETURNS jsonb immutable strict language sql as $$
		SELECT jsonb_agg(                                               
			jsonb_build_object(                                     
				'amount', x.d->'amount',                        
				'unit', jsonb_build_object(                     
					'id', x.d->'unit'->>'id',               
					'labels', wikidata.labels               
				)                                               
			)                                                       
		)                                                               
		FROM jsonb_array_elements($1) AS x(d)                           
		LEFT JOIN wikidata on (x.d->'unit'->>'id') = wikidata.id       
	$$; 
	`

	/*
		NOTE: using 2 slashes as below e.g. 'https://upload…'
		will result in a 301 redirect to 'https:/upload….'
		then a 200 response….trying to cut out the redirect
		by using just 1 slash  will result in an invalid signature...
		we have to have the redirect for some reason.
	*/
	buildImage := `
	CREATE OR REPLACE FUNCTION build_image(jsonb) 
	RETURNS jsonb immutable strict language sql as $$
	SELECT jsonb_agg(  
		'https://upload.wikimedia.org/wikipedia/commons/' || LEFT(item.m, 1) || '/' || LEFT(item.m, 2) || '/' || s                                  
	)                                                               
	FROM (
		SELECT 
			md5(REPLACE(x.d::text, ' ', '_')) AS m,
			REPLACE(x.d::text, ' ', '_') AS s
		FROM jsonb_array_elements_text($1) AS x(d) 
	) item  
	$$;
	`

	for _, f := range []string{buildItem, buildDateTime, buildQuantity, buildImage} {
		if _, err := p.DB.Exec(f); err != nil {
			return err
		}
	}

	return nil
}
