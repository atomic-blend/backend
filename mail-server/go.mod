module github.com/atomic-blend/backend/mail-server

go 1.24.5

require (
	connectrpc.com/connect v1.16.0
	github.com/atomic-blend/backend/grpc v0.0.0-00010101000000-000000000000
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6
	github.com/emersion/go-smtp v0.23.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emersion/go-message v0.18.2 // indirect
	github.com/emersion/go-msgauth v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.33.0
	github.com/streadway/amqp v1.1.0
	github.com/stretchr/testify v1.10.0
	go.mongodb.org/mongo-driver v1.17.4
	golang.org/x/sys v0.31.0 // indirect
)

replace github.com/atomic-blend/backend/grpc => ../grpc
