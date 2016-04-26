package utils

import (
	"net"
	"net/http"
	"time"
)

func NewHTTPTimeoutClient(timeout_milliseconds int64) *http.Client {
	td := time.Duration(time.Duration(timeout_milliseconds) * time.Millisecond)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, td)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(td))
				return conn, nil
			},
			ResponseHeaderTimeout: td,
		},
	}
	return client
}
