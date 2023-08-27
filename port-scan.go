package main

import (
	"fmt"
	"net"
	"sort"
)

func portsacq(ports, results chan int, address string) {
	for p := range ports {
		finad := fmt.Sprintf(address+":"+"%d", p)
		conn, err := net.Dial("tcp", finad)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
func main() {
	// acquiring the target address
	var address string
	fmt.Println("Type the target address (Domain/IP)")
	fmt.Scanln(&address)

	results := make(chan int)
	ports := make(chan int, 100)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go portsacq(ports, results, address)
	}

	go func() {
		for i := 0; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("The %d port is open\n", port)
	}
}
