package jsonx

import (
	"encoding/json"
)

// OrderedKey has no documentation
type OrderedKey string

// MarshalText has no documentation
func (ok *OrderedKey) MarshalText() (text []byte, err error) {
	return []byte(*ok), nil
}

// UnmarshalText has no documentation
func (ok *OrderedKey) UnmarshalText(text []byte) error {
	*ok = OrderedKey(string(text))
	return nil
}

// OrderedMap has no documentation
type OrderedMap map[OrderedKey]json.RawMessage
