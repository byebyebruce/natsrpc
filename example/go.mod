module github.com/byebyebruce/natsrpc/example

go 1.17

require (
	github.com/golang/protobuf v1.5.2
	github.com/nats-io/nats-server/v2 v2.1.4
	github.com/nats-io/nats.go v1.9.1
	github.com/stretchr/testify v1.7.0
	github.com/byebyebruce/natsrpc v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/byebyebruce/natsrpc v0.0.0 => ../
