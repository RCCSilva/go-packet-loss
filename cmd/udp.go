package main

import (
	"bufio"
	"context"
	"encoding/binary"
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
	rWriteUDP := WriteUDP(ctx, rListenUDP.Address)
	if r := <-rWriteUDP; !r.isOk() {
		log.Fatalf("udp: failed to start writer: %v", r.Error)
	}
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
		log.Print("udp writer: writer started!")
		writer := bufio.NewWriter(c)
		counter := 2001
		for {
			// log.Printf("udp writer: writing %v", counter)
			binary.Write(writer, binary.LittleEndian, int16(counter))
			writer.Flush()
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

		var current, last int16
		log.Print("udp listener: started!")

		for {
			binary.Read(l, binary.LittleEndian, &current)
			if last > 0 && last > current {
				log.Printf("udp listener: received \"%v\" - OUT OF ORDER!", current)
			} else {

				log.Printf("udp listener: received \"%v\"", current)
			}
			last = current
		}
	}()

	return result
}
