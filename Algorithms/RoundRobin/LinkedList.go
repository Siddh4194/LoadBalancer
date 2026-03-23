package roundrobin

import (
	"fmt"
	"net/url"
)

type Node struct {
	server   string
	proxyUrl url.URL
	next     *Node
}

type CircularLinkedList struct {
	head *Node
	tail *Node
}

func (cll *CircularLinkedList) Add(server string,proxyUrl url.URL) {
	newNode := &Node{server: server, proxyUrl: proxyUrl}

	if cll.head == nil {
		cll.head = newNode
		cll.tail = newNode
		newNode.next = cll.head
		return
	}
	cll.tail.next = newNode
	newNode.next = cll.head
	cll.tail = newNode
}

func (cll *CircularLinkedList) Remove(server string) {
	if cll.head == nil {
		return
	}

	prev := cll.tail
	curr := cll.head

	for {
		if curr.server == server {

			if curr == cll.head && curr == cll.tail {
				cll.head = nil
				cll.tail = nil
				return
			}

			if curr == cll.head {
				cll.head = curr.next
				cll.tail.next = cll.head
			}

			if curr == cll.tail {
				cll.tail = prev
				cll.tail.next = cll.head
			}

			return
		}

		prev = curr
		curr = curr.next

		if curr == cll.head {
			break
		}
	}
}

// load balancer next iteration for the operations
type LoadBalancer struct {
	Servers *CircularLinkedList
	current *Node
}

func (lb *LoadBalancer) GetNextServer() (*url.URL, error) {
	if lb.Servers.head == nil {
		return nil, fmt.Errorf("No servers available")
	}

	if lb.current == nil {
		lb.current = lb.Servers.head
		return &lb.current.proxyUrl, nil
	}

	lb.current = lb.current.next
	return &lb.current.proxyUrl,nil
}