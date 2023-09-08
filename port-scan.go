package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

const (
	porterror = "Invalid Port Format"
)

type ConnDetail struct {
	addr   string
	prange string
	lrange int
	hrange int
	open   []int
}

func (x ConnDetail) pscan(capa, results chan int) {
	for p := range capa {
		finad := fmt.Sprintf(x.addr+":"+"%d", p)
		log.Println("scanning port", p)
		conn, err := net.Dial("tcp", finad)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
func (x ConnDetail) splitrange(ports *[]int) error {
	ranges := strings.Split(x.prange, "-")
	x.lrange, _ = strconv.Atoi(ranges[0])
	x.hrange, _ = strconv.Atoi(ranges[1])
	if x.lrange > x.hrange || x.lrange < 1 || x.hrange > 65535 {
		return errors.New(porterror)
	}
	for ; x.lrange <= x.hrange; x.lrange++ {
		*ports = append(*ports, x.lrange)
	}
	return nil
}
func main() {
	var details ConnDetail
	fmt.Println("Type the target address (Domain/IP):")
	fmt.Scanln(&details.addr)
	fmt.Println("The range of ports to scan (example: 100-200): ")
	fmt.Scanln(&details.prange)
	ports := []int{}
	details.splitrange(&ports)
	results := make(chan int)
	capacity := make(chan int, 100)

	for i := 0; i < cap(capacity); i++ {
		go details.pscan(capacity, results)
	}

	go func() {
		for _, p := range ports {
			capacity <- p
		}
	}()

	for i := 0; i < len(ports); i++ {
		port := <-results
		if port != 0 {
			details.open = append(details.open, port)
		}
	}
	close(capacity)
	close(results)
	sort.Ints(details.open)
	for _, port := range details.open {
		log.Printf("Port %d is open\n", port)
	}
}
