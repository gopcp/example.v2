package pprof

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"testing"
	"time"
)

const (
	tempWorkDir      = "../../../pprof"
	cpuProfilePath   = tempWorkDir + "/cpu.out"
	memProfilePath   = tempWorkDir + "/mem.out"
	blockProfilePath = tempWorkDir + "/block.out"
)

var newTempWorkDir bool
var removeTempWorkDir bool = false

func TestCommonProf(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	makeWorkDir()
	*cpuProfile = cpuProfilePath
	*blockProfile = blockProfilePath
	*memProfile = memProfilePath
	Start()
	doSomething()
	Stop()
	cleanWorkDir()
}

func doSomething() {
	size := 10000000
	randomNumbers := make([]int, size)
	max := int32(size)
	for i := 0; i < size; i++ {
		target := rand.Int31n(max)
		randomNumbers[target] = i
	}
}

func makeWorkDir() {
	absTempWorkDir := getAbsFilePath(tempWorkDir)
	f, err := os.Open(absTempWorkDir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			panic(err)
		} else {
			newTempWorkDir = true
		}
	}
	if f != nil {
		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}
		if !fi.IsDir() {
			panic("There are name of a file conflict with temp work dir name!")
		} else {
			fmt.Printf("Temp work dir: '%s'\n", absTempWorkDir)
		}
	}
	if f == nil {
		fmt.Printf("Make temp work dir '%s'...\n", absTempWorkDir)
		err = os.Mkdir(absTempWorkDir, os.ModeDir|os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	_ = f
}

func cleanWorkDir() {
	if removeTempWorkDir {
		err := remove(cpuProfilePath, 5, 2)
		if err != nil {
			fmt.Errorf("Error: Couldn't remove cpu profile '%s': %s\n",
				getAbsFilePath(cpuProfilePath), err)
		}
		err = remove(blockProfilePath, 5, 2)
		if err != nil {
			fmt.Errorf("Error: Couldn't remove block profile '%s': %s\n",
				getAbsFilePath(blockProfilePath), err)
		}
		err = remove(memProfilePath, 5, 2)
		if err != nil {
			fmt.Errorf("Error: Couldn't remove mem profile '%s': %s\n",
				getAbsFilePath(memProfilePath), err)
		}
	}
	if newTempWorkDir && removeTempWorkDir {
		err := remove(tempWorkDir, 5, 3)
		if err != nil {
			fmt.Errorf("Error: Couldn't remove temp work dir '%s': %s\n",
				getAbsFilePath(tempWorkDir), err)
		}
	}
}

func remove(path string, repeat int, intervalSeconds int) error {
	if repeat < 0 {
		repeat = 5
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 2
	}
	loopNumber := repeat + 1
	absPath := getAbsFilePath(path)
	for i := 1; i <= loopNumber; i++ {
		fmt.Printf("Try to remove file/dir '%s'...", absPath)
		err := os.Remove(absPath)
		if err == nil {
			fmt.Println("ok.")
			break
		} else {
			fmt.Println("failing!")
		}
		if err != nil && i == loopNumber {
			return err
		}
		if err != nil && i < loopNumber {
			time.Sleep(time.Duration(intervalSeconds) * time.Second)
		}
	}
	return nil
}

func TestRuntimeProf(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	for _, v := range []int{0, 1, 2} {
		for _, w := range []string{"goroutine", "threadcreate", "heap", "block"} {
			profileName := fmt.Sprintf("%s.%d", w, v)
			fmt.Printf("Try to save %s profile to file '%s' ...\n", w, profileName)
			SaveProfile(tempWorkDir, profileName, ProfileType(w), v)
		}
	}
}
