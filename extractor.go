package mylittlesalesman_test

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/goccy/go-json"
)

func MarshalExtractor(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("not a struct")
	}

	result := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonKey := field.Tag.Get("json")
		extractor := field.Tag.Get("extractor")
		if jsonKey != "" && extractor != "" {
			result[jsonKey] = extractor
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(result); err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}
