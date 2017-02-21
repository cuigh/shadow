package main

import (
	"flag"
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
		cfg     = new(Config)
		err     error
		config  string
		address string
		proxy   string
		timeout int64
		verbose bool
		version bool
	)
	registerBoolArg(&version, "version", "v", false, "Prints the shadow version")
	registerStringArg(&config, "config", "c", "", "Configuration file")
	registerStringArg(&address, "addr", "a", ":1080", "Listen address")
	registerStringArg(&proxy, "proxy", "p", "", "Parent proxy address")
	registerInt64Arg(&timeout, "timeout", "t", 30000, "Timeout for response header of milliseconds")
	registerBoolArg(&verbose, "verbose", "", false, "Verbose output")
	flag.Parse()

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
			Address: address,
			Proxy:   proxy,
			Timeout: timeout,
			Verbose: verbose,
		}
	}

	return cfg
}

func registerStringArg(p *string, name, shortName, value, usage string) {
	flag.StringVar(p, name, value, usage)
	if shortName != "" {
		flag.StringVar(p, shortName, value, usage+" (shorthand)")
	}
}

func registerInt64Arg(p *int64, name, shortName string, value int64, usage string) {
	flag.Int64Var(p, name, value, usage)
	if shortName != "" {
		flag.Int64Var(p, shortName, value, usage+" (shorthand)")
	}
}

func registerBoolArg(p *bool, name, shortName string, value bool, usage string) {
	flag.BoolVar(p, name, value, usage)
	if shortName != "" {
		flag.BoolVar(p, shortName, value, usage+" (shorthand)")
	}
}
