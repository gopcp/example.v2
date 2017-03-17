package pkgtool

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
)

func TestGetGoroot(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	expGoroot := runtime.GOROOT()
	actGoroot := GetGoroot()
	if actGoroot != expGoroot {
		t.Errorf("Error: The goroot should be '%s'' but '%s'.\n", expGoroot, actGoroot)
	}
}

func TestGetAllGopath(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	expGopath := os.Getenv("GOPATH")
	actGopaths := GetAllGopath()
	var sep string
	if runtime.GOOS == "windows" {
		sep = ";"
	} else {
		sep = ":"
	}
	expGopathLen := len(expGopath)
	tempString := expGopath
	tempIndex := -1
	sepNumber := 0
	for {
		if tempString == "" {
			break
		}
		tempIndex = strings.Index(tempString, sep)
		if tempIndex < 0 {
			sepNumber++
			break
		} else {
			sepNumber++
		}
		if (tempIndex + 1) >= expGopathLen {
			break
		}
		tempString = tempString[tempIndex+1:]
	}
	actLen := len(actGopaths)
	if actLen != sepNumber {
		t.Errorf("Error: The length of gopaths should be %d but %d.\n", sepNumber, actLen)
	}
}

func TestGetSrcDirs(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	actSrcDirs := GetSrcDirs(true)
	actSrcDirLen := len(actSrcDirs)
	expSrcDirLen := len(GetAllGopath()) + 1
	if actSrcDirLen != expSrcDirLen {
		t.Errorf("Error: The length of gopaths should be %d but %d.\n", expSrcDirLen, actSrcDirLen)
	}
	fileSep := string(filepath.Separator)
	for _, v := range actSrcDirs {
		if !strings.HasSuffix(v, (fileSep+"src")) &&
			!(strings.HasPrefix(v, GetGoroot()) &&
				strings.HasSuffix(v, fileSep+"src"+fileSep+"pkg")) {
			t.Errorf("Error: The src dir '%s' is incorrect.\n", v)
		}
	}
}
