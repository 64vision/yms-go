module tools

go 1.24.0

replace gollux/utils => ../utils

require (
	github.com/go-pg/pg v8.0.7+incompatible
	gollux/account v0.0.0-00010101000000-000000000000
	gollux/sms v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/sms => ../sms

replace gollux/account => ../account
