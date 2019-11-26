package minecraft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func (m Minecraft) Ping(resolve bool) (PingInfo, error) {
	if m.Conn == nil {
		return PingInfo{}, fmt.Errorf("no connection")
	}

	if resolve {
		err := m.Resolve()
		if err != nil {
			return PingInfo{}, fmt.Errorf("Ping Resolve: %w", err)
		}
	}

	return m.pingQuery()
}

func (m Minecraft) pingQuery() (PingInfo, error) {
	var buf, buff bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, []byte{00, 04})
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery magic: %w", err)
	}
	err = binary.Write(&buf, binary.LittleEndian, []byte(string(len(m.Address))))
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery address length: %w", err)
	}
	err = binary.Write(&buf, binary.LittleEndian, []byte(m.Address))
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery address: %w", err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(m.Port))
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery port: %w", err)
	}
	err = binary.Write(&buf, binary.LittleEndian, []byte{01})
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery status: %w", err)
	}
	err = binary.Write(&buff, binary.LittleEndian, []byte(string(buf.Len())))
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery buffer length: %w", err)
	}
	err = binary.Write(&buff, binary.LittleEndian, buf.Bytes())
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery buffer bytes: %w", err)
	}
	err = binary.Write(&buff, binary.LittleEndian, []byte{01, 00})
	if err != nil {
		return PingInfo{}, fmt.Errorf("pingQuery status ping: %w", err)
	}

	ret := make(chan []byte, 1)
	e := make(chan error, 1)

	go func() {
		_, err = m.Conn.Write(buff.Bytes())
		if err != nil {
			e <- fmt.Errorf("handshake write: %w", err)
			ret <- []byte{}
			return
//			return PingInfo{}, fmt.Errorf("handshake write: %w", err)
		}
		buff.Reset()
		buf.Reset()

		r := bufio.NewReader(m.Conn)
		z, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				e <- fmt.Errorf("read byte buffer: %w", err)
				ret <- []byte{}
				return
//				return PingInfo{}, fmt.Errorf("read byte buffer: %w", err)
			}
		}

		e <- nil
		ret <- z
	}()

	j := PingInfo{}
	select {
		case r := <- ret:
			cerr := <- e
			if cerr != nil {
				return PingInfo{}, fmt.Errorf("ping chan: %w", cerr)
			}
			err = json.Unmarshal(r[3:], &j)
			if err != nil {
				return PingInfo{}, fmt.Errorf("unmarshal json: %w", err)
			}

			return j, nil
		case <- time.After(time.Duration(m.Timeout) * time.Second):
			return PingInfo{}, fmt.Errorf("ping timeout")
	}

	return j, nil
}
