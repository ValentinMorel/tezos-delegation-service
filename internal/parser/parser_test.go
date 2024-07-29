package parser

import (
	"testing"
	"tezos-delegation-service/internal/db"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int64
		hasError bool
	}{
		{input: float64(123), expected: 123, hasError: false},
		{input: "456", expected: 456, hasError: false},
		{input: int(789), expected: 789, hasError: false},
		{input: uint64(1011), expected: 1011, hasError: false},
		{input: "abc", expected: 0, hasError: true},
	}

	for _, test := range tests {
		result, err := ParseID(test.input)
		if test.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int64
		hasError bool
	}{
		{input: float64(123), expected: 123, hasError: false},
		{input: "456", expected: 456, hasError: false},
		{input: int(789), expected: 789, hasError: false},
		{input: uint64(1011), expected: 1011, hasError: false},
		{input: "abc", expected: 0, hasError: true},
	}

	for _, test := range tests {
		result, err := ParseInt64(test.input)
		if test.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestParseDelegationParameters(t *testing.T) {
	tests := []struct {
		input    map[string]interface{}
		expected *db.Delegation
		hasError bool
	}{
		{
			input: map[string]interface{}{
				"id":        "1",
				"sender":    map[string]interface{}{"address": "tz1abc"},
				"timestamp": "2024-07-29T12:34:56Z",
				"amount":    "1000",
				"level":     "1",
			},
			expected: &db.Delegation{
				ID:        1,
				Delegator: "tz1abc",
				Timestamp: time.Date(2024, 7, 29, 12, 34, 56, 0, time.UTC),
				Amount:    1000,
				Level:     1,
			},
			hasError: false,
		},
		{
			input: map[string]interface{}{
				"id":        "abc",
				"sender":    map[string]interface{}{"address": "tz1abc"},
				"timestamp": "2024-07-29T12:34:56Z",
				"amount":    "1000",
				"level":     "1",
			},
			expected: nil,
			hasError: true,
		},
		{
			input: map[string]interface{}{
				"id":        "1",
				"sender":    map[string]interface{}{"address": "tz1abc"},
				"timestamp": "invalid",
				"amount":    "1000",
				"level":     "1",
			},
			expected: nil,
			hasError: true,
		},
		{
			input: map[string]interface{}{
				"id":        "1",
				"sender":    map[string]interface{}{"address": 123},
				"timestamp": "2024-07-29T12:34:56Z",
				"amount":    "1000",
				"level":     "1",
			},
			expected: nil,
			hasError: true,
		},
		{
			input: map[string]interface{}{
				"id":        "1",
				"sender":    "invalid_sender",
				"timestamp": "2024-07-29T12:34:56Z",
				"amount":    "1000",
				"level":     "1",
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := ParseDelegationParameters(test.input)
		if test.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}
