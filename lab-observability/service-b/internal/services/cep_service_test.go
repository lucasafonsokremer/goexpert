package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCEP(t *testing.T) {
	service := NewCEPService()

	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "valid CEP",
			cep:      "12345678",
			expected: true,
		},
		{
			name:     "invalid CEP with letters",
			cep:      "1234567a",
			expected: false,
		},
		{
			name:     "invalid CEP too short",
			cep:      "1234567",
			expected: false,
		},
		{
			name:     "invalid CEP too long",
			cep:      "123456789",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ValidateCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}
