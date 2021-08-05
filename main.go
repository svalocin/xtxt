package main

import (
	"flag"
	"fmt"

	"github.com/slyerr/xtxt/job/dapenti"
)

var out = flag.String("o", "out", "输出目录")

func main() {
	flag.Parse()

	if err := dapenti.Run(*out); err != nil {
		fmt.Printf("执行喷嚏网 RRS 作业出错：%v", err)
		return
	}

	fmt.Println("ok")
}
