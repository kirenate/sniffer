package searchEngines

import (
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

type VirustotalEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (virustotal *VirustotalEnum) Init() {
	virustotal.BaseURL = "https://www.virustotal.com/gui/domain/{domain}/relations"
	virustotal.EngineName = "Virustotal"
	virustotal.MaxDomains = 10
	virustotal.MaxPages = 10
}

func (virustotal *VirustotalEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(virustotal.Domain, ".")
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
		if link != "" && !slices.Contains(virustotal.Subdomains, link) && link != virustotal.Domain {
			link = strings.Join(strings.Fields(link), " ")
			virustotal.Subdomains = append(virustotal.Subdomains, link)
		}
	}
	return LinksList, nil
}

func (virustotal *VirustotalEnum) GenerateQuery() string {
	if !slices.Equal(virustotal.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(virustotal.Subdomains, " -")
		query := utils.Format(format, virustotal.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", virustotal.Domain, "{domain}")
		return query
	}
}

func (virustotal *VirustotalEnum) CheckResponseBlock(resp string) bool {
	if strings.Contains(resp, "Our systems have detected unusual traffic") {
		return false
	}
	return true
}

func (virustotal *VirustotalEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}
func (virustotal *VirustotalEnum) GetSubdomains() []string {
	return virustotal.Subdomains
}
func (virustotal *VirustotalEnum) GetDomain() string {
	return virustotal.Domain
}
func (virustotal *VirustotalEnum) GetBaseURL() string {
	return virustotal.BaseURL
}
func (virustotal *VirustotalEnum) GetMaxDomains() int {
	return virustotal.MaxDomains
}
func (virustotal *VirustotalEnum) GetMaxPages() int {
	return virustotal.MaxPages
}
func (virustotal *VirustotalEnum) GetEngineName() string {
	return virustotal.EngineName
}
