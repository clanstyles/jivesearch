package instant

import (
	"fmt"
	"jivesearch/instant/contributors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Potus is an instant answer
type Potus struct {
	Answer
}

// POTUS represents a single President & their Vice President(s)
type POTUS struct {
	Ordinal   string
	President Person
	Party     string
	Terms     int
	Vice      []Person
}

// Person is a President or Vice President
type Person struct {
	Name  string
	Start string
	End   string
}

func (p *Potus) setQuery(r *http.Request) answerer {
	p.Answer.setQuery(r)
	return p
}

func (p *Potus) setUserAgent(r *http.Request) answerer {
	return p
}

func (p *Potus) setType() answerer {
	p.Type = "potus"
	return p
}

func (p *Potus) setContributors() answerer {
	p.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return p
}

func (p *Potus) setTriggers() answerer {
	p.triggers = []string{
		"president of the united states", "potus",
	}
	return p
}

func (p *Potus) setTriggerFuncs() answerer {
	p.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return p
}

func (p *Potus) setSolution() answerer {
	// maybe a better solution is to have
	// a set of non-trigger funcs???
	if strings.Contains(p.query, "vice") {
		p.Solution = Solution{}
		return p
	}

	// Find POTUS
	re := regexp.MustCompile("[0-9]+")
	num := re.FindAllString(p.remainder, -1)

	var i int
	if len(num) >= 1 { // we just take the first number in their query
		i, _ = strconv.Atoi(num[0])
	}

	// current POTUS
	if i == 0 {
		i = 45
	}

	// for now we just return the President's name
	// we have the data for dates, VP's, etc but
	// until we get further along with the HTML output
	// then we'll just do this.
	data, found := presidents(i)
	if !found {
		p.Solution = Solution{}
		return p
	}

	p.Text = data.President.Name

	return p
}

func (p *Potus) setCache() answerer {
	p.Cache = true
	return p
}

func (p *Potus) tests() []test {
	// there is an obvious flaw in the tests below:
	// e.g. "2st", etc. Also, we need to support the
	// numbers spelled out ("first", "second", etc.)
	typ := "potus"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: "current POTUS",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Donald Trump",
					Cache:        true,
				},
			},
		},
		test{
			query: "46th POTUS",
			expected: []Solution{
				Solution{},
			},
		},
		test{
			query: "32nd vice POTUS",
			expected: []Solution{
				Solution{},
			},
		},
	}

	for _, q := range []string{
		"%dst president of the united states",
		"who was the %dnd POTUS",
		"%d president of the united states",
	} {
		for i, pres := range []string{
			"George Washington",
			"John Adams",
			"Thomas Jefferson",
			"James Madison",
			"James Monroe",
			"John Quincy Adams",
			"Andrew Jackson",
			"Martin Van Buren",
			"William Henry Harrison",
			"John Tyler",
			"James K. Polk",
			"Zachary Taylor",
			"Millard Fillmore",
			"Franklin Pierce",
			"James Buchanan",
			"Abraham Lincoln",
			"Andrew Johnson",
			"Ulysses S. Grant",
			"Rutherford B. Hayes",
			"James A. Garfield",
			"Chester A. Arthur",
			"Grover Cleveland",
			"Benjamin Harrison",
			"Grover Cleveland",
			"William McKinley",
			"Theodore Roosevelt",
			"William Howard Taft",
			"Woodrow Wilson",
			"Warren G. Harding",
			"Calvin Coolidge",
			"Herbert Hoover",
			"Franklin D. Roosevelt",
			"Harry S. Truman",
			"Dwight D. Eisenhower",
			"John F. Kennedy",
			"Lyndon B. Johnson",
			"Richard Nixon",
			"Gerald Ford",
			"Jimmy Carter",
			"Ronald Reagan",
			"George H. W. Bush",
			"Bill Clinton",
			"George W. Bush",
			"Barack Obama",
			"Donald Trump",
		} {
			t := test{
				query: fmt.Sprintf(q, i+1),
				expected: []Solution{
					Solution{
						Type:         typ,
						Triggered:    true,
						Contributors: contrib,
						Text:         pres,
						Cache:        true,
					},
				},
			}

			tests = append(tests, t)
		}
	}

	return tests
}

func presidents(i int) (POTUS, bool) {
	p := POTUS{}
	found := true

	switch i {
	case 1:
		p = POTUS{
			Ordinal:   "1st",
			Party:     "Non-partisan",
			Terms:     2,
			President: Person{Name: "George Washington", Start: "4/30/1789", End: "3/4/1797"},
			Vice: []Person{
				Person{Name: "John Adams", Start: "4/30/1789", End: "3/4/1797"},
			},
		}
	case 2:
		p = POTUS{
			Ordinal:   "2nd",
			Party:     "Federalist",
			Terms:     1,
			President: Person{Name: "John Adams", Start: "3/4/1797", End: "3/4/1801"},
			Vice: []Person{
				Person{Name: "Thomas Jefferson", Start: "3/4/1797", End: "3/4/1801"},
			},
		}
	case 3:
		p = POTUS{
			Ordinal:   "3rd",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Thomas Jefferson", Start: "3/4/1801", End: "3/4/1809"},
			Vice: []Person{
				Person{Name: "Aaron Burr", Start: "3/4/1801", End: "3/4/1805"},
				Person{Name: "George Clinton", Start: "3/4/1805", End: "3/4/1809"},
			},
		}
	case 4:
		p = POTUS{
			Ordinal:   "4th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "James Madison", Start: "3/4/1809", End: "3/4/1817"},
			Vice: []Person{
				Person{Name: "George Clinton", Start: "3/4/1809", End: "4/20/1812"},
				Person{Name: "Elbridge Gerry", Start: "3/4/1813", End: "11/23/1814"},
			},
		}
	case 5:
		p = POTUS{
			Ordinal:   "5th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "James Monroe", Start: "3/4/1817", End: "3/4/1825"},
			Vice: []Person{
				Person{Name: "Daniel Tomkins", Start: "3/4/1817", End: "3/4/1825"},
			},
		}
	case 6:
		p = POTUS{
			Ordinal:   "6th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "John Quincy Adams", Start: "3/4/1825", End: "3/4/1829"},
			Vice: []Person{
				Person{Name: "John C. Calhoun", Start: "3/4/1825", End: "3/4/1829"},
			},
		}
	case 7:
		p = POTUS{
			Ordinal:   "7th",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Andrew Jackson", Start: "3/4/1829", End: "3/4/1837"},
			Vice: []Person{
				Person{Name: "John C. Calhoun", Start: "3/4/1829", End: "12/28/1832"},
				Person{Name: "Martin Van Buren", Start: "3/4/1833", End: "3/4/1837"},
			},
		}
	case 8:
		p = POTUS{
			Ordinal:   "8th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "Martin Van Buren", Start: "3/4/1837", End: "3/4/1841"},
			Vice: []Person{
				Person{Name: "Richard Mentor Johnson", Start: "3/4/1837", End: "3/4/1841"},
			},
		}
	case 9:
		p = POTUS{
			Ordinal:   "9th",
			Party:     "Whig",
			Terms:     1,
			President: Person{Name: "William Henry Harrison", Start: "3/4/1841", End: "4/4/1841"},
			Vice: []Person{
				Person{Name: "John Tyler", Start: "3/4/1841", End: "4/4/1841"},
			},
		}
	case 10:
		p = POTUS{
			Ordinal:   "10th",
			Party:     "Whig / Independent",
			Terms:     1,
			President: Person{Name: "John Tyler", Start: "4/4/1841", End: "3/4/1845"},
		}
	case 11:
		p = POTUS{
			Ordinal:   "11th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "James K. Polk", Start: "3/4/1845", End: "3/4/1849"},
			Vice: []Person{
				Person{Name: "George M. Dallas", Start: "3/4/1845", End: "3/4/1849"},
			},
		}
	case 12:
		p = POTUS{
			Ordinal:   "12th",
			Party:     "Whig",
			Terms:     1,
			President: Person{Name: "Zachary Taylor", Start: "3/4/1849", End: "7/9/1850"},
			Vice: []Person{
				Person{Name: "Millard Fillmore", Start: "3/4/1849", End: "7/9/1850"},
			},
		}
	case 13:
		p = POTUS{
			Ordinal:   "13th",
			Party:     "Whig",
			Terms:     1,
			President: Person{Name: "Millard Fillmore", Start: "7/9/1850", End: "3/4/1853"},
		}
	case 14:
		p = POTUS{
			Ordinal:   "14th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "Franklin Pierce", Start: "3/4/1853", End: "3/4/1857"},
			Vice: []Person{
				Person{Name: "William R. King", Start: "3/4/1853", End: "4/18/1853"},
			},
		}
	case 15:
		p = POTUS{
			Ordinal:   "15th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "James Buchanan", Start: "3/4/1857", End: "3/4/1861"},
			Vice: []Person{
				Person{Name: "John C. Breckinridge", Start: "3/4/1857", End: "3/4/1861"},
			},
		}
	case 16:
		p = POTUS{
			Ordinal:   "16th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Abraham Lincoln", Start: "3/4/1861", End: "4/15/1865"},
			Vice: []Person{
				Person{Name: "Hannibal Hamlin", Start: "3/4/1861", End: "3/4/1865"},
				Person{Name: "Andrew Johnson", Start: "3/4/1865", End: "4/15/1865"},
			},
		}
	case 17:
		p = POTUS{
			Ordinal:   "17th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "Andrew Johnson", Start: "4/15/1865", End: "3/4/1869"},
		}
	case 18:
		p = POTUS{
			Ordinal:   "18th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Ulysses S. Grant", Start: "3/4/1869", End: "3/4/1877"},
			Vice: []Person{
				Person{Name: "Schuyler Colfax", Start: "3/4/1869", End: "3/4/1873"},
				Person{Name: "Henry Wilson", Start: "3/4/1873", End: "11/22/1875"},
			},
		}
	case 19:
		p = POTUS{
			Ordinal:   "19th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Rutherford B. Hayes", Start: "3/4/1877", End: "3/4/1881"},
			Vice: []Person{
				Person{Name: "William A. Wheeler", Start: "3/4/1877", End: "3/4/1881"},
			},
		}
	case 20:
		p = POTUS{
			Ordinal:   "20th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "James A. Garfield", Start: "3/4/1881", End: "9/19/1881"},
			Vice: []Person{
				Person{Name: "Chester A. Arthur", Start: "3/4/1881", End: "9/19/1881"},
			},
		}
	case 21:
		p = POTUS{
			Ordinal:   "21st",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Chester A. Arthur", Start: "9/19/1881", End: "3/4/1885"},
		}
	case 22:
		p = POTUS{
			Ordinal:   "22nd",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Grover Cleveland", Start: "3/4/1885", End: "3/4/1889"},
			Vice: []Person{
				Person{Name: "Thomas A. Hendricks", Start: "3/4/1885", End: "11/25/1885"},
			},
		}
	case 23:
		p = POTUS{
			Ordinal:   "23rd",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Benjamin Harrison", Start: "3/4/1889", End: "3/4/1893"},
			Vice: []Person{
				Person{Name: "Levi P. Morton", Start: "3/4/1889", End: "3/4/1893"},
			},
		}
	case 24:
		p = POTUS{
			Ordinal:   "24th",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Grover Cleveland", Start: "3/4/1893", End: "3/4/1897"},
			Vice: []Person{
				Person{Name: "Adlai Stevenson", Start: "3/4/1893", End: "3/4/1897"},
			},
		}
	case 25:
		p = POTUS{
			Ordinal:   "25th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "William McKinley", Start: "3/4/1897", End: "9/14/1901"},
			Vice: []Person{
				Person{Name: "Garret Hobart", Start: "3/4/1897", End: "11/21/1899"},
				Person{Name: "Theodore Roosevelt", Start: "3/4/1901", End: "9/14/1901"},
			},
		}
	case 26:
		p = POTUS{
			Ordinal:   "26th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Theodore Roosevelt", Start: "9/14/1901", End: "3/4/1909"},
			Vice: []Person{
				Person{Name: "Charles W. Fairbanks", Start: "3/4/1905", End: "3/4/1909"},
			},
		}
	case 27:
		p = POTUS{
			Ordinal:   "27th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "William Howard Taft", Start: "3/4/1909", End: "3/4/1913"},
			Vice: []Person{
				Person{Name: "James S. Sherman", Start: "3/4/1909", End: "10/30/1912"},
			},
		}
	case 28:
		p = POTUS{
			Ordinal:   "28th",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Woodrow Wilson", Start: "3/4/1913", End: "3/4/1921"},
			Vice: []Person{
				Person{Name: "Thomas R. Marshall", Start: "3/4/1913", End: "3/4/1921"},
			},
		}
	case 29:
		p = POTUS{
			Ordinal:   "29th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Warren G. Harding", Start: "3/4/1921", End: "8/2/1923"},
			Vice: []Person{
				Person{Name: "Calvin Coolidge", Start: "3/4/1921", End: "8/2/1923"},
			},
		}
	case 30:
		p = POTUS{
			Ordinal:   "30th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Calvin Coolidge", Start: "8/2/1923", End: "3/4/1929"},
			Vice: []Person{
				Person{Name: "Charles G. Dawes", Start: "3/4/1925", End: "3/4/1929"},
			},
		}
	case 31:
		p = POTUS{
			Ordinal:   "31st",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Herbert Hoover", Start: "3/4/1929", End: "3/4/1933"},
			Vice: []Person{
				Person{Name: "Charles Curtis", Start: "3/4/1929", End: "3/4/1933"},
			},
		}
	case 32:
		p = POTUS{
			Ordinal:   "32nd",
			Party:     "Democrat",
			Terms:     4,
			President: Person{Name: "Franklin D. Roosevelt", Start: "3/4/1933", End: "4/12/1945"},
			Vice: []Person{
				Person{Name: "John Nance Garner", Start: "3/4/1933", End: "1/20/1941"},
				Person{Name: "Henry A. Wallace", Start: "1/20/1941", End: "1/20/1945"},
				Person{Name: "Harry S. Truman", Start: "1/20/1945", End: "4/12/1945"},
			},
		}
	case 33:
		p = POTUS{
			Ordinal:   "33rd",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Harry S. Truman", Start: "4/12/1945", End: "1/20/1953"},
			Vice: []Person{
				Person{Name: "Alben W. Barkley", Start: "1/20/1949", End: "1/20/1953"},
			},
		}
	case 34:
		p = POTUS{
			Ordinal:   "34th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Dwight D. Eisenhower", Start: "1/20/1953", End: "1/20/1961"},
			Vice: []Person{
				Person{Name: "Richard Nixon", Start: "1/20/1953", End: "1/20/1961"},
			},
		}
	case 35:
		p = POTUS{
			Ordinal:   "35th",
			Party:     "Democrat",
			Terms:     1,
			President: Person{Name: "John F. Kennedy", Start: "1/20/1961", End: "11/22/1963"},
			Vice: []Person{
				Person{Name: "Lyndon B. Johnson", Start: "1/20/1961", End: "11/22/1963"},
			},
		}
	case 36:
		p = POTUS{
			Ordinal:   "36th",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Lyndon B. Johnson", Start: "11/22/1963", End: "1/20/1969"},
			Vice: []Person{
				Person{Name: "Hubert Humphrey", Start: "1/20/1965", End: "1/20/1969"},
			},
		}
	case 37:
		p = POTUS{
			Ordinal:   "37th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Richard Nixon", Start: "1/20/1969", End: "8/9/1974"},
			Vice: []Person{
				Person{Name: "Spiro Agnew", Start: "1/20/1969", End: "10/10/1973"},
				Person{Name: "Gerald Ford", Start: "12/6/1973", End: "8/9/1974"},
			},
		}
	case 38:
		p = POTUS{
			Ordinal:   "38th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Gerald Ford", Start: "8/9/1974", End: "1/20/1977"},
			Vice: []Person{
				Person{Name: "Nelson Rockefeller", Start: "12/19/1974", End: "1/20/1977"},
			},
		}
	case 39:
		p = POTUS{
			Ordinal: "39th",
			Party:   "Democrat",
			Terms:   1,
			President: Person{
				Name: "Jimmy Carter", Start: "1/20/1977", End: "1/20/1981",
			},
			Vice: []Person{Person{Name: "Walter Mondale", Start: "1/20/1977", End: "1/20/1981"}},
		}
	case 40:
		p = POTUS{
			Ordinal:   "40th",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "Ronald Reagan", Start: "1/20/1981", End: "1/20/1989"},
			Vice: []Person{
				Person{Name: "George H. W. Bush", Start: "1/20/1981", End: "1/20/1989"},
			},
		}
	case 41:
		p = POTUS{
			Ordinal:   "41st",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "George H. W. Bush", Start: "1/20/1989", End: "1/20/1993"},
			Vice: []Person{
				Person{Name: "Dan Quayle", Start: "1/20/1989", End: "1/20/1993"},
			},
		}
	case 42:
		p = POTUS{
			Ordinal:   "42nd",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Bill Clinton", Start: "1/20/1993", End: "1/20/2001"},
			Vice: []Person{
				Person{Name: "Al Gore", Start: "1/20/1993", End: "1/20/2001"},
			},
		}
	case 43:
		p = POTUS{
			Ordinal:   "43rd",
			Party:     "Republican",
			Terms:     2,
			President: Person{Name: "George W. Bush", Start: "1/20/2001", End: "1/20/2009"},
			Vice: []Person{
				Person{Name: "Dick Cheney", Start: "1/20/2001", End: "1/20/2009"},
			},
		}
	case 44:
		p = POTUS{
			Ordinal:   "44th",
			Party:     "Democrat",
			Terms:     2,
			President: Person{Name: "Barack Obama", Start: "1/20/2009", End: "1/20/2017"},
			Vice: []Person{
				Person{Name: "Joe Biden", Start: "1/20/2009", End: "1/20/2017"},
			},
		}
	case 45:
		p = POTUS{
			Ordinal:   "45th",
			Party:     "Republican",
			Terms:     1,
			President: Person{Name: "Donald Trump", Start: "1/20/2017", End: ""},
			Vice: []Person{
				Person{Name: "Mike Pence", Start: "1/20/2017", End: ""},
			},
		}
	default:
		found = false
	}
	return p, found
}
