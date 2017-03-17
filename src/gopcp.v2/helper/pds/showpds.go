// Show the dependency structure of specified package
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"gopcp.v2/helper/pkgtool"
)

// ARROWS 代表有导入关系的代码包之间的分隔符。
const ARROWS = "->"

var (
	pkgImportPathFlag string
)

func init() {
	flag.StringVar(&pkgImportPathFlag, "p", "", "The path of target package.")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("FATAL ERROR: %s", err)
			debug.PrintStack()
		}
	}()
	flag.Parse()
	pkgImportPath := getPkgImportPath()
	pn := pkgtool.NewPkgNode(pkgImportPath)
	fmt.Printf("The package node of '%s': %v\n", pkgImportPath, *pn)
	err := pn.Grow()
	if err != nil {
		fmt.Printf("GROW ERROR: %s\n", err)
	}
	fmt.Printf("The dependency structure of package '%s':\n", pkgImportPath)
	ShowDepStruct(pn, 0, "")
}

// 展示的序号。
var sn = 1

// ShowDepStruct 用于展示代码包依赖结构。
func ShowDepStruct(pnode *pkgtool.PkgNode, depth int, prefix string) {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	importPath := pnode.ImportPath()
	buf.WriteString(importPath)
	deps := pnode.ImportedNodes()
	// fmt.Printf("P_NODE: '%s', DEP_LEN: %d\n", importPath, len(deps))
	if len(deps) == 0 {
		fmt.Printf("%d[%d]: %s\n", sn, depth, buf.String())
		sn++
		return
	}
	buf.WriteString(ARROWS)
	depth++
	for _, v := range deps {
		ShowDepStruct(v, depth, buf.String())
	}
}

// getPkgImportPath 会返回指定目录的代码包导入路径。
func getPkgImportPath() string {
	if len(pkgImportPathFlag) > 0 {
		return pkgImportPathFlag
	}
	fmt.Printf("The flag p is invalid, use current dir as package import path.\n")
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	srcDirs := pkgtool.GetSrcDirs(false)
	var importPath string
	for _, v := range srcDirs {
		if strings.HasPrefix(currentDir, v) {
			importPath = currentDir[len(v)+1:]
			break
		}
	}
	if strings.TrimSpace(importPath) == "" {
		panic(errors.New("Couldn't parse the import path!"))
	}
	return importPath
}
