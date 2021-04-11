package mur_tcp

import (
	"encoding/hex"
	"strings"
	"time"
)

// Binary-Coded Decimal (Hex string) to float64
func (r *MURResponse) parseBCDIntString(s string) string {
	sum := ""

	// +2 because of Hex string like "aabb"
	for i := 0; i < len(s); i += 2 {
		sum = s[i:i+2] + sum
	}

	return sum
}

// Binary-Coded Float (Hex string) to usual float string
func (r *MURResponse) parseBCDFloatString(s string) string {
	sum := r.parseBCDIntString(s[:len(s)-2])

	dst := make([]byte, 1)
	_, _ = hex.Decode(dst, []byte(s[len(s)-2:]))
	exp := int(dst[0]) - 127

	if exp > len(s)-2 {
		sum += strings.Repeat("0", exp)
	} else if exp > 0 {
		sum = sum[:exp] + "." + sum[exp:]
	} else {
		if exp < -1 {
			sum = strings.Repeat("0", -exp-1) + sum
		}
		sum = "." + sum
	}

	sum = strings.Trim(sum, "0")
	// sum = strings.TrimSuffix(sum, "0")

	if sum[0] == '.' {
		sum = "0" + sum
	}

	return sum
}

func (r *MURResponse) readBCDByteString(b []byte, offset *int) string {
	return r.parseBCDIntString(r.readNext(b, offset, 1))
}

func (r *MURResponse) readBCDIntString(b []byte, offset *int, length int) string {
	return r.parseBCDIntString(r.readNext(b, offset, length))
}

func (r *MURResponse) readDateString(b []byte, offset *int) string {
	t := r.parseBCDInt(r.readNext(b, offset, 3))

	day := t % 100
	t = t / 100

	month := t % 100
	t = t / 100

	year := t%100 + 2000

	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC).String()
}

func (r *MURResponse) readBCDFloatString(b []byte, offset *int, length int) string {
	return r.parseBCDFloatString(r.readNext(b, offset, length))
}
