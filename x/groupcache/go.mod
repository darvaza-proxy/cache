module darvaza.org/cache/x/groupcache

go 1.21

toolchain go1.22.6

require (
	darvaza.org/cache v0.2.10
	darvaza.org/core v0.15.2
	darvaza.org/slog v0.5.14
)

require github.com/mailgun/groupcache/v2 v2.5.0

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

replace darvaza.org/cache => ../../
