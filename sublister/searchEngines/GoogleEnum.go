package searchEngines

import (
	"net/url"
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

type GoogleEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}

func (google *GoogleEnum) Init() {
	google.BaseURL = "https://www.google.com/search?q=m{query}&amp;btnG=Search&amp;hl=en-US&amp;biw=&amp;bih=&amp;gbv=1&amp;start={page_no}&amp;filter=0"
	google.EngineName = "Google"
	google.MaxDomains = 10
	google.MaxPages = 200
}

func (google *GoogleEnum) ExtractDomains(resp string) ([]string, error) {
	var LinksList []string
	LinkRegexp, err := regexp.Compile("'<cite.*?>(.*?)<\\/cite>'")
	if err != nil {
		return LinksList, err
	}
	LinksList = LinkRegexp.FindAllString(resp, -1)
	for _, link := range LinksList {
		re, err := regexp.Compile("'<span.*>'")
		if err != nil {
			return LinksList, err
		}
		link = re.ReplaceAllString(link, "")
		if !strings.HasPrefix(link, "http") {
			link = "http://" + link
		}
		parsedURL, err := url.ParseRequestURI(link)
		if err != nil {
			return LinksList, err
		}
		subdomain := parsedURL.Host
		if subdomain != "" && !slices.Contains(google.Subdomains, subdomain) && subdomain != google.Domain {
			subdomain = strings.Join(strings.Fields(subdomain), "")
			google.Subdomains = append(google.Subdomains, subdomain)
		}
	}
	return LinksList, nil
}

func (google *GoogleEnum) GenerateQuery() string {
	if !slices.Equal(google.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(google.Subdomains[:google.MaxDomains-1], " -")
		query := utils.Format(format, google.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", google.Domain, "{domain}")
		return query
	}
}

func (google *GoogleEnum) CheckResponseBlock(resp string) bool {
	if strings.Contains(resp, "Our systems have detected unusual traffic") {
		return false
	}
	return true
}

func (google *GoogleEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}
func (google *GoogleEnum) GetSubdomains() []string {
	return google.Subdomains
}
func (google *GoogleEnum) GetDomain() string {
	return google.Domain
}
func (google *GoogleEnum) GetBaseURL() string {
	return google.BaseURL
}
func (google *GoogleEnum) GetMaxDomains() int {
	return google.MaxDomains
}
func (google *GoogleEnum) GetMaxPages() int {
	return google.MaxPages
}
func (google *GoogleEnum) GetEngineName() string {
	return google.EngineName
}
