module github.com/byebyebruce/natsrpc/example

go 1.17

replace github.com/byebyebruce/natsrpc => ../

require (
	github.com/byebyebruce/natsrpc v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats-server/v2 v2.5.0
	github.com/nats-io/nats.go v1.12.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/minio/highwayhash v1.0.1 // indirect
	github.com/nats-io/jwt/v2 v2.0.3 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
