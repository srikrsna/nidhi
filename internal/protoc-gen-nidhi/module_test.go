package pgn_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/spf13/afero"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	pgn "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi"
)

func TestModule(t *testing.T) {
	req, err := os.Open("./test_data/code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	fs := afero.NewMemMapFs()

	var out bytes.Buffer
	pgs.Init(
		pgs.ProtocInput(req),
		pgs.ProtocOutput(&out),
		pgs.FileSystem(fs),
	).RegisterModule(pgn.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()

	var res plugin_go.CodeGeneratorResponse
	if err := proto.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	fmt.Print(&res)
}
