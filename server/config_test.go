package server_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/telnet_chat/server"
)

func TestParseConf(t *testing.T) {
	type conf struct {
		telnetPort string
		telnetIp   string
		httpPort   string
		httpIp     string
		logPath    string
	}

	var parseTests = []struct {
		data     string
		errorMsg string
		config   conf
	}{
		{
			`telnetIp = 0.0.0.1`,
			"config failed: Near line 1 (last key parsed 'telnetIp'): Invalid float value: \"0.0.0.1\"",
			conf{},
		}, {
			`telnetPort = "9002"
      telnetIp = "0.0.0.2"
			httpPort = "90020"
      httpIp = "0.0.0.2"
      logPath = "/tmp/test/awesome/2"`,
			"",
			conf{
				"9002", "0.0.0.2", "90020", "0.0.0.2", "/tmp/test/awesome/2",
			},
		}, {
			`telnetPort = "9003"`,
			"",
			conf{
				"9003", "", "", "", "",
			},
		},
	}

	for _, pt := range parseTests {
		data := bytes.NewReader([]byte(pt.data))
		config, err := server.ParseConf(data)
		if pt.errorMsg != "" {
			assert.EqualError(t, err, pt.errorMsg)
		} else {
			assert.Equal(t, pt.config.telnetPort, config.TelnetPort)
			assert.Equal(t, pt.config.telnetIp, config.TelnetIP)
			assert.Equal(t, pt.config.logPath, config.LogPath)
		}

	}
}
