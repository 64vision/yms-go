module zerasuite/bookings

go 1.26.4

replace gollux/utils => ../utils

require (
	github.com/go-pg/pg v8.0.7+incompatible
	gollux/dbconfig v0.0.0-00010101000000-000000000000
	gollux/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace gollux/dbconfig => ../dbconfig
