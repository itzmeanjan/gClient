package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/itzmeanjan/pub0sub/ops"
)

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

func LogMsg(fd *os.File, buf *bytes.Buffer, sent uint64, received uint64, topic string) error {
	defer func() {
		buf.Reset()
	}()

	n, err := buf.WriteString(fmt.Sprintf("%d; %d; %s", sent, received, topic))
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
