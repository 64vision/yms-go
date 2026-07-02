module zerasuite/bookings

go 1.26.4

replace gollux/utils => ../utils

require (
	github.com/go-pg/pg v8.0.7+incompatible
	gollux/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.42.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	mellium.im/sasl v0.3.2 // indirect
)
