package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	pgn "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi"
)

func main() {
	pgs.Init().RegisterModule(pgn.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()
}
