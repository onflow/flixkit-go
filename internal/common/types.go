package common

// Logger interface for consistent logging across packages
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

// NetworkConfig holds network connection configuration
type NetworkConfig struct {
	Name string
	Host string
	Key  string
}
