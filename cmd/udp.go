package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func TestUDP() {
	ctx := context.Background()

	rcListenUDP := ListenUDP(ctx)
	rListenUDP := <-rcListenUDP
	if !rListenUDP.isOk() {
		log.Fatalf("udp: failed to start udp listener: %v", rListenUDP.Error)
	}
	log.Print("udp: listener started!")
	rWriteUDP := WriteUDP(ctx, rListenUDP.Address)
	if r := <-rWriteUDP; !r.isOk() {
		log.Fatalf("udp: failed to start writer: %v", r.Error)
	}
	log.Printf("udp: writer started!")
	<-ctx.Done()
}

func WriteUDP(ctx context.Context, addr string) <-chan Result {
	result := make(chan Result)

	go func() {
		c, err := net.Dial("udp", addr)
		if err != nil {
			result <- Result{Error: err}
			return
		}
		defer c.Close()
		log.Print("udp: writer started!")
		counter := 2001
		for {
			c.Write([]byte(fmt.Sprint(counter)))
			time.Sleep(1 * time.Second)
			counter++
		}
	}()

	return result
}

func ListenUDP(ctx context.Context) <-chan Result {
	result := make(chan Result)

	go func() {
		addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1")}
		l, err := net.ListenUDP("udp", addr)
		if err != nil {
			result <- Result{Error: err}
			return
		}
		defer l.Close()

		result <- Result{
			Address: l.LocalAddr().String(),
		}

		for {
			b := make([]byte, 16)
			rlen, _, err := l.ReadFromUDP(b)
			if err != nil {
				continue
			}
			log.Printf("udp: %q", b[:rlen])
		}
	}()

	return result
}
