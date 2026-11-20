// launcher renders an xray config from environment variables and execs xray.
// It is the image ENTRYPOINT and must itself be a static PIE binary,
// like the xray binary it launches.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type client struct {
	ID string `json:"id"`
}

type vlessSettings struct {
	Clients    []client `json:"clients"`
	Decryption string   `json:"decryption"`
}

type wsSettings struct {
	Path string `json:"path"`
}

type streamSettings struct {
	Network    string     `json:"network"`
	WSSettings wsSettings `json:"wsSettings"`
}

type inbound struct {
	Port           int            `json:"port"`
	Listen         string         `json:"listen"`
	Protocol       string         `json:"protocol"`
	Settings       vlessSettings  `json:"settings"`
	StreamSettings streamSettings `json:"streamSettings"`
}

type outbound struct {
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

type config struct {
	Log struct {
		LogLevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds  []inbound  `json:"inbounds"`
	Outbounds []outbound `json:"outbounds"`
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	uuid := os.Getenv("UUID")
	if uuid == "" {
		fmt.Fprintln(os.Stderr, "launcher: UUID environment variable is required")
		os.Exit(1)
	}
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "launcher: invalid PORT %q\n", v)
			os.Exit(1)
		}
		port = p
	}
	wsPath := envOr("WS_PATH", "/ws")

	var cfg config
	cfg.Log.LogLevel = "warning"
	cfg.Inbounds = []inbound{{
		Port:     port,
		Listen:   "0.0.0.0",
		Protocol: "vless",
		Settings: vlessSettings{
			Clients:    []client{{ID: uuid}},
			Decryption: "none",
		},
		StreamSettings: streamSettings{
			Network:    "ws",
			WSSettings: wsSettings{Path: wsPath},
		},
	}}
	cfg.Outbounds = []outbound{{Protocol: "freedom", Tag: "direct"}}

	data, err := json.Marshal(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "launcher: marshal config:", err)
		os.Exit(1)
	}
	if err := os.WriteFile("/config.json", data, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "launcher: write config:", err)
		os.Exit(1)
	}
	fmt.Printf("launcher: starting xray on 0.0.0.0:%d, ws path %s\n", port, wsPath)
	if err := syscall.Exec("/xray", []string{"xray", "run", "-c", "/config.json"}, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, "launcher: exec xray:", err)
		os.Exit(1)
	}
}
