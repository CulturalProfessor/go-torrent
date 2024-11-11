package metainfo

import (
	"net"
	"io"	
	"bytes"
	"encoding/binary"
	"errors"
	"encoding"
	"math/rand"
)
var zeroHash Hash

const HashSize = 20

type Hash [HashSize]byte

func NewHash(b []byte) (h Hash) {
	copy(h[:], b[:HashSize])
	return
}

type CompactAddr struct {
	IP   net.IP // For IPv4, its length must be 4.
	Port uint16
}

func (a CompactAddr) WriteBinary(w io.Writer) (n int, err error) {
	if n, err = w.Write(a.IP); err == nil {
		if err = binary.Write(w, binary.BigEndian, a.Port); err == nil {
			n += 2
		}
	}
	return
}

var (
	_ encoding.BinaryMarshaler   = new(CompactAddr)
	_ encoding.BinaryUnmarshaler = new(CompactAddr)
)

// MarshalBinary implements the interface encoding.BinaryMarshaler,
func (a CompactAddr) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	buf.Grow(18)
	if _, err = a.WriteBinary(buf); err == nil {
		data = buf.Bytes()
	}
	return
}

// UnmarshalBinary implements the interface encoding.BinaryUnmarshaler.
func (a *CompactAddr) UnmarshalBinary(data []byte) error {
	_len := len(data) - 2
	switch _len {
	case net.IPv4len, net.IPv6len:
	default:
		return errors.New("invalid compact ip-address/port info")
	}

	a.IP = make(net.IP, _len)
	copy(a.IP, data[:_len])
	a.Port = binary.BigEndian.Uint16(data[_len:])
	return nil
}
func NewRandomHash() (h Hash) {
	rand.Read(h[:])
	return
}
func (h Hash) IsZero() bool {
	return h == zeroHash
}
