package ui

import (
	"crypto/tls"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type proxy struct {
	addr   string
	prefix string
}

func (p *proxy) proxy(req *http.Request) error {
	destURL, _ := url.Parse(p.addr + req.URL.Path)
	destURL.RawQuery = req.URL.RawQuery
	req.URL = destURL
	return nil
}

func New(addr, prefix string) http.Handler {
	p := &proxy{
		addr:   addr,
		prefix: prefix,
	}
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if err := p.proxy(req); err != nil {
				zap.S().Infof("Failed to proxy: %v", err)
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func NewProxy(addr, prefix string) http.Handler {
	p := &proxy{
		addr:   addr,
		prefix: prefix,
	}
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if err := p.proxy(req); err != nil {
				zap.S().Infof("Failed to proxy: %v", err)
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
