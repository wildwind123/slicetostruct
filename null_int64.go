package slicetostruct

import (
	"reflect"
	"strconv"

	"github.com/go-faster/errors"
)

type ConvertNullInt64 struct {
	Value *int64
}

func (c *ConvertNullInt64) Set(value *ConvertValueParams) error {
	if value.Items[value.Index] == "" {
		return nil
	}
	v, err := strconv.ParseInt(value.Items[value.Index], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "cant ParseInt, %s", value.Items[value.Index])
	}
	c.Value = &v
	value.ReflectValue.Set(reflect.ValueOf(c.Value))
	return nil
}
