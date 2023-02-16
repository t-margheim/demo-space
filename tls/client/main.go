package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func main() {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: &tr,
		Timeout:   time.Second,
	}
	httpResp, err := client.Get("http://localhost:8100")
	if err != nil {
		fmt.Println("http: failed to get response:", err)
		return
	}
	fmt.Println("http: request completed:", httpResp.Status)
	httpsResp, err := client.Get("https://localhost:443/v1/metrics")
	if err != nil {
		fmt.Println("https: failed to get response:", err)
		return
	}
	fmt.Println("https: request completed:", httpsResp.Status)
}
