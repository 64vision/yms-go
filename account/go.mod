module gollux/account

go 1.26.4

replace gollux/utils => ../utils

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-pg/pg v8.0.7+incompatible
	gollux/dbconfig v0.0.0-00010101000000-000000000000
	gollux/sms v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
	zerasuite/bookings v0.0.0-00010101000000-000000000000
	zerasuite/shippinglines v0.0.0-00010101000000-000000000000
	zerasuite/yards v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/text v0.1.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/sms => ../sms

replace gollux/game => ../game

replace gollux/reports => ../reports

replace zerasuite/yards => ../yards

replace zerasuite/shippinglines => ../shippinglines

replace zerasuite/bookings => ../bookings

replace gollux/dbconfig => ../dbconfig
