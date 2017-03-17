package module

import (
	"strconv"
	"testing"
)

var legalNetworkMap = map[string]struct{}{
	"http":  struct{}{},
	"https": struct{}{},
}

var legalIPMap = map[string]struct{}{
	"192.0.2.1":    struct{}{},
	"2001:db8::68": struct{}{},
}

func TestAddrNew(t *testing.T) {
	var network string
	var ip string
	port := uint64(8080)
	for network = range legalNetworkMap {
		for ip = range legalIPMap {
			addr, err := NewAddr(network, ip, port)
			if err != nil {
				t.Fatalf("An error occurs when creating an address: %s (network: %s, ip: %s, port: %d)",
					err, network, ip, port)
			}
			if addr == nil {
				t.Fatal("Couldn't create address!")
			}
			if addr.Network() != network {
				t.Fatalf("Inconsistent network for address: expected: %s, actual: %s",
					network, addr.Network())
			}
			expectedIPPort := ip + ":" + strconv.FormatUint(port, 10)
			if addr.String() != expectedIPPort {
				t.Fatalf("Inconsistent IP and/or port for address: expected: %s, actual: %s",
					expectedIPPort, addr.String())
			}
		}
	}
	network = "tcp"
	_, err := NewAddr(network, ip, port)
	if err == nil {
		t.Fatalf("No error when create a buffer with illegal network %q!", network)
	}
	network = "http"
	ip = "127.0.0.0.1"
	_, err = NewAddr(network, ip, port)
	if err == nil {
		t.Fatalf("No error when create a buffer with illegal network %q!", network)
	}
}
