package cow

import (
	"sync"
	"testing"
)

func TestConcurrentIntArray(t *testing.T) {
	arrayLen := 1000
	t.Run("all2", func(t *testing.T) {
		array := NewConcurrentIntArray(arrayLen)
		if array == nil {
			t.Fatalf("Unnormal array!")
		}
		if array.Len() != arrayLen {
			t.Fatalf("Incorrect array length!")
		}
		maxI := 2000
		t.Run("Set2", func(t *testing.T) {
			testSet2(array, maxI, t)
		})
		t.Run("Get2", func(t *testing.T) {
			testGet2(array, maxI, t)
		})
		t.Run("SetAndGet2", func(t *testing.T) {
			for i := 0; i < arrayLen; i++ {
				testSetAndGet2(arrayLen, t)
			}
		})
	})
}

func testSet2(array ConcurrentIntArray, maxI int, t *testing.T) {
	arrayLen := array.Len()
	var wg sync.WaitGroup
	wg.Add(maxI)
	errChan := make(chan error, maxI)
	for i := 0; i < maxI; i++ {
		go func(i int) {
			defer wg.Done()
			var err error
			defer func() {
				errChan <- err
			}()
			for j := 0; j < arrayLen; j++ {
				_, err = array.Set(j, j*i)
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

func testGet2(array ConcurrentIntArray, maxI int, t *testing.T) {
	arrayLen := array.Len()
	intMax := (maxI - 1) * (arrayLen - 1)
	for i := 0; i < arrayLen; i++ {
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

func testSetAndGet2(maxI int, t *testing.T) {
	array := NewConcurrentIntArray(maxI)
	var wg sync.WaitGroup
	errChan := make(chan error, maxI)
	for i := 0; i < maxI; i++ {
		wg.Add(1)
		go func(index int, t1 *testing.T) {
			defer wg.Done()
			var err error
			defer func() {
				errChan <- err
			}()
			_, err = array.Set(index, int(index))
		}(i, t)
	}
	wg.Wait()
	for j := 0; j < maxI; j++ {
		item, err := array.Get(j)
		if err != nil {
			t.Fatal(err)
		}
		if item != int(j) {
			t.Fatalf("Fail to set array[%d] = %d", j, item)
		}
	}
}
