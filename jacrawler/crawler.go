package jacrawler

import (
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/gocolly/colly/v2"
	"regexp"
	"time"
)

const (
	pYear  = `Year`
	pMonth = `Month`
	pDay   = `Day`
)

type (
	common struct {
		uri               string
		pattern           string
		htmlSearchElement string
		dateLayout        string
		regex             *regexp.Regexp
	}

	covidCrawl struct {
		crawler *colly.Collector

		emb IEmbassyData
	}

	ICovidCrawler interface {
		CrawlPage() (*model.Embassy, error)
	}

	IEmbassyData interface {
		GetUri() string
		GetPattern() string
		GetDateLayout() string
		GetHtmlSearchElement() string
		GetRegex() *regexp.Regexp
		GetDateValue(map[string]string) string
		YearModifier() int
		IsDateUpdated(*model.Embassy, ports.IConnection) bool
		UpdateDate(*model.Embassy, ports.IConnection) error
		GetISO() string
		GetUpdatedDateFromText(*colly.HTMLElement) (time.Time, bool, error)
	}
)

func NewCovidCrawler(emb IEmbassyData) ICovidCrawler {
	return &covidCrawl{
		crawler: colly.NewCollector(),
		emb:     emb,
	}
}

func (c *covidCrawl) CrawlPage() (data *model.Embassy, err error) {
	data = new(model.Embassy)
	var errCrawler error

	c.crawler.OnHTML(c.emb.GetHtmlSearchElement(), func(h *colly.HTMLElement) {
		d, isCritical, err := c.emb.GetUpdatedDateFromText(h)
		if err != nil && isCritical {
			errCrawler = err
		} else {
			if !d.IsZero() {
				data.UpdatedDate = d
			}
		}
	})

	if err := c.crawler.Visit(c.emb.GetUri()); err != nil {
		return nil, err
	}

	return data, errCrawler
}
