module github.com/etcd-monitor/pingap-go

go 1.21

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.1
	github.com/prometheus/client_golang v1.18.0
	go.etcd.io/etcd/client/v3 v3.5.11
	go.opentelemetry.io/otel v1.22.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.22.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.18.0
	golang.org/x/net v0.20.0
	golang.org/x/time v0.5.0
)
