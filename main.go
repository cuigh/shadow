package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	Version = "v0.6.1"
)

func main() {
	cfg := loadConfig()
	s, err := NewShadow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("shadow started at", s.cfg.Address)
	if cfg.Proxy != "" {
		log.Println("using parent proxy", s.cfg.Proxy)
	}
	http.ListenAndServe(cfg.Address, s)
}

func loadConfig() *Config {
	var (
		cfg         = new(Config)
		err         error
		config      string
		address     string
		proxy       string
		dialTimeout int64
		readTimeout int64
		verbose     bool
		version     bool
	)

	fs := NewFlags(`shadow is a simple and fast http/https proxy written by Go.

    shadow -a[-addr] :7777
    shadow -a[-addr] :7777 -p[-proxy] :7778 -verbose`)
	fs.Bool(&version, "version", "v", false, "Prints the shadow version")
	fs.String(&config, "config", "c", "", "Configuration file")
	fs.String(&address, "addr", "a", ":1080", "Listen address")
	fs.String(&proxy, "proxy", "p", "", "Parent proxy address")
	fs.Int64(&dialTimeout, "dial_timeout", "dt", 30000, "Timeout for dial proxy in milliseconds")
	fs.Int64(&readTimeout, "read_timeout", "rt", 30000, "Timeout for read response header in milliseconds")
	fs.Bool(&verbose, "verbose", "", false, "Verbose output")
	fs.Parse()

	if version {
		fmt.Println("shadow " + Version)
		os.Exit(0)
	}

	if config != "" {
		cfg, err = NewConfig(config)
		if err != nil {
			log.Fatalf("Invalid config arg '%s': %s", config, err)
		}
	} else {
		cfg = &Config{
			Address:     address,
			Proxy:       proxy,
			DialTimeout: dialTimeout,
			ReadTimeout: readTimeout,
			Verbose:     verbose,
		}
	}
	return cfg
}
