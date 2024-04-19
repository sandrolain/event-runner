package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack"
)

func Marshal(t string, v any) ([]byte, error) {
	t = strings.ToLower(t)
	switch t {
	case "json":
		return json.Marshal(v)
	case "msgpack":
		return msgpack.Marshal(v)
	case "gob":
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(v)
		return buf.Bytes(), err
	}
	return nil, fmt.Errorf("unsupported type: %s", t)
}

func Unmarshal(t string, b []byte, v any) error {
	t = strings.ToLower(t)
	switch t {
	case "json":
		return json.Unmarshal(b, v)
	case "msgpack":
		return msgpack.Unmarshal(b, v)
	case "gob":
		dec := gob.NewDecoder(bytes.NewBuffer(b))
		return dec.Decode(v)
	}
	return fmt.Errorf("unsupported type: %s", t)
}
