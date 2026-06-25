module hyperball.com/router

go 1.24.0

replace hyperball.com/account => ../account

require (
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	hyperball.com/account v0.0.0-00010101000000-000000000000
	hyperball.com/auth v0.0.0-00010101000000-000000000000
	hyperball.com/game v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-pg/pg v8.0.7+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	hyperball.com/reports v0.0.0-00010101000000-000000000000 // indirect
	hyperball.com/sms v0.0.0-00010101000000-000000000000 // indirect
	hyperball.com/utils v0.0.0-00010101000000-000000000000 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace hyperball.com/database => ../database

replace hyperball.com/utils => ../utils

replace hyperball.com/config => ../config

replace hyperball.com/auth => ../auth

replace hyperball.com/game => ../game

replace hyperball.com/sms => ../sms

replace hyperball.com/reports => ../reports
