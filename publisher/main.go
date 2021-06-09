package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itzmeanjan/gClient/utils"
	"github.com/itzmeanjan/pub0sub/publisher"
)

var (
	proto  = "tcp"
	addr   = flag.String("addr", utils.GetAddr(), "Connect to address")
	port   = flag.Uint64("port", utils.GetPort(), "Connect to port")
	client = flag.Uint64("client", utils.GetClientCount(), "#-of concurrent publishers to use")
	repeat = flag.Uint64("repeat", utils.GetRepeatCount(), "Repeat publish ( = 0 :-> infinite )")
	delay  = flag.Duration("delay", utils.GetDelay(), "Gap between subsequent message publish")
	out    = flag.Bool("out", utils.GetLoggingPreference(), "Persist publisher log")
	topics utils.TopicList
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

	ctx, cancel := context.WithCancel(context.Background())
	fullAddr := fmt.Sprintf("%s:%d", *addr, *port)
	pubs := make([]*publisher.Publisher, 0, *client)

	var logHandles []*os.File
	if *out {
		logHandles = make([]*os.File, 0, *client)
	}

	for i := 0; i < int(*client); i++ {
		pub, err := publisher.New(ctx, proto, fullAddr)
		if err != nil {
			log.Printf("[gClient] Error : %s\n", err.Error())
			return
		}

		pubs = append(pubs, pub)
		if *out {
			fd, err := os.OpenFile(fmt.Sprintf("log_pub_%d.csv", i), os.O_CREATE|os.O_RDWR, 0x1b6)
			if err != nil {
				log.Printf("[gClient] Error : %s\n", err.Error())
				return
			}

			defer func() {
				if err := fd.Close(); err != nil {
					log.Printf("[gClient] Error : %s\n", err.Error())
				}
			}()

			logHandles = append(logHandles, fd)
		}
	}

	publishers := utils.Publishers{Handles: pubs, Logs: logHandles, Buffer: new(bytes.Buffer)}
	log.Printf("[gClient] Connected to %s [ %d client(s) ] ✅\n", fullAddr, *client)

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

				start := time.Now()
				if err := publishers.PublishMsg(topics); err != nil {
					log.Printf("[gClient] Error : %s\n", err.Error())
					break OUT_1
				}

				log.Printf("[gClient] Publish iteration : %d [ in %s ]\n", i+1, time.Since(start))
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

				start := time.Now()
				if err := publishers.PublishMsg(topics); err != nil {
					log.Printf("[gClient] Error : %s\n", err.Error())
					break OUT_2
				}

				log.Printf("[gClient] Publish iteration : %d [ in %s ]\n", i+1, time.Since(start))
			}

		}
	}

	cancel()
	<-time.After(time.Second)

	log.Printf("[gClient] Graceful shutdown\n")
}
