package pgn_test

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/spf13/afero"

	pgn "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi"
)

var (
	update = flag.Bool("update", false, "True will replace the golden files output with latest output")
)

func TestModule(t *testing.T) {
	req, err := os.Open("./test_data/code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer req.Close()

	fs := afero.NewMemMapFs()

	var out bytes.Buffer
	pgs.Init(
		pgs.ProtocInput(req),
		pgs.ProtocOutput(&out),
		pgs.FileSystem(fs),
	).RegisterModule(pgn.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()

	if *update {
		res, err := os.Create("./test_data/code_generator_response.pb.bin")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Close()

		if _, err := io.Copy(res, &out); err != nil {
			t.Error(err)
		}
	} else {
		exp, err := ioutil.ReadFile("./test_data/code_generator_response.pb.bin")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(exp, out.Bytes()) {
			t.Fatal("Output mismatch")
		}
	}
}
