package app

type Config struct {
	Port        int
	Host        string
	Master_host string
	Master_port int
}

func NewConfig(host string, port int, replica_host string, replica_port int) *Config {
	return &Config{
		Host:        host,
		Port:        port,
		Master_host: replica_host,
		Master_port: replica_port,
	}
}
