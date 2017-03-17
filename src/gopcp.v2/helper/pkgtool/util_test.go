package pkgtool

import (
	"runtime/debug"
	"testing"
)

func TestAppendIfAbsent(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	var s []string
	e1 := "a"
	s = appendIfAbsent(s, e1)
	expLen := 1
	actLen := len(s)
	if actLen != expLen {
		t.Errorf("Error: The length of slice should be %d but %d.\n", expLen, actLen)
	}
	if !contains(s, e1) {
		t.Errorf("Error: The slice should contains '%s' but not.\n", e1)
	}
	s = appendIfAbsent(s, e1)
	expLen = 1
	actLen = len(s)
	if actLen != expLen {
		t.Errorf("Error: The length of slice should be %d but %d.\n", expLen, actLen)
	}
	e2 := "b"
	s = appendIfAbsent(s, e2)
	expLen = 2
	actLen = len(s)
	if actLen != expLen {
		t.Errorf("Error: The length of slice should be %d but %d.\n", expLen, actLen)
	}
	expS := []string{"a", "b"}
	if !sliceEquels(s, expS) {
		t.Errorf("Error: The slice should be %v but %v.\n", expS, s)
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func sliceEquels(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if v1 != v2 {
			return false
		}
	}
	return true
}
