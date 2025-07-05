package searchEngines

import (
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

type CrtEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (crt *CrtEnum) Init() {
	crt.BaseURL = "https://crt.sh/?q=%25.{query}"
	crt.EngineName = "SSL Certificates"
	crt.MaxDomains = 10
	crt.MaxPages = 10
}

func (crt *CrtEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(crt.Domain, ".")
	reg := "([a-z0-9]+\\.)+" + split[0] + "\\." + split[1] //([a-z0-9]+\.)+vk\.com
	if len(split) > 2 {
		for i := 2; i < len(split); i++ {
			reg += "\\." + split[i]
		}
	}
	linksRegexp, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	linksList := linksRegexp.FindAllString(resp, -1)
	linksList = utils.DeleteRepetitions(linksList)
	for _, link := range linksList {
		if link != "" && !slices.Contains(crt.Subdomains, link) && link != crt.Domain {
			link = strings.Join(strings.Fields(link), " ")
			crt.Subdomains = append(crt.Subdomains, link)
		}
	}
	return crt.Subdomains, nil
}

func (crt *CrtEnum) GenerateQuery() string {
	if !slices.Equal(crt.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(crt.Subdomains, " -")
		query := utils.Format(format, crt.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", crt.Domain, "{domain}")
		return query
	}
}

func (crt *CrtEnum) CheckResponseBlock(resp string) bool {
	if strings.Contains(resp, "Our systems have detected unusual traffic") {
		return false
	}
	return true
}

func (crt *CrtEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}
func (crt *CrtEnum) GetSubdomains() []string {
	return crt.Subdomains
}
func (crt *CrtEnum) GetDomain() string {
	return crt.Domain
}
func (crt *CrtEnum) GetBaseURL() string {
	return crt.BaseURL
}
func (crt *CrtEnum) GetMaxDomains() int {
	return crt.MaxDomains
}
func (crt *CrtEnum) GetMaxPages() int {
	return crt.MaxPages
}
func (crt *CrtEnum) GetEngineName() string {
	return crt.EngineName
}
