package slicetostruct

import (
	"reflect"
	"strconv"

	"github.com/go-faster/errors"
)

type ConvertInt struct {
	Value int
}

func (c *ConvertInt) Set(value *ConvertValueParams) error {
	v, err := strconv.ParseInt(value.Items[value.Index], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "cant ParseInt, %s", value.Items[value.Index])
	}
	c.Value = int(v)
	value.ReflectValue.Set(reflect.ValueOf(c.Value))
	return nil
}
