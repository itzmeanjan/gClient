package main

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"github.com/itzmeanjan/pub0sub/ops"
	"github.com/itzmeanjan/pub0sub/publisher"
)

type topicList []string

func (t *topicList) String() string {
	if t == nil {
		return ""
	}

	return strings.Join(*t, "")
}

func (t *topicList) Set(val string) error {
	*t = append(*t, val)
	return nil
}

func prepareMsg(topics topicList) (*ops.Msg, error) {
	now := uint64(time.Now().UTC().Nanosecond() / 1_000_000)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, now); err != nil {
		return nil, err
	}

	return &ops.Msg{Topics: topics, Data: buf.Bytes()}, nil
}

type publishers []*publisher.Publisher

func (p *publishers) publishMsg(topics topicList) error {
	msg, err := prepareMsg(topics)
	if err != nil {
		return err
	}

	for _, pub := range *p {
		if pub.Connected() {

			if _, err := pub.Publish(msg); err != nil {
				return err
			}

		}
	}

	return nil
}
