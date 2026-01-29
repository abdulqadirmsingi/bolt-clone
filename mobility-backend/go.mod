module github.com/zeng/mobility-backend

go 1.23

require (
	github.com/jackc/pgx/v5 v5.7.2
	github.com/redis/go-redis/v9 v9.4.0
	github.com/uber/h3-go/v4 v4.1.0
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.62.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

// Fix: revision 7065a208a1fe was removed from jackc/pgservicefile; use a valid pseudo-version.
replace github.com/jackc/pgservicefile => github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761
