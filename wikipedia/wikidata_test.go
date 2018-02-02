package wikipedia

import (
	"fmt"
	"reflect"
	"testing"
)

var shaqRawLabels = `{"de": {"value": "Shaquille O’Neal", "language": "de"}, "en": {"value": "Shaquille O'Neal", "language": "en"}}`
var shaqLabels = Labels{
	"de": Text{Text: "Shaquille O’Neal", Language: "de"},
	"en": Text{Text: "Shaquille O'Neal", Language: "en"},
}

var shaqRawAliases = `{"de": [{"value": "Shaquille Rashaun O'Neal", "language": "de"}, {"value": "Shaq", "language": "de"}, {"value": "Shaquille O'Neal", "language": "de"}, {"value": "Shaquille Rashaun O’Neal", "language": "de"}], "en": [{"value": "Shaquille Rashaun O'Neal", "language": "en"}, {"value": "Shaq", "language": "en"}, {"value": "The Diesel", "language": "en"}, {"value": "Superman", "language": "en"}, {"value": "The Big Aristotle", "language": "en"}, {"value": "M.D.E. (Most Dominant Ever)", "language": "en"}, {"value": "L.C.L. (Last Center Left)", "language": "en"}, {"value": "The Big Shakespeare", "language": "en"}, {"value": "Wilt Chamberneezy", "language": "en"}, {"value": "Officer Shaq", "language": "en"}, {"value": "The Big Deporter", "language": "en"}, {"value": "Shaq Daddy", "language": "en"}, {"value": "Shaq-Fu", "language": "en"}, {"value": "Big Daddy", "language": "en"}, {"value": "The Big Baryshnikov", "language": "en"}, {"value": "Dr. Shaq", "language": "en"}, {"value": "Osama Bin Shaq", "language": "en"}, {"value": "The Big Cactus", "language": "en"}, {"value": "Warrior", "language": "en"}, {"value": "The Eighth Wonder", "language": "en"}, {"value": "Big Fella", "language": "en"}]}`
var shaqAliases = Aliases{
	"de": []Text{
		Text{Text: "Shaquille Rashaun O'Neal", Language: "de"}, Text{Text: "Shaq", Language: "de"},
		Text{Text: "Shaquille O'Neal", Language: "de"}, Text{Text: "Shaquille Rashaun O’Neal", Language: "de"},
	},
	"en": []Text{
		Text{Text: "Shaquille Rashaun O'Neal", Language: "en"}, Text{Text: "Shaq", Language: "en"},
		Text{Text: "The Diesel", Language: "en"}, Text{Text: "Superman", Language: "en"},
		Text{Text: "The Big Aristotle", Language: "en"}, Text{Text: "M.D.E. (Most Dominant Ever)", Language: "en"},
		Text{Text: "L.C.L. (Last Center Left)", Language: "en"}, Text{Text: "The Big Shakespeare", Language: "en"},
		Text{Text: "Wilt Chamberneezy", Language: "en"}, Text{Text: "Officer Shaq", Language: "en"},
		Text{Text: "The Big Deporter", Language: "en"}, Text{Text: "Shaq Daddy", Language: "en"},
		Text{Text: "Shaq-Fu", Language: "en"}, Text{Text: "Big Daddy", Language: "en"},
		{Text: "The Big Baryshnikov", Language: "en"}, {Text: "Dr. Shaq", Language: "en"},
		{Text: "Osama Bin Shaq", Language: "en"}, {Text: "The Big Cactus", Language: "en"},
		{Text: "Warrior", Language: "en"}, {Text: "The Eighth Wonder", Language: "en"},
		{Text: "Big Fella", Language: "en"},
	},
}

var shaqRawDescriptions = `{"de": {"value": "US-amerikanischer Basketballspieler", "language": "de"}, "en": {"value": "American basketball player", "language": "en"}}`
var shaqDescriptions = Descriptions{
	"de": Text{Text: "US-amerikanischer Basketballspieler", Language: "de"},
	"en": Text{Text: "American basketball player", Language: "en"},
}

var shaqRawClaims = `{
	"P18": [{
		"mainsnak": {
			"snaktype": "value",
			"property": "P18",
			"datavalue": {
				"value": "Shaqmiami.jpg",
				"type": "string"
			},
			"datatype": "commonsMedia"
		},
		"type": "statement",
		"id": "Q169452$40D3BCAD-5D62-467F-8AF1-7D162638B5A5",
		"rank": "normal",
		"references": [{
			"hash": "732ec1c90a6f0694c7db9a71bf09fe7f2b674172",
			"snaks": {
				"P143": [{
					"snaktype": "value",
					"property": "P143",
					"datavalue": {
						"value": {
							"entity-type": "item",
							"numeric-id": 10000,
							"id": "Q10000"
						},
						"type": "wikibase-entityid"
					},
					"datatype": "wikibase-item"
				}]
			},
			"snaks-order": [
				"P143"
			]
		}]
	}],
	"P21": [{
		"mainsnak": {
			"snaktype": "value",
			"property": "P21",
			"datavalue": {
				"value": {
					"entity-type": "item",
					"numeric-id": 6581097,
					"id": "Q6581097"
				},
				"type": "wikibase-entityid"
			},
			"datatype": "wikibase-item"
		},
		"type": "statement",
		"id": "q169452$3E7B6D0C-D465-4681-BFD3-CCAA213C1813",
		"rank": "normal",
		"references": [{
				"hash": "bd49d3e4f67bc460ce7a06b6ac3027347cf5ee55",
				"snaks": {
					"P143": [{
						"snaktype": "value",
						"property": "P143",
						"datavalue": {
							"value": {
								"entity-type": "item",
								"numeric-id": 169514,
								"id": "Q169514"
							},
							"type": "wikibase-entityid"
						},
						"datatype": "wikibase-item"
					}]
				},
				"snaks-order": [
					"P143"
				]
			},
			{
				"hash": "d5847b9b6032aa8b13dae3c2dfd9ed5d114d21b3",
				"snaks": {
					"P143": [{
						"snaktype": "value",
						"property": "P143",
						"datavalue": {
							"value": {
								"entity-type": "item",
								"numeric-id": 11920,
								"id": "Q11920"
							},
							"type": "wikibase-entityid"
						},
						"datatype": "wikibase-item"
					}]
				},
				"snaks-order": [
					"P143"
				]
			},
			{
				"hash": "1e95bd1cdd89a73487b69026d14b67bdf322e862",
				"snaks": {
					"P248": [{
						"snaktype": "value",
						"property": "P248",
						"datavalue": {
							"value": {
								"entity-type": "item",
								"numeric-id": 20666306,
								"id": "Q20666306"
							},
							"type": "wikibase-entityid"
						},
						"datatype": "wikibase-item"
					}],
					"P813": [{
						"snaktype": "value",
						"property": "P813",
						"datavalue": {
							"value": {
								"time": "+2015-10-10T00:00:00Z",
								"timezone": 0,
								"before": 0,
								"after": 0,
								"precision": 11,
								"calendarmodel": "http:\/\/www.wikidata.org\/entity\/Q1985727"
							},
							"type": "time"
						},
						"datatype": "time"
					}],
					"P854": [{
						"snaktype": "value",
						"property": "P854",
						"datavalue": {
							"value": "http:\/\/data.bnf.fr\/ark:\/12148\/cb13972273t",
							"type": "string"
						},
						"datatype": "url"
					}]
				},
				"snaks-order": [
					"P248",
					"P813",
					"P854"
				]
			}
		]
	}],
	"P26": [{
		"mainsnak": {
			"snaktype": "value",
			"property": "P26",
			"datavalue": {
				"value": {
					"entity-type": "item",
					"numeric-id": 7491045,
					"id": "Q7491045"
				},
				"type": "wikibase-entityid"
			},
			"datatype": "wikibase-item"
		},
		"type": "statement",
		"qualifiers": {
			"P582": [{
				"snaktype": "value",
				"property": "P582",
				"hash": "9c0f3eb5795ad1f49b9febe37b921ca66cabec98",
				"datavalue": {
					"value": {
						"time": "+2009-11-10T00:00:00Z",
						"timezone": 0,
						"before": 0,
						"after": 0,
						"precision": 11,
						"calendarmodel": "http:\/\/www.wikidata.org\/entity\/Q1985727"
					},
					"type": "time"
				},
				"datatype": "time"
			}],
			"P580": [{
				"snaktype": "value",
				"property": "P580",
				"hash": "d3ab4bb8dc6bce5ab9900c09830f7f5fa5eee710",
				"datavalue": {
					"value": {
						"time": "+2002-12-26T00:00:00Z",
						"timezone": 0,
						"before": 0,
						"after": 0,
						"precision": 11,
						"calendarmodel": "http:\/\/www.wikidata.org\/entity\/Q1985727"
					},
					"type": "time"
				},
				"datatype": "time"
			}],
			"P276": [{
				"snaktype": "value",
				"property": "P276",
				"hash": "c39c07ee9b7bcd11304302ac0ee0babba63de5ea",
				"datavalue": {
					"value": {
						"entity-type": "item",
						"numeric-id": 2021533,
						"id": "Q2021533"
					},
					"type": "wikibase-entityid"
				},
				"datatype": "wikibase-item"
			}],
			"P2842": [{
				"snaktype": "value",
				"property": "P2842",
				"hash": "9ac00672af441ed3f87aaab0ddc35081b38db98a",
				"datavalue": {
					"value": {
						"entity-type": "item",
						"numeric-id": 2021533,
						"id": "Q2021533"
					},
					"type": "wikibase-entityid"
				},
				"datatype": "wikibase-item"
			}]
		},
		"qualifiers-order": [
			"P582",
			"P580",
			"P276",
			"P2842"
		],
		"id": "Q169452$08BFA5B0-1414-45EB-B978-DE0F55DABBB1",
		"rank": "normal"
	}],
	"P2048": [{
		"mainsnak": {
			"snaktype": "value",
			"property": "P2048",
			"datavalue": {
				"value": {
					"amount": "+2.16",
					"unit": "http:\/\/www.wikidata.org\/entity\/Q11573",
					"upperBound": "+2.17",
					"lowerBound": "+2.15"
				},
				"type": "quantity"
			},
			"datatype": "quantity"
		},
		"type": "statement",
		"id": "Q169452$2ab9184a-4b78-7344-51b7-4658604bf594",
		"rank": "normal",
		"references": [{
			"hash": "676146481b933e27e3de51d02eb6f244ec5ef635",
			"snaks": {
				"P813": [{
					"snaktype": "value",
					"property": "P813",
					"datavalue": {
						"value": {
							"time": "+2016-01-28T00:00:00Z",
							"timezone": 0,
							"before": 0,
							"after": 0,
							"precision": 11,
							"calendarmodel": "http:\/\/www.wikidata.org\/entity\/Q1985727"
						},
						"type": "time"
					},
					"datatype": "time"
				}],
				"P143": [{
					"snaktype": "value",
					"property": "P143",
					"datavalue": {
						"value": {
							"entity-type": "item",
							"numeric-id": 328,
							"id": "Q328"
						},
						"type": "wikibase-entityid"
					},
					"datatype": "wikibase-item"
				}]
			},
			"snaks-order": [
				"P813",
				"P143"
			]
		}]
	}],
	"P1477": [{
		"mainsnak": {
			"snaktype": "value",
			"property": "P1477",
			"hash": "",
			"datavalue": {
				"value": {
					"text": "This is not the real name for Shaq",
					"language": "en"
				},
				"type": "monolingualtext"
			},
			"datatype": "monolingualtext"
		},
		"type": "statement",
		"id": "",
		"rank": "preferred",
		"references": [{
			"hash": "",
			"snaks": {
				"P143": [{
					"snaktype": "value",
					"property": "P143",
					"hash": "",
					"datavalue": {
						"value": {
							"entity-type": "item"
						},
						"type": "wikibase-entityid"
					},
					"datatype": "wikibase-item"
				}]
			},
			"snaks-order": [
				"P143"
			]
		}]
	}]
}`

var shaqClaims = &Claims{
	Image: []string{"Shaqmiami.jpg"},
	Sex: []Wikidata{
		Wikidata{
			ID: "Q6581097",
		},
	},
	Spouse: []Spouse{
		Spouse{
			Item: []Wikidata{
				Wikidata{
					ID: "Q7491045",
				},
			},
			Start: []DateTime{
				DateTime{
					Value: "2002-12-26T00:00:00Z", Calendar: Wikidata{ID: "Q1985727"},
				},
			},
			End: []DateTime{
				DateTime{
					Value: "2009-11-10T00:00:00Z", Calendar: Wikidata{ID: "Q1985727"},
				},
			},
			Place: []Wikidata{
				Wikidata{
					ID: "Q2021533",
				},
			},
		},
	},
	Height: []Quantity{
		Quantity{
			Amount: "2.16",
			Unit:   Wikidata{ID: "Q11573"},
		},
	},
	BirthName: []Text{
		Text{
			Text:     "This is not the real name for Shaq",
			Language: "en",
		},
	},
}

var shaqRawJSON = []byte(fmt.Sprintf(
	`{"type": "item", "id": "Q169452", "labels": %v, "aliases": %v, "descriptions": %v, "claims": %v}`,
	shaqRawLabels, shaqRawAliases, shaqRawDescriptions, shaqRawClaims,
))

var italyRawClaims = `{
	"P625": [
	{
		"mainsnak": {
		"snaktype": "value",
		"property": "P625",
		"datavalue": {
			"value": {
			"latitude": 43,
			"longitude": 12,
			"altitude": 1,
			"precision": 1,
			"globe": "http:\/\/www.wikidata.org\/entity\/Q2"
			},
			"type": "globecoordinate"
		},
		"datatype": "globe-coordinate"
		},
		"type": "statement",
		"id": "q38$6be5dd88-4e91-89b3-92d4-2cc199659ce7",
		"rank": "normal"
	}
	]
}
`

var shaqWikidata = &Wikidata{
	ID:           "Q169452",
	Labels:       shaqLabels,
	Aliases:      shaqAliases,
	Descriptions: shaqDescriptions,
	Claims:       shaqClaims,
}

var italyRawJSON = []byte(fmt.Sprintf(`{"type": "item", "id": "Q38", "claims": %v}`, italyRawClaims))

var italyClaims = &Claims{ // coordinate not implemented yet
/*
	Coordinate: []Coordinate{
		Coordinate{
			Latitude:  []float64{float64(43)},
			Longitude: []float64{float64(12)},
			Altitude:  []float64{float64(1)},
			Precision: []float64{float64(1)},
			Globe:     []Wikidata{Wikidata{ID: "Q2"}},
		},
	},
*/
}

func TestWikiData_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want *Wikidata
	}{
		{
			"shaq",
			args{shaqRawJSON},
			shaqWikidata,
		},
		{
			"italy",
			args{italyRawJSON},
			&Wikidata{
				ID:     "Q38",
				Claims: italyClaims,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Wikidata{}

			if err := got.UnmarshalJSON(tt.args.b); err != nil {
				t.Errorf("Wikidata.UnmarshalJSON() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v; want %+v", got, tt.want)
			}
		})
	}
}
