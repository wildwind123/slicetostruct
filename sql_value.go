package slicetostruct

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/go-faster/errors"
)

type ConvertSqlValue struct {
	Value sql.NullInt64
}

func (c *ConvertSqlValue) Set(value *ConvertValueParams) error {
	if value.Items[value.Index] == "" {
		return nil
	}

	var err error
	switch value.FieltType {
	case "sql.NullInt64":
		v := sql.NullInt64{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullFloat64":
		v := sql.NullFloat64{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullString":
		v := sql.NullString{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullInt32":
		v := sql.NullInt32{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullInt16":
		v := sql.NullInt16{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullByte":
		v := sql.NullByte{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullBool":
		v := sql.NullBool{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	case "sql.NullTime":
		v := sql.NullTime{}
		err = v.Scan(value.Items[value.Index])
		if err != nil {
			return errors.Wrap(err, "cant c.Value.Scan")
		}
		value.ReflectValue.Set(reflect.ValueOf(v))
	default:
		return errors.New(fmt.Sprintf("field type unknown = %s", value.FieltType))
	}

	return nil
}
