package scheduler

import "testing"

func TestGetPrimaryDomain(t *testing.T) {
	host := "127.0.0.1"
	pd, err := getPrimaryDomain(host)
	if err != nil {
		t.Fatalf("An error occurs when getting primary domain: %s (host: %s)",
			err, host)
	}
	if pd != host {
		t.Fatalf("Inconsistent primary domain: expected: %s, actual: %s",
			host, pd)
	}
	host = "cn.bing.com"
	pd, err = getPrimaryDomain(host)
	if err != nil {
		t.Fatalf("An error occurs when getting primary domain: %s (host: %s)",
			err, host)
	}
	expectedPD := "bing.com"
	if pd != expectedPD {
		t.Fatalf("Inconsistent primary domain: expected: %s, actual: %s",
			expectedPD, pd)
	}
	_, err = getPrimaryDomain("")
	if err == nil {
		t.Fatal("It still can get primary domain for a empty host!")
	}
	host = "123.abc"
	_, err = getPrimaryDomain(host)
	if err == nil {
		t.Fatalf("It still can get primary domain for a unrecognized host %q!", host)
	}
}
