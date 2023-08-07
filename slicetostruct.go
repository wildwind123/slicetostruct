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
	ReturnErrIndexDoesNotExist bool
	FieldNames                 []string
}

func New[T any](params Params) *SliceToStruct[T] {
	sTS := &SliceToStruct[T]{
		Params: params,
	}
	sTS.SetFieldNames(params.FieldNames)
	return sTS
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

		switch fieldType {
		case "int64":
			v, err := strconv.ParseInt(items[fieldIndex], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.SetInt(v)
		case "*int64":
			if items[fieldIndex] == "" {
				continue
			}
			v, err := strconv.ParseInt(items[fieldIndex], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(&v))
		case "int":
			v, err := strconv.ParseInt(items[fieldIndex], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(int(v)))
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
			v, err := strconv.ParseFloat(items[fieldIndex], 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			field.Set(reflect.ValueOf(v))
		case "*float64":
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
	res := []string{}

	res = strings.Split(tagStr, ",")
	for i := range res {
		if len(res) >= (i+1) && strings.HasSuffix(res[i], `\`) {
			res[i] = res[i][:len(res[i])-1] + res[i+1]

			fmt.Println("sss")
		}
	}

	return res
}
