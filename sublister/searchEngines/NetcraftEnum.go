package searchEngines

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"regexp"
	"slices"
	"sniffer/sublister/utils"
	"strings"
	"time"
	"unicode/utf8"
)

type NetcraftEnum struct {
	Subdomains []string
	Domain     string
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
	Cookie     http.Cookie
}

func (net *NetcraftEnum) Init() {
	net.BaseURL = "https://searchdns.netcraft.com/?restriction=site+ends+with&host={domain}"
	net.EngineName = "Netcraft"
	net.MaxDomains = 0
	net.MaxPages = 0
}

func (net *NetcraftEnum) ExtractDomains(resp string) ([]string, error) {
	split := strings.Split(net.Domain, ".")
	reg := "([a-z]+\\.)+(" + split[0] + "\\." + split[1] + ")+"
	LinksRegexp, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	LinksList := LinksRegexp.FindAllString(resp, -1)
	LinksList = utils.DeleteRepetitions(LinksList)
	for _, link := range LinksList {
		if link != "" && !slices.Contains(net.Subdomains, link) && link != net.Domain {
			link = strings.Join(strings.Fields(link), " ")
			net.Subdomains = append(net.Subdomains, link)
		}
	}
	return LinksList, nil
}

func (net *NetcraftEnum) GenerateQuery() string {
	if !slices.Equal(net.Subdomains, []string{}) {
		format := "site:{domain} -www.{domain} -{found}"
		found := strings.Join(net.Subdomains, " -")
		query := utils.Format(format, net.Domain, "{domain}")
		query = utils.Format(query, found, "{found}")
		return query
	} else {
		query := utils.Format("site:{domain}", net.Domain, "{domain}")
		return query
	}
}

func (net *NetcraftEnum) CheckResponseBlock(resp string) bool {
	return true
}

func (net *NetcraftEnum) ShouldSleep() {
	time.Sleep(5 * time.Millisecond)
}

func (net *NetcraftEnum) CreateCookies(cookieList []*http.Cookie) {
	for _, respCookie := range cookieList {
		valid := utf8.ValidString(respCookie.Value)
		if valid {
			hasher := sha1.New()
			hasher.Write([]byte(respCookie.Value))
			sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

			reqCookie := http.Cookie{
				Name:   "netcraft_js_verification_response",
				Value:  sha,
				Domain: "searchdns.netcraft.com",
				Path:   "/",
			}
			net.Cookie = reqCookie
		}
	}
}

func (net *NetcraftEnum) SendReqWithCookie(query string, pageNum int) (*http.Response, error) {
	BaseURL := net.GetBaseURL()
	BaseURL = utils.Format(BaseURL, query, "{query}")
	BaseURL = utils.Format(BaseURL, pageNum, "{page_no}")
	var buf []byte
	req, _ := http.NewRequest("GET", BaseURL, bytes.NewBuffer(buf))
	req.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	req.Header.Add("Cookie", net.Cookie.Value)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func (net *NetcraftEnum) GetSubdomains() []string {
	return net.Subdomains
}
func (net *NetcraftEnum) GetDomain() string {
	return net.Domain
}
func (net *NetcraftEnum) GetBaseURL() string {
	return net.BaseURL
}
func (net *NetcraftEnum) GetMaxDomains() int {
	return net.MaxDomains
}
func (net *NetcraftEnum) GetMaxPages() int {
	return net.MaxPages
}
func (net *NetcraftEnum) GetEngineName() string {
	return net.EngineName
}
func (net *NetcraftEnum) GetCookie() http.Cookie {
	return net.Cookie
}
