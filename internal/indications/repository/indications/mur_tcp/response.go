package mur_tcp

import (
	"fmt"
	"io"
	"strings"
)

type DataType byte

type MURResponse struct {
	Address    byte
	Text       string
	Adapter    string
	Indication map[string]string // *indication.Indication
	zeroLength byte
}

// UnmarshalIndications unmarshals a byte slice into a MURResponse.
func (r *MURResponse) UnmarshalIndications(b []byte) error {

	// // Правила
	// // 1b 07 - перевод каретки на 07 знаков вперёд
	//
	// // 0000   1b 01 01 5c 00 4c 10 95 49 01 40 00 28 01 21 00
	// // 0010   00 87 25 92 66 01 82 00 40 56 65 01 84 00 60 24
	// // 0020   29 02 83 00 10 29 00 05 83 1b 03 63 03 80 1b 03
	// // 0030   20 01 7f 00 00 06 7f 00 30 01 81 00 10 07 82 00
	// // 0040   11 05 82 00 00 90 51 05 85 1b 03 03 08 84 00 a1
	// // 0050   1b 00
	//
	// // 1b01015c004c10954901400028012100008725926601820040566501840060242902830010290005831b036303801b0320017f0000067f0030018100100782001105820000905105851b0303088400a11b00
	//
	// example1 := map[string]string{
	// 	"Adapter number":        "10 95 49 01",
	// 	"Adapter version":       "40",
	// 	"Status (error code)":   "00",
	// 	"Date":                  "28 01 21",
	// 	"Energy, GCal":          "00 00 87 25 92 66 01 82",
	// 	"Volume, cu. m.":        "00 40 56 65 01 84",
	// 	"Add. volume 1, cu. m.": "00 60 24 29 02 83",
	// 	"Add. volume 2, cu. m.": "00 10 29 00 05 83",
	// 	"Add. volume 3, cu. m.": "00 00 00 63 03 80",
	// 	"Add. volume 4, cu. m.": "00 00 00 20 01 7f",
	// 	"Consumption, cu. m./h": "00 00 06 7f",
	// 	"Power, kW":             "00 30 01 81",
	// 	"In temperature, C":     "00 10 07 82",
	// 	"Out temperature, C":    "00 11 05 82",
	// 	"Runtime, h":            "00 00 90 51 05 85",
	// 	"Runtime w/ error, h":   "00 00 00 03 08 84",
	// }
	//
	// // 0000   1b 01 04 5c 00 4c 14 94 49 01 40 00 29 01 21 00
	// // 0010   00 55 56 33 71 01 82 00 34 00 63 03 84 1b 07 50
	// // 0020   37 74 04 83 1b 0d 26 02 80 00 80 01 81 00 16 07
	// // 0030   82 00 45 06 82 00 00 16 52 05 85 00 00 70 21 08
	// // 0040   84 00 03 1b 00
	//
	// // 1b01045c004c14944901400029012100005556337101820034006303841b0750377404831b0d26028000800181001607820045068200001652058500007021088400031b00
	//
	// example2 := map[string]string{
	// 	"Adapter number":        "14 94 49 01",
	// 	"Adapter version":       "40",
	// 	"Status (error code)":   "00",
	// 	"Date":                  "29 01 21",
	// 	"Energy, GCal":          "00 00 55 56 33 71 01 82",
	// 	"Volume, cu. m.":        "00 34 00 63 03 84",
	// 	"Add. volume 1, cu. m.": "00 00 00 00 00 00",
	// 	"Add. volume 2, cu. m.": "00 50 37 74 04 83",
	// 	"Add. volume 3, cu. m.": "00 00 00 00 00 00",
	// 	"Add. volume 4, cu. m.": "00 00 00 00 00 00",
	// 	"Consumption, cu. m./h": "00 26 02 80",
	// 	"Power, kW":             "00 80 01 81",
	// 	"In temperature, C":     "00 16 07 82",
	// 	"Out temperature, C":    "00 45 06 82",
	// 	"Runtime, h":            "00 00 16 52 05 85",
	// 	"Runtime w/ error, h":   "00 00 70 21 08 84",
	// }
	//
	// // 0000   1b 01 00 5c 00 4c 97 98 49 01 40 00 18 12 20 00
	// // 0010   00 88 74 75 25 04 82 00 12 34 70 02 84 00 10 13
	// // 0020   27 02 83 00 30 31 30 05 83 1b 03 40 02 7f 1b 07
	// // 0030   70 06 7f 00 10 02 81 00 17 07 82 00 36 04 82 00
	// // 0040   00 97 41 05 85 00 00 70 40 09 84 00 44 1b 00
	//
	// // 1b01005c004c97984901400018122000008874752504820012347002840010132702830030313005831b0340027f1b0770067f00100281001707820036048200009741058500007040098400441b00
	//
	// example3 := map[string]string{
	// 	"Adapter number":        "97 98 49 01",
	// 	"Adapter version":       "40",
	// 	"Status (error code)":   "00",
	// 	"Date":                  "18 12 20",
	// 	"Energy, GCal":          "00 00 88 74 75 25 04 82",
	// 	"Volume, cu. m.":        "00 12 34 70 02 84",
	// 	"Add. volume 1, cu. m.": "00 10 13 27 02 83",
	// 	"Add. volume 2, cu. m.": "00 30 31 30 05 83",
	// 	"Add. volume 3, cu. m.": "00 00 00 40 02 7f",
	// 	"Add. volume 4, cu. m.": "00 00 00 00 00 00",
	// 	"Consumption, cu. m./h": "00 70 06 7f",
	// 	"Power, kW":             "00 10 02 81",
	// 	"In temperature, C":     "00 17 07 82",
	// 	"Out temperature, C":    "00 36 04 82",
	// 	"Runtime, h":            "00 00 97 41 05 85",
	// 	"Runtime w/ error, h":   "00 00 70 40 09 84",
	// }
	//
	// fmt.Printf("%+v\n%+v\n%+v\n", example1, example2, example3)

	// Verify that header is present
	if len(b) < 7 ||
		b[0] != ByteStart ||
		b[1] != 0x01 {
		return io.ErrUnexpectedEOF
	}

	_ = b[2]

	// Track offset in packet for reading data
	n := 3

	r.Address = b[n]
	n += 1

	if b[n] != 0x00 {
		return io.ErrUnexpectedEOF
	}
	n += 1

	dataLength := b[n]
	if dataLength != 0x4c {
		if dataLength == 0x01 {
			return nil
		}
		return io.ErrUnexpectedEOF
	}
	n += 1

	// tn := n
	//
	// // TODO: Make certain struct for indications
	indications := make(map[string]string)

	// r.Adapter = r.readNext(b, &n, 4)
	r.Adapter = r.readBCDIntString(b, &n, 4)
	// indications["Adapter version"] = r.readNext(b, &n, 1)
	// indications["Status (error code)"] = r.readNext(b, &n, 1)
	// indications["Date"] = r.readNext(b, &n, 3)
	// indications["Energy, GCal"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 8)))
	// indications["Volume, cu. m."] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Add. volume 1, cu. m."] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Add. volume 2, cu. m."] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Add. volume 3, cu. m."] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Add. volume 4, cu. m."] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Consumption, cu. m./h"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["Power, kW"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["In temperature, C"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["Out temperature, C"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["Runtime, h"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Runtime w/ error, h"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))

	// indications["AdapterVersion"] = r.readNext(b, &n, 1)
	// indications["Status"] = r.readNext(b, &n, 1)
	// indications["CurrentDate"] = r.readNext(b, &n, 3)
	// indications["Energy"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 8)))
	// indications["Volume0"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Volume1"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Volume2"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Volume3"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Volume4"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["Consumption"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["Power"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["TemperatureIn"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["TemperatureOut"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 4)))
	// indications["Runtime"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))
	// indications["RuntimeWithError"] = fmt.Sprintf("%.6f", r.parseBCDFloat(r.readNext(b, &n, 6)))

	indications["AdapterVersion"] = r.readBCDByteString(b, &n)
	indications["Status"] = r.readBCDByteString(b, &n)
	indications["CurrentDate"] = r.readDateString(b, &n)
	indications["Energy"] = r.readBCDFloatString(b, &n, 8)
	indications["Volume"] = r.readBCDFloatString(b, &n, 6)
	indications["Volume1"] = r.readBCDFloatString(b, &n, 6)
	indications["Volume2"] = r.readBCDFloatString(b, &n, 6)
	indications["Volume3"] = r.readBCDFloatString(b, &n, 6)
	indications["Volume4"] = r.readBCDFloatString(b, &n, 6)
	indications["Consumption"] = r.readBCDFloatString(b, &n, 4)
	indications["Power"] = r.readBCDFloatString(b, &n, 4)
	indications["TemperatureIn"] = r.readBCDFloatString(b, &n, 4)
	indications["TemperatureOut"] = r.readBCDFloatString(b, &n, 4)
	indications["Runtime"] = r.readBCDFloatString(b, &n, 6)
	indications["RuntimeWithError"] = r.readBCDFloatString(b, &n, 6)

	r.Indication = indications

	// fmt.Printf("%v\n", indications)
	// n = tn

	// r.Indication = &indication.Indication{
	// 	AdapterNumber:    r.readBCDInt(b, &n, 4),
	// 	AdapterVersion:   r.readBCDByte(b, &n),
	// 	Status:           r.readBCDByte(b, &n),
	// 	CurrentDate:      r.readDate(b, &n),
	// 	Energy:           r.readBCDFloat(b, &n, 8),
	// 	Volume:           r.readBCDFloat(b, &n, 6),
	// 	Volume1:          r.readBCDFloat(b, &n, 6),
	// 	Volume2:          r.readBCDFloat(b, &n, 6),
	// 	Volume3:          r.readBCDFloat(b, &n, 6),
	// 	Volume4:          r.readBCDFloat(b, &n, 6),
	// 	Consumption:      r.readBCDFloat(b, &n, 4),
	// 	Power:            r.readBCDFloat(b, &n, 4),
	// 	TemperatureIn:    r.readBCDFloat(b, &n, 4),
	// 	TemperatureOut:   r.readBCDFloat(b, &n, 4),
	// 	Runtime:          r.readBCDFloat(b, &n, 6),
	// 	RuntimeWithError: r.readBCDFloat(b, &n, 6),
	// }

	_ = b[n]
	_ = b[n+1]
	n += 2

	if b[n] != ByteStart ||
		b[n+1] != ByteEnd {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func (r *MURResponse) readNext(b []byte, offset *int, length int) string {
	var result string

	if r.zeroLength > 0 { // Saved zeros from previous indication
		neededZeros := min(length, int(r.zeroLength))

		result = strings.Repeat("00", neededZeros)

		r.zeroLength -= byte(neededZeros)
		length -= neededZeros

	} else if b[*offset] == ByteStart { // Skip some zeros
		zeros := int(b[*offset+1])
		*offset += 2

		neededZeros := min(length, zeros)

		result = strings.Repeat("00", neededZeros)

		length -= neededZeros
		zeros -= neededZeros

		// It's need to reed some more bytes
		if length > 0 {
			result += fmt.Sprintf("%x", b[*offset:*offset+length])
			*offset += length
		}

		// Some more zeros are left
		if zeros > 0 {
			r.zeroLength = byte(zeros - length)
		}

		return result
	}

	result += fmt.Sprintf("%x", b[*offset:*offset+length])
	*offset += length

	return result
}

// Function to find minimum of x and y
func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

// Function to find maximum of x and y
func max(x, y int) int {
	if x < y {
		return y
	}

	return x
}
