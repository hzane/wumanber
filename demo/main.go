package main

import (
	"fmt"

	"github.com/heartszhang/wumanber"
)

func main() {
	patterns := []string{"tst", "only", "hello", "his", "中文"}
	wm := wumanber.New(patterns)
	x := wm.Search("hello this is only english 中文 test")
	for _, idx := range x {
		fmt.Println(patterns[idx])
	}
}
