module admin

go 1.26.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	gollux/account v0.0.0-00010101000000-000000000000
	gollux/auth v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-pg/pg v8.0.7+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	gollux/email v0.0.0-00010101000000-000000000000 // indirect
	gollux/sms v0.0.0-00010101000000-000000000000 // indirect
	gollux/utils v0.0.0-00010101000000-000000000000 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/auth => ../auth

replace gollux/account => ../account

replace gollux/utils => ../utils

replace gollux/game => ../game

replace gollux/reports => ../reports

replace gollux/sms => ../sms

replace gollux/tools => ../tools

replace gollux/email => ../email
