package cow

import (
	"sync"
	"testing"
)

func TestConcurrentArray(t *testing.T) {
	arrayLen := uint32(1000)
	t.Run("all", func(t *testing.T) {
		array := NewConcurrentArray(arrayLen)
		if array == nil {
			t.Fatalf("Unnormal array!")
		}
		if array.Len() != arrayLen {
			t.Fatalf("Incorrect array length!")
		}
		maxI := uint32(2000)
		t.Run("Set", func(t *testing.T) {
			testSet(array, maxI, t)
		})
		t.Run("Get", func(t *testing.T) {
			testGet(array, maxI, t)
		})
	})
}

func testSet(array ConcurrentArray, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	var wg sync.WaitGroup
	wg.Add(int(maxI))
	for i := uint32(0); i < maxI; i++ {
		go func(i uint32) {
			defer wg.Done()
			for j := uint32(0); j < arrayLen; j++ {
				err := array.Set(j, int(j*i))
				if uint32(j) >= arrayLen && err == nil {
					t.Fatalf("Unexpected nil error! (index: %d)", j)
				} else {
					if err != nil {
						t.Fatalf("Unexpected error: %s (index: %d)", err, j)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func testGet(array ConcurrentArray, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	intMax := int((maxI - 1) * (arrayLen - 1))
	for i := uint32(0); i < arrayLen; i++ {
		elem, err := array.Get(i)
		if err != nil {
			t.Fatalf("Unexpected error: %s (index: %d)", err, i)
		}
		if elem < 0 || elem > intMax {
			t.Fatalf("Incorrect element: %d! (index: %d, expect max: %d)",
				elem, i, intMax)
		}
	}
}
