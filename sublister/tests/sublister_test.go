package tests

import (
	"fmt"
	"github.com/coreos/etcd/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"slices"
	"sniffer/sublister"
	"sniffer/sublister/searchEngines"
	"sniffer/sublister/utils"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	s := utils.Format("https://crt.sh/?q=%25.{domain}", "misis.ru", "{domain}")
	s = utils.Format(s, 1, "{page_no}")

	res := "https://crt.sh/?q=%25.misis.ru"
	testutil.AssertEqual(t, res, s, "test successful")
	fmt.Println(s)
}

func TestSendReq(t *testing.T) {
	en := searchEngines.CrtEnum{}
	en.Domain = "vk.com"
	en.Init()
	resp, err := sublister.SendReq(&en, "GET", "vk.com", 1)
	require.NoError(t, err)
	body, err := sublister.GetResponseGzipBody(resp)
	require.NoError(t, err)
	if resp.StatusCode == 500 || strings.Contains(string(body), "error") {
		t.Fail()
	}
	fmt.Println(resp)
	fmt.Println(string(body))
}

func TestEnumerate(t *testing.T) {
	en := searchEngines.CrtEnum{}
	en.Init()
	en.Domain = "misis.ru"
	fmt.Println("enums struct: ", en)
	subdomains, err := sublister.Enumerate(&en)
	require.NoError(t, err)
	for _, v := range subdomains {
		fmt.Println("subdomain found: ", v)
	}
}

func TestGenerateQuery(t *testing.T) {
	en := searchEngines.YahooEnum{}
	en.Init()
	en.Domain = "misis.ru"
	en.Subdomains = append(en.Subdomains, "en.misis.ru")
	query := en.GenerateQuery()
	fmt.Println(query)
}

func TestExtractDomains(t *testing.T) {
	en := searchEngines.CrtEnum{}
	en.Init()
	en.Domain = "mts.ru"
	resp, errs := sublister.SendReq(&en, "GET", en.Domain, 1)
	if errs != nil {
		fmt.Println(errs)
	}
	body, err := sublister.GetResponseGzipBody(resp)
	require.NoError(t, err)
	//fmt.Println(body)
	domainList, err := en.ExtractDomains(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(domainList)
	assert.NotEmpty(t, domainList)
}

func TestDeleteRepetitions(t *testing.T) {
	s := []string{"a", "a", "a", "a", "a", "b", "b", "b", "c", "c"}
	res := DeleteRepetitions(s)
	assert.Equal(t, []string{"a", "b", "c"}, res)
}

func DeleteRepetitions(s []string) []string {
	slices.Sort(s)
	var newSlice []string
	for _, v := range s {
		if slices.Contains(newSlice, v) {
			continue
		}
		newSlice = append(newSlice, v)
	}
	return newSlice
}

func TestCreateCooikes(t *testing.T) {
	net := searchEngines.NetcraftEnum{
		Subdomains: nil,
		Domain:     "misis.ru",
		BaseURL:    "https://searchdns.netcraft.com/?restriction=site+ends+with&host={domain}",
		MaxDomains: 10,
		MaxPages:   10,
		EngineName: "Netcraft",
		Cookie:     http.Cookie{},
	}
	resp, _ := sublister.SendReq(&net, "GET", "misis.ru", 1)
	cookies := resp.Cookies()
	net.CreateCookies(cookies)
	fmt.Println(net.Cookie.Value)
	assert.NotEmpty(t, net.Cookie.Value)
}

func TestSendReqWithCookie(t *testing.T) {
	net := searchEngines.NetcraftEnum{
		Subdomains: nil,
		Domain:     "misis.ru",
		BaseURL:    "https://searchdns.netcraft.com/?restriction=site+ends+with&host={domain}",
		MaxDomains: 10,
		MaxPages:   10,
		EngineName: "Netcraft",
		Cookie:     http.Cookie{},
	}
	if sublister.CheckIfCookieNeeded(&net) {
		resp, _ := sublister.SendReq(&net, "GET", "misis.ru", 1)
		cookies := resp.Cookies()
		fmt.Println(cookies)
		net.CreateCookies(cookies)
		resp, _ = net.SendReqWithCookie(net.Domain, 1)
		fmt.Println(resp)
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), "does not support") {
			t.Fail()
		}
		fmt.Println(string(body))
	}
}

func TestDNSdumpsterSendReq(t *testing.T) {
	dump := searchEngines.DNSdumpsterEnum{}
	dump.Init()
	dump.Domain = "misis.ru"
	resp, _ := sublister.SendReq(&dump, "POST", dump.Domain, 1)
	fmt.Println(resp)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if strings.Contains(string(body), "Sorry, you have been blocked") {
		t.Fail()
	}
}
