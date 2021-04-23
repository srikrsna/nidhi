package main

import (
	"go/format"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"

	pgn "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi"
)

func main() {
	pgs.Init().RegisterModule(pgn.New()).RegisterPostProcessor(goFmt{}).Render()
}

type goFmt struct{}

func (p goFmt) Match(a pgs.Artifact) bool {
	var n string

	switch a := a.(type) {
	case pgs.GeneratorFile:
		n = a.Name
	case pgs.GeneratorTemplateFile:
		n = a.Name
	case pgs.CustomFile:
		n = a.Name
	case pgs.CustomTemplateFile:
		n = a.Name
	case pgs.GeneratorAppend:
		n = a.FileName
	case pgs.GeneratorTemplateAppend:
		n = a.FileName
	case pgs.GeneratorInjection:
		n = a.FileName
	case pgs.GeneratorTemplateInjection:
		n = a.FileName
	default:
		return false
	}

	return strings.HasSuffix(n, ".go")
}

func (p goFmt) Process(in []byte) ([]byte, error) { return format.Source(in) }

var _ pgs.PostProcessor = goFmt{}
