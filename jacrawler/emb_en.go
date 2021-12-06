package jacrawler

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"log"
	"regexp"
)

type (
	english struct {
		*common
		iso string
	}
)

func NewEnglishEmbassy() IEmbassyData {
	eng := &english{
		common: &common{
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

func (e *english) IsDateUpdated(data *model.Embassy, db ports.IConnection) bool {
	c := new(model.Embassy)
	c.EmbassyISO = e.iso
	if err := db.FetchLatestDateFromEmbassy(c); err != nil {
		log.Printf("There was an error fetching data from db: %s\n", err)
		return true
	}

	return c.UpdatedDate.After(data.UpdatedDate) || c.UpdatedDate.Equal(data.UpdatedDate)
}

func (e *english) UpdateDate(data *model.Embassy, db ports.IConnection) error {
	data.EmbassyISO = e.iso
	return db.Save(data)
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
