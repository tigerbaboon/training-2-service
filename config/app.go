package config

// Config is a struct that contains all the configuration of the application.
type Config struct {
	Database Database

	AppName string
	AppKey  string
	AppEnv  string
	Debug   bool

	Port           int
	HttpPrefork    bool
	HttpJsonNaming string

	SslCaPath      string
	SslPrivatePath string
	SslCertPath    string

	OtelEnable            bool
	OtelMetricMode        string
	OtelTraceMode         string
	OtelLogMode           string
	OtelCollectorEndpoint string
	OtelTraceRatio        float64
	MySecret              string
}

var App = Config{
	Database: database,

	AppName: "go_app",
	Port:    8080,
	AppKey:  "secret",
	AppEnv:  "development",
	Debug:   false,

	HttpPrefork:    false,
	HttpJsonNaming: "pascal_case",

	SslCaPath:      "storage/cert/ca.pem",
	SslPrivatePath: "storage/cert/server.pem",
	SslCertPath:    "storage/cert/server-key.pem",

	OtelCollectorEndpoint: "", // localhost:4317
	OtelEnable:            false,
	OtelMetricMode:        "none", // fallback to stdout if not none and collector is not available
	OtelTraceMode:         "none", // fallback to stdout if not none and collector is not available
	OtelLogMode:           "stdout",
	OtelTraceRatio:        1,
	MySecret:              "",
}
