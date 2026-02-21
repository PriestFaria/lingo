package main

import (
	"github.com/PriestFaria/lingo/internal/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}










