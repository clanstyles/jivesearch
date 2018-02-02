package frontend

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jivesearch/jivesearch/log"
	"github.com/jivesearch/jivesearch/wikipedia"
	"golang.org/x/text/language"
)

var funcMap = template.FuncMap{
	"Commafy":       commafy,
	"SafeHTML":      safeHTML,
	"Truncate":      truncate,
	"HMACKey":       hmacKey,
	"WikiCanonical": wikiCanonical,
	"WikiDateTime":  wikiDateTime,
	"WikiLabel":     wikiLabel,
	"WikiJoin":      wikiJoin,
	"WikiAmount":    wikiAmount,
}

// where did this come from?
func commafy(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)

		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

func safeHTML(value string) template.HTML {
	return template.HTML(value)
}

// Preserving words is a crude translation from the python answer:
// http://stackoverflow.com/questions/250357/truncate-a-string-without-ending-in-the-middle-of-a-word
func truncate(txt string, max int, preserve bool) string {
	if len(txt) <= max {
		return txt
	}

	if preserve {
		c := strings.Fields(txt[:max+1])
		return strings.Join(c[0:len(c)-1], " ") + " ..."
	}

	return txt[:max] + "..."
}

var hmacSecret = func() string {
	secret := os.Getenv("hmac_secret")
	if secret == "" {
		log.Info.Println(`hmac secret for image proxy is blank. Please set the "hmac_secret" env variable`)
	}
	return secret
}

// hmacKey generates an hmac key for our reverse image proxy
func hmacKey(u string) string {
	secret := hmacSecret()

	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(u)); err != nil {
		log.Info.Println(err)
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// wikiCanonical returns the canonical form of a wikipedia title.
// if this breaks Wikidata dumps have "sitelinks"
func wikiCanonical(t string) string {
	return strings.Replace(t, " ", "_", -1)
}

// wikiDateTime tries to format a date with optional time.
// We assume Gregorian calendar below. (Julian calendar TODO).
// Note: Wikidata only uses Gregorian and Julian calendars.
func wikiDateTime(dt wikipedia.DateTime, includeAge bool) string {
	for j, f := range []string{time.RFC3339Nano, "2006"} {
		var ff string

		switch j {
		case 1:
			dt.Value = dt.Value[:4]
			ff = f
		default:
			ff = "January 2, 2006"
		}

		t, err := time.Parse(f, dt.Value)
		if err != nil {
			log.Debug.Println(err)
			continue
		}

		formatted := t.Format(ff)

		if includeAge {
			now := time.Now()
			years := now.Year() - t.Year()
			if now.YearDay() < t.YearDay() {
				years--
			}

			return fmt.Sprintf("%v (age %d)", formatted, years)
		}

		return formatted
	}

	return ""
}

// wikiLabel extracts the closest label for a Wikipedia Item using a language matcher
func wikiLabel(labels map[string]wikipedia.Text, preferred []language.Tag) string {
	// create a matcher based on the available labels
	langs := []language.Tag{}

	for k := range labels {
		t, err := language.Parse(k)
		if err != nil { // sr-el doesn't parse
			continue
		}

		langs = append(langs, t)
	}

	m := language.NewMatcher(langs)
	lang, _, _ := m.Match(preferred...)

	label := labels[lang.String()]
	return label.Text
}

// wikiJoin joins a slice of Wikidata items
func wikiJoin(items []wikipedia.Wikidata, preferred []language.Tag) string {
	sl := []string{}
	for _, item := range items {
		sl = append(sl, wikiLabel(item.Labels, preferred))
	}

	return strings.Join(sl, ", ")
}

// wikiAmount displays a unit in meters if non-US.
func wikiAmount(q wikipedia.Quantity, r language.Region) string {
	var f string

	amt, err := strconv.ParseFloat(q.Amount, 64)
	if err != nil {
		log.Debug.Println(err)
		return ""
	}

	switch r.String() {
	case "US", "LR", "MM": // only 3 countries that don't use metric system
		switch q.Unit.ID {
		case "Q11573", "Q174728", "Q218593":
			if q.Unit.ID == "Q11573" { // 1 meter = 39.3701 inches
				amt = amt * 39.3701
			} else if q.Unit.ID == "Q174728" { // 1 cm = 0.393701 inches
				amt = amt * .393701
			}

			if amt < 12 {
				f = fmt.Sprintf(`%f"`, amt)
			} else {
				f = fmt.Sprintf(`%d'%d"`, int(amt)/int(12), int(math.Mod(amt, 12)))
			}

		case "Q11570": // 1 kilogram = 2.20462 lbs
			amt = amt * 2.20462
			f = fmt.Sprintf("%d lbs", int(amt+.5))

		default:
			log.Debug.Printf("unknown unit %v\n", q.Unit.ID)
		}
	default:
		s := strconv.FormatFloat(amt, 'f', -1, 64)

		switch q.Unit.ID {
		case "Q11573":
			f = fmt.Sprintf("%v %v", s, "m")
		case "Q174728":
			f = fmt.Sprintf("%v %v", s, "cm")
		case "Q218593":
			amt = amt / .393701
			f = fmt.Sprintf("%v %v", int(amt+.5), "cm")
		case "Q11570":
			f = fmt.Sprintf("%v %v", s, "kg")
		default:
			log.Debug.Printf("unknown unit %v\n", q.Unit.ID)
		}
	}

	return f
}
