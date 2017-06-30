package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	Version = "v0.6"
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
		cfg               = new(Config)
		err               error
		config            string
		address           string
		proxy             string
		timeout           int64
		connectionTimeout int64
		deadlineTimeout   int64
		verbose           bool
		version           bool
	)

	fs := NewFlags(`shadow is a simple and fast http/https proxy written by Go.

    shadow -a[-addr] :7777
    shadow -a[-addr] :7777 -p[-proxy] :7778 -verbose`)
	fs.Bool(&version, "version", "v", false, "Prints the shadow version")
	fs.String(&config, "config", "c", "", "Configuration file")
	fs.String(&address, "addr", "a", ":1080", "Listen address")
	fs.String(&proxy, "proxy", "p", "", "Parent proxy address")
	fs.Int64(&timeout, "timeout", "t", 30000, "Timeout for response header of milliseconds")
	fs.Int64(&connectionTimeout, "connection timeout", "ct", 5000, "Timeout for Dial connection of milliseconds")
	fs.Int64(&deadlineTimeout, "deadline timeout", "dt", 10000, "Timeout for recevie after send of milliseconds")
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
			Address:           address,
			Proxy:             proxy,
			Timeout:           timeout,
			Verbose:           verbose,
			ConnectionTimeout: connectionTimeout,
			DeadlineTimeout:   deadlineTimeout,
		}
	}
	return cfg
}
