package utils

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"

	"github.com/itzmeanjan/pub0sub/ops"
	"github.com/itzmeanjan/pub0sub/publisher"
)

type Publishers struct {
	Handles []*publisher.Publisher
	Logs    []*os.File
	Buffer  *bytes.Buffer
	Idx     uint64
}

func (p *Publishers) PublishMsg(topics TopicList) error {
	p.Idx++

	for i, pub := range p.Handles {
		if pub.Connected() {
			start := uint64(time.Now().UnixNano() / 1_000_000)
			msg, err := p.prepareMsg(topics)
			if err != nil {
				return err
			}

			if _, err := pub.Publish(msg); err != nil {
				return err
			}
			end := uint64(time.Now().UnixNano() / 1_000_000)

			if p.Logs != nil {
				LogMsg(p.Logs[i], p.Buffer, 0, start, end, "na")
			}
		}
	}

	return nil
}

func (p *Publishers) prepareMsg(topics TopicList) (*ops.Msg, error) {
	now := uint64(time.Now().UnixNano() / 1_000_000)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, p.Idx); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, now); err != nil {
		return nil, err
	}

	return &ops.Msg{Topics: topics, Data: buf.Bytes()}, nil
}
