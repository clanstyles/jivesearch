package wikipedia

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Wikidata is a Wikidata item
type Wikidata struct {
	ID           string `json:"id,omitempty"`
	Labels       `json:"labels,omitempty"`
	Aliases      `json:"aliases,omitempty"`
	Descriptions `json:"descriptions,omitempty"`
	*Claims
}

// Labels holds the labels for an Item
type Labels map[string]Text

// Aliases holds the alternative names for an Item
type Aliases map[string][]Text

// Descriptions holds the descriptions for an Item
type Descriptions map[string]Text

type claim struct {
	MainSnak        property              `json:"mainsnak"`
	Qualifiers      map[string][]property `json:"qualifiers"`
	QualifiersOrder []string              `json:"qualifiers-order"`
}

type property struct {
	Property  string `json:"property"`
	DataValue struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	} `json:"datavalue"`
	DataType string `json:"datatype"`
}

// DateTime is the raw, unformatted version of a datetime
// Note: Wikidata only uses Gregorian and Julian calendars
type DateTime struct {
	Value    string   `json:"value,omitempty"`
	Calendar Wikidata `json:"calendar,omitempty"`
}

// Text is a language and value
type Text struct {
	Text     string `json:"value,omitempty"`
	Language string `json:"language,omitempty"`
}

// Quantity is a Wikipedia quantity
type Quantity struct {
	Amount string   `json:"amount,omitempty"`
	Unit   Wikidata `json:"unit,omitempty"`
}

// Coordinate is a Wikipedia coordinate
type Coordinate struct {
	Latitude  []float64  `json:"latitude,omitempty"`
	Longitude []float64  `json:"longitude,omitempty"`
	Altitude  []float64  `json:"altitude,omitempty"`
	Precision []float64  `json:"precision,omitempty"`
	Globe     []Wikidata `json:"globe,omitempty"`
}

// Spouse represents a person's spouse or partner
type Spouse struct {
	Item  []Wikidata `json:"item,omitempty"`
	Start []DateTime `json:"start,omitempty" property:"P580"`
	End   []DateTime `json:"end,omitempty" property:"P582"`       // do we also need P585 as we do for Partner?
	Place []Wikidata `json:"location,omitempty" property:"P2842"` // AKA Location P276
}

// Team represents a team on which a person played
type Team struct {
	Item     []Wikidata `json:"item,omitempty"`
	Start    []DateTime `json:"start,omitempty" property:"P580"`
	End      []DateTime `json:"end,omitempty" property:"P582"`
	Position []Wikidata `json:"position,omitempty" property:"P413"`
	Number   []string   `json:"number,omitempty" property:"P1618"`
}

// Education represents the education of a person
type Education struct {
	Item   []Wikidata `json:"item,omitempty"`
	Start  []DateTime `json:"start,omitempty" property:"P580"`
	End    []DateTime `json:"end,omitempty" property:"P582"`
	Degree []Wikidata `json:"degree,omitempty" property:"P512"`
	Major  []Wikidata `json:"major,omitempty" property:"P812"`
}

// Interment is the place a person was buried
type Interment struct {
	Item  []Wikidata `json:"item,omitempty"`
	Start []DateTime `json:"start,omitempty" property:"P580"`
	End   []DateTime `json:"end,omitempty" property:"P582"`
}

// Award is an award someone won
type Award struct {
	Item []Wikidata `json:"item,omitempty"`
	Date []DateTime `json:"date,omitempty" property:"P585"`
}

// Military is a person's history in the military
type Military struct {
	Item  []Wikidata `json:"item,omitempty"`
	Start []DateTime `json:"start,omitempty" property:"P580"`
	End   []DateTime `json:"end,omitempty" property:"P582"`
}

// Member is a part of a group (band, etc)
type Member struct {
	Item  []Wikidata `json:"item,omitempty"`
	Start []DateTime `json:"start,omitempty" property:"P580"`
	End   []DateTime `json:"end,omitempty" property:"P582"`
	Date  []DateTime `json:"date,omitempty" property:"P585"` // some don't have start/end time just a point-in-time.
}

// Population is a point-in-time value of a country's population
type Population struct {
	Value []Quantity `json:"value,omitempty"`
	Date  []DateTime `json:"date,omitempty" property:"P585"`
}

// Instrument is a musical instrument (guitar, drums, etc)
type Instrument struct {
	Item         []Wikidata `json:"item,omitempty"`
	Manufacturer []Wikidata `json:"manufacturer,omitempty" property:"P176"`
}

// Nomination is a nomination for an award
type Nomination struct {
	Item []Wikidata `json:"item,omitempty"`
	For  []Wikidata `json:"for,omitempty" property:"P1686"`
	Date []DateTime `json:"date,omitempty" property:"P585"`
}

// claims is the raw version of the wikidata claims
type claims struct {
	Image       []claim `json:"P18"`
	BirthPlace  []claim `json:"P19"`
	Sex         []claim `json:"P21"`
	Father      []claim `json:"P22"`
	Mother      []claim `json:"P25"`
	Spouse      []claim `json:"P26"`
	Country     []claim `json:"P27"`
	Instance    []claim `json:"P31"` // eg person, book, etc...
	Capital     []claim `json:"P36"`
	Currency    []claim `json:"P38"`
	Flag        []claim `json:"P41"`
	Teams       []claim `json:"P54"`
	Education   []claim `json:"P69"`
	Occupation  []claim `json:"P106"`
	Signature   []claim `json:"P109"`
	Interment   []claim `json:"P119"`
	Genre       []claim `json:"P136"`
	Religion    []claim `json:"P140"`
	Awards      []claim `json:"P166"`
	Ethnicity   []claim `json:"P172"`
	Military    []claim `json:"P241"`
	RecordLabel []claim `json:"P264"`
	Discography []claim `json:"P358"`
	Position    []claim `json:"P413"`
	Partner     []claim `json:"P451"` // non-married spouse
	Origin      []claim `json:"P495"`
	DeathCause  []claim `json:"P509"`
	Members     []claim `json:"P527"`
	Residence   []claim `json:"P551"`
	Hand        []claim `json:"P552"`
	//Coordinate  []claim `json:"P625"`
	Birthday    []claim `json:"P569"`
	Death       []claim `json:"P570"`
	Start       []claim `json:"P571"`
	Sport       []claim `json:"P641"`
	Drafted     []claim `json:"P647"`
	GivenName   []claim `json:"P735"`
	Influences  []claim `json:"P737"`
	Location    []claim `json:"P740"` // location where a group was formed
	Website     []claim `json:"P856"`
	Population  []claim `json:"P1082"`
	Instrument  []claim `json:"P1303"`
	Participant []claim `json:"P1344"`
	Nominations []claim `json:"P1411"`
	Languages   []claim `json:"P1412"`
	BirthName   []claim `json:"P1477"`
	Spotify     []claim `json:"P1902"`
	Twitter     []claim `json:"P2002"`
	Instagram   []claim `json:"P2003"`
	Facebook    []claim `json:"P2013"`
	YouTube     []claim `json:"P2397"`
	WorkStart   []claim `json:"P2031"`
	Height      []claim `json:"P2048"`
	Weight      []claim `json:"P2067"`
	Siblings    []claim `json:"P3373"`
}

// Claims are the formatted and condensed version of the Wikidata claims
type Claims struct {
	Image       []string    `json:"image,omitempty"`
	BirthPlace  []Wikidata  `json:"birthplace,omitempty"`
	Sex         []Wikidata  `json:"sex,omitempty"`
	Father      []Wikidata  `json:"father,omitempty"`
	Mother      []Wikidata  `json:"mother,omitempty"`
	Spouse      []Spouse    `json:"spouse,omitempty"`
	Country     []Wikidata  `json:"country,omitempty"` // country of residence
	Instance    []Wikidata  `json:"instance,omitempty"`
	Capital     []Wikidata  `json:"capital,omitempty"`
	Currency    []Wikidata  `json:"currency,omitempty"`
	Flag        []string    `json:"flag,omitempty"`
	Teams       []Team      `json:"teams,omitempty"` // sports teams
	Education   []Education `json:"education,omitempty"`
	Occupation  []Wikidata  `json:"occupation,omitempty"`
	Signature   []string    `json:"signature,omitempty"`
	Interment   []Interment `json:"interment,omitempty"` // burial/ashes location
	Genre       []Wikidata  `json:"genre,omitempty"`
	Religion    []Wikidata  `json:"religion,omitempty"`
	Awards      []Award     `json:"awards,omitempty"`
	Ethnicity   []Wikidata  `json:"ethnicity,omitempty"`
	Military    []Military  `json:"military,omitempty"` // military branch
	RecordLabel []Wikidata  `json:"record_label,omitempty"`
	Discography []Wikidata  `json:"discography,omitempty"`
	Position    []Wikidata  `json:"position,omitempty"` // e.g. position on team...forward, center, etc..
	Partner     []Spouse    `json:"partner,omitempty"`
	Origin      []Wikidata  `json:"origin,omitempty"`         // country of origin
	DeathCause  []Wikidata  `json:"cause_of_death,omitempty"` // there is also P1196 "manner of death"
	Members     []Member    `json:"members,omitempty"`
	Residence   []Wikidata  `json:"residence,omitempty"`
	Hand        []Wikidata  `json:"hand,omitempty"` // left or right-handed
	//Coordinate  []Coordinate `json:"coordinate,omitempty"`
	Birthday    []DateTime   `json:"birthday,omitempty"`
	Death       []DateTime   `json:"death,omitempty"`
	Start       []DateTime   `json:"start,omitempty"`
	Sport       []Wikidata   `json:"sport,omitempty"`
	Drafted     []Wikidata   `json:"drafted,omitempty"`
	GivenName   []Wikidata   `json:"given_name,omitempty"`
	Influences  []Wikidata   `json:"influences,omitempty"`
	Location    []Wikidata   `json:"location,omitempty"`
	Website     []string     `json:"website,omitempty"`
	Population  []Population `json:"population,omitempty"`
	Instrument  []Instrument `json:"instrument,omitempty"` // Jimi Hendrix Fender Stratocaster
	Participant []Wikidata   `json:"participant,omitempty"`
	Nominations []Nomination `json:"nominations,omitempty"`
	Languages   []Wikidata   `json:"languages,omitempty"` // languages spoken and/or written proficiency
	BirthName   []Text       `json:"birth_name,omitempty"`
	Spotify     []string     `json:"spotify,omitempty"`
	Twitter     []string     `json:"twitter,omitempty"`
	Instagram   []string     `json:"instagram,omitempty"`
	Facebook    []string     `json:"facebook,omitempty"`
	YouTube     []string     `json:"youtube,omitempty"`
	WorkStart   []DateTime   `json:"work_start,omitempty"` // better name??? P571 is similar tag
	Height      []Quantity   `json:"height,omitempty"`
	Weight      []Quantity   `json:"weight,omitempty"`
	Siblings    []Wikidata   `json:"siblings,omitempty"`
}

// UnmarshalJSON formats and extracts only the info we need from claims
func (w *Wikidata) UnmarshalJSON(b []byte) error {
	type alias Wikidata

	raw := struct {
		*claims `json:"claims"`
		*alias
	}{
		claims: &claims{},
		alias:  (*alias)(w),
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	w.Claims = &Claims{}

	r := reflect.Indirect(reflect.ValueOf(raw.claims))

	for i := 0; i < r.NumField(); i++ {
		field := reflect.Indirect(reflect.ValueOf(w.Claims)).Field(i)

		for _, c := range r.Field(i).Interface().([]claim) {
			if c.MainSnak.DataValue.Value == nil {
				continue
			}

			switch field.Interface().(type) {
			case []DateTime, []Quantity, []string, []Text, []Wikidata, []Coordinate:
				if err := setValue(field, c.MainSnak); err != nil {
					return err
				}
			default: // eg has qualifiers
				var elem reflect.Value

				typ := field.Type().Elem()
				if typ.Kind() == reflect.Ptr {
					elem = reflect.New(typ.Elem())
				}
				if typ.Kind() == reflect.Struct {
					elem = reflect.New(typ).Elem()
				}

				for j := 0; j < reflect.Indirect(elem).NumField(); j++ {
					f := elem.Field(j)

					switch j {
					case 0:
						if err := setValue(f, c.MainSnak); err != nil {
							return err
						}
					default:
						// only set the qualifiers we want
						tag := strings.Split(elem.Type().Field(j).Tag.Get("property"), ",")[0]
						for k, qual := range c.Qualifiers {
							if tag == k {
								for _, q := range qual {
									if q.DataValue.Value == nil {
										continue
									}
									if err := setValue(f, q); err != nil {
										return err
									}
								}
							}
						}
					}
				}
				field.Set(reflect.Append(field, elem))
			}
		}
	}

	return nil
}

// reduce tries to reduce the number of items we store
// https://www.wikidata.org/wiki/Special:ListDatatypes
func setValue(field reflect.Value, p property) error {
	switch p.DataType {
	case "string", "commonsMedia", "url", "external-id":
		v := p.DataValue.Value.(string)
		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	case "monolingualtext":
		val := p.DataValue.Value.(map[string]interface{})
		v := Text{
			Text:     val["text"].(string),
			Language: val["language"].(string),
		}

		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	case "time":
		// We don't want to alter the timestamp...
		// eg don't turn '1965-00-00T00:00:00Z' to '1965-01-01T00:00:00Z'.
		// Also, we might need to add precision, etc...
		val := p.DataValue.Value.(map[string]interface{})
		c := val["calendarmodel"].(string)

		w := Wikidata{
			ID: strings.TrimPrefix(c, "http://www.wikidata.org/entity/"),
		}

		v := DateTime{
			Value:    strings.TrimPrefix(val["time"].(string), "+"),
			Calendar: w,
		}

		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	case "quantity":
		val := p.DataValue.Value.(map[string]interface{})
		w := Wikidata{
			ID: strings.TrimPrefix(val["unit"].(string), "http://www.wikidata.org/entity/"),
		}

		v := Quantity{
			Amount: strings.TrimPrefix(val["amount"].(string), "+"),
			Unit:   w,
		}

		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	case "globe-coordinate":
		val := p.DataValue.Value.(map[string]interface{})

		v := Coordinate{
			Globe: []Wikidata{
				Wikidata{
					ID: strings.TrimPrefix(val["globe"].(string), "http://www.wikidata.org/entity/"),
				},
			},
		}

		switch val["latitude"].(type) {
		case float64:
			v.Latitude = []float64{val["latitude"].(float64)}
		}

		switch val["longitude"].(type) {
		case float64:
			v.Longitude = []float64{val["longitude"].(float64)}
		}

		switch val["altitude"].(type) {
		case float64:
			v.Altitude = []float64{val["altitude"].(float64)}
		}

		switch val["precision"].(type) {
		case float64:
			v.Precision = []float64{val["precision"].(float64)}
		}

		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	case "wikibase-item":
		val := p.DataValue.Value.(map[string]interface{})
		v := Wikidata{
			ID: val["id"].(string),
		}

		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	default: // "math", "geo-shape", and "tabular-data" don't seem to be used
		return fmt.Errorf("unknown type %v", p.DataType)
	}

	return nil
}
