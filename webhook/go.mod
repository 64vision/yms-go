module zerasuite/webhook

go 1.26.4

require (
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	gollux/account v0.0.0-00010101000000-000000000000
	gollux/auth v0.0.0-00010101000000-000000000000
	gollux/dbconfig v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/text v0.1.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	gollux/email v0.0.0-00010101000000-000000000000 // indirect
	gollux/sms v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	mellium.im/sasl v0.3.2 // indirect
	zerasuite/bookings v0.0.0-00010101000000-000000000000 // indirect
	zerasuite/shippinglines v0.0.0-00010101000000-000000000000 // indirect
	zerasuite/yards v0.0.0-00010101000000-000000000000 // indirect
)

replace gollux/utils => ../utils

replace gollux/auth => ../auth

replace gollux/account => ../account

replace gollux/sms => ../sms

replace zerasuite/bookings => ../bookings

replace zerasuite/shippinglines => ../shippinglines

replace zerasuite/yards => ../yards

replace gollux/email => ../email

replace gollux/dbconfig => ../dbconfig
