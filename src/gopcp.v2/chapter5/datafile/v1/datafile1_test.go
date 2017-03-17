package v1

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func removeFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	file.Close()
	return os.Remove(path)
}

func TestIDataFile(t *testing.T) {
	t.Run("v1/all", func(t *testing.T) {
		dataLen := uint32(3)
		path1 := filepath.Join(os.TempDir(), "data_file_test_new.txt")
		defer func() {
			if err := removeFile(path1); err != nil {
				t.Errorf("Open file error: %s\n", err)
			}
		}()
		t.Run("New", func(t *testing.T) {
			testNew(path1, dataLen, t)
		})
		path2 := filepath.Join(os.TempDir(), "data_file_test.txt")
		defer func() {
			if err := removeFile(path2); err != nil {
				t.Fatalf("Open file error: %s\n", err)
			}
		}()
		max := 100000
		t.Run("WriteAndRead", func(t *testing.T) {
			testRW(path2, dataLen, max, t)
		})
	})
}

func testNew(path string, dataLen uint32, t *testing.T) {
	t.Logf("New a data file (path: %s, dataLen: %d)...\n",
		path, dataLen)
	dataFile, err := NewDataFile(path, dataLen)
	if err != nil {
		t.Logf("Couldn't new a data file: %s", err)
		t.FailNow()
	}
	if dataFile == nil {
		t.Log("Unnormal data file!")
		t.FailNow()
	}
	defer dataFile.Close()
	if dataFile.DataLen() != dataLen {
		t.Fatalf("Incorrect data length!")
	}
}

func testRW(path string, dataLen uint32, max int, t *testing.T) {
	t.Logf("New a data file (path: %s, dataLen: %d)...\n",
		path, dataLen)
	dataFile, err := NewDataFile(path, dataLen)
	if err != nil {
		t.Logf("Couldn't new a data file: %s", err)
		t.FailNow()
	}
	defer dataFile.Close()
	var wg sync.WaitGroup
	wg.Add(5)
	// 写入。
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			var prevWSN int64 = -1
			for j := 0; j < max; j++ {
				data := Data{
					byte(rand.Int31n(256)),
					byte(rand.Int31n(256)),
					byte(rand.Int31n(256)),
				}
				wsn, err := dataFile.Write(data)
				if err != nil {
					t.Fatalf("Unexpected writing error: %s\n", err)
				}
				if prevWSN >= 0 && wsn <= prevWSN {
					t.Fatalf("Incorect WSN %d! (lt %d)\n", wsn, prevWSN)
				}
				prevWSN = wsn
			}
		}()
	}
	// 读取。
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			var prevRSN int64 = -1
			for i := 0; i < max; i++ {
				rsn, date, err := dataFile.Read()
				if err != nil {
					t.Fatalf("Unexpected writing error: %s\n", err)
				}
				if date == nil {
					t.Fatalf("Unnormal data!")
				}
				if prevRSN >= 0 && rsn <= prevRSN {
					t.Fatalf("Incorect RSN %d! (lt %d)\n", rsn, prevRSN)
				}
				prevRSN = rsn
			}
		}()
	}
	wg.Wait()
}
