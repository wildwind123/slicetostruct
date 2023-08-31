package slicetostruct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/errors"
)

var ErrIndexDoesNotExist = fmt.Errorf("index of field does not exist")

const keyTag = "ss"
const defaultTimeLayout = "02.01.2006"

type SliceToStruct[T any] struct {
	Params
	fieldNames map[string]int
}

type Params struct {
	// replace comma to dot for float value
	ReplaceCommaToDot          bool
	ReturnErrIndexDoesNotExist bool
	FieldNames                 []string
	converters                 *converters
}

func New[T any](params Params) *SliceToStruct[T] {
	if params.converters == nil {
		params.converters = &converters{}
		params.converters.SetConverter("int64", &ConvertInt64{})
		params.converters.SetConverter("*int64", &ConvertNullInt64{})
		params.converters.SetConverter("int", &ConvertInt{})
		convertSqlValue := ConvertSqlValue{
			params: &params,
		}
		params.converters.SetConverter("sql.NullInt64", &convertSqlValue)
		params.converters.SetConverter("sql.NullFloat64", &convertSqlValue)
		params.converters.SetConverter("sql.NullString", &convertSqlValue)
		params.converters.SetConverter("sql.NullInt32", &convertSqlValue)
		params.converters.SetConverter("sql.NullInt16", &convertSqlValue)
		params.converters.SetConverter("sql.NullByte", &convertSqlValue)
		params.converters.SetConverter("sql.NullBool", &convertSqlValue)
		params.converters.SetConverter("sql.NullTime", &convertSqlValue)
	}

	sTS := &SliceToStruct[T]{
		Params: params,
	}
	sTS.SetFieldNames(params.FieldNames)
	return sTS
}

func (sTS *SliceToStruct[T]) SetConverter(name string, converter Converter) {
	sTS.converters.SetConverter(name, converter)
}

func (sTS *SliceToStruct[T]) SetFieldNames(fieldNames []string) {
	if len(fieldNames) == 0 {
		sTS.fieldNames = nil
		return
	}

	fieldNamesMap := make(map[string]int, len(fieldNames))
	for i := range fieldNames {
		fieldNamesMap[fieldNames[i]] = i
	}
	sTS.fieldNames = fieldNamesMap
}

func (sTS *SliceToStruct[T]) ToStruct(items []string) (*T, error) {
	if len(sTS.fieldNames) > 0 && len(sTS.fieldNames) < len(items) {
		return nil, errors.New("count items greater then fieldNames")
	}

	var val interface{} = new(T)

	pointToStruct := reflect.ValueOf(val)
	curStruct := pointToStruct.Elem()
	kind := curStruct.Kind()
	if kind != reflect.Struct {
		return nil, errors.New("generic type does not struct")
	}

	structType := pointToStruct.Elem().Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldInfo := structType.Field(i)
		sliceFieldName := fieldInfo.Name
		tags := getTags(fieldInfo.Tag.Get(keyTag))
		if len(tags) > 0 && tags[0] != "" {
			sliceFieldName = tags[0]
		}
		if sliceFieldName == "-" {
			continue
		}

		fieldIndex, err := sTS.GetSliceIndexForField(sliceFieldName, i, len(items))
		if err != nil && !errors.Is(err, ErrIndexDoesNotExist) {
			return nil, errors.Wrap(err, "")
		}
		if errors.Is(err, ErrIndexDoesNotExist) {
			if sTS.ReturnErrIndexDoesNotExist {
				return nil, ErrIndexDoesNotExist
			}
			continue
		}

		fieldType := fieldInfo.Type.String()
		field := curStruct.FieldByName(fieldInfo.Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if fieldType[0] == '*' && items[fieldIndex] == "" {
			continue
		}
		if len(tags) > 1 && tags[1] == "omitempty" && items[fieldIndex] == "" {
			continue
		}

		converter, err := sTS.converters.GetConverter(fieldType)
		if err != nil && !errors.Is(err, ErrConverterDoesNotExist) {
			return nil, errors.Wrap(err, "cant sTS.converters.GetConverter")
		}
		if err == nil {
			err = converter.Set(&ConvertValueParams{
				Items:        items,
				Index:        fieldIndex,
				ReflectValue: &field,
				Tags:         tags,
				FieldName:    &sliceFieldName,
				FieltType:    fieldType,
			})
			if err != nil {
				return nil, errors.Wrap(err, "cant converter.Set")
			}
			continue
		}

		switch fieldType {
		case "*int":
			v, err := strconv.ParseInt(items[fieldIndex], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			vInt := int(v)
			field.Set(reflect.ValueOf(&vInt))
		case "string":
			field.SetString(items[fieldIndex])
		case "*string":
			v := items[fieldIndex]
			field.Set(reflect.ValueOf(&v))
		case "float64":
			if sTS.Params.ReplaceCommaToDot {
				items[fieldIndex] = strings.Replace(items[fieldIndex], ",", ".", 1)
			}
			v, err := strconv.ParseFloat(items[fieldIndex], 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(v))
		case "*float64":
			if sTS.Params.ReplaceCommaToDot {
				items[fieldIndex] = strings.Replace(items[fieldIndex], ",", ".", 1)
			}
			v, err := strconv.ParseFloat(items[fieldIndex], 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(&v))
		case "time.Time":
			timeLayout := defaultTimeLayout
			if len(tags) > 2 {
				timeLayout = tags[2]
			}
			t, err := time.Parse(timeLayout, items[fieldIndex])
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(t))
		case "*time.Time":
			timeLayout := defaultTimeLayout
			if len(tags) > 2 {
				timeLayout = tags[2]
			}
			t, err := time.Parse(timeLayout, items[fieldIndex])
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(&t))
		default:
			return nil, errors.New(fmt.Sprintf("type not implement %s", fieldType))
		}

	}
	v := curStruct.Interface().(T)
	return &v, nil
}

func (sTS *SliceToStruct[T]) GetSliceIndexForField(fieldName string, fieldIndex int, lenSlice int) (int, error) {
	if len(sTS.fieldNames) > 0 {
		v, ok := sTS.fieldNames[fieldName]
		if !ok {
			return 0, errors.Errorf("fieldName does not exist on fieldNames. fieldName = %s, fieldNames = %v", fieldName, sTS.fieldNames)
		}
		if v > (lenSlice - 1) {
			return 0, errors.Errorf("fieldName index does not exist on slice, fieldName = %s, index = %d", fieldName, v)
		}

		return v, nil
	}

	if lenSlice < fieldIndex+1 {
		return 0, ErrIndexDoesNotExist
	}
	return fieldIndex, nil
}

func getTags(tagStr string) []string {

	res := strings.Split(tagStr, ",")
	for i := range res {
		if len(res) >= (i+1) && strings.HasSuffix(res[i], `#`) {
			res[i] = res[i][:len(res[i])-1] + "," + res[i+1]

			if len(res) >= (i + 2) {
				res = append(res[:i+1], res[i+2:]...)
				continue
			}
			res = res[:len(res)-1]
		}
	}

	return res
}
