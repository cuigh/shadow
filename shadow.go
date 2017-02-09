package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Address string `json:"addr"`
	Proxy   string `json:"proxy"`
	Timeout int64  `json:"timeout"`
	Verbose bool   `json:"verbose"`
}

func NewConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

type Shadow struct {
	cfg *Config
	tp  *http.Transport
}

func NewShadow(cfg *Config) (*Shadow, error) {
	tp := &http.Transport{
		ResponseHeaderTimeout: time.Duration(cfg.Timeout) * time.Millisecond,
	}
	if cfg.Proxy != "" {
		u, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, err
		}
		tp.Proxy = func(*http.Request) (*url.URL, error) {
			return u, nil
		}
	}

	return &Shadow{
		cfg: cfg,
		tp:  tp,
	}, nil
}

func (s *Shadow) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.cfg.Verbose {
		log.Println(r.Method, ">", r.URL)
	}

	resp, err := s.tp.RoundTrip(r)
	if err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
