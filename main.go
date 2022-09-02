package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"go-curl/conf"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Tropicana33/common/golog"
	"github.com/pkg/errors"
)

func init() {
	flag.StringVar(&conf_file, "f", "conf.ini", "config file path!")
	flag.Parse()
}

var conf_file string
var exiting = false

var gCount = struct {
	count int
	sync.Mutex
}{
	count: 0,
}

var clientPool = struct {
	sync.Mutex
	ClientMap map[int]*http.Client
}{
	ClientMap: make(map[int]*http.Client, 30),
}

func getClient(proxyAddr string, senderID int) (client *http.Client) {
	clientPool.Lock()
	client = clientPool.ClientMap[senderID]
	if client != nil {
		clientPool.Unlock()
		return client
	}

	client = &http.Client{}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	client.Transport = &http.Transport{
		Dial: dialer.Dial,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			addr = proxyAddr
			return dialer.DialContext(ctx, network, addr)
		},
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}

	clientPool.ClientMap[senderID] = client
	clientPool.Unlock()

	return client
}

func fetchByProxy(count int, proxyAddr string, senderID int) (*http.Response, error) {
	var body bytes.Buffer

	r, err := http.NewRequest("GET", conf.GetConfig().GetBaseUrl()+strconv.Itoa(count), nil)
	fmt.Println(r.URL)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	client := getClient(proxyAddr, senderID)

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	if _, err = io.Copy(&body, resp.Body); err != nil {
		golog.Error("io.Copy err", err.Error())
		return nil, nil
	}

	return resp, nil
}

func reqSender(senderID int) {
	count := -1

	for {
		gCount.Lock()
		count = gCount.count + 1
		gCount.count = count
		gCount.Unlock()

		if count > conf.GetConfig().GetReqLimit() {
			return
		}

		_, err := fetchByProxy(count, conf.GetConfig().GetUpsServer(), senderID)
		if err != nil {
			golog.Error("fetchByProxy err", err.Error())
			continue
		}
	}
}

func handler() int {
	gCount.count = conf.GetConfig().GetStartIndex()

	for i := 0; i < conf.GetConfig().GetWorkerCount(); i++ {
		go reqSender(i)
	}

	for {
		time.Sleep(time.Millisecond * 500)

		if exiting {
			break
		}
	}

	return 0
}

func main() {
	conf.InitConfig(conf_file)
	os.Exit(handler())
}
