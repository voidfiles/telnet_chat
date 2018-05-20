package server

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// Config contains configuration information for the telenet chat.
type Config struct {
	TelnetPort string
	TelnetIP   string
	HTTPPort   string
	HTTPIP     string
	LogPath    string
}

// ParseConf will parse a Reader for config information.
func ParseConf(rawConfig io.Reader) (Config, error) {
	var config Config
	if _, err := toml.DecodeReader(rawConfig, &config); err != nil {
		return config, fmt.Errorf("config failed: %s", err)
	}

	return config, nil
}

// ReadConfig parses a file at a path for configuration information.
func ReadConfig(configfile string) (Config, error) {
	_, err := os.Stat(configfile)
	if err != nil {
		return Config{}, fmt.Errorf("Config file is missing: %s %s", configfile, err)
	}

	f, err := os.Open(configfile)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to open config file: %s %s", configfile, err)
	}

	return ParseConf(f)
}
