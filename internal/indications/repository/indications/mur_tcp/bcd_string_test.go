package mur_tcp

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestParseBCDIntString(t *testing.T) {
	type TestInOut struct {
		BCD   string
		Human string
	}

	table := []TestInOut{
		{"97984901", "01499897"},
		{"181220", "201218"},
		{"00008874752504", "04257574880000"},
		{"0012347002", "0270341200"},
	}

	t.Run("TestParseBCDIntString", func(t *testing.T) {
		r := MURResponse{}
		for _, inOut := range table {
			result := 		r.parseBCDIntString(inOut.BCD)
			assert.Equal(t, result, inOut.Human)
		}
	})
}

func TestParseBCDFloatString(t *testing.T) {
	type TestInOut struct {
		BCD   string
		Human string
	}

	table := []TestInOut{
		{"0000554524920382", "39.2244555"},
		{"002341180384", "3184.123"},
		{"000010620181", "1.621"},
		{"000020000381", "3.002"},
		{"00000040037f", "0.034"},
		{"000040040181", "1.044"},
	}

	t.Run("TestParseBCDFloatString", func(t *testing.T) {
		r := MURResponse{}
		for _, inOut := range table {
			result := 		r.parseBCDFloatString(inOut.BCD)
			assert.Equal(t, result, inOut.Human)
		}
	})
}
