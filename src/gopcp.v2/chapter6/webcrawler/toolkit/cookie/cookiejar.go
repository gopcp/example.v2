package cookie

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// NewCookiejar 用于创建http.CookieJar类型的实例。
func NewCookiejar() http.CookieJar {
	options := &cookiejar.Options{PublicSuffixList: &myPublicSuffixList{}}
	cj, _ := cookiejar.New(options)
	return cj
}

// myPublicSuffixList 代表cookiejar.PublicSuffixList接口的实现类型。
type myPublicSuffixList struct{}

func (psl *myPublicSuffixList) PublicSuffix(domain string) string {
	suffix, _ := publicsuffix.PublicSuffix(domain)
	return suffix
}

func (psl *myPublicSuffixList) String() string {
	return "Web crawler - public suffix list (rev 1.0) power by \"golang.org/x/net/publicsuffix\""
}
