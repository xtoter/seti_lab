// Иногда наши программы Go должны порождать другие, не
// Go процессы. Например, подсветка синтаксиса на этом
// сайте [реализуется](https://github.com/mmcgrana/gobyexample/blob/master/tools/generate.go)
// путем запуска [`pygmentize`](http://pygments.org/)
// процесса из программы Go. Давайте рассмотрим несколько
// примеров порождающих процессов из Go.

package main

import (
	"fmt"
	"os/exec"
)

func main() {
	dateCmd := exec.Command("traceroute", "ya.ru")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dateOut))
}
