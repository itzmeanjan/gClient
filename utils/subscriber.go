package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/itzmeanjan/pub0sub/ops"
)

func DeserialiseMsg(msg *ops.PushedMessage) (uint64, uint64, error) {
	if msg == nil {
		return 0, 0, errors.New("nil message")
	}

	var (
		id  uint64
		ts  uint64
		buf = bytes.NewReader(msg.Data)
	)

	if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
		return 0, 0, err
	}
	if err := binary.Read(buf, binary.BigEndian, &ts); err != nil {
		return 0, 0, err
	}

	return id, ts, nil
}

func LogMsg(fd *os.File, buf *bytes.Buffer, id, sent, received uint64, topic string) error {
	defer func() {
		buf.Reset()
	}()

	var template string
	if id != 0 {
		template = fmt.Sprintf("%d; %d; %d; %s\n", sent, received, id, topic)
	} else {
		template = fmt.Sprintf("%d; %d; %s\n", sent, received, topic)
	}

	n, err := buf.WriteString(template)
	if err != nil {
		return err
	}

	m, err := fd.Write(buf.Bytes())
	if err != nil {
		return err
	}

	if n != m {
		return errors.New("incomplete write")
	}

	return nil
}
