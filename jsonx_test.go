package jsonx_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/dtcookie/jsonx"
	"github.com/dtcookie/opt"
)

type verysimple struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (vs *verysimple) MarshalJSON() ([]byte, error) {
	rawProperties := jsonx.RawProperties{}
	if err := rawProperties.Marshal("a", vs.A); err != nil {
		return nil, err
	}
	if err := rawProperties.Marshal("b", vs.B); err != nil {
		return nil, err
	}
	return json.Marshal(rawProperties)
}

func (vs *verysimple) UnmarshalJSON(data []byte) error {
	rawProperties := jsonx.RawProperties{}
	if err := json.Unmarshal(data, &rawProperties); err != nil {
		return err
	}
	if err := rawProperties.Unmarshal("a", &vs.A); err != nil {
		return err
	}
	if err := rawProperties.Unmarshal("b", &vs.B); err != nil {
		return err
	}
	return nil
}

func TestWithoutUnknownProperties(t *testing.T) {
	input := `{ "a": "0123", "b": "abcd", "c": "foo" }`
	obj := new(verysimple)
	if err := json.Unmarshal([]byte(input), obj); err != nil {
		t.Fatal(err)
	}
	if obj.A != "0123" {
		t.Fatal(fmt.Errorf("simple.A didn't get unmarshalled properly. expected: %s, actual: %s", "0123", obj.A))
	}
	if obj.B != "abcd" {
		t.Fatal(fmt.Errorf("simple.B didn't get unmarshalled properly. expected: %s, actual: %s", "abcd", obj.B))
	}
	var data []byte
	var err error
	if data, err = json.Marshal(obj); err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{}
	if err = json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatal(fmt.Errorf("number of marshalled properties doesn't match. expected: %d, actual: %d", 2, len(m)))
	}
	if value, found := m["a"]; !found {
		t.Fatal(fmt.Errorf("property %s didn't get serialized", "a"))
	} else if stringValue, ok := value.(string); !ok {
		t.Fatal(fmt.Errorf("property %s was expected to be a string. actual: %T", "a", value))
	} else if stringValue != "0123" {
		t.Fatal(fmt.Errorf("property %s didn't get marshalled correctly. expected: %s, actual: %s", "a", "0123", stringValue))
	}
	if value, found := m["b"]; !found {
		t.Fatal(fmt.Errorf("property %s didn't get serialized", "b"))
	} else if stringValue, ok := value.(string); !ok {
		t.Fatal(fmt.Errorf("property %s was expected to be a string. actual: %T", "b", value))
	} else if stringValue != "abcd" {
		t.Fatal(fmt.Errorf("property %s didn't get marshalled correctly. expected: %s, actual: %s", "b", "abcd", stringValue))
	}
}

type simple struct {
	Value    *string        `json:"value"`
	Unknowns jsonx.Unknowns `json:"-"`
}

func (s *simple) MarshalJSON() ([]byte, error) {
	rawProperties := jsonx.RawProperties{}
	if err := rawProperties.Marshal("value", s.Value); err != nil {
		return nil, err
	}
	rawProperties.MarshalUnknowns(s.Unknowns)
	return json.Marshal(rawProperties)
}

func (s *simple) UnmarshalJSON(data []byte) error {
	rawProperties := jsonx.RawProperties{}
	if err := json.Unmarshal(data, &rawProperties); err != nil {
		return err
	}
	if err := rawProperties.Unmarshal("value", &s.Value); err != nil {
		return err
	}
	rawProperties.UnmarshalUnknowns(&s.Unknowns)
	return nil
}

func TestNoUnknownPropertiesSpecified(t *testing.T) {
	input := `{ "value": "0123" }`
	obj := new(simple)
	if err := json.Unmarshal([]byte(input), obj); err != nil {
		t.Fatal(err)
	}
	if *obj.Value != "0123" {
		t.Fatal(fmt.Errorf("obj.Value didn't get unmarshalled properly. expected: %s, actual: %s", "0123", opt.String(obj.Value)))
	}
	if obj.Unknowns != nil {
		t.Fatal("obj.Unknowns was expected to be nil")
	}
}

func TestUnknownProperties(t *testing.T) {
	input := `{ "value": "0123", "addProperty": "abcd" }`
	obj := new(simple)
	if err := json.Unmarshal([]byte(input), obj); err != nil {
		t.Fatal(err)
	}
	if *obj.Value != "0123" {
		t.Fatal(fmt.Errorf("obj.Value didn't get unmarshalled properly. expected: %s, actual: %s", "0123", opt.String(obj.Value)))
	}
	if _, ok := obj.Unknowns["addProperty"]; !ok {
		t.Fatal(fmt.Errorf("expected to find unknown property '%s'. none found.", "addProperty"))
	}
	if string(obj.Unknowns["addProperty"]) != "\"abcd\"" {
		t.Fatal(fmt.Errorf("unknown property '%s' didn't get unmarshalled property. expected: %s, actual: %s", "addProperty", "\"abcd\"", string(obj.Unknowns["addProperty"])))
	}
}

type Foo interface {
	String() string
}

type Bar struct {
	Value string
}

func (b *Bar) String() string {
	return b.Value
}

func TestASDF(t *testing.T) {
	var foo Foo
	foo = &Bar{Value: "asdf"}
	var data []byte
	var err error
	if data, err = json.Marshal(foo); err != nil {
		t.Fatal(err)
	}
	foo = new(Bar)
	if err = json.Unmarshal(data, foo); err != nil {
		t.Fatal(err)
	}
	log.Println("value", foo.String())
}
