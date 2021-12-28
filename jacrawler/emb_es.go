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
	spanish struct {
		*common
		iso          string
		yearModifier int
	}
)

func (s *spanish) GetUri() string {
	return s.uri
}

func (s *spanish) GetPattern() string {
	return s.pattern
}

func (s *spanish) GetDateLayout() string {
	return s.dateLayout
}

func (s *spanish) GetHtmlSearchElement() string {
	return s.htmlSearchElement
}

func (s *spanish) GetRegex() *regexp.Regexp {
	return s.regex
}

func (s *spanish) GetDateValue(data map[string]string) string {
	return fmt.Sprintf(`%s/%s/%s`, data[pYear], data[pMonth], data[pDay])
}

func (s *spanish) YearModifier() int {
	return s.yearModifier
}

func (s *spanish) IsDateUpdated(embassy *model.Embassy, connection ports.IConnection) bool {
	c := new(model.Embassy)
	c.EmbassyISO = s.iso
	if err := connection.FetchLatestDateFromEmbassy(c); err != nil {
		log.Printf("There was an error fetching data from db: %s\n", err)
		return true
	}

	return c.UpdatedDate.After(embassy.UpdatedDate) || c.UpdatedDate.Equal(embassy.UpdatedDate)
}

func (s *spanish) UpdateDate(embassy *model.Embassy, connection ports.IConnection) error {
	embassy.EmbassyISO = s.iso
	return connection.Save(embassy)
}

func (s *spanish) GetISO() string {
	return s.iso
}

func (s *spanish) GetUpdatedDateFromText(element *colly.HTMLElement) (time.Time, bool, error) {
	if strings.TrimSpace(element.Text) == "" {
		return time.Time{}, false, fmt.Errorf("no element found in page")
	}

	data := getParams(s, element.Text)

	y, err := strconv.Atoi(data[pYear])
	if err != nil {
		return time.Time{}, true, err
	}

	data[pYear] = strconv.Itoa(s.YearModifier() + y)

	pt, err := time.Parse(s.GetDateLayout(), s.GetDateValue(data))
	return pt, false, err
}

func NewSpanishEmbassy() IEmbassyData {
	esp := &spanish{
		yearModifier: 0,
		iso:          "es",
		common: &common{
			uri:               "https://www.es.emb-japan.go.jp/itpr_es/00_001125.html",
			pattern:           fmt.Sprintf(`(?P<%s>\d{4})\/(?P<%s>\d{1,2})\/(?P<%s>\d{1,2})`, pYear, pMonth, pDay),
			htmlSearchElement: "div[class=rightalign]",
			dateLayout:        "2006/01/2",
		},
	}
	esp.regex = regexp.MustCompile(esp.pattern)

	return esp
}
