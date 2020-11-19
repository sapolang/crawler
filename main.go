package main

import (
	"crawler/rules/gushiwen"

	"github.com/henrylee2cn/pholcus/exec"
)

func main() {
	gushiwen.Gushiwen.Register()
	exec.DefaultRun("web")
}
