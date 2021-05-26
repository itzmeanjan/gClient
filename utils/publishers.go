package utils

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/itzmeanjan/pub0sub/ops"
	"github.com/itzmeanjan/pub0sub/publisher"
)

type Publishers []*publisher.Publisher

func (p *Publishers) PublishMsg(topics TopicList) error {
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

func prepareMsg(topics TopicList) (*ops.Msg, error) {
	now := uint64(time.Now().UnixNano() / 1_000_000)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, now); err != nil {
		return nil, err
	}

	return &ops.Msg{Topics: topics, Data: buf.Bytes()}, nil
}
