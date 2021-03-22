module github.com/srikrsna/nidhi

go 1.16

require (
	github.com/elgris/sqrl v0.0.0-20190909141434-5a439265eeec
	github.com/gertd/go-pluralize v0.1.7
	github.com/google/go-cmp v0.5.5
	github.com/google/gofuzz v1.2.0
	github.com/google/uuid v1.1.2
	github.com/json-iterator/go v1.1.10
	github.com/lyft/protoc-gen-star v0.5.2
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/srikrsna/protoc-gen-fuzz/wkt v0.0.0-20210321095126-38865cd101ba
	gocloud.dev v0.20.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/srikrsna/protoc-gen-fuzz/wkt => ../protoc-gen-fuzz/wkt
