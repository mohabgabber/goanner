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

func validate(lrange, hrange int, addr string) {
	if lrange > hrange || hrange > 65535 || lrange < 1 {
		flag.PrintDefaults()
		log.Fatal("Error: Invalid Port Range Formatting")
	}
	if len(addr) == 0 {
		flag.PrintDefaults()
		log.Fatal("Error: Invalid Target Address")
	}
}

func splitrange(ports *[]int, prange, addr string) error {
	ranges := strings.Split(prange, "-")
	lrange, _ := strconv.Atoi(ranges[0])
	hrange, _ := strconv.Atoi(ranges[1])
	fmt.Println(lrange, hrange, addr)
	validate(lrange, hrange, addr)
	for ; lrange <= hrange; lrange++ {
		*ports = append(*ports, lrange)
	}
	return nil
}

func main() {

	// Target Address And Port Range (By The User)
	var addr string
	var prange string
	flag.StringVar(&addr, "t", "127.0.0.1", "Target Address (Domain/IP)")
	flag.StringVar(&prange, "p", "1-1024", "Port Range (From-To)")
	flag.Parse()

	// Stores The Port Range
	ports := []int{}
	splitrange(&ports, prange, addr)

	// Stores The Open Ports
	var open []int

	// Stores Open Ports
	results := make(chan int)

	// Limits Execution
	capacity := make(chan int, 100)

	// Starting Of The Port Scan
	log.Printf("Scanning The Port Range: %s \t On Target: %s", prange, addr)
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

	// Printing The Results
	log.Printf("Discovered %d Open Ports", len(open))
	for _, port := range open {
		log.Printf("Port %d is open\n", port)
	}
}
