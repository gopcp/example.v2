package buffer

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBufferNew(t *testing.T) {
	size := uint32(10)
	buf, err := NewBuffer(size)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, size)
	}
	if buf == nil {
		t.Fatal("Couldn't create buffer!")
	}
	if buf.Cap() != size {
		t.Fatalf("Inconsistent buffer cap: expected: %d, actual: %d",
			size, buf.Cap())
	}
	buf, err = NewBuffer(0)
	if err == nil {
		t.Fatal("No error when new a buffer with zero size!")
	}
}

func TestBufferPut(t *testing.T) {
	size := uint32(10)
	buf, err := NewBuffer(size)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, size)
	}
	data := make([]uint32, size)
	for i := uint32(0); i < size; i++ {
		data[i] = i
	}
	var count uint32
	var datum uint32
	for _, datum = range data {
		ok, err := buf.Put(datum)
		if err != nil {
			t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
				err, datum)
		}
		if !ok {
			t.Fatalf("Couldn't put datum to the buffer! (datum: %d)",
				datum)
		}
		count++
		if buf.Len() != count {
			t.Fatalf("Inconsistent buffer len: expected: %d, actual: %d",
				count, buf.Len())
		}
	}
	datum = size
	ok, err := buf.Put(datum)
	if err != nil {
		t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
			err, datum)
	}
	if ok {
		t.Fatalf("It still can put datum to the full buffer! (datum: %d)",
			datum)
	}
	buf.Close()
	_, err = buf.Put(datum)
	if err == nil {
		t.Fatalf("It still can put datum to the closed buffer! (datum: %d)", datum)
	}
}

func TestBufferPutInParallel(t *testing.T) {
	size := uint32(22)
	bufferSize := uint32(20)
	buf, err := NewBuffer(bufferSize)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, size)
	}
	data := make([]uint32, size)
	for i := uint32(0); i < size; i++ {
		data[i] = i
	}
	testingFunc := func(datum interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := buf.Put(datum)
			if err != nil {
				t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
					err, datum)
			}
			if !ok && buf.Len() < buf.Cap() {
				t.Fatalf("Couldn't put datum to the buffer! (datum: %d)",
					datum)
			}
		}
	}
	t.Run("Put in parallel(1)", func(t *testing.T) {
		for _, datum := range data[:size/2] {
			t.Run(fmt.Sprintf("Datum=%d", datum), testingFunc(datum, t))
		}
	})
	t.Run("Put in parallel(2)", func(t *testing.T) {
		for _, datum := range data[size/2:] {
			t.Run(fmt.Sprintf("Datum=%d", datum), testingFunc(datum, t))
		}
	})
	if buf.Len() != buf.Cap() {
		t.Fatalf("Inconsistent buffer len: expected: %d, actual: %d",
			buf.Cap(), buf.Len())
	}
}

func TestBufferGet(t *testing.T) {
	size := uint32(10)
	buf, err := NewBuffer(size)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, size)
	}
	for i := uint32(0); i < size; i++ {
		buf.Put(i)
	}
	count := size
	var datum uint32
	var ok bool
	for i := uint32(0); i < size; i++ {
		d, err := buf.Get()
		if err != nil {
			t.Fatalf("An error occurs when getting a datum from the buffer: %s",
				err)
		}
		datum, ok = d.(uint32)
		if !ok {
			t.Fatalf("Inconsistent datum type: expected: %T, actual: %T",
				datum, d)
		}
		if datum != i {
			t.Fatalf("Inconsistent datum: expected: %#v, actual: %#v",
				i, datum)
		}
		count--
		if buf.Len() != count {
			t.Fatalf("Inconsistent buffer len: expected: %d, actual: %d",
				count, buf.Len())
		}
	}
	d, err := buf.Get()
	if err != nil {
		t.Fatalf("An error occurs when getting a datum from the buffer: %s",
			err)
	}
	if d != nil {
		t.Fatal("It still can get a datum from the empty buffer!")
	}
	datum = 0
	buf.Put(datum)
	buf.Close()
	_, err = buf.Get()
	if err != nil {
		t.Fatalf("An error occurs when getting a datum from the buffer: %s",
			err)
	}
	_, err = buf.Get()
	if err == nil {
		t.Fatal("It still can get datum from the closed buffer!")
	}
}

func TestBufferGetInParallel(t *testing.T) {
	bufferSize := uint32(30)
	buf, err := NewBuffer(bufferSize)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, bufferSize)
	}
	for i := uint32(0); i < bufferSize; i++ {
		buf.Put(i)
	}
	marks := make([]uint8, bufferSize)
	var lock sync.Mutex
	testingFunc := func(t *testing.T) {
		t.Parallel()
		d, err := buf.Get()
		if err != nil {
			t.Fatalf("An error occurs when getting a datum from the buffer: %s",
				err)
		}
		if d == nil && buf.Len() != 0 {
			t.Fatalf("Get an empty datum! (len: %d)", buf.Len())
		}
		if d != nil {
			datum := d.(uint32)
			lock.Lock()
			marks[int(datum)]++
			lock.Unlock()
		}
	}
	t.Run("Get in parallel(1)", func(t *testing.T) {
		num := bufferSize/2 + 1
		for i := uint32(0); i < num; i++ {
			t.Run(fmt.Sprintf("Index=%d", i), testingFunc)
		}
	})
	t.Run("Get in parallel(2)", func(t *testing.T) {
		num := bufferSize/2 + 2
		for i := uint32(0); i < num; i++ {
			t.Run(fmt.Sprintf("Index=%d", i), testingFunc)
		}
	})
	if buf.Len() != 0 {
		t.Fatalf("Inconsistent buffer len: expected: %d, actual: %d",
			0, buf.Len())
	}
	for i, m := range marks {
		if m == 0 && i != 0 {
			t.Fatalf("Haven’t got the number: %d", i)
		}
		if m > 1 {
			t.Fatalf("Got the number more than once: %d", i)
		}
	}
}

func TestBufferPutAndGetInParallel(t *testing.T) {
	bufferSize := uint32(50)
	buf, err := NewBuffer(bufferSize)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, bufferSize)
	}
	maxPuttingNumber := bufferSize + uint32(rand.Int31n(20))
	maxGettingNumber := bufferSize + uint32(rand.Int31n(20))
	puttingCount := maxPuttingNumber
	gettingCount := maxGettingNumber
	marks := make([]uint8, maxPuttingNumber)
	var lock sync.Mutex
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			begin := uint32(0)
			end := maxPuttingNumber / 2
			for i := begin; i < end; i++ {
				ok, err := buf.Put(i)
				if err != nil {
					t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
						err, i)
				}
				if !ok &&
					atomic.LoadUint32(&gettingCount) == 0 &&
					buf.Len() < buf.Cap() {
					t.Fatalf("Couldn't put datum to the buffer! (datum: %d)",
						i)
				}
				atomic.AddUint32(&puttingCount, ^uint32(0))
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			begin := maxPuttingNumber / 2
			end := maxPuttingNumber
			for i := begin; i < end; i++ {
				ok, err := buf.Put(i)
				if err != nil {
					t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
						err, i)
				}
				if !ok &&
					atomic.LoadUint32(&gettingCount) == 0 &&
					buf.Len() < buf.Cap() {
					t.Fatalf("Couldn't put datum to the buffer! (datum: %d)",
						i)
				}
				atomic.AddUint32(&puttingCount, ^uint32(0))
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			max := bufferSize/2 + 1
			for i := uint32(0); i < max; i++ {
				d, err := buf.Get()
				if err != nil {
					t.Fatalf("An error occurs when getting a datum from the buffer: %s",
						err)
				}
				if d == nil &&
					atomic.LoadUint32(&puttingCount) == 0 &&
					buf.Len() != 0 {
					t.Fatalf("Get an empty datum! (len: %d)", buf.Len())
				}
				atomic.AddUint32(&gettingCount, ^uint32(0))
				if d != nil {
					datum := d.(uint32)
					lock.Lock()
					marks[int(datum)]++
					lock.Unlock()
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			max := bufferSize/2 + 2
			for i := uint32(0); i < max; i++ {
				d, err := buf.Get()
				if err != nil {
					t.Fatalf("An error occurs when getting a datum from the buffer: %s",
						err)
				}
				if d == nil &&
					atomic.LoadUint32(&puttingCount) == 0 &&
					buf.Len() != 0 {
					t.Fatalf("Get an empty datum! (len: %d)", buf.Len())
				}
				atomic.AddUint32(&gettingCount, ^uint32(0))
				if d != nil {
					datum := d.(uint32)
					lock.Lock()
					marks[int(datum)]++
					lock.Unlock()
				}
			}
		})
	})
	for i, m := range marks {
		if m > 1 {
			t.Fatalf("Got the number more than once: %d", i)
		}
	}
}

func TestBufferCloseInParallel(t *testing.T) {
	bufferSize := uint32(100)
	buf, err := NewBuffer(bufferSize)
	if err != nil {
		t.Fatalf("An error occurs when new a buffer: %s (size: %d)",
			err, bufferSize)
	}
	maxNumber := bufferSize + uint32(rand.Int31n(100))
	t.Run("Put", func(t *testing.T) {
		t.Parallel()
		for i := uint32(0); i < maxNumber; i++ {
			_, err := buf.Put(i)
			if err != nil && !buf.Closed() {
				t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
					err, i)
			}
			if err == nil && buf.Closed() {
				t.Fatalf("It still can put datum to the closed buffer! (datum: %d)", i)
			}
		}
	})
	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		max := bufferSize/2 + 1
		for i := uint32(0); i < max; i++ {
			_, err := buf.Get()
			if err != nil && !buf.Closed() {
				t.Fatalf("An error occurs when getting a datum from the buffer: %s (datum: %d)",
					err, i)
			}
			if buf.Closed() {
				if _, err = buf.Get(); err == nil {
					t.Fatalf("It still can get datum from the closed buffer! (datum: %d)", i)
				}
			}
		}
	})
	t.Run("Close", func(t *testing.T) {
		t.Parallel()
		time.Sleep(time.Millisecond)
		ok := buf.Close()
		if !ok {
			t.Fatal("Couldn't close the buffer!")
		}
		if !buf.Closed() {
			t.Fatalf("Inconsistent buffer status: expected closed: %v, actual closed: %v",
				true, buf.Closed())
		}
		ok = buf.Close()
		if ok {
			t.Fatal("It still can close the closed buffer!")
		}
		if !buf.Closed() {
			t.Fatalf("Inconsistent buffer status: expected closed: %v, actual closed: %v",
				true, buf.Closed())
		}
	})
}
