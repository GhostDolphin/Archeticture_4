package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"hash/fnv"

	"github.com/GhostDolphin/Architecture_4/httptools"
	"github.com/GhostDolphin/Architecture_4/signal"
)

var (
	port = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
)

var (
	timeout = time.Duration(*timeoutSec) * time.Second
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
)

func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func forward(dst string, rw http.ResponseWriter, r *http.Request) error {
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

type IsHealthy struct {
	status map[string]bool
	health func(dst string) bool
}

type Balancer struct {
	isHealthy *IsHealthy
}

func (hp *IsHealthy) CheckAll() {
	for _, server := range serversPool {
		if hp.health(server) {
			hp.status[server] = true
		} else {
			hp.status[server] = false
		}
	}
}

func (bal *Balancer) verifyServer() {
	for {
		bal.isHealthy.CheckAll()
		time.Sleep(5 * time.Second)
	}
}

func (hp *IsHealthy) AllHealthy() []string {
	var allHealthy []string
	for _, server := range serversPool {
		if hp.status[server] {
			allHealthy = append(allHealthy, server)
		}
	}
	return allHealthy
}

func hash(input string) uint32 {
	hsh := fnv.New32a()
	hsh.Write([]byte(input))
	return hsh.Sum32()
}

func (bal *Balancer) doBalancer(url string) string {
	allHealthy := bal.isHealthy.AllHealthy()

	if len(allHealthy) == 0 {
		log.Println("There are no healthy servers")
		return "There are no healthy servers"
	}
	return ""
}

func main() {
	flag.Parse()
	isHealthy := &IsHealthy{}
	isHealthy.status = map[string]bool{}
	isHealthy.health = health
	balance := &Balancer{}
	balance.isHealthy = isHealthy

	go balance.verifyServer()

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		server := balance.doBalancer(r.URL.Path)
		forward(server, rw, r)
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
