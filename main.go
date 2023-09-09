package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

func pscan(capa, results chan int, addr string) {
	for p := range capa {
		finad := fmt.Sprintf(addr+":"+"%d", p)
		conn, err := net.Dial("tcp", finad)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
func splitrange(ports *[]int, lrange, hrange int, prange string) error {
	ranges := strings.Split(prange, "-")
	lrange, _ = strconv.Atoi(ranges[0])
	hrange, _ = strconv.Atoi(ranges[1])
	if lrange > hrange || lrange < 1 || hrange > 65535 {
		log.Fatal("Error: Wrong Port Format")
	}
	for ; lrange <= hrange; lrange++ {
		*ports = append(*ports, lrange)
	}
	return nil
}

func main() {
	// Target Address And Port Range (By The User)
	var addr string
	var prange string
	flag.StringVar(&addr, "t", "127.0.0.1", "Target Address")
	flag.StringVar(&prange, "p", "1-1024", "Port Range")
	flag.Parse()
	log.Printf("Scanning The Port Range: %s \t On Target: %s", prange, addr)
	// Stores The Port Range
	var lrange int
	var hrange int
	ports := []int{}
	splitrange(&ports, lrange, hrange, prange)
	// Stores The Open Ports
	var open []int
	// Stores Open Ports
	results := make(chan int)
	// Limits Execution
	capacity := make(chan int, 100)

	for i := 0; i < cap(capacity); i++ {
		go pscan(capacity, results, addr)
	}

	go func() {
		for _, p := range ports {
			capacity <- p
		}
	}()

	for i := 0; i < len(ports); i++ {
		port := <-results
		if port != 0 {
			open = append(open, port)
		}
	}
	close(capacity)
	close(results)
	sort.Ints(open)
	log.Printf("Discovered %d Open Ports", len(open))
	for _, port := range open {
		log.Printf("Port %d is open\n", port)
	}
}
