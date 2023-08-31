package slicetostruct

import (
	"fmt"
	"reflect"
	"sync"
)

type ConvertValueParams struct {
	StringValue  string
	ReflectValue *reflect.Value
}

type Converter interface {
	Set(value *ConvertValueParams) error
}

var ErrConverterDoesNotExist = fmt.Errorf("converter does not exist")

type converters struct {
	converters map[string]Converter
	mu         sync.Mutex
}

func (converters *converters) SetConverter(name string, converter Converter) {
	converters.mu.Lock()
	defer converters.mu.Unlock()
	if converters.converters == nil {
		converters.converters = make(map[string]Converter)
	}
	converters.converters[name] = converter
}

func (converters *converters) GetConverter(name string) (Converter, error) {
	converters.mu.Lock()
	defer converters.mu.Unlock()
	converter, ok := converters.converters[name]
	if !ok {
		return nil, ErrConverterDoesNotExist
	}
	return converter, nil
}
