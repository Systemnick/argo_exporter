package mur_tcp

import (
	"bufio"
	"fmt"
)

const (
	ByteStart byte = 0x1b
	Byte01    byte = 0x01
	ByteEnd   byte = 0x00

	AddressBroadcast = byte(0xff)

	Channel0 byte = 0
	Channel3 byte = 3
	Channel4 byte = 4
)

var (
	CommandGetVersion     byte = 0x8a
	CommandInstantValues  byte = 0x06
	CommandAdapterRequest byte = 0x07
)

type MURPacket struct {
	Channel byte
	Adapter uint16
	Command byte
}

type MURRawPacket struct {
	Command  byte
	Data     []byte
	sequence *byte
}

type Sequence byte

var sequence Sequence = 255

func (s *Sequence) next() byte {
	*s += 1
	if *s > 4 {
		*s = 0
	}

	return byte(*s)
}

func (p *MURRawPacket) checksum() byte {
	p.ensureSequence()

	checksum := *p.sequence ^ p.Command
	for i := 0; i < len(p.Data); i++ {
		checksum = checksum ^ p.Data[i]
	}

	return checksum
}

func (p *MURRawPacket) ensureSequence() {
	if p.sequence == nil {
		seq := sequence.next()
		p.sequence = &seq
	}
}

func (p *MURRawPacket) escapeData() {
	for i := 0; i < len(p.Data); i++ {
		if p.Data[i] == byte(0x1b) && p.Data[i+1] == byte(0x00) {
			// Ensure slice capacity
			p.Data = append(p.Data, 0)
			copy(p.Data[i+2:], p.Data[i+1:len(p.Data)-1])
			p.Data[i+1] = ByteStart
			i += 2
		}
	}
}

// MarshalBinary allocates a byte slice and marshals a MURPacket into binary form.
func (p *MURRawPacket) MarshalBinary() ([]byte, error) {
	p.ensureSequence()

	checksum := p.checksum()

	p.escapeData()

	b := make([]byte, p.length())

	n := 0

	copy(b[n:n+4], []byte{ByteStart, Byte01, AddressBroadcast, *p.sequence})
	n += 4

	b[n] = p.Command
	n += 1

	// binary.BigEndian.PutUint16(b[n:n+2], uint16(p.Command))
	copy(b[n:n+len(p.Data)], p.Data)
	n += len(p.Data)

	if checksum == ByteStart {
		copy(b[n:n+4], []byte{checksum, ByteStart, ByteStart, ByteEnd})
		n += 4
	} else {
		copy(b[n:n+3], []byte{checksum, ByteStart, ByteEnd})
		n += 3
	}

	return b, nil
}

func (p *MURRawPacket) Send(conn *ReConn) error {
	binary, err := p.MarshalBinary()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	_, err = fmt.Fprint(conn.conn, string(binary))
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	return nil
}

func (p *MURRawPacket) Read(conn *ReConn) ([]byte, error) {
	length := 0
	response := [1024]byte{}

	for {
		buff := make([]byte, 128)
		n, err := bufio.NewReader(conn.conn).Read(buff)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			return nil, err
		}

		copy(response[length:length+n], buff[:n])
		length += n

		if response[length-2] == byte(0x1b) && response[length-1] == byte(0x00) {
			break
		}
	}

	return response[:length], nil
}

func (p *MURRawPacket) length() int {
	add := 0
	if p.checksum() == ByteStart {
		add++
	}

	return 4 + 1 + len(p.Data) + 3 + add
}

// MarshalBinary allocates a byte slice and marshals a MURPacket into binary form.
func (p *MURPacket) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.length())
	_, err := p.read(b)
	return b, err
}

// read reads data from a MURPacket into b. read is used to marshal a MURPacket
// into binary form, but does not allocate on its own.
func (p *MURPacket) read(b []byte) (int, error) {
	// tcp.port == 5000 && ip.addr == 10.0.0.78 && tcp.payload[0] == 0x1b

	n := 0

	copy(b[n:n+4], []byte{ByteStart, 1, AddressBroadcast, p.Channel})
	n += 4

	// binary.BigEndian.PutUint16(b[n:n+2], uint16(p.Command))
	copy(b[n:n+2], []byte{p.Command, byte(p.Adapter % 256)})
	n += 2

	switch p.Command {
	case CommandInstantValues:
		copy(b[n:n+2], []byte{0x00, 0x28})
		n += 2
	case CommandAdapterRequest:
		copy(b[n:n+3], []byte{0x0b, 0x07, 0xcf})
		n += 3
	}

	copy(b[n:n+2], []byte{ByteStart, ByteEnd})
	n += 2

	return n, nil
}

// UnmarshalBinary unmarshals a byte slice into a MURPacket.
func (p *MURPacket) UnmarshalBinary(b []byte) error {
	// // Verify that both hardware addresses and a single EtherType are present
	// if len(b) < 14 {
	// 	return io.ErrUnexpectedEOF
	// }
	//
	// // Track offset in packet for reading data
	// n := 14
	//
	// // Continue looping and parsing VLAN tags until no more VLAN EtherType
	// // values are detected
	// et := EtherType(binary.BigEndian.Uint16(b[n-2 : n]))
	// switch et {
	// case EtherTypeServiceVLAN, EtherTypeVLAN:
	// 	// VLAN type is hinted for further parsing.  An index is returned which
	// 	// indicates how many bytes were consumed by VLAN tags.
	// 	nn, err := p.unmarshalVLANs(et, b[n:])
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	n += nn
	// default:
	// 	// No VLANs detected.
	// 	p.EtherType = et
	// }
	//
	// // Allocate single byte slice to store destination and source hardware
	// // addresses, and payload
	// bb := make([]byte, 6+6+len(b[n:]))
	// copy(bb[0:6], b[0:6])
	// p.Destination = bb[0:6]
	// copy(bb[6:12], b[6:12])
	// p.Source = bb[6:12]
	//
	// // There used to be a minimum payload length restriction here, but as
	// // long as two hardware addresses and an EtherType are present, it
	// // doesn't really matter what is contained in the payload.  We will
	// // follow the "robustness principle".
	// copy(bb[12:], b[n:])
	// p.Payload = bb[12:]

	return nil
}

// length calculates the number of bytes required to store a MURPacket.
func (p *MURPacket) length() int {
	// // If payload is less than the required minimum length, we zero-pad up to
	// // the required minimum length
	// pl := len(p.Payload)
	// if pl < minPayload {
	// 	pl = minPayload
	// }
	//
	// // Add additional length if VLAN tags are needed.
	// var vlanLen int
	// switch {
	// case p.ServiceVLAN != nil && p.VLAN != nil:
	// 	vlanLen = 8
	// case p.VLAN != nil:
	// 	vlanLen = 4
	// }
	//
	// // 6 bytes: destination hardware address
	// // 6 bytes: source hardware address
	// // N bytes: VLAN tags (if present)
	// // 2 bytes: EtherType
	// // N bytes: payload length (may be padded)
	// return 6 + 6 + vlanLen + 2 + pl
	if p.Command == CommandGetVersion {
		return 8
	} else if p.Command == CommandInstantValues {
		return 10
	}

	switch p.Command {
	case CommandGetVersion:
		return 8
	case CommandInstantValues:
		return 10
	case CommandAdapterRequest:
		return 11
	}
	return 0
}
