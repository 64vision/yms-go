module gollux/client

go 1.26.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	gollux/account v0.0.0-00010101000000-000000000000
	gollux/auth v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-pg/pg v8.0.7+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	gollux/sms v0.0.0-00010101000000-000000000000 // indirect
	gollux/utils v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/account => ../account

replace gollux/auth => ../auth

replace gollux/utils => ../utils

replace gollux/sms => ../sms
