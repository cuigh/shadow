package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
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
		//TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	// set http proxy
	if cfg.Proxy != "" {
		u, err := url.Parse("http://" + cfg.Proxy)
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

	if r.Method == http.MethodConnect {
		s.handleHTTPS(w, r)
		return
	} else {
		s.handleHTTP(w, r)
	}
}

func (s *Shadow) handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := s.tp.RoundTrip(r)
	if err != nil {
		log.Println("RoundTrip failed: ", err)
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

func (s *Shadow) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	h, ok := w.(http.Hijacker)
	if !ok {
		log.Println("HTTP server does not support hijacking")
		return
	}

	client, _, err := h.Hijack()
	if err != nil {
		log.Println("Cannot hijack connection: ", err)
		return
	}

	var remote net.Conn
	if s.cfg.Proxy == "" {
		_, err = client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		if err != nil {
			client.Close()
			log.Println("Write response to client failed: ", err)
			return
		}

		remote, err = net.Dial("tcp", r.URL.Host)
		if err != nil {
			client.Close()
			log.Println("Dial remote host failed: ", err)
			return
		}
	} else {
		remote, err = net.Dial("tcp", s.cfg.Proxy)
		if err != nil {
			client.Close()
			log.Println("Dial proxy failed: ", err)
			return
		}

		err = r.WriteProxy(remote)
		if err != nil {
			client.Close()
			remote.Close()
			log.Println("WriteProxy failed: ", err)
			return
		}
	}

	go s.writeToRemote(client, remote)
	go s.readFromRemote(client, remote)
}

func (s *Shadow) writeToRemote(client, remote net.Conn) {
	_, err := io.Copy(remote, client)
	if err != nil && s.cfg.Verbose {
		log.Println("Write to remote failed: ", err)
	}

	client.Close()
	remote.Close()
}

func (s *Shadow) readFromRemote(client, remote net.Conn) {
	_, err := io.Copy(client, remote)
	if err != nil && s.cfg.Verbose {
		log.Println("Read from remote failed: ", err)
	}

	client.Close()
	remote.Close()
}
