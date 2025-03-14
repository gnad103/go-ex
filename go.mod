module github.com/gnad103/go-ex

go 1.24

require (
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

replace (
	github.com/gnad103/go-ex/proto => ./proto
	github.com/gnad103/go-ex/product-service => ./product-service
	github.com/gnad103/go-ex/user-service => ./user-service
)
