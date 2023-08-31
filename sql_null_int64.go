package slicetostruct

import (
	"database/sql"
	"reflect"

	"github.com/go-faster/errors"
)

type ConvertSqlNullInt64 struct {
	Value sql.NullInt64
}

func (c *ConvertSqlNullInt64) Set(value *ConvertValueParams) error {
	if value.Items[value.Index] == "" {
		return nil
	}
	err := c.Value.Scan(value.Items[value.Index])
	if err != nil {
		return errors.Wrap(err, "cant c.Value.Scan")
	}
	value.ReflectValue.Set(reflect.ValueOf(c.Value))
	return nil
}
