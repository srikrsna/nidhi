// Code generated by protoc-gen-nidhi. DO NOT EDIT.
//
// Source: test/v1/test.proto

package testv1nidhi

import (
	context "context"
	sql "database/sql"
	nidhi "github.com/srikrsna/nidhi"
	v1 "github.com/srikrsna/nidhi/internal/gen/test/v1"
)

// TestUpdates is an updates type that can be passed to
// [*nidhi.Store.Update] and [*nidhi.Store.UpdateMany]
type TestUpdates struct {
	Title    *string            `json:"title,omitempty"`
	SubTest  **v1.SubTest       `json:"subTest,omitempty"`
	SubTests *[]*v1.SubTest     `json:"subTests,omitempty"`
	M        *map[string]string `json:"m,omitempty"`
}

// TestSchema is the set of path selectors for the Test type.
var TestSchema = struct {
	Id       testId
	Title    testTitle
	SubTest  testSubTest
	SubTests testSubTests
	M        testM
}{}

// TestQuery is an alias for nidhi.Query[TestField]
type TestQuery = nidhi.Query[TestField]

// TestConj is an alias for nidhi.Conj[TestField]
type TestConj = nidhi.Conj[TestField]

// NewTestStore is a document store for Test
func NewTestStore(
	ctx context.Context,
	db *sql.DB,
	opt nidhi.StoreOptions,
) (*nidhi.Store[v1.Test], error) {
	return nidhi.NewStore(
		ctx,
		db,
		"test_v1",
		"test",
		[]string{
			"id",
			"title",
			"subTest",
			"subTests",
			"m",
		},
		func(x *v1.Test) string { return x.Id },
		func(x *v1.Test, id string) { x.Id = id },
		opt,
	)
}

type TestField interface {
	nidhi.Field
	testField()
}
type baseTestField struct{}

func (baseTestField) testField() {}

type testId struct {
	baseTestField
}

func (testId) Selector() string                              { return `JSON_VALUE('$.id' RETURNING TEXT DEFAULT '' ON EMPTY)` }
func (f testId) Is(c *nidhi.StringCond) (testId, nidhi.Cond) { return f, c }

type testTitle struct {
	baseTestField
}

func (testTitle) Selector() string                                 { return `JSON_VALUE('$.title' RETURNING TEXT DEFAULT '' ON EMPTY)` }
func (f testTitle) Is(c *nidhi.StringCond) (testTitle, nidhi.Cond) { return f, c }

type testSubTest struct {
	baseTestField
	Name  testSubTestName
	Inner testSubTestInner
}

func (testSubTest) Selector() string {
	return `JSON_VALUE('$.subTest' RETURNING JSONB DEFAULT '{}' ON EMPTY)`
}
func (f testSubTest) Is(c *nidhi.JsonCond) (testSubTest, nidhi.Cond) { return f, c }

type testSubTestName struct {
	baseTestField
}

func (testSubTestName) Selector() string {
	return `JSON_VALUE('$.subTest.name' RETURNING TEXT DEFAULT '' ON EMPTY)`
}
func (f testSubTestName) Is(c *nidhi.StringCond) (testSubTestName, nidhi.Cond) { return f, c }

type testSubTestInner struct {
	baseTestField
	Yes testSubTestInnerYes
}

func (testSubTestInner) Selector() string {
	return `JSON_VALUE('$.subTest.inner' RETURNING JSONB DEFAULT '{}' ON EMPTY)`
}
func (f testSubTestInner) Is(c *nidhi.JsonCond) (testSubTestInner, nidhi.Cond) { return f, c }

type testSubTestInnerYes struct {
	baseTestField
}

func (testSubTestInnerYes) Selector() string {
	return `JSON_VALUE('$.subTest.inner.yes' RETURNING TEXT DEFAULT '' ON EMPTY)`
}
func (f testSubTestInnerYes) Is(c *nidhi.StringCond) (testSubTestInnerYes, nidhi.Cond) { return f, c }

type testSubTests struct {
	baseTestField
	Name  testSubTestsName
	Inner testSubTestsInner
}

func (testSubTests) Selector() string {
	return `JSON_VALUE('$.subTests' RETURNING JSONB DEFAULT '{}' ON EMPTY)`
}
func (f testSubTests) Is(c *nidhi.JsonCond) (testSubTests, nidhi.Cond) { return f, c }

type testSubTestsName struct {
	baseTestField
}

func (testSubTestsName) Selector() string {
	return `JSON_VALUE('$.subTests[*].name' RETURNING TEXT[] DEFAULT '{}' ON EMPTY)`
}
func (f testSubTestsName) Is(c *nidhi.StringSliceCond) (testSubTestsName, nidhi.Cond) { return f, c }

type testSubTestsInner struct {
	baseTestField
	Yes testSubTestsInnerYes
}

func (testSubTestsInner) Selector() string {
	return `JSON_VALUE('$.subTests[*].inner' RETURNING JSONB DEFAULT '{}' ON EMPTY)`
}
func (f testSubTestsInner) Is(c *nidhi.JsonCond) (testSubTestsInner, nidhi.Cond) { return f, c }

type testSubTestsInnerYes struct {
	baseTestField
}

func (testSubTestsInnerYes) Selector() string {
	return `JSON_VALUE('$.subTests[*].inner.yes' RETURNING TEXT[] DEFAULT '{}' ON EMPTY)`
}
func (f testSubTestsInnerYes) Is(c *nidhi.StringSliceCond) (testSubTestsInnerYes, nidhi.Cond) {
	return f, c
}

type testM struct {
	baseTestField
	Key   testMKey
	Value testMValue
}

func (testM) Selector() string                           { return `JSON_VALUE('$.m' RETURNING JSONB DEFAULT '{}' ON EMPTY)` }
func (f testM) Is(c *nidhi.JsonCond) (testM, nidhi.Cond) { return f, c }

type testMKey struct {
	baseTestField
}

func (testMKey) Selector() string                                { return `JSON_VALUE('$.m.key' RETURNING TEXT DEFAULT '' ON EMPTY)` }
func (f testMKey) Is(c *nidhi.StringCond) (testMKey, nidhi.Cond) { return f, c }

type testMValue struct {
	baseTestField
}

func (testMValue) Selector() string {
	return `JSON_VALUE('$.m.value' RETURNING TEXT DEFAULT '' ON EMPTY)`
}
func (f testMValue) Is(c *nidhi.StringCond) (testMValue, nidhi.Cond) { return f, c }
