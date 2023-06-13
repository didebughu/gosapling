package main

import (
	"github.com/didebughu/gosapling/analysis/pointer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(pointer.Analyzer) }
