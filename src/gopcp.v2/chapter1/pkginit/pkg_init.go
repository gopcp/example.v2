package main // 命令源码文件必须在这里声明自己属于main包。

import ( // 引入了代码包fmt和runtime。
	"fmt"
	"runtime"
)

func init() { // 代码包初始化函数。
	fmt.Printf("Map: %v\n", m) // 格式化的打印。
	// 通过调用runtime包的代码获取当前机器的操作系统和计算架构。
	// 而后通过fmt包的Sprintf方法进行格式化字符串生成并赋值给变量info。
	info = fmt.Sprintf("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
}

// 非局部变量，map类型，且已初始化。
var m = map[int]string{1: "A", 2: "B", 3: "C"}

// 非局部变量，string类型，未被初始化。
var info string

func main() { // 命令源码文件必须有的入口函数，也称主函数。
	fmt.Println(info) // 打印变量info。
}
