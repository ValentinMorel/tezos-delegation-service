package parser

import (
	"errors"
	"fmt"
	"strconv"
	"tezos-delegation-service/internal/db"
	"time"
)

// Error messages
var (
	ErrInvalidIDType          = errors.New("invalid type for id")
	ErrInvalidInt64Type       = errors.New("invalid type for int64")
	ErrInvalidSenderType      = errors.New("invalid type for sender")
	ErrMissingSenderAddress   = errors.New("missing sender address")
	ErrInvalidTimestampFormat = errors.New("invalid timestamp format")
	ErrMissingField           = errors.New("missing field")
)

func ParseID(id interface{}) (int64, error) {
	switch v := id.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int, int32, int64:
		return int64(v.(int)), nil
	case uint, uint32, uint64:
		return int64(v.(uint64)), nil
	default:
		return 0, fmt.Errorf("%w: %T", ErrInvalidIDType, id)
	}
}

func ParseInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int, int32, int64:
		return int64(v.(int)), nil
	case uint, uint32, uint64:
		return int64(v.(uint64)), nil
	default:
		return 0, fmt.Errorf("%w: %T", ErrInvalidInt64Type, val)
	}
}

func checkField(raw map[string]interface{}, key string) (interface{}, error) {
	val, ok := raw[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingField, key)
	}
	return val, nil
}

func ParseDelegationParameters(raw map[string]interface{}) (*db.Delegation, error) {
	idVal, err := checkField(raw, "id")
	if err != nil {
		return nil, err
	}
	id, err := ParseID(idVal)
	if err != nil {
		return nil, err
	}

	senderVal, err := checkField(raw, "sender")
	if err != nil {
		return nil, err
	}
	sender, ok := senderVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrInvalidSenderType, senderVal)
	}

	delegator, ok := sender["address"].(string)
	if !ok {
		return nil, ErrMissingSenderAddress
	}

	timestampVal, err := checkField(raw, "timestamp")
	if err != nil {
		return nil, err
	}
	timestampStr, ok := timestampVal.(string)
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrInvalidTimestampFormat, timestampVal)
	}
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidTimestampFormat, timestampStr)
	}

	amountVal, err := checkField(raw, "amount")
	if err != nil {
		return nil, err
	}
	amount, err := ParseInt64(amountVal)
	if err != nil {
		return nil, err
	}

	levelVal, err := checkField(raw, "level")
	if err != nil {
		return nil, err
	}
	level, err := ParseInt64(levelVal)
	if err != nil {
		return nil, err
	}

	delegation := db.Delegation{
		ID:        id,
		Delegator: delegator,
		Timestamp: timestamp,
		Amount:    amount,
		Level:     level,
	}
	return &delegation, nil
}
