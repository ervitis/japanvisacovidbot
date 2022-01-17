package jacrawler

import (
	"context"
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
	english struct {
		*common
		iso string
	}
)

func NewEnglishEmbassy(db ports.IConnection) IEmbassyData {
	eng := &english{
		common: &common{
			db:                db,
			uri:               "https://www.mofa.go.jp/ca/cp/page22e_000925.html",
			htmlSearchElement: "div[class=rightalign]",
			dateLayout:        "January 02, 2006",
			pattern:           fmt.Sprintf(`(?P<%s>\w+)\s(?P<%s>\d{1,2}),\s(?P<%s>\d{4})`, pMonth, pDay, pYear),
		},
		iso: "en",
	}

	eng.regex = regexp.MustCompile(eng.pattern)
	return eng
}

func (e *english) IsDateUpdated(ctx context.Context, data *model.Embassy) bool {
	c := new(model.Embassy)
	c.EmbassyISO = e.iso
	if err := e.db.FetchLatestDateFromEmbassy(ctx, c); err != nil {
		log.Printf("There was an error fetching data from db: %s\n", err)
		return true
	}

	return c.UpdatedDate.After(data.UpdatedDate) || c.UpdatedDate.Equal(data.UpdatedDate)
}

func (e *english) UpdateDate(ctx context.Context, data *model.Embassy) error {
	data.EmbassyISO = e.iso
	return e.db.Save(ctx, data)
}

func (e *english) GetUri() string {
	return e.uri
}

func (e *english) GetPattern() string {
	return e.pattern
}

func (e *english) GetDateLayout() string {
	return e.dateLayout
}

func (e *english) GetHtmlSearchElement() string {
	return e.htmlSearchElement
}

func (e *english) GetRegex() *regexp.Regexp {
	return e.regex
}

func (e *english) YearModifier() int {
	return 0
}

func (e *english) GetDateValue(data map[string]string) string {
	return fmt.Sprintf(`%s %s, %s`, data[pMonth], data[pDay], data[pYear])
}

func (e *english) GetISO() string {
	return e.iso
}

func (e *english) GetUpdatedDateFromText(el *colly.HTMLElement) (time.Time, bool, error) {
	if strings.TrimSpace(el.Text) == "" {
		return time.Time{}, false, nil
	}

	data := getParams(e, el.Text)

	y, err := strconv.Atoi(data[pYear])
	if err != nil {
		return time.Time{}, true, err
	}

	data[pYear] = strconv.Itoa(e.YearModifier() + y)

	pt, err := time.Parse(e.GetDateLayout(), e.GetDateValue(data))
	return pt, false, err
}
