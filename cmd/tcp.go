package main

import (
	"context"
	"encoding/binary"
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
	rWriteTCP := WriteTCP(ctx, rListenTCP.Address)
	if r := <-rWriteTCP; !r.isOk() {
		log.Fatalf("failed to start tcp writer: %v", r.Error)
	}
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
		log.Printf("tcp writer: started")

		for {
			binary.Write(c, binary.LittleEndian, int16(counter))
			time.Sleep(1 * time.Second)
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
		log.Print("tcp listener: started!")

		for {
			c, err := listener.Accept()

			if err != nil {
				log.Printf("tcp listener: failed to establish connection: %v", err)
				continue
			}

			go func() {
				defer c.Close()
				var current, last int16

				binary.Read(c, binary.LittleEndian, &current)

				for {
					err := binary.Read(c, binary.LittleEndian, &current)
					if err != nil {
						continue
					}
					if last > 0 && current < last {
						log.Printf("tcp listener: received %v - OUT OF ORDER!", current)
					} else {
						log.Printf("tcp listener: received %v", current)
					}
					last = current
				}
			}()
		}
	}()

	return result
}
