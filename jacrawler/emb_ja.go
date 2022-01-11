package jacrawler

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	japanese struct {
		*common
		yearModifier int
		iso          string
	}
)

const (
	reiwaBegin = 2018
)

func NewJapaneseEmbassy() IEmbassyData {
	jap := &japanese{
		yearModifier: reiwaBegin,
		iso:          "ja",
		common: &common{
			uri:               "https://www.mofa.go.jp/mofaj/ca/cp/page22_003380.html",
			htmlSearchElement: "div[class=rightalign]",
			dateLayout:        "2006-01-02",
			pattern:           fmt.Sprintf(`.{2}(?P<%s>\d)年(?P<%s>\d{2})月(?P<%s>\d{1,2})日`, pYear, pMonth, pDay),
		},
	}

	jap.regex = regexp.MustCompile(jap.pattern)
	return jap
}

func (j *japanese) IsDateUpdated(data *model.Embassy, db ports.IConnection) bool {
	c := new(model.Embassy)
	c.EmbassyISO = j.iso
	if err := db.FetchLatestDateFromEmbassy(c); err != nil {
		log.Printf("There was an error fetching data from db: %s\n", err)
		return true
	}

	return c.UpdatedDate.After(data.UpdatedDate) || c.UpdatedDate.Equal(data.UpdatedDate)
}

func (j *japanese) UpdateDate(data *model.Embassy, db ports.IConnection) error {
	data.EmbassyISO = j.iso
	return db.Save(data)
}

func (j *japanese) GetUri() string {
	return j.uri
}

func (j *japanese) GetPattern() string {
	return j.pattern
}

func (j *japanese) GetDateLayout() string {
	return j.dateLayout
}

func (j *japanese) GetHtmlSearchElement() string {
	return j.htmlSearchElement
}

func (j *japanese) GetRegex() *regexp.Regexp {
	return j.regex
}

func (j *japanese) YearModifier() int {
	return reiwaBegin
}

func (j *japanese) GetDateValue(data map[string]string) string {
	return fmt.Sprintf(`%s-%s-%s`, data[pYear], data[pMonth], data[pDay])
}

func (j *japanese) GetISO() string {
	return j.iso
}

func (j *japanese) GetUpdatedDateFromText(el *colly.HTMLElement) (time.Time, bool, error) {
	if strings.TrimSpace(el.Text) == "" {
		return time.Time{}, false, nil
	}

	data := getParams(j, el.Text)
	if data[pYear] == "" {
		return time.Time{}, false, nil
	}

	y, err := strconv.Atoi(data[pYear])
	if err != nil {
		return time.Time{}, true, err
	}

	data[pYear] = strconv.Itoa(j.YearModifier() + y)

	pt, err := time.Parse(j.GetDateLayout(), j.GetDateValue(data))
	return pt, false, err
}
