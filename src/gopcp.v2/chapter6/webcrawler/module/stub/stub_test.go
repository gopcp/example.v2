package stub

import (
	"testing"

	"gopcp.v2/chapter6/webcrawler/module"
)

// addrStr 代表测试用的网络地址。
var addrStr = "127.0.0.1:8080"

// mid 代表测试用的MID。
var mid = module.MID("D1|" + addrStr)

func TestNew(t *testing.T) {
	mi, err := NewModuleInternal(mid, module.CalculateScoreSimple)
	if err != nil {
		t.Fatalf("An error occurs when creating an internal module: %s (mid: %s)",
			err, mid)
	}
	if mi == nil {
		t.Fatal("Couldn't create internal module!")
	}
	if mi.ID() != mid {
		t.Fatalf("Inconsistent MID for internal module: expected: %s, actual: %s",
			mid, mi.ID())
	}
	if mi.Addr() != addrStr {
		t.Fatalf("Inconsistent addr for internal module: expected: %s, actual: %s",
			addrStr, mi.Addr())
	}
	if mi.Score() != 0 {
		t.Fatalf("Inconsistent score for internal module: expected: %d, actual: %d",
			0, mi.Score())
	}
	if mi.ScoreCalculator() == nil {
		t.Fatalf("Inconsistent score calculator for internal module: expected: %p (%T), actual: %p (%T)",
			module.CalculateScoreSimple, module.CalculateScoreSimple,
			mi.ScoreCalculator(), mi.ScoreCalculator())
	}
	if mi.CalledCount() != 0 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			0, mi.CalledCount())
	}
	if mi.AcceptedCount() != 0 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			0, mi.AcceptedCount())
	}
	if mi.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, mi.CompletedCount())
	}
	if mi.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, mi.HandlingNumber())
	}
	illegalMID := module.MID("D127.0.0.1")
	mi, err = NewModuleInternal(illegalMID, nil)
	if err == nil {
		t.Fatalf("No error when create an internal module with illegal MID %q!", illegalMID)
	}
}

func TestScore(t *testing.T) {
	number := uint64(10)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.SetScore(i)
		score := mi.Score()
		if score != i {
			t.Fatalf("Inconsistent score for internal module: expected: %d, actual: %d",
				i, score)
		}
	}
}

func TestCalledCount(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.IncrCalledCount()
		if mi.CalledCount() != i {
			t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
				i, mi.CalledCount())
		}
	}
}

func TestAcceptedCount(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.IncrAcceptedCount()
		if mi.AcceptedCount() != i {
			t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
				i, mi.AcceptedCount())
		}
	}
}

func TestCompletedCount(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.IncrCompletedCount()
		if mi.CompletedCount() != i {
			t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
				i, mi.CompletedCount())
		}
	}
}

func TestHandlingNumber(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.IncrHandlingNumber()
		mi.IncrHandlingNumber()
		mi.IncrHandlingNumber()
		if mi.HandlingNumber() != i+2 {
			t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
				i+2, mi.HandlingNumber())
		}
		mi.DecrHandlingNumber()
		mi.DecrHandlingNumber()
		if mi.HandlingNumber() != i {
			t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
				i, mi.HandlingNumber())
		}
	}
}

func TestClearAndCounts(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	mod := uint64(17)
	for i := uint64(1); i < number; i++ {
		mi.IncrCalledCount()
		mi.IncrAcceptedCount()
		mi.IncrCompletedCount()
		mi.IncrHandlingNumber()
		if i%mod == 0 {
			mi.Clear()
		}
		counts := mi.Counts()
		if counts.CalledCount != mi.CalledCount() {
			t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
				mi.CalledCount(), counts.CalledCount)
		}
		if counts.AcceptedCount != mi.AcceptedCount() {
			t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
				mi.AcceptedCount(), counts.AcceptedCount)
		}
		if counts.CompletedCount != mi.CompletedCount() {
			t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
				mi.CompletedCount(), counts.CompletedCount)
		}
		if counts.HandlingNumber != mi.HandlingNumber() {
			t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
				mi.HandlingNumber(), counts.HandlingNumber)
		}
		expectedCount := i % 17
		if counts.CalledCount != expectedCount {
			t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
				expectedCount, counts.CalledCount)
		}
		if counts.AcceptedCount != expectedCount {
			t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
				expectedCount, counts.AcceptedCount)
		}
		if counts.CompletedCount != expectedCount {
			t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
				expectedCount, counts.CompletedCount)
		}
		if counts.HandlingNumber != expectedCount {
			t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
				expectedCount, counts.HandlingNumber)
		}
	}
}

func TestSummary(t *testing.T) {
	number := uint64(10000)
	mi, _ := NewModuleInternal(mid, nil)
	for i := uint64(1); i < number; i++ {
		mi.IncrCalledCount()
		mi.IncrAcceptedCount()
		mi.IncrCompletedCount()
		mi.IncrHandlingNumber()
		if i%17 == 0 {
			mi.Clear()
		}
		counts := mi.Counts()
		expectedSummary := module.SummaryStruct{
			ID:        mi.ID(),
			Called:    counts.CalledCount,
			Accepted:  counts.AcceptedCount,
			Completed: counts.CompletedCount,
			Handling:  counts.HandlingNumber,
		}
		summary := mi.Summary()
		if summary != expectedSummary {
			t.Fatalf("Inconsistent summary for internal module: expected: %#v, actual: %#v",
				expectedSummary, summary)
		}
	}
}

func TestAllInParallel(t *testing.T) {
	number := uint64(100000)
	mi, _ := NewModuleInternal(mid, nil)
	t.Run("CalledCount in Parallel", func(t *testing.T) {
		t.Parallel()
		for i := uint64(1); i < number; i++ {
			mi.IncrCalledCount()
			if mi.CalledCount() != i {
				t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
					i, mi.CalledCount())
			}
		}
	})
	t.Run("AcceptedCount in Parallel", func(t *testing.T) {
		t.Parallel()
		for i := uint64(1); i < number; i++ {
			mi.IncrAcceptedCount()
			if mi.AcceptedCount() != i {
				t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
					i, mi.AcceptedCount())
			}
		}
	})
	t.Run("CompletedCount in Parallel", func(t *testing.T) {
		t.Parallel()
		for i := uint64(1); i < number; i++ {
			mi.IncrCompletedCount()
			if mi.CompletedCount() != i {
				t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
					i, mi.CompletedCount())
			}
		}
	})
	t.Run("HandlingNumber in Parallel", func(t *testing.T) {
		t.Parallel()
		for i := uint64(1); i < number; i++ {
			mi.IncrHandlingNumber()
			mi.IncrHandlingNumber()
			mi.IncrHandlingNumber()
			if mi.HandlingNumber() != i+2 {
				t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
					i+2, mi.HandlingNumber())
			}
			mi.DecrHandlingNumber()
			mi.DecrHandlingNumber()
			if mi.HandlingNumber() != i {
				t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
					i, mi.HandlingNumber())
			}
		}
	})
}
