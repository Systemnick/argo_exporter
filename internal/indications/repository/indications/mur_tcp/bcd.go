package mur_tcp

import (
	"encoding/hex"
	"math"
	"time"
)

// Binary-Coded Decimal (Hex string) to float64
func (r *MURResponse) parseBCDInt(s string) uint64 {
	sum := uint64(0)
	multiplier := 0

	// +2 because of Hex string like "aabb"
	for i := 0; i < len(s); i += 2 {
		sum += uint64(s[i+1]-0x30) * uint64(math.Pow10(multiplier))
		multiplier++

		sum += uint64(s[i]-0x30) * uint64(math.Pow10(multiplier))
		multiplier++
	}

	return sum
}

// Binary-Coded Decimal (Hex string) to float64
func (r *MURResponse) parseBCDFloat(s string) float64 {
	f := float64(0)
	sum := r.parseBCDInt(s[:len(s)-2])

	dst := make([]byte, 1)
	_, _ = hex.Decode(dst, []byte(s[len(s)-2:]))
	exp := int(dst[0]) - 127 - (len(s) - 2)

	f = float64(sum) * math.Pow10(exp)

	return f
}

func (r *MURResponse) readBCDByte(b []byte, offset *int) uint8 {
	return uint8(r.parseBCDInt(r.readNext(b, offset, 1)))
}

func (r *MURResponse) readBCDInt(b []byte, offset *int, length int) uint32 {
	return uint32(r.parseBCDInt(r.readNext(b, offset, length)))
}

func (r *MURResponse) readDate(b []byte, offset *int) time.Time {
	t := r.parseBCDInt(r.readNext(b, offset, 3))

	day := t % 100
	t = t / 100

	month := t % 100
	t = t / 100

	year := t%100 + 2000

	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
}

func (r *MURResponse) readBCDFloat(b []byte, offset *int, length int) float64 {
	return r.parseBCDFloat(r.readNext(b, offset, length))
}
