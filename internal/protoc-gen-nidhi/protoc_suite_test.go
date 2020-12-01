package pgn_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	pb "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data"
)

func TestProtocNidhi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protoc Nidhi Suite")
}

type byId []*pb.All

// Len is the number of elements in the collection.
func (a byId) Len() int           { return len(a) }
func (a byId) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a byId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func ProtoEqual(expected proto.Message) types.GomegaMatcher {
	return &protoEqual{
		expected: expected,
	}
}

type protoEqual struct {
	expected proto.Message
}

func (matcher *protoEqual) Match(actual interface{}) (success bool, err error) {
	response, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("ProtoEqual matcher expects a proto.Message")
	}

	return proto.Equal(response, matcher.expected), nil
}

func (matcher *protoEqual) FailureMessage(actual interface{}) (message string) {
	msg := actual.(proto.Message)
	return fmt.Sprintf("Actual\n\t%v\nExpected\n\t%v", prototext.Format(msg), prototext.Format(matcher.expected))
}

func (matcher *protoEqual) NegatedFailureMessage(actual interface{}) (message string) {
	return matcher.FailureMessage(actual)
}
