package pgn_test

import (
	"fmt"
	"sort"
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

type byInt32Field []*pb.All

// Len is the number of elements in the collection.
func (a byInt32Field) Len() int           { return len(a) }
func (a byInt32Field) Less(i, j int) bool { return a[i].Int32Field < a[j].Int32Field }
func (a byInt32Field) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

type allSliceEqual struct {
	expected []*pb.All
}

func AllSliceEqual(exp []*pb.All) types.GomegaMatcher {
	return &allSliceEqual{
		expected: exp,
	}
}

func (matcher *allSliceEqual) Match(actual interface{}) (success bool, err error) {
	response, ok := actual.([]*pb.All)
	if !ok {
		return false, fmt.Errorf("ProtoEqual matcher expects a proto.Message")
	}

	sort.Sort(byId(response))
	sort.Sort(byId(matcher.expected))

	if len(response) != len(matcher.expected) {
		return false, nil
	}

	for i := range response {
		if !proto.Equal(response[i], matcher.expected[i]) {
			return false, nil
		}
	}

	return true, nil
}

func (matcher *allSliceEqual) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Actual\n\t%v\nExpected\n\t%v", actual, matcher.expected)
}

func (matcher *allSliceEqual) NegatedFailureMessage(actual interface{}) (message string) {
	return matcher.FailureMessage(actual)
}
