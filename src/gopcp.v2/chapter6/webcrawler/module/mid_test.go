package module

import (
	"net"
	"strconv"
	"testing"
)

var legalMIDs = []MID{}

func init() {
	for _, mt := range legalTypes {
		for mip := range legalIPMap {
			addr, _ := NewAddr("http", mip, 8080)
			mid, _ := GenMID(mt, DefaultSNGen.Get(), addr)
			legalMIDs = append(legalMIDs, mid)
		}
	}
}

var illegalMIDs = []MID{
	MID("D"),
	MID("DZ"),
	MID("D1|"),
	MID("D1|127.0.0.1:-1"),
	MID("D1|127.0.0.1:"),
	MID("D1|127.0.0.1"),
	MID("D1|127.0.0."),
	MID("D1|127"),
	MID("D1|127.0.0.0.1:8080"),
	MID("DZ|127.0.0.1:8080"),
	MID("A"),
	MID("AZ"),
	MID("A1|"),
	MID("A1|127.0.0.1:-1"),
	MID("A1|127.0.0.1:"),
	MID("A1|127.0.0.1"),
	MID("A1|127.0.0."),
	MID("A1|127"),
	MID("A1|127.0.0.0.1:8080"),
	MID("AZ|127.0.0.1:8080"),
	MID("P"),
	MID("PZ"),
	MID("P1|"),
	MID("P1|127.0.0.1:-1"),
	MID("P1|127.0.0.1:"),
	MID("P1|127.0.0.1"),
	MID("P1|127.0.0."),
	MID("P1|127"),
	MID("P1|127.0.0.0.1:8080"),
	MID("PZ|127.0.0.1:8080"),
	MID("M1|127.0.0.1:8080"),
}

func TestMIDGenAndSplit(t *testing.T) {
	addr, _ := NewAddr("http", "127.0.0.1", 8080)
	addrs := []net.Addr{nil, addr}
	for _, addr := range addrs {
		for _, mt := range legalTypes {
			expectedSN := DefaultSNGen.Get()
			mid, err := GenMID(mt, expectedSN, addr)
			if err != nil {
				t.Fatalf("An error occurs when generating module ID: %s (type: %s, sn: %d, addr: %s)",
					err, mt, expectedSN, addr)
			}
			expectedLetter := legalTypeLetterMap[mt]
			var expectedAddrStr string
			if addr != nil {
				expectedAddrStr = addr.String()
			}
			parts, err := SplitMID(mid)
			if err != nil {
				t.Fatalf("An error occurs when splitting MID %q: %s", mid, err)
			}
			letter, snStr, addrStr := parts[0], parts[1], parts[2]
			if letter != expectedLetter {
				t.Fatalf("Inconsistent type letter in MID: expected: %s, actual: %s",
					expectedLetter, letter)
			}
			sn, err := strconv.ParseUint(snStr, 10, 64)
			if err != nil {
				t.Fatalf("An error occurs when parsing SN: %s (snStr: %s)", err, snStr)
			}
			if sn != expectedSN {
				t.Fatalf("Inconsistent SN in MID: expected: %d, actual: %d",
					expectedSN, sn)
			}
			if addrStr != expectedAddrStr {
				t.Fatalf("Inconsistent address string in MID: expected: %s, actual: %s",
					expectedAddrStr, addrStr)
			}
		}
	}
	for _, addr := range addrs {
		for _, mt := range illegalTypes {
			mid, err := GenMID(mt, DefaultSNGen.Get(), addr)
			if err == nil {
				t.Fatalf("It still can generate module ID with illegal type %q!",
					mt)
			}
			if string(mid) != "" {
				t.Fatalf("It still can generate module ID %q with illegal type %q!",
					mid, mt)
			}
		}
	}
	for _, illegalMID := range illegalMIDs {
		if _, err := SplitMID(illegalMID); err == nil {
			t.Fatalf("It still can split illegal module ID %q", illegalMID)
		}
	}
}

func TestMIDLegal(t *testing.T) {
	var addr net.Addr
	var mid MID
	var err error
	for _, mt := range legalTypes {
		for mip := range legalIPMap {
			sn := DefaultSNGen.Get()
			addr, err = NewAddr("http", mip, 8080)
			if err == nil {
				mid, err = GenMID(mt, sn, addr)
			}
			if err != nil {
				t.Fatalf("An error occurs when judging the legality for MID: %s (type: %s, sn: %d, addr: %s)",
					err, mt, sn, addr)
			}
			if !LegalMID(mid) {
				t.Fatalf("The generated MID %q is legal, but do not be detected!", mid)
			}
		}
	}
	for _, illegalMID := range illegalMIDs {
		if LegalMID(illegalMID) {
			t.Fatalf("The generated MID %q is illegal, but do not be detected!", illegalMID)
		}
	}
}
