package cookie

import "testing"

func TestCookiejar(t *testing.T) {
	cookiejar := NewCookiejar()
	if cookiejar == nil {
		t.Fatal("Couldn't create cookiejar!")
	}
}

func TestMyPublicSuffixList(t *testing.T) {
	domains := []string{
		"golang.org",
		"cn.bing.com",
		"zhihu.sogou.com",
		"www.beijing.gov.cn",
	}
	expectedPSs := []string{
		"org",
		"com",
		"com",
		"gov.cn",
	}
	psl := &myPublicSuffixList{}
	domainsLen := len(domains)
	for i := 0; i < domainsLen; i++ {
		domain := domains[i]
		expectedPS := expectedPSs[i]
		suffix := psl.PublicSuffix(domain)
		if suffix != expectedPS {
			t.Fatalf("Inconsistent publice suffix for domain %q: expected: %s, actual: %s",
				domain, expectedPS, suffix)
		}
	}
	expectedString := "Web crawler - public suffix list (rev 1.0) power by \"golang.org/x/net/publicsuffix\""
	actualString := psl.String()
	if actualString != expectedString {
		t.Fatalf("Inconsistent string for myPublicSuffixList: expected: %s, actual: %s",
			expectedString, actualString)
	}
}
