package pkgtool

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

// getImportsFromGoSource 会返回指定源码文件中导入的代码包的导入路径。
func getImportsFromGoSource(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	r := bufio.NewReader(f)
	var isMultiImport bool
	var importPaths []string
	for {
		lineBytes, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line := strings.TrimSpace(string(lineBytes))
		// Ignore when build ignore
		if line == "// +build ignore" {
			return importPaths, nil
		}
		// Ignore command source
		// if strings.HasPrefix(line, "package") && strings.Split(line, " ")[1] == "main" {
		// 	return importPaths, nil
		// }
		if strings.HasPrefix(line, "import") {
			if strings.Contains(line, "(") {
				isMultiImport = true
			} else {
				importPath := line[strings.Index(line, "\"")+1 : strings.LastIndex(line, "\"")]
				importPaths = appendIfAbsent(importPaths, importPath)
				break
			}
		} else {
			if isMultiImport {
				if strings.HasPrefix(line, ")") {
					break
				} else {
					// Ignore irregular import
					if !strings.HasPrefix(line, "\"") || !strings.HasSuffix(line, "\"") {
						continue
					}
					importPath := strings.Replace(line, "\"", "", 2)
					importPaths = appendIfAbsent(importPaths, importPath)
				}
			}
		}
	}
	sort.Strings(importPaths)
	return importPaths, nil
}

// getImportsFromPackage 会返回指定代码包中导入的所有代码包的导入路径。
func getImportsFromPackage(importPath string, containsTestFile bool) ([]string, error) {
	var importPaths []string
	packageAbsPath := getAbsPathOfPackage(importPath)
	// return empty slice when the import path is invalid.
	if packageAbsPath == "" {
		return importPaths, nil
	}
	srcFileAbsPaths, err := getGoSourceFileAbsPaths(packageAbsPath, containsTestFile)
	if err != nil {
		return nil, err
	}
	for _, v := range srcFileAbsPaths {
		currentImportPaths, err := getImportsFromGoSource(v)
		if err != nil {
			return nil, err
		}
		importPaths = appendIfAbsent(importPaths, currentImportPaths...)
	}
	return importPaths, nil
}
