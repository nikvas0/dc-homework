module product_upload_inserter

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/jinzhu/gorm v1.9.12
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	lib v0.0.0-00010101000000-000000000000
)

replace lib => ../lib
