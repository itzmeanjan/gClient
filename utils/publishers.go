package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	now := uint64(time.Now().UTC().Nanosecond() / 1_000_000)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, now); err != nil {
		return nil, err
	}

	return &ops.Msg{Topics: topics, Data: buf.Bytes()}, nil
}

func DeserialiseMsg(msg *ops.PushedMessage) (uint64, error) {
	if msg == nil {
		return 0, errors.New("nil message")
	}

	var ts uint64
	buf := bytes.NewReader(msg.Data)
	if err := binary.Read(buf, binary.BigEndian, &ts); err != nil {
		return 0, err
	}

	return ts, nil
}
