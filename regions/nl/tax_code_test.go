package nl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyTaxCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{
			name: "empty",
			code: "",
			err:  "invalid VAT number",
		},
		{
			name: "too long",
			code: "a really really long string that's way too long",
			err:  "invalid VAT number",
		},
		{
			name: "too short",
			code: "shorty",
			err:  "invalid VAT number",
		},
		{
			name: "valid",
			code: "NL000099995B57",
		},
		{
			name: "lowercase",
			code: "nl000099995b57",
		},
		{
			name: "no B",
			code: "NL000099998X57",
			err:  "invalid VAT number",
		},
		{
			name: "non numbers",
			code: "NL000099998B5a",
			err:  "invalid VAT number",
		},
		{
			name: "invalid checksum",
			code: "NL123456789B12",
			err:  "checkusum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyTaxCode(tt.code)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
		})
	}
}