module gollux/account

go 1.26.4

replace gollux/utils => ../utils

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-pg/pg v8.0.7+incompatible
	gollux/sms v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
	zerasuite/shippinglines v0.0.0-00010101000000-000000000000
	zerasuite/yards v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.36.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/sms => ../sms

replace gollux/game => ../game

replace gollux/reports => ../reports

replace zerasuite/yards => ../yards

replace zerasuite/shippinglines => ../shippinglines
