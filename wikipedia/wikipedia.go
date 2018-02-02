// Package wikipedia fetches Wikipedia articles
package wikipedia

import (
	"encoding/json"
	"regexp"
	"strings"

	"golang.org/x/text/language"
)

// Fetcher outlines the methods used to retrieve Wikipedia snippets
type Fetcher interface {
	Setup() error
	Fetch(query string, lang language.Tag) (*Item, error)
}

// Item is the text portion of a wikipedia article
type Item struct {
	Wikipedia
	*Wikidata
}

// Wikipedia holds the summary text of an article
type Wikipedia struct {
	ID       string `json:"wikibase_item,omitempty"`
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
	truncate int
	//Popularity float32 `json:"popularity_score"` // I can't seem to find any documentation for this
}

var reParen = regexp.MustCompile(`\s?\((.*?)\)`) // replace parenthesis

// UnmarshalJSON truncates the text
func (w *Wikipedia) UnmarshalJSON(data []byte) error {
	// copy the fields of Wikipedia but not the
	// methods so we don't recursively call UnmarshalJSON
	type Alias Wikipedia
	a := &struct {
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	w.Text = reParen.ReplaceAllString(w.Text, "")
	w.Text = strings.Replace(w.Text, "\u00a0", "", -1) // otherwise causes a panic below

	if len(w.Text) > w.truncate { // truncates while preserving words.
		c := strings.Fields(w.Text[:w.truncate+1])
		w.Text = strings.Join(c[0:len(c)-1], " ")
		if !strings.HasSuffix(w.Text, ".") {
			w.Text = w.Text + " ..."
		}
	}

	return nil
}

// Languages verifies languages based on Wikipedia's supported languages.
// An empty slice of supported languages implies you support every language available.
func Languages(supported []language.Tag) ([]language.Tag, []language.Tag) {
	// make sure supported languages are supported by Wikipedia
	s := []language.Tag{}
	unsupported := []language.Tag{}

	switch len(supported) {
	case 0:
		for lang := range Available {
			s = append(s, lang)
		}
	default:
		for _, lang := range supported {
			if _, ok := Available[lang]; !ok {
				unsupported = append(unsupported, lang)
				continue
			}

			s = append(s, lang)
		}
	}

	return s, unsupported
}

// Available is a map of all languages that Wikipedia supports.
// https://en.wikipedia.org/wiki/List_of_Wikipedias
// We sort their table by # of Articles descending.
var Available = map[language.Tag]struct{}{
	language.MustParse("en"):         struct{}{}, // english is fallback
	language.MustParse("ceb"):        struct{}{},
	language.MustParse("sv"):         struct{}{},
	language.MustParse("de"):         struct{}{},
	language.MustParse("nl"):         struct{}{},
	language.MustParse("fr"):         struct{}{},
	language.MustParse("ru"):         struct{}{},
	language.MustParse("it"):         struct{}{},
	language.MustParse("es"):         struct{}{},
	language.MustParse("war"):        struct{}{},
	language.MustParse("pl"):         struct{}{},
	language.MustParse("vi"):         struct{}{},
	language.MustParse("ja"):         struct{}{},
	language.MustParse("pt"):         struct{}{},
	language.MustParse("zh"):         struct{}{},
	language.MustParse("uk"):         struct{}{},
	language.MustParse("ca"):         struct{}{},
	language.MustParse("fa"):         struct{}{},
	language.MustParse("ar"):         struct{}{},
	language.MustParse("no"):         struct{}{},
	language.MustParse("sh"):         struct{}{},
	language.MustParse("fi"):         struct{}{},
	language.MustParse("hu"):         struct{}{},
	language.MustParse("id"):         struct{}{},
	language.MustParse("ro"):         struct{}{},
	language.MustParse("cs"):         struct{}{},
	language.MustParse("ko"):         struct{}{},
	language.MustParse("sr"):         struct{}{},
	language.MustParse("tr"):         struct{}{},
	language.MustParse("ms"):         struct{}{},
	language.MustParse("eu"):         struct{}{},
	language.MustParse("eo"):         struct{}{},
	language.MustParse("bg"):         struct{}{},
	language.MustParse("da"):         struct{}{},
	language.MustParse("min"):        struct{}{},
	language.MustParse("kk"):         struct{}{},
	language.MustParse("sk"):         struct{}{},
	language.MustParse("hy"):         struct{}{},
	language.MustParse("zh-min-nan"): struct{}{},
	language.MustParse("he"):         struct{}{},
	language.MustParse("lt"):         struct{}{},
	language.MustParse("hr"):         struct{}{},
	language.MustParse("ce"):         struct{}{},
	language.MustParse("sl"):         struct{}{},
	language.MustParse("et"):         struct{}{},
	language.MustParse("gl"):         struct{}{},
	language.MustParse("nn"):         struct{}{},
	language.MustParse("uz"):         struct{}{},
	language.MustParse("el"):         struct{}{},
	language.MustParse("be"):         struct{}{},
	language.MustParse("la"):         struct{}{},
	//language.MustParse("simple"):struct{}{}, // Simple English...does not parse
	language.MustParse("vo"):        struct{}{},
	language.MustParse("hi"):        struct{}{},
	language.MustParse("ur"):        struct{}{},
	language.MustParse("th"):        struct{}{},
	language.MustParse("az"):        struct{}{},
	language.MustParse("ka"):        struct{}{},
	language.MustParse("ta"):        struct{}{},
	language.MustParse("cy"):        struct{}{},
	language.MustParse("mk"):        struct{}{},
	language.MustParse("mg"):        struct{}{},
	language.MustParse("oc"):        struct{}{},
	language.MustParse("lv"):        struct{}{},
	language.MustParse("bs"):        struct{}{},
	language.MustParse("new"):       struct{}{},
	language.MustParse("tt"):        struct{}{},
	language.MustParse("tg"):        struct{}{},
	language.MustParse("te"):        struct{}{},
	language.MustParse("tl"):        struct{}{},
	language.MustParse("sq"):        struct{}{},
	language.MustParse("pms"):       struct{}{},
	language.MustParse("ky"):        struct{}{},
	language.MustParse("br"):        struct{}{},
	language.MustParse("be-tarask"): struct{}{},
	language.MustParse("zh-yue"):    struct{}{},
	language.MustParse("ht"):        struct{}{},
	language.MustParse("jv"):        struct{}{},
	language.MustParse("ast"):       struct{}{},
	language.MustParse("bn"):        struct{}{},
	language.MustParse("lb"):        struct{}{},
	language.MustParse("ml"):        struct{}{},
	language.MustParse("mr"):        struct{}{},
	language.MustParse("af"):        struct{}{},
	language.MustParse("pnb"):       struct{}{},
	language.MustParse("sco"):       struct{}{},
	language.MustParse("is"):        struct{}{},
	language.MustParse("ga"):        struct{}{},
	language.MustParse("cv"):        struct{}{},
	language.MustParse("ba"):        struct{}{},
	language.MustParse("fy"):        struct{}{},
	language.MustParse("sw"):        struct{}{},
	language.MustParse("my"):        struct{}{},
	language.MustParse("lmo"):       struct{}{},
	language.MustParse("an"):        struct{}{},
	language.MustParse("yo"):        struct{}{},
	language.MustParse("ne"):        struct{}{},
	language.MustParse("io"):        struct{}{},
	language.MustParse("gu"):        struct{}{},
	language.MustParse("nds"):       struct{}{},
	language.MustParse("scn"):       struct{}{},
	language.MustParse("bpy"):       struct{}{},
	language.MustParse("pa"):        struct{}{},
	language.MustParse("ku"):        struct{}{},
	language.MustParse("als"):       struct{}{},
	language.MustParse("bar"):       struct{}{},
	language.MustParse("kn"):        struct{}{},
	language.MustParse("qu"):        struct{}{},
	language.MustParse("ia"):        struct{}{},
	language.MustParse("su"):        struct{}{},
	language.MustParse("ckb"):       struct{}{},
	language.MustParse("mn"):        struct{}{},
	language.MustParse("arz"):       struct{}{},
	language.MustParse("bat-smg"):   struct{}{},
	language.MustParse("azb"):       struct{}{},
	language.MustParse("nap"):       struct{}{},
	language.MustParse("wa"):        struct{}{},
	language.MustParse("gd"):        struct{}{},
	language.MustParse("bug"):       struct{}{},
	language.MustParse("yi"):        struct{}{},
	language.MustParse("am"):        struct{}{},
	language.MustParse("map-bms"):   struct{}{},
	language.MustParse("si"):        struct{}{},
	language.MustParse("fo"):        struct{}{},
	language.MustParse("mzn"):       struct{}{},
	language.MustParse("or"):        struct{}{},
	language.MustParse("li"):        struct{}{},
	language.MustParse("sah"):       struct{}{},
	language.MustParse("hsb"):       struct{}{},
	language.MustParse("vec"):       struct{}{},
	language.MustParse("sa"):        struct{}{},
	language.MustParse("os"):        struct{}{},
	language.MustParse("mai"):       struct{}{},
	language.MustParse("ilo"):       struct{}{},
	language.MustParse("mrj"):       struct{}{},
	language.MustParse("hif"):       struct{}{},
	language.MustParse("mhr"):       struct{}{},
	language.MustParse("xmf"):       struct{}{},
	//language.MustParse("roa-tara"):struct{}{}, // Does not parse
	language.MustParse("nah"): struct{}{},
	//language.MustParse("eml"):struct{}{}, // Does not parse
	language.MustParse("bh"):      struct{}{},
	language.MustParse("pam"):     struct{}{},
	language.MustParse("ps"):      struct{}{},
	language.MustParse("nso"):     struct{}{},
	language.MustParse("diq"):     struct{}{},
	language.MustParse("hak"):     struct{}{},
	language.MustParse("sd"):      struct{}{},
	language.MustParse("se"):      struct{}{},
	language.MustParse("mi"):      struct{}{},
	language.MustParse("bcl"):     struct{}{},
	language.MustParse("nds-nl"):  struct{}{},
	language.MustParse("gan"):     struct{}{},
	language.MustParse("glk"):     struct{}{},
	language.MustParse("vls"):     struct{}{},
	language.MustParse("rue"):     struct{}{},
	language.MustParse("bo"):      struct{}{},
	language.MustParse("wuu"):     struct{}{},
	language.MustParse("szl"):     struct{}{},
	language.MustParse("fiu-vro"): struct{}{},
	language.MustParse("sc"):      struct{}{},
	language.MustParse("co"):      struct{}{},
	language.MustParse("vep"):     struct{}{},
	language.MustParse("lrc"):     struct{}{},
	language.MustParse("tk"):      struct{}{},
	language.MustParse("csb"):     struct{}{},
	//language.MustParse("zh-classical"):struct{}{}, // Does not parse
	language.MustParse("crh"):     struct{}{},
	language.MustParse("km"):      struct{}{},
	language.MustParse("gv"):      struct{}{},
	language.MustParse("kv"):      struct{}{},
	language.MustParse("frr"):     struct{}{},
	language.MustParse("as"):      struct{}{},
	language.MustParse("lad"):     struct{}{},
	language.MustParse("zea"):     struct{}{},
	language.MustParse("so"):      struct{}{},
	language.MustParse("cdo"):     struct{}{},
	language.MustParse("ace"):     struct{}{},
	language.MustParse("ay"):      struct{}{},
	language.MustParse("udm"):     struct{}{},
	language.MustParse("kw"):      struct{}{},
	language.MustParse("stq"):     struct{}{},
	language.MustParse("nrm"):     struct{}{},
	language.MustParse("ie"):      struct{}{},
	language.MustParse("lez"):     struct{}{},
	language.MustParse("myv"):     struct{}{},
	language.MustParse("koi"):     struct{}{},
	language.MustParse("rm"):      struct{}{},
	language.MustParse("pcd"):     struct{}{},
	language.MustParse("ug"):      struct{}{},
	language.MustParse("lij"):     struct{}{},
	language.MustParse("mt"):      struct{}{},
	language.MustParse("fur"):     struct{}{},
	language.MustParse("gn"):      struct{}{},
	language.MustParse("dsb"):     struct{}{},
	language.MustParse("gom"):     struct{}{},
	language.MustParse("dv"):      struct{}{},
	language.MustParse("cbk-zam"): struct{}{},
	language.MustParse("ext"):     struct{}{},
	language.MustParse("ang"):     struct{}{},
	language.MustParse("kab"):     struct{}{},
	language.MustParse("mwl"):     struct{}{},
	language.MustParse("ksh"):     struct{}{},
	language.MustParse("ln"):      struct{}{},
	language.MustParse("gag"):     struct{}{},
	language.MustParse("sn"):      struct{}{},
	language.MustParse("nv"):      struct{}{},
	language.MustParse("frp"):     struct{}{},
	language.MustParse("pag"):     struct{}{},
	language.MustParse("pi"):      struct{}{},
	language.MustParse("av"):      struct{}{},
	language.MustParse("lo"):      struct{}{},
	language.MustParse("dty"):     struct{}{},
	language.MustParse("xal"):     struct{}{},
	language.MustParse("pfl"):     struct{}{},
	language.MustParse("krc"):     struct{}{},
	language.MustParse("haw"):     struct{}{},
	language.MustParse("kaa"):     struct{}{},
	language.MustParse("olo"):     struct{}{},
	language.MustParse("bxr"):     struct{}{},
	language.MustParse("rw"):      struct{}{},
	language.MustParse("pdc"):     struct{}{},
	language.MustParse("pap"):     struct{}{},
	language.MustParse("bjn"):     struct{}{},
	language.MustParse("to"):      struct{}{},
	language.MustParse("nov"):     struct{}{},
	language.MustParse("kl"):      struct{}{},
	language.MustParse("arc"):     struct{}{},
	language.MustParse("jam"):     struct{}{},
	language.MustParse("kbd"):     struct{}{},
	language.MustParse("ha"):      struct{}{},
	language.MustParse("tet"):     struct{}{},
	language.MustParse("tyv"):     struct{}{},
	language.MustParse("tpi"):     struct{}{},
	language.MustParse("ki"):      struct{}{},
	language.MustParse("ig"):      struct{}{},
	language.MustParse("na"):      struct{}{},
	language.MustParse("ab"):      struct{}{},
	language.MustParse("lbe"):     struct{}{},
	language.MustParse("roa-rup"): struct{}{},
	language.MustParse("jbo"):     struct{}{},
	language.MustParse("ty"):      struct{}{},
	language.MustParse("kg"):      struct{}{},
	language.MustParse("za"):      struct{}{},
	language.MustParse("lg"):      struct{}{},
	language.MustParse("wo"):      struct{}{},
	language.MustParse("mdf"):     struct{}{},
	language.MustParse("srn"):     struct{}{},
	language.MustParse("zu"):      struct{}{},
	language.MustParse("bi"):      struct{}{},
	language.MustParse("ltg"):     struct{}{},
	language.MustParse("chr"):     struct{}{},
	language.MustParse("tcy"):     struct{}{},
	language.MustParse("sm"):      struct{}{},
	language.MustParse("om"):      struct{}{},
	language.MustParse("tn"):      struct{}{},
	language.MustParse("chy"):     struct{}{},
	language.MustParse("xh"):      struct{}{},
	language.MustParse("tw"):      struct{}{},
	language.MustParse("cu"):      struct{}{},
	language.MustParse("rmy"):     struct{}{},
	language.MustParse("tum"):     struct{}{},
	language.MustParse("pih"):     struct{}{},
	language.MustParse("rn"):      struct{}{},
	language.MustParse("got"):     struct{}{},
	language.MustParse("pnt"):     struct{}{},
	language.MustParse("ss"):      struct{}{},
	language.MustParse("ch"):      struct{}{},
	language.MustParse("bm"):      struct{}{},
	language.MustParse("ady"):     struct{}{},
	language.MustParse("mo"):      struct{}{},
	language.MustParse("ts"):      struct{}{},
	language.MustParse("iu"):      struct{}{},
	language.MustParse("st"):      struct{}{},
	language.MustParse("ny"):      struct{}{},
	language.MustParse("fj"):      struct{}{},
	language.MustParse("ee"):      struct{}{},
	language.MustParse("ak"):      struct{}{},
	language.MustParse("ks"):      struct{}{},
	language.MustParse("sg"):      struct{}{},
	language.MustParse("ik"):      struct{}{},
	language.MustParse("ve"):      struct{}{},
	language.MustParse("dz"):      struct{}{},
	language.MustParse("ff"):      struct{}{},
	language.MustParse("ti"):      struct{}{},
	language.MustParse("cr"):      struct{}{},
	language.MustParse("ng"):      struct{}{},
	language.MustParse("cho"):     struct{}{},
	language.MustParse("kj"):      struct{}{},
	language.MustParse("mh"):      struct{}{},
	language.MustParse("ho"):      struct{}{},
	language.MustParse("ii"):      struct{}{},
	language.MustParse("aa"):      struct{}{},
	language.MustParse("mus"):     struct{}{},
	language.MustParse("hz"):      struct{}{},
	language.MustParse("kr"):      struct{}{},
	language.MustParse("hil"):     struct{}{},
	language.MustParse("kbp"):     struct{}{},
	language.MustParse("din"):     struct{}{},
}
