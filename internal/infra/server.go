package infra

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// InitDefaultServer initializes and returns a default HTTP server using configuration from viper.
func InitDefaultServer() *http.Server {
	host := viper.GetString("server.host")
	if host == "" {
		host = "localhost"
	}

	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080
	}

	addr := host + ":" + fmt.Sprintf("%d", port)
	return &http.Server{
		Addr: addr,
	}
}
