module zerasuite/webhook

go 1.26.4

require (
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	gollux/auth v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.42.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/utils => ../utils

replace gollux/auth => ../auth
