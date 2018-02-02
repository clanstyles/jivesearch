package vote

import (
	"database/sql"
	"errors"
	"fmt"
)

// PostgreSQL contains our client and database info
type PostgreSQL struct {
	*sql.DB
	Table string
}

// ErrScoreFnExists indicates an issue setting up our default score() function
var ErrScoreFnExists = errors.New(`pq: function "score" already exists with same argument types`)

// Get retrieves the urls & vote tallies for a given query
// score($1,$2,$3) is a PostgreSQL stored procedure to return
// votes for the matching query.
func (p *PostgreSQL) Get(query string, limit int) ([]Result, error) {
	votes := []Result{}

	// Pagination here won't work...???
	rows, err := p.DB.Query("SELECT * FROM score($1) LIMIT $2", query, limit)
	if err != nil {
		return votes, err
	}

	defer rows.Close()
	for rows.Next() {
		res := Result{}
		if err := rows.Scan(&res.URL, &res.Votes); err != nil {
			return votes, err
		}
		votes = append(votes, res)
	}

	if err := rows.Err(); err != nil {
		return votes, err
	}

	return votes, err
}

// Insert saves a vote to PostgreSQL
// using %s here for table name s/b safe from sql injection???
func (p *PostgreSQL) Insert(v *Vote) error {
	_, err := p.DB.Exec(fmt.Sprintf(
		`INSERT INTO %s (query, url, domain, vote, date)
		VALUES ($1, $2, $3, $4, $5)`, p.Table), v.Query, v.URL, v.Domain, v.Vote, v.Date,
	)

	return err
}

// Setup creates our table
func (p *PostgreSQL) Setup() error {
	_, err := p.DB.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			query text,
			url text,
			domain text,
			vote smallint,
			date date
		);`, p.Table),
	)

	if err != nil {
		return err
	}

	if _, err = p.DB.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS query_idx ON %s (query);", p.Table)); err != nil {
		return err
	}

	// Create a default score() function.
	// Do not replace an existing fn if it already exists.
	// This function is meant to be overidden.
	scoreFn := `
		CREATE FUNCTION score(q text)
		RETURNS TABLE(url text, votes bigint)
		AS 
		$$
			SELECT  url, sum(vote) AS votes
			FROM votes
			WHERE lower(query)=lower(q) 
			GROUP BY url
			ORDER BY votes DESC  
		$$ 
		LANGUAGE sql;
	`

	_, err = p.DB.Exec(scoreFn)

	return err
}
