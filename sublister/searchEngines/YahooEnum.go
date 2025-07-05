package searchEngines

import (
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

type YahooEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (yahoo *YahooEnum) Init() {
	yahoo.BaseURL = "https://search.yahoo.com/search?p={query}&b={page_no}"
	yahoo.EngineName = "Yahoo"
	yahoo.MaxDomains = 10
	yahoo.MaxPages = 10
}

func (yahoo *YahooEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(yahoo.Domain, ".")
	reg1 := "(>[a-z]+\\.)+<b>" + split[0] + "<\\/b>\\.<b>" + split[1] + "<\\/b>" //(>[a-z]+\.)+<b>misis<\/b>\.<b>ru<\/b>
	LinksRegexp1, err := regexp.Compile(reg1)
	if err != nil {
		return nil, err
	}
	LinksList1 := LinksRegexp1.FindAllString(resp, -1)
	LinksList1 = utils.DeleteRepetitions(LinksList1)
	for i, link := range LinksList1 {
		if strings.Contains(link, ">") {
			LinksList1[i] = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(link, "</b>", ""), "<b>", ""), ">", "")
		}
	}
	reg2 := "%2f([a-z])+\\." + split[0] + "\\." + split[1]
	LinksRegexp2, err := regexp.Compile(reg2)
	if err != nil {
		return nil, err
	}
	LinksList2 := LinksRegexp2.FindAllString(resp, -1)
	LinksList2 = utils.DeleteRepetitions(LinksList2)
	for i, link := range LinksList2 {
		if strings.Contains(link, "%2f") {
			LinksList2[i] = strings.ReplaceAll(link, "%2f", "")
		}
	}
	LinksList := slices.Concat(LinksList1, LinksList2)
	LinksList = utils.DeleteRepetitions(LinksList)
	for _, link := range LinksList {
		if link != "" && !slices.Contains(yahoo.Subdomains, link) && link != yahoo.Domain {
			link = strings.Join(strings.Fields(link), " ")
			yahoo.Subdomains = append(yahoo.Subdomains, link)
		}
	}
	return LinksList, nil
}

func (yahoo *YahooEnum) GenerateQuery() string {
	if !slices.Equal(yahoo.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(yahoo.Subdomains, " -")
		query := utils.Format(format, yahoo.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", yahoo.Domain, "{domain}")
		return query
	}
}

func (yahoo *YahooEnum) CheckResponseBlock(resp string) bool {
	if strings.Contains(resp, "Our systems have detected unusual traffic") {
		return false
	}
	return true
}

func (yahoo *YahooEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}
func (yahoo *YahooEnum) GetSubdomains() []string {
	return yahoo.Subdomains
}
func (yahoo *YahooEnum) GetDomain() string {
	return yahoo.Domain
}
func (yahoo *YahooEnum) GetBaseURL() string {
	return yahoo.BaseURL
}
func (yahoo *YahooEnum) GetMaxDomains() int {
	return yahoo.MaxDomains
}
func (yahoo *YahooEnum) GetMaxPages() int {
	return yahoo.MaxPages
}
func (yahoo *YahooEnum) GetEngineName() string {
	return yahoo.EngineName
}
