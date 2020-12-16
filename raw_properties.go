package jsonx

import (
	"encoding/json"
	"reflect"
)

// RawProperties has no documentation
type RawProperties map[string]json.RawMessage

// Unmarshal has no documentation
func (rp *RawProperties) Unmarshal(key string, target interface{}) error {
	// parameter target needs to be a non nil pointer
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(target)}
	}

	var err error
	if rawMessage, found := (*rp)[key]; found {
		delete(*rp, key)
		if err = json.Unmarshal(rawMessage, target); err != nil {
			return err
		}
	}
	return nil
}

// MarshalAll has no documentation
func (rp *RawProperties) MarshalAll(obj interface{}) error {
	if obj == nil {
		return nil
	}
	var err error
	var data []byte
	if data, err = json.Marshal(obj); err != nil {
		return err
	}
	if err = json.Unmarshal(data, rp); err != nil {
		return err
	}
	return nil
}

// MarshalUnknowns has no documentation
func (rp *RawProperties) MarshalUnknowns(p Unknowns) error {
	var err error
	var data []byte
	for k, v := range p {
		if data, err = json.Marshal(v); err != nil {
			return err
		}
		(*rp)[k] = data
	}
	return nil
}

// UnmarshalUnknowns has no documentation
func (rp *RawProperties) UnmarshalUnknowns(p *Unknowns) {
	if len(*rp) == 0 {
		return
	}
	*p = Unknowns{}
	for k, v := range *rp {
		(*p)[k] = v
	}
}

// Marshal has no documentation
func (rp *RawProperties) Marshal(key string, obj interface{}) error {
	var err error
	var data []byte
	if data, err = json.Marshal(obj); err != nil {
		return err
	}
	(*rp)[key] = data
	return nil
}
