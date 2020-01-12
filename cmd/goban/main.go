package main

import (
	"github.com/utrack/goban/goban"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(goban.Analyzer())
}
