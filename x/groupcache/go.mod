module github.com/darvaza-proxy/cache/x/groupcache

go 1.19

replace github.com/darvaza-proxy/cache => ../../

require (
	darvaza.org/core v0.9.0
	darvaza.org/slog v0.5.0
	github.com/darvaza-proxy/cache v0.0.3
	github.com/mailgun/groupcache/v2 v2.4.2
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/protobuf v1.29.1 // indirect
)
