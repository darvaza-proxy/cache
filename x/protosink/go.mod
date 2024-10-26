module darvaza.org/cache/x/protosink

go 1.21

toolchain go1.22.6

require (
	darvaza.org/cache v0.3.0
	darvaza.org/core v0.15.2
	darvaza.org/slog v0.5.14 // indirect
)

require google.golang.org/protobuf v1.35.1

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace darvaza.org/cache => ../../
