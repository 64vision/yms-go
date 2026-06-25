module hyperball.com/account

go 1.24.0

replace hyperball.com/utils => ../utils

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-pg/pg v8.0.7+incompatible
	hyperball.com/sms v0.0.0-00010101000000-000000000000
	hyperball.com/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.36.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace hyperball.com/sms => ../sms

replace hyperball.com/game => ../game

replace hyperball.com/reports => ../reports
