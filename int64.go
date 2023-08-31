package slicetostruct

import (
	"strconv"

	"github.com/go-faster/errors"
)

type ConvertInt64 struct {
	Value int64
}

func (c *ConvertInt64) Set(value *ConvertValueParams) error {
	v, err := strconv.ParseInt(value.StringValue, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "cant ParseInt, %s", value.StringValue)
	}
	c.Value = v
	value.ReflectValue.SetInt(c.Value)
	return nil
}
