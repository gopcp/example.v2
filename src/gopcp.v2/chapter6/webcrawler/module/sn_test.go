package module

import (
	"math"
	"testing"
)

func TestGenerator(t *testing.T) {
	// 测试最大序列号的自动修正。
	maxmax := uint64(math.MaxUint64)
	start := uint64(1)
	max := uint64(0)
	snGen := NewSNGenertor(start, max)
	if snGen == nil {
		t.Fatalf("Couldn't create SN generator! (start: %d, max: %d)",
			start, max)
	}
	if snGen.Start() != start {
		t.Fatalf("Inconsistent start for SN: expected: %d, actual: %d",
			start, snGen.Start())
	}
	if snGen.Max() != maxmax {
		t.Fatalf("Inconsistent max for SN: expected: %d, actual: %d",
			maxmax, snGen.Max())
	}
	// 测试循环序列号生成器。
	max = uint64(7)
	max = uint64(101)
	snGen = NewSNGenertor(start, max)
	if snGen == nil {
		t.Fatalf("Couldn't create SN generator! (start: %d, max: %d)",
			start, max)
	}
	if snGen.Max() != max {
		t.Fatalf("Inconsistent max for SN: expected: %d, actual: %d",
			max, snGen.Max())
	}
	end := snGen.Max()*5 + 11
	expectedSN := start
	var expectedNext uint64
	var expectedCycleCount uint64
	for i := start; i < end; i++ {
		sn := snGen.Get()
		if expectedSN > snGen.Max() {
			expectedSN = start
		}
		if sn != expectedSN {
			t.Fatalf("Inconsistent ID: expected: %d, actual: %d (index: %d)",
				expectedSN, sn, i)
		}
		expectedNext = expectedSN + 1
		if expectedNext > snGen.Max() {
			expectedNext = start
		}
		if snGen.Next() != expectedNext {
			t.Fatalf("Inconsistent next ID: expected: %d, actual: %d (sn: %d, index: %d)",
				expectedNext, snGen.Next(), sn, i)
		}
		if sn == snGen.Max() {
			expectedCycleCount++
		}
		if snGen.CycleCount() != expectedCycleCount {
			t.Fatalf("Inconsistent cycle count: expected: %d, actual: %d (sn: %d, index: %d)",
				expectedCycleCount, snGen.CycleCount(), sn, i)
		}
		expectedSN++
	}
}
