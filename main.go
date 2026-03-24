package main

import (
	roundrobin "LoadBalancer/Algorithms/RoundRobin"
	"LoadBalancer/lib"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)


func roundRobinAlGoToSelectServer(server int) *httputil.ReverseProxy {
	server1, err1 := url.Parse("http://localhost:3000")
	server2, err2 := url.Parse("http://localhost:3001")
	server3, err3 := url.Parse("http://localhost:3002")

	if err1 != nil ||  err2 != nil || err3 != nil {
		log.Fatal("Error parsing server URLs")
	} 
	var selectedServer *url.URL

	if server == 1 {
		selectedServer = server1
	} else if server == 2 {
		selectedServer = server2
	} else if server == 3 {
		selectedServer = server3
	}

	proxy := httputil.NewSingleHostReverseProxy(selectedServer)
	return proxy
}

func hashIP(ip string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(ip))      // get hash
	return h.Sum32()
}

func parseUrl(urlStr string) *url.URL {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}
	return parsedUrl
}

func initiateRoundRobin() *roundrobin.LoadBalancer {
	cll := &roundrobin.CircularLinkedList{}
	cll.Add("1", *parseUrl("http://localhost:3000"),"http://localhost:3000")
	cll.Add("2", *parseUrl("http://localhost:3001"),"http://localhost:3001")
	cll.Add("3", *parseUrl("http://localhost:3002"),"http://localhost:3002")

	lb := &roundrobin.LoadBalancer{Servers: cll}

	return lb;
}

func main() {
	proxyUrl, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	// ip hashing
	ipHashing := false
	roundRobin := true
	lb := initiateRoundRobin()

	lb.ForEachServer(func(server *roundrobin.Node) {
		go lib.HealthCheck(server)
	})

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	fmt.Println(proxy)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.URL.Path)
		if(ipHashing){
			hashedIp := hashIP(r.RemoteAddr)
			fmt.Printf("Hashed IP: %d\n", hashedIp % 3)
			roundRobinAlGoToSelectServer(int(hashedIp % 3)).ServeHTTP(w, r)
		}
		if(roundRobin){
			server,err := lb.GetNextServer()
			if err != nil {
				http.Error(w, "No servers available", http.StatusServiceUnavailable)
				return
			}
			server.ServeHTTP(w,r)
		}
		// fmt.Fprintf(w, "Hello, World from load balancer!")
	})


	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}