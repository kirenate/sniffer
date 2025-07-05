package searchEngines

import (
	"regexp"
	"slices"
	"sniffer/logger"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

var log = logger.MakeLogger()

type DNSdumpsterEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (DNSdumpster *DNSdumpsterEnum) Init() {
	DNSdumpster.BaseURL = "https://dnsdumpster.com/"
	DNSdumpster.EngineName = "DNSdumpster"
	DNSdumpster.MaxDomains = 10
	DNSdumpster.MaxPages = 10
}

func (DNSdumpster *DNSdumpsterEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(DNSdumpster.Domain, ".")
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
		if link != "" && !slices.Contains(DNSdumpster.Subdomains, link) && link != DNSdumpster.Domain {
			link = strings.Join(strings.Fields(link), " ")
			DNSdumpster.Subdomains = append(DNSdumpster.Subdomains, link)
		}
	}
	return LinksList, nil
}

func (DNSdumpster *DNSdumpsterEnum) GenerateQuery() string {
	if !slices.Equal(DNSdumpster.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(DNSdumpster.Subdomains, " -")
		query := utils.Format(format, DNSdumpster.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", DNSdumpster.Domain, "{domain}")
		return query
	}
}

func (DNSdumpster *DNSdumpsterEnum) CheckResponseBlock(resp string) bool {
	if strings.Contains(resp, "Our systems have detected unusual traffic") {
		return false
	}
	return true
}

func (DNSdumpster *DNSdumpsterEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}
func (DNSdumpster *DNSdumpsterEnum) GetSubdomains() []string {
	return DNSdumpster.Subdomains
}
func (DNSdumpster *DNSdumpsterEnum) GetDomain() string {
	return DNSdumpster.Domain
}
func (DNSdumpster *DNSdumpsterEnum) GetBaseURL() string {
	return DNSdumpster.BaseURL
}
func (DNSdumpster *DNSdumpsterEnum) GetMaxDomains() int {
	return DNSdumpster.MaxDomains
}
func (DNSdumpster *DNSdumpsterEnum) GetMaxPages() int {
	return DNSdumpster.MaxPages
}
func (DNSdumpster *DNSdumpsterEnum) GetEngineName() string {
	return DNSdumpster.EngineName
}
