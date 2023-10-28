package main

import (
	"discomfort_index/api"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var apiPath, port, ipAddress string

type DiscomfortIndexes []DiscomfortIndex

type DiscomfortIndex struct {
	City        string  `json:"city"`
	Feeling     string  `json:"feeling"`
	Humidity    int     `json:"humidity"`
	Index       float64 `json:"index"`
	Temperature float64 `json:"temperature"`
}

func init() {
	// Initialize flag
	// 引数は変数のポインタ(メモリのアドレス値)、フラグの名前、デフォルト値、使い方の説明
	flag.StringVar(&apiPath, "apipath", "/di", "API path")
	flag.StringVar(&port, "port", "18080", "Port")
	flag.StringVar(&ipAddress, "ipaddress", "", "IP Address")
}

func main() {
	flag.Parse()

	handleRequests(apiPath, port)
}

func handleRequests(apiPath, port string) {
	http.HandleFunc(apiPath, entryPoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	// If IP Address is not specified in arg
	if ipAddress == "" {
		// Get IP address
		ip, err := getIP(r)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("No valid ip"))
		}

		if ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192") || strings.HasPrefix(ip, "172") {
			url := "http://ifconfig.me"
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			ipAddress = string(body)
		} else {
			ipAddress = ip
		}
	}
	fmt.Printf("IP Address: %s\n", ipAddress)

	apiKey := os.Getenv("WEATHERAPI_API_KEY")

	rtn := api.EntryPoint(ipAddress, apiKey)

	var discomfortIndex DiscomfortIndex
	discomfortIndex.Feeling = rtn.Feeling
	discomfortIndex.Index = rtn.Index
	discomfortIndex.City = rtn.City
	discomfortIndex.Humidity = rtn.Humidity
	discomfortIndex.Temperature = rtn.Temperature

	json.NewEncoder(w).Encode(discomfortIndex)
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
