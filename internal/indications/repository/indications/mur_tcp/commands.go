package mur_tcp

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// 1b 01 ff 00 8a 8a 1b 00
func (m *MUR) GetVersion() string {
	var version string
	{
		packet := MURRawPacket{
			Command: CommandGetVersion,
		}

		if err := m.sendReconnect(packet); err != nil {
			return ""
		}

		buff := m.readReconnect(packet)

		// Check response length
		if len(buff) < 7 {
			return ""
		}

		// Track offset in packet for reading data
		n := 4

		length := int(binary.BigEndian.Uint16(buff[n : n+2]))
		n += 2

		version = string(buff[n : n+length])
	}

	return version
}

type RegistrarInfo struct {
	NetworkAddress byte
	AdapterCount   uint16
	SerialNumber   uint16
}

func (m *MUR) GetRegistrarInfo() *RegistrarInfo {
	var networkAddress byte
	var count uint16
	var serialNumber uint16
	{
		packet := MURRawPacket{
			Command: byte(0x07),
			Data:    []byte{0x00, 0x00, 0xf8},
		}

		if err := m.sendReconnect(packet); err != nil {
			return nil
		}

		buff := m.readReconnect(packet)
		if len(buff) < 23 {
			return nil
		}

		networkAddress = buff[3]
		serialNumber = binary.LittleEndian.Uint16(buff[6:8])
		count = binary.LittleEndian.Uint16(buff[20:22])
	}

	return &RegistrarInfo{
		NetworkAddress: networkAddress,
		AdapterCount:   count,
		SerialNumber:   serialNumber,
	}
}

// 1b 01 ff 00 8a 8a 1b 00        - GetVersion
// 1b 01 ff 00 06 2e 00 28 1b 00
// 1b 01 ff 04 07 c0 0b 07 cf 1b 00 - SetInstantMode
func (m *MUR) SetInstantMode() {
	// req := MURPacket{
	// 	Channel: Channel4,
	// 	Adapter: 0xc0,
	// 	Command: CommandAdapterRequest,
	// }
	// binary, err := req.MarshalBinary()
	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return ""
	// }

	// Originally sent by CfgWin2RCX.exe after GetVersion - Check status?
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x01, 0xff, 0xaa, 0xbb, 0xcc, 0x23, 0x1b, 0x00})
	// if err := m.sendPlainCheckResponse([]byte{0x1b, 0x01, 0xff, 0x01, 0xff, 0xaa, 0xbb, 0xcc, 0x23, 0x1b, 0x00},
	// 	[]byte{0x1b, 0x01, 0x01, 0x5c, 0x00, 0xa2, 0x1b, 0x00}); err != nil {
	// 	fmt.Printf("Check status after GetVersion error: %+v", err)
	// }

	// Way 1
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x03, 0xff, 0xaa, 0xbb, 0xcc, 0x21, 0x1b, 0x00})
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x04, 0x09, 0x00, 0x05, 0x00, 0x03, 0x0b, 0x1b, 0x00})
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x01, 0x07, 0xc0, 0x0b, 0x07, 0xca, 0x1b, 0x00})

	// Way 2
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x02, 0xff, 0xaa, 0xbb, 0xcc, 0x20, 0x1b, 0x00})
	// m.sendPlainIgnoreResponse([]byte{0x1b, 0x01, 0xff, 0x03, 0x07, 0x20, 0x06, 0x07, 0x25, 0x1b, 0x00})

	// packet := MURRawPacket{
	// 	Command: byte(0xff),
	// 	Data:    []byte{0xaa, 0xbb, 0xcc},
	// }
	//
	// if err := m.sendReconnect(packet); err != nil {
	// 	return
	// }
	//
	// _ = m.readReconnect(packet)

	m.sendIgnore([]byte{0xff, 0xaa, 0xbb, 0xcc})
}

func (m *MUR) Foreplay() {
	m.sendIgnore([]byte{0x09, 0x00, 0x05, 0x00, 0x03})
	m.sendIgnore([]byte{0x2b})
}

func (m *MUR) sendIgnore(cmd []byte) {
	if len(cmd) < 1 {
		return
	}

	packet := MURRawPacket{
		Command: cmd[0],
		Data:    cmd[1:],
	}

	if err := m.sendReconnect(packet); err != nil {
		return
	}

	_ = m.readReconnect(packet)
}

func (m *MUR) sendReconnect(packet MURRawPacket) error {
	m.conn.SetTimeout(15 * time.Second)

	if err := packet.Send(m.conn); err != nil {
		err = m.conn.Reconnect()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		m.SetInstantMode()
		return err
	}

	return nil
}

func (m *MUR) readReconnect(packet MURRawPacket) []byte {
	m.conn.SetTimeout(15 * time.Second)

	buff, err := packet.Read(m.conn)
	if err != nil {
		err = m.conn.Reconnect()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		m.SetInstantMode()
		return nil
	}

	// Debug
	// fmt.Printf("%d: %#v\n", len(buff), buff)

	return buff
}

// Returns adapter number and map of indications.
//
// 1b 01 ff 00 8a 8a 1b 00        - GetVersion
// 1b 01 ff 00 06 2e 00 28 1b 00
// 1b 01 ff 00 07 00 06 07 06 1b 00
// 1b 01 ff 04 06 00 00 02 1b 00 - Get adapter 0
func (m *MUR) GetAdapterIndications(number uint16) (string, map[string]string) {
	ada1 := byte(number % 256)
	ada2 := byte(number / 256)
	unk1 := byte(0x07)
	pada := uint16(0x0600) + number*0x20 // Pre-command adapter number
	pad1 := byte(pada % 256)
	pad2 := byte(pada / 256)

	// Pre-request
	{
		prePacket := MURRawPacket{
			Command: CommandAdapterRequest,
			Data:    []byte{pad1, pad2, unk1},
		}

		if err := m.sendReconnect(prePacket); err != nil {
			return "", nil
		}

		// Ignore pre-request results
		_ = m.readReconnect(prePacket)
	}

	// Main request
	res := MURResponse{}
	{
		packet := MURRawPacket{
			Command: CommandInstantValues,
			Data:    []byte{ada1, ada2},
		}

		if err := m.sendReconnect(packet); err != nil {
			return "", nil
		}

		buff := m.readReconnect(packet)

		// Result parsing
		err := res.UnmarshalIndications(buff)
		if err != nil {
			fmt.Printf("indication unmarshaling error: %v\n", err)
			return "", nil
		}
	}

	return res.Adapter, res.Indication
}

func (m *MUR) GetAllIndications() map[string]map[string]string {
	res := make(map[string]map[string]string)

	for i := uint16(0); i < m.meterCount; i++ {
		timeStart := time.Now()
		adapter, in := m.GetAdapterIndications(i)
		if in == nil {
			fmt.Printf("Adapter %d timeout\n", i)
			continue
		}
		in["AdapterID"] = strconv.Itoa(int(i))
		in["RegistrarIP"] = m.ip.String()
		in["PollDuration"] = fmt.Sprintf("%d", time.Since(timeStart))
		in["LastIndicationsTime"] = fmt.Sprintf("%d", time.Now().UTC().Unix())
		fmt.Printf("Adapter %s: %+v\n", adapter, in)

		res[adapter] = in
	}

	return res
}
