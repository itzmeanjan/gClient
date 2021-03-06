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
	"github.com/itzmeanjan/pub0sub/subscriber"
)

var (
	proto    = "tcp"
	addr     = flag.String("addr", utils.GetAddr(), "Connect to address")
	port     = flag.Uint64("port", utils.GetPort(), "Connect to port")
	client   = flag.Uint64("client", utils.GetClientCount(), "#-of concurrent subscribers to use")
	capacity = flag.Uint64("capacity", utils.GetCapacity(), "Pending message queue capacity")
	out      = flag.Bool("out", utils.GetLoggingPreference(), "Persist subscriber log")
	topics   utils.TopicList
)

func main() {
	flag.Var(&topics, "topic", "Topic to subscribe")
	flag.Parse()

	if len(topics) == 0 {
		log.Printf("[0sub] Error : no topics specified\n")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	fullAddr := fmt.Sprintf("%s:%d", *addr, *port)
	subs := make([]*subscriber.Subscriber, 0, *client)

	for i := 0; i < int(*client); i++ {
		sub, err := subscriber.New(ctx, proto, fullAddr, *capacity, topics...)
		if err != nil {
			log.Printf("[gClient] Error : %s\n", err.Error())
			return
		}

		subs = append(subs, sub)
	}

	log.Printf("[gClient] Connected to %s [ %d clients ] ✅\n", fullAddr, *client)

	for id, sub := range subs {
		func(id int, sub *subscriber.Subscriber) {

			go func() {
				defer func() {
					if err := sub.Disconnect(); err != nil {
						log.Printf("[gClient] Error : %s\n", err.Error())
					}
				}()

				var buf = new(bytes.Buffer)
				var fd *os.File

				if *out {
					handle, err := os.OpenFile(fmt.Sprintf("log_sub_%d.csv", id), os.O_CREATE|os.O_RDWR, 0x1b6)
					if err != nil {
						log.Printf("[gClient] Error : %s\n", err.Error())
						return
					}

					fd = handle
				}

				defer func() {
					if err := fd.Close(); err != nil {
						log.Printf("[gClient] Error : %s\n", err.Error())
					}
				}()

				for {
					select {
					case <-ctx.Done():
						return

					case <-sub.Watch():
						received := uint64(time.Now().UnixNano() / 1_000_000)
						msg := sub.Next()
						id, sent, err := utils.DeserialiseMsg(msg)
						if err != nil {
							log.Printf("[gClient] Error : %s\n", err.Error())
							return
						}

						if *out {
							if err := utils.LogMsg(fd, buf, id, sent, received, msg.Topic); err != nil {
								log.Printf("[gClient] Error : %s\n", err.Error())
								break
							}
						}

						log.Printf("[gClient] Received : `%d : %d` from `%s`\n", sent, id, msg.Topic)
					}
				}
			}()

		}(id, sub)
	}

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT)

	<-interruptChan
	cancel()
	<-time.After(time.Second)

	log.Printf("[gClient] Graceful shutdown\n")
}
