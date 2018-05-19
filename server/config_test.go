package server_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/telnet_chat/server"
)

func TestParseConf(t *testing.T) {
	type conf struct {
		port    string
		ip      string
		logPath string
	}

	var parseTests = []struct {
		data     string
		errorMsg string
		config   conf
	}{
		{
			`ip = 0.0.0.1`,
			"config failed: Near line 1 (last key parsed 'ip'): Invalid float value: \"0.0.0.1\"",
			conf{},
		}, {
			`port = "9002"
      ip = "0.0.0.2"
      logPath = "/tmp/test/awesome/2"`,
			"",
			conf{
				"9002", "0.0.0.2", "/tmp/test/awesome/2",
			},
		}, {
			`port = "9003"`,
			"",
			conf{
				"9003", "", "",
			},
		},
	}

	for _, pt := range parseTests {
		data := bytes.NewReader([]byte(pt.data))
		config, err := server.ParseConf(data)
		if pt.errorMsg != "" {
			assert.EqualError(t, err, pt.errorMsg)
		} else {
			assert.Equal(t, pt.config.port, config.Port)
			assert.Equal(t, pt.config.ip, config.IP)
			assert.Equal(t, pt.config.logPath, config.LogPath)
		}

	}
}
