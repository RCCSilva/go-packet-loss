package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func TestTCP() {
	ctx := context.Background()

	rcListenTCP := ListenTCP(ctx)
	rListenTCP := <-rcListenTCP
	if !rListenTCP.isOk() {
		log.Fatalf("failed to start tcp listener: %v", rListenTCP.Error)
	}
	log.Print("tcp listener started!")
	rWriteTCP := WriteTCP(ctx, rListenTCP.Address)
	if r := <-rWriteTCP; !r.isOk() {
		log.Fatalf("failed to start tcp writer: %v", r.Error)
	}
	log.Printf("tcp writer started!")
	<-ctx.Done()
}

func WriteTCP(ctx context.Context, addr string) <-chan Result {
	result := make(chan Result)

	go func() {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			result <- Result{Error: err}
			return
		}

		defer c.Close()
		result <- Result{}
		counter := 1001

		for {
			time.Sleep(1 * time.Second)
			_, err := c.Write([]byte(fmt.Sprint(counter)))
			if err != nil {
				log.Printf("failed to write %v", err)
			}
			counter++
		}
	}()

	return result
}

func ListenTCP(ctx context.Context) <-chan Result {
	result := make(chan Result)

	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:0")

		if err != nil {
			result <- Result{Error: err}
			return
		}

		defer listener.Close()

		result <- Result{
			Address: listener.Addr().String(),
		}

		for {
			c, err := listener.Accept()

			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			go func() {
				defer c.Close()
				for {
					b := make([]byte, 16)
					n, err := c.Read(b)
					if err != nil {
						continue
					}
					log.Printf("tcp: %q", b[:n])
				}
			}()
		}
	}()

	return result
}
