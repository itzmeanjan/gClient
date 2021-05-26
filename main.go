package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itzmeanjan/pub0sub/publisher"
)

var (
	proto  = "tcp"
	addr   = flag.String("addr", "127.0.0.1", "Connect to address")
	port   = flag.Uint64("port", 13000, "Connect to port")
	client = flag.Uint64("client", 1, "#-of concurrent publishers to use")
	repeat = flag.Uint64("repeat", 1, "Repeat publish ( = 0 :-> infinite )")
	delay  = flag.Duration("delay", time.Duration(100)*time.Millisecond, "Gap between subsequent message publish")
	topics topicList
)

func main() {
	flag.Var(&topics, "topic", "Topic to publish data on")
	flag.Parse()

	if len(topics) == 0 {
		log.Printf("[gClient] Error : no topics specified\n")
		return
	}

	if *client < 1 {
		*client = 1
	}
	if *repeat < 1 {
		*repeat = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	fullAddr := fmt.Sprintf("%s:%d", *addr, *port)
	pubs := make(publishers, 0, *client)

	for i := 0; i < int(*client); i++ {

		pub, err := publisher.New(ctx, proto, fullAddr)
		if err != nil {
			log.Printf("[gClient] Error : %s\n", err.Error())
			return
		}

		pubs = append(pubs, pub)

	}

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT)

	if *repeat == 0 {
		var i uint64

	OUT_1:
		for ; ; i++ {

			select {
			case <-interruptChan:
				break OUT_1

			default:
				<-time.After(*delay)

				if err := pubs.publishMsg(topics); err != nil {
					log.Printf("[gClient] Error : %s\n", err.Error())
					break OUT_1
				}
			}

		}
	} else {
		var i uint64

	OUT_2:
		for ; i < *repeat; i++ {

			select {
			case <-interruptChan:
				break OUT_2

			default:
				<-time.After(*delay)

				if err := pubs.publishMsg(topics); err != nil {
					log.Printf("[gClient] Error : %s\n", err.Error())
					break OUT_2
				}
			}

		}
	}

	cancel()
	<-time.After(time.Second)

	log.Printf("[gClient] Graceful shutdown\n")
}
