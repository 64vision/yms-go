module hyperball.com/wallet

go 1.24.0

replace hyperball.com/utils => ../utils

require (
	github.com/go-pg/pg v8.0.7+incompatible
	google.golang.org/grpc v1.70.0
	hyperball.com/utils v0.0.0-00010101000000-000000000000
	hyperball.com/wallet/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.36.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace wallet/proto => ../wallet/wallet/proto

replace hyperball.com/wallet/proto => ../wallet/wallet/proto
