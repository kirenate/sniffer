package sublister

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sniffer/logger"
	"sniffer/sublister/searchEngines"
	"sniffer/sublister/utils"
	"strings"
	"time"
)

var log = logger.MakeLogger()

type Enumerator interface {
	Init()
	ExtractDomains(resp string) ([]string, error)
	CheckResponseBlock(resp string) bool
	ShouldSleep()
	GenerateQuery() string
	GetSubdomains() []string
	GetDomain() string
	GetBaseURL() string
	GetMaxDomains() int
	GetMaxPages() int
	GetEngineName() string
}

type Config struct {
	BaseURL    string
	MaxDomains int
	MaxPages   int
	EngineName string
}
type EnumeratorContainer struct {
	Config     Config
	Enumerator Enumerator
}

func GetEnumerators() []EnumeratorContainer {
	return []EnumeratorContainer{
		{Enumerator: &searchEngines.CrtEnum{}, Config: Config{BaseURL: "...", MaxDomains: 10}},
	}
}

var client = http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       25 * time.Second,
}

func SendReq(en Enumerator, method string, query string, pageNum int) (resp *http.Response, err error) {
	BaseURL := en.GetBaseURL()
	BaseURL = utils.Format(BaseURL, query, "{query}")
	BaseURL = utils.Format(BaseURL, pageNum, "{page_no}")
	var buf []byte
	req, _ := http.NewRequest(method, BaseURL, bytes.NewBuffer(buf))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err = client.Do(req)
	return resp, err
}

func GetResponseGzipBody(resp *http.Response) (string, error) {
	gzreader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", nil
	}
	body, err := io.ReadAll(gzreader)
	if err != nil {
		return "", nil
	}
	return string(body), nil

}

func CheckIfCookieNeeded(en Enumerator) bool {
	domain := en.GetDomain()
	resp, errs := SendReq(en, "GET", domain, 0)
	if errs != nil {
		fmt.Println(errs)
	}
	cookiesList := resp.Cookies()
	if len(cookiesList) == 0 {
		return false
	}
	return true
}

func CheckMaxSubdomains(en Enumerator, count int) bool {
	MaxDomains := en.GetMaxDomains()
	if count >= MaxDomains {
		return true
	}
	return false
}

func CheckMaxPages(en Enumerator, count int) bool {
	MaxPages := en.GetMaxPages()
	if count >= MaxPages {
		return true
	}
	return false
}
func GetPage(pageNum int) int {
	return pageNum + 10
}

func Enumerate(en Enumerator) ([]string, error) {
	PageNum := 0
	PrevLinks := []string{}
	Retries := 0
	Domain := en.GetDomain()
	for {
		query := en.GenerateQuery()
		count := strings.Count(query, Domain)

		if CheckMaxSubdomains(en, count) {
			PageNum = GetPage(PageNum)
		}

		if CheckMaxPages(en, PageNum) {

			return en.GetSubdomains(), nil
		}
		resp, err := SendReq(en, "GET", en.GetDomain(), PageNum)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			return en.GetSubdomains(), err
		}
		body, err := GetResponseGzipBody(resp)
		if err != nil {
			return en.GetSubdomains(), err
		}
		links, err := en.ExtractDomains(string(body))
		if err != nil {
			log.Error().Stack().Err(err).Msgf("failed to extract domains from %s", en.GetEngineName())
			return en.GetSubdomains(), err
		}
		if slices.Equal(links, PrevLinks) {
			Retries++
			PageNum = GetPage(PageNum)
			if Retries >= 3 {
				return links, nil
			}
		}
		PrevLinks = links
		en.ShouldSleep()
	}
}

func Sublister(domain string) []string {
	//var yahoo searchEngines.YahooEnum
	//yahoo.Init()
	//yahoo.Domain = domain
	//var err error
	//yahoo.Subdomains, err = Enumerate(&yahoo)
	//if err != nil {
	//	log.Error().Stack().Err(err).Send()
	//}
	//return yahoo.GetSubdomains()
	var crt searchEngines.CrtEnum
	crt.Init()
	crt.Domain = domain
	var err error
	crt.Subdomains, err = Enumerate(&crt)
	if err != nil {
		log.Error().Stack().Err(err).Send()
	}
	return crt.Subdomains
}
