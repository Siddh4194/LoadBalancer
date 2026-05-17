package roundrobin

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Node struct {
	Server    string
	ProxyUrl  url.URL
	Url       string
	Proxy     httputil.ReverseProxy
	Next      *Node
	FailCount int
	Healthy   bool
}

type CircularLinkedList struct {
	head *Node
	tail *Node
}

func (cll *CircularLinkedList) Add(server string, proxyUrl url.URL, url string) {
	newNode := &Node{Server: server, ProxyUrl: proxyUrl, Url: url, Proxy: *httputil.NewSingleHostReverseProxy(&proxyUrl), FailCount: 0, Healthy: true}

	if cll.head == nil {
		cll.head = newNode
		cll.tail = newNode
		newNode.Next = cll.head
		return
	}
	cll.tail.Next = newNode
	newNode.Next = cll.head
	cll.tail = newNode
}

func (cll *CircularLinkedList) Remove(server string) {
	if cll.head == nil {
		return
	}

	prev := cll.tail
	curr := cll.head

	for {
		if curr.Server == server {

			if curr == cll.head && curr == cll.tail {
				cll.head = nil
				cll.tail = nil
				return
			}

			if curr == cll.head {
				cll.head = curr.Next
				cll.tail.Next = cll.head
			}

			if curr == cll.tail {
				cll.tail = prev
				cll.tail.Next = cll.head
			}

			return
		}

		prev = curr
		curr = curr.Next

		if curr == cll.head {
			break
		}
	}
}

// load balancer next iteration for the operations
type LoadBalancer struct {
	Servers *CircularLinkedList
	current *Node
	mu      sync.RWMutex
}

func (lb *LoadBalancer) GetNextServer() (*httputil.ReverseProxy, error) {
	if lb.Servers.head == nil {
		return nil, fmt.Errorf("No servers available")
	}

	if lb.current == nil {
		lb.current = lb.Servers.head
	} else {
		lb.current = lb.current.Next
	}

	start := lb.current

	for {
		if lb.current.Healthy {
			fmt.Printf("Selected server: %s (Healthy: %t)\n", lb.current.Server, lb.current.Healthy)
			return &lb.current.Proxy, nil
		}

		lb.current = lb.current.Next

		if lb.current == start {
			break
		}
	}

	return nil, fmt.Errorf("Couldn't found healthy servers")
}

func (lb *LoadBalancer) ForEachServer(fn func(server *Node)) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	current := lb.Servers.head

	if current == nil {
		return
	}

	for {
		fn(current)
		current = current.Next

		if current == lb.Servers.head {
			break
		}
	}
}
