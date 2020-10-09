package cow

import (
	"sync"
	"testing"
)

func TestConcurrentArray2(t *testing.T) {
	arrayLen := uint32(500)
	t.Run("all2", func(t *testing.T) {
		array := NewConcurrentArray2(arrayLen)
		if array == nil {
			t.Fatalf("Unnormal array!")
		}
		if array.Len() != arrayLen {
			t.Fatalf("Incorrect array length!")
		}
		maxI := uint32(2000)
		t.Run("Set2", func(t *testing.T) {
			testSet2(array, maxI, t)
		})
		t.Run("Get2", func(t *testing.T) {
			testGet2(array, maxI, t)
		})
		t.Run("SetAndGet2", func(t *testing.T) {
			for i := uint32(0); i < arrayLen; i++ {
				testSetAndGet2(arrayLen, t)
			}
		})
	})
}

func testSet2(array ConcurrentArray2, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	var wg sync.WaitGroup
	wg.Add(int(maxI))
	errChan := make(chan error, maxI)
	for i := uint32(0); i < maxI; i++ {
		go func(i uint32) {
			defer wg.Done()
			var err error
			defer func() {
				errChan <- err
			}()
			for j := uint32(0); j < arrayLen; j++ {
				_, err = array.Set(j, int(j*i))
				if err != nil {
					break
				}
			}
		}(i)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func testGet2(array ConcurrentArray2, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	intMax := int((maxI - 1) * (arrayLen - 1))
	for i := uint32(0); i < arrayLen; i++ {
		elem, err := array.Get(i)
		if err != nil {
			t.Fatalf("Unexpected error: %s (index: %d)", err, i)
		}
		if elem < 0 || elem > intMax {
			t.Fatalf("Incorect element: %d! (index: %d, expect max: %d)",
				elem, i, intMax)
		}
	}
}

func testSetAndGet2(maxI uint32, t *testing.T) {
	array := NewConcurrentArray2(maxI)
	var wg sync.WaitGroup
	errChan := make(chan error, maxI)
	for i := uint32(0); i < maxI; i++ {
		wg.Add(1)
		go func(index uint32, t1 *testing.T) {
			defer wg.Done()
			var err error
			defer func() {
				errChan <- err
			}()
			_, err = array.Set(index, int(index))
		}(i, t)
	}
	wg.Wait()
	for j := uint32(0); j < maxI; j++ {
		item, err := array.Get(j)
		if err != nil {
			t.Fatal(err)
		}
		if item != int(j) {
			t.Fatalf("Fail to set array[%d] = %d", j, item)
		}
	}
}
