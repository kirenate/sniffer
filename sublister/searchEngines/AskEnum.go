package searchEngines

import (
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

type AskEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (ask *AskEnum) Init() {
	ask.BaseURL = "https://www.ask.com/web?q={query}&page={page_no}&qid=8D6EE6BF52E0C04527E51F64F22C4534&o=0&l=dir&qsrc=998&qo=pagination"
	ask.EngineName = "Ask"
	ask.MaxDomains = 11
	ask.MaxPages = 10
}

func (ask *AskEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(ask.Domain, ".")
	reg := "([a-z]+\\.)+(" + split[0] + "\\." + split[1] + ")+"
	LinksRegexp, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	LinksList := LinksRegexp.FindAllString(resp, -1)
	LinksList = utils.DeleteRepetitions(LinksList)
	for _, link := range LinksList {
		if link != "" && !slices.Contains(ask.Subdomains, link) && link != ask.Domain {
			link = strings.Join(strings.Fields(link), " ")
			ask.Subdomains = append(ask.Subdomains, link)
		}
	}
	return LinksList, nil
}

func (ask *AskEnum) CheckResponseBlock(resp string) bool {
	return true
}

func (ask *AskEnum) ShouldSleep() {
	time.Sleep(1 * time.Millisecond)
}
func (ask *AskEnum) GenerateQuery() string {
	if !slices.Equal(ask.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(ask.Subdomains, " -")
		query := utils.Format(format, ask.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", ask.Domain, "{domain}")
		return query
	}
}

func (ask *AskEnum) GetSubdomains() []string {
	return ask.Subdomains
}

func (ask *AskEnum) GetDomain() string {
	return ask.Domain
}

func (ask *AskEnum) GetBaseURL() string {
	return ask.BaseURL
}

func (ask *AskEnum) GetMaxDomains() int {
	return ask.MaxDomains
}

func (ask *AskEnum) GetMaxPages() int {
	return ask.MaxPages
}

func (ask *AskEnum) GetEngineName() string {
	return ask.EngineName
}
