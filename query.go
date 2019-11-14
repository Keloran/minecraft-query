package minecraft

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/eknkc/basex"
	"io"
)

func (m Minecraft) Query() error {
	if m.Conn == nil {
		return fmt.Errorf("no connection")
	}

	challenge, err := m.GetChallenge()
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	status, err := m.GetStatus(challenge)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	fmt.Println(fmt.Sprintf("Status: %v", status))

	return nil
}

func (m Minecraft) GetChallenge() (int32, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, []byte("\xFE\xFD"))
	binary.Write(&buf, binary.BigEndian, []byte("\x09"))
	binary.Write(&buf, binary.BigEndian, uint32(1))
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		return 0, fmt.Errorf("challange write: %w", err)
	}
	fmt.Println(fmt.Sprintf(hex.EncodeToString(buf.Bytes())))

	out := make([]byte, 12)
	_, err = io.ReadFull(m.Conn, out)
	if err != nil {
		if err != io.EOF {
			return 0, fmt.Errorf("challenge read: %w", err)
		}
	}

	e, err := basex.NewEncoding("0123456789abcdef")
	if err != nil {
		return 0, fmt.Errorf("challenge encode: %w", err)
	}
	f := e.Encode(out[5:])

	fmt.Println(fmt.Sprintf("F: %v", f))

	return 0, nil
}

func (m Minecraft) GetStatus(challenge int32) (interface{}, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, []byte("\xFE\xFD"))
	binary.Write(&buf, binary.BigEndian, []byte("\x00"))
	binary.Write(&buf, binary.BigEndian, uint32(1))
	binary.Write(&buf, binary.BigEndian, challenge)
	binary.Write(&buf, binary.BigEndian, []byte("\x00\x00\x00\x00"))
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("status write: %w", err)
	}
	fmt.Println(fmt.Sprintf("status write: %v", hex.EncodeToString(buf.Bytes())))

	out := make([]byte, 4)
	_, err = io.ReadFull(m.Conn, out)
	fmt.Println(fmt.Sprintf("out %v, hex: %v", out, hex.EncodeToString(out)))

	return nil, nil
}

