package pprof

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
)

// ProfileType 代表概要类型。
type ProfileType string

var ( // profile flags
	memProfile       = flag.String("memprofile", "", "write a memory profile to the named file after execution")
	memProfileRate   = flag.Int("memprofilerate", 0, "if > 0, sets runtime.MemProfileRate")
	cpuProfile       = flag.String("cpuprofile", "", "write a cpu profile to the named file during execution")
	blockProfile     = flag.String("blockprofile", "", "write a goroutine blocking profile to the named file after execution")
	blockProfileRate = flag.Int("blockprofilerate", 1, "if > 0, calls runtime.SetBlockProfileRate()")
)

var running bool
var lock = new(sync.Mutex)

func parseProfFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}
	*cpuProfile = getAbsFilePath(*cpuProfile)
	*blockProfile = getAbsFilePath(*blockProfile)
	*memProfile = getAbsFilePath(*memProfile)
}

func getAbsFilePath(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.FromSlash(path)
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("Couldn't get current work dir: %s\n", err))
		}
		path = filepath.Join(baseDir, path)
	}
	return path
}

// Start 用于开始采样并生成概要文件。
func Start() {
	lock.Lock()
	defer lock.Unlock()
	parseProfFlags()
	startBlockProfile()
	startCPUProfile()
	startMemProfile()
	running = true
}

func startBlockProfile() {
	if *blockProfile != "" && *blockProfileRate > 0 {
		runtime.SetBlockProfileRate(*blockProfileRate)
	}
}

func startCPUProfile() {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create cpu profile output file: %s\n", err)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't start cpu profile: %s\n", err)
			f.Close()
			return
		}
	}
}

func startMemProfile() {
	if *memProfile != "" && *memProfileRate > 0 {
		runtime.MemProfileRate = *memProfileRate
	}
}

// Stop 用于停止采样并生成概要文件。
func Stop() {
	lock.Lock()
	defer lock.Unlock()
	stopBlockProfile()
	stopCPUProfile()
	stopMemProfile()
	running = false
}

func stopBlockProfile() {
	if *blockProfile != "" && *blockProfileRate >= 0 {
		f, err := os.Create(*blockProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create block profile output file: %s\n", err)
			return
		}
		if err = pprof.Lookup("block").WriteTo(f, 0); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't write %s: %s\n", *blockProfile, err)
		}
		f.Close()
	}
}

func stopCPUProfile() {
	if *cpuProfile != "" {
		pprof.StopCPUProfile() // 把记录的概要信息写到已指定的文件
	}
}

func stopMemProfile() {
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create mem profile output file: %s\n", err)
			return
		}
		if err = pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't write %s: %s\n", *memProfile, err)
		}
		f.Close()
	}
}

// SaveProfile 用于存储概要文件。
func SaveProfile(workDir string, profileName string, ptype ProfileType, debug int) {
	absWorkDir := getAbsFilePath(workDir)
	if profileName == "" {
		profileName = string(ptype)
	}
	profileName += ".out"
	profilePath := filepath.Join(absWorkDir, profileName)
	f, err := os.Create(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create profile output file: %s\n", err)
		return
	}
	if err = pprof.Lookup(string(ptype)).WriteTo(f, debug); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write %s: %s\n", profilePath, err)
	}
	f.Close()
}
