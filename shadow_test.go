package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	const (
		proxy = "http://localhost:7777"
		// reqURL = "http://headers.jsontest.com/"
		reqURL = "http://www.baidu.com/"
	)

	u, err := url.Parse(proxy)
	if err != nil {
		t.Fatal(err)
	}

	tp := &http.Transport{
		ResponseHeaderTimeout: time.Second * 10,
		Proxy: func(*http.Request) (*url.URL, error) {
			return u, nil
		},
	}
	c := http.Client{
		Transport: tp,
	}

	resp, err := c.Get(reqURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode, string(data))
}

func BenchmarkDirect(b *testing.B) {
	tp := &http.Transport{
		ResponseHeaderTimeout: time.Second * 10,
	}
	c := &http.Client{
		Transport: tp,
	}

	b.N = 100
	for i := 0; i < b.N; i++ {
		err := get(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProxy(b *testing.B) {
	u, _ := url.Parse("http://localhost:7777")
	tp := &http.Transport{
		ResponseHeaderTimeout: time.Second * 10,
		Proxy: func(*http.Request) (*url.URL, error) {
			return u, nil
		},
	}
	c := &http.Client{
		Transport: tp,
	}

	b.N = 100
	for i := 0; i < b.N; i++ {
		err := get(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func get(c *http.Client) error {
	resp, err := c.Get("http://www.baidu.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
