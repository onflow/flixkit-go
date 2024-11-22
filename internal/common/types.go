package common

// Logger interface for consistent logging across packages
type Logger interface {
	Debug(string)
	Info(string)
	Error(string)
}

// NetworkConfig holds network connection configuration
type NetworkConfig struct {
	Name string
	Host string
	Key  string
}
