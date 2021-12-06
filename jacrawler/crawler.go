package jacrawler

import (
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
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
	}
)

func NewCovidCrawler(emb IEmbassyData) ICovidCrawler {
	return &covidCrawl{
		crawler: colly.NewCollector(),
		emb:     emb,
	}
}

func (c *covidCrawl) getParams(text string) (paramsMap map[string]string) {
	match := c.emb.GetRegex().FindStringSubmatch(text)

	paramsMap = make(map[string]string)
	for i, name := range c.emb.GetRegex().SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func (c *covidCrawl) getUpdatedDateFromText(text string) (time.Time, error) {
	if strings.TrimSpace(text) == "" {
		return time.Time{}, nil
	}

	data := c.getParams(text)

	y, err := strconv.Atoi(data[pYear])
	if err != nil {
		return time.Time{}, err
	}

	data[pYear] = strconv.Itoa(c.emb.YearModifier() + y)

	pt, err := time.Parse(c.emb.GetDateLayout(), c.emb.GetDateValue(data))
	return pt, err
}

func (c *covidCrawl) CrawlPage() (data *model.Embassy, err error) {
	data = new(model.Embassy)
	var errCrawler error

	c.crawler.OnHTML(c.emb.GetHtmlSearchElement(), func(h *colly.HTMLElement) {
		d, err := c.getUpdatedDateFromText(h.Text)
		if err != nil {
			errCrawler = err
		} else {
			data.UpdatedDate = d
		}
	})

	if err := c.crawler.Visit(c.emb.GetUri()); err != nil {
		return nil, err
	}

	return data, errCrawler
}
