package netlib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func HttpPost(url string, data string, connTimeoutMs int, serveTimeoutMs int, traceid string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	body := strings.NewReader(data)
	reqest, _ := http.NewRequest("POST", url, body)
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("didi-header-rid", traceid)
	response, err := client.Do(reqest)
	if err != nil {
		err = errors.New(fmt.Sprintf("http failed, POST url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("http status code error, POST url:%s, code:%d", url, response.StatusCode))
		return nil, err
	}

	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("cannot read http response, POST url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	return res_body, nil
}

func HttpGet(url string, connTimeoutMs int, serveTimeoutMs int, traceid string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("didi-header-rid", traceid)
	response, err := client.Do(reqest)
	if err != nil {
		err = errors.New(fmt.Sprintf("http failed, GET url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("http status code error, GET url:%s, code:%d", url, response.StatusCode))
		return nil, err
	}

	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("cannot read http response, GET url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	return res_body, nil
}
