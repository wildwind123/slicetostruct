package slicetostruct

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/go-faster/errors"
)

type T1 struct {
	ID int64
}

type T2 struct {
	ID  int64
	ID2 int64
}

type T3 struct {
	ID  int64 `ss:"id"`
	ID2 int64 `ss:"id2"`
}

type T4 struct {
	ID  int64 `ss:""`
	ID2 int64 `ss:",omitempty"`
	ID3 int64 `ss:",omitempty"`
}

type T5 struct {
	ID  int64 `ss:"id"`
	ID2 int64 `ss:"-"`
	ID3 int64 `ss:"id3"`
}

type T6 struct {
	ID    int64 `ss:"id"`
	ID2   int64 `ss:"id2"`
	ID1_2 int64 `ss:"-"`
}

type T7 struct {
	ID    int64 `ss:"'is, id',omitempty"`
	ID2   int64 `ss:"id2"`
	ID1_2 int64 `ss:"-"`
}

type T8 struct {
	ID    int64 `ss:"'is, id',omitempty"`
	ID2   int64 `ss:"id2"`
	ID1_2 int64 `ss:"-"`
}

type T9 struct {
	ID sql.NullInt64
}

type T10 struct {
	ID int64 `ss:"Id_alias"`
}

type T11 struct {
	ID int64 `ss:"ID_alias"`
}

type T12 struct {
	ID int64 `ss:"id_alias"`
}

type T13 struct {
	ID int64
}

type T14 struct {
	Date sql.NullTime `ss:"date"`
}

type TAll struct {
	ID         int64      `ss:"id"`
	Name       string     `ss:"name"`
	Int        int        `ss:"int"`
	IDNil      *int64     `ss:"id_nil"`
	NameNil    *string    `ss:"name_nil"`
	IntNil     *int       `ss:"int_nil"`
	Float64    float64    `ss:"float_64"`
	Float64Nil *float64   `ss:"float_64_nil"`
	Time       time.Time  `ss:"time"`
	TimeNil    *time.Time `ss:"time_nil"`
}

func TestToStruct(t *testing.T) {
	// test T1
	sliceToStruct := New[T1](Params{})
	res, err := sliceToStruct.ToStruct([]string{"123"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 {
		t.Error("wrong id")
	}
	// test T2
	sliceToStruct2 := New[T2](Params{})
	res2, err := sliceToStruct2.ToStruct([]string{"111"})
	if err != nil {
		t.Error(err)
		return
	}
	if res2.ID != 111 || res2.ID2 != 0 {
		t.Error("wrong result")
		return
	}

	sliceToStruct2.ReturnErrIndexDoesNotExist = true
	_, err = sliceToStruct2.ToStruct([]string{"111"})
	if err == nil {
		t.Error("should has error")
		return
	}

	res2, err = sliceToStruct2.ToStruct([]string{"111", "222"})
	if err != nil {
		t.Error(err)
		return
	}
	if res2.ID != 111 || res2.ID2 != 222 {
		t.Error("wrong result")
		return
	}

}

func TestToStructTag(t *testing.T) {
	// test T1
	sliceToStruct := New[T3](Params{})
	res, err := sliceToStruct.ToStruct([]string{"123"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 {
		t.Error("wrong result")
	}
	if res.ID2 != 0 {
		t.Error("wrong result")
	}
	sliceToStruct.SetFieldNames([]string{"id", "id2"})
	_, err = sliceToStruct.ToStruct([]string{"123"})
	if err == nil {
		t.Error("should has error, different fieldNames and count slice")
		return
	}

	res, err = sliceToStruct.ToStruct([]string{"123", "33"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 || res.ID2 != 33 {
		t.Error("wrong result")
		return
	}

	sliceToStruct.SetFieldNames([]string{"fake", "id", "id2", "ss"})
	res, err = sliceToStruct.ToStruct([]string{"0", "123", "33"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 || res.ID2 != 33 {
		t.Error("wrong result")
		return
	}

}

func TestAllTypes(t *testing.T) {
	sliceToStruct := New[TAll](Params{
		FieldNames: []string{
			"id", "name", "int", "id_nil", "name_nil", "int_nil", "float_64", "float_64_nil", "time", "time_nil",
		},
	})
	res, err := sliceToStruct.ToStruct([]string{"123", "name test", "1", "1232", "name test_2", "12", "23.1", "23.2", "01.01.2012", "03.03.2003"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 ||
		res.Name != "name test" ||
		res.Int != 1 ||
		*res.IDNil != 1232 ||
		*res.NameNil != "name test_2" ||
		*res.IntNil != 12 {
		t.Error("wrong result")
	}
	if res.Time.Unix() != 1325376000 ||
		res.TimeNil.Unix() != 1046649600 {
		t.Error("wrong time")
	}
	fmt.Println(res)
}

func TestOmitEmpty(t *testing.T) {
	sliceToStruct := New[T4](Params{})
	t4, err := sliceToStruct.ToStruct([]string{"1", "", "4"})
	if err != nil {
		t.Error(err)
		return
	}
	if t4.ID != 1 || t4.ID2 != 0 || t4.ID3 != 4 {
		t.Error("wrong result")
		return
	}
	fmt.Println(t4)
}

func TestSkip(t *testing.T) {
	sliceToStruct := New[T5](Params{})
	t4, err := sliceToStruct.ToStruct([]string{"1", "2", "4"})
	if err != nil {
		t.Error(err)
		return
	}
	if t4.ID != 1 || t4.ID2 != 0 || t4.ID3 != 4 {
		t.Error("wrong result")
		return
	}

}

func TestSkip2(t *testing.T) {
	sliceToStruct := New[T6](Params{})

	sliceToStruct.SetFieldNames([]string{"fake", "id", "id2", "ss"})
	res, err := sliceToStruct.ToStruct([]string{"d", "123", "33"})
	if err != nil {
		t.Error(err)
		return
	}
	if res.ID != 123 || res.ID2 != 33 || res.ID1_2 != 0 {
		t.Error("wrong result")
		return
	}
}

func TestGetTags(t *testing.T) {
	res := getTags("test,test1,dddd323")
	if res[0] != "test" || res[1] != "test1" || res[2] != "dddd323" || len(res) != 3 {
		t.Error("wrong tags")
	}
	res = getTags("test")
	if res[0] != "test" || len(res) != 1 {
		t.Error("wrong result")
	}
	res = getTags(`test#,test1,dddd323`)
	if res[0] != `test,test1` || res[1] != "dddd323" || len(res) != 2 {
		t.Error("wrong result")
	}
	res = getTags(`Организация#, у которой прибор учета находится на праве собственности или на ином законном основании,test,123`)
	if len(res) != 3 ||
		res[0] != `Организация, у которой прибор учета находится на праве собственности или на ином законном основании` ||
		res[1] != `test` ||
		res[2] != `123` {
		t.Error("wrong result")
	}
	res = getTags(`test#,test1`)
	if res[0] != `test,test1` || len(res) != 1 {
		t.Error("wrong result")
	}
	res = getTags(`test1,test#,test1`)
	if res[0] != `test1` || res[1] != `test,test1` || len(res) != 2 {
		t.Error("wrong result")
	}
}

type Int64Test struct {
}

func (k *Int64Test) Set(value *ConvertValueParams) error {
	if *value.FieldName == "id" {
		value.Items[value.Index] = "333"
	}
	d := ConvertInt64{}
	err := d.Set(value)
	if err != nil {
		return errors.Wrap(err, "Int64Test")
	}
	return nil
}

func TestSetConverter(t *testing.T) {
	sliceToStruct := New[T6](Params{})

	res, err := sliceToStruct.ToStruct([]string{"1", "123", "33"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res.ID != 1 || res.ID1_2 != 0 || res.ID2 != 123 {
		t.Error("wrong result")
	}
	sliceToStruct.SetConverter("int64", &Int64Test{})
	res, _ = sliceToStruct.ToStruct([]string{"1", "123", "33"})
	if res.ID != 333 || res.ID1_2 != 0 || res.ID2 != 123 {
		t.Error("wrong result")
	}
}

func TestSqlNullInt64(t *testing.T) {
	sliceToStruct := New[T9](Params{})
	res, err := sliceToStruct.ToStruct([]string{"1", "123", "33"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res.ID.Int64 != 1 || !res.ID.Valid {
		t.Error("wrong result")
	}
	fmt.Println(res)
}

func TestNoneCaseSensitive(t *testing.T) {
	sliceToStruct := New[T10](Params{
		NotCaseSensitive: true,
	})
	sliceToStruct.SetFieldNames([]string{"iD_alias"})
	res, err := sliceToStruct.ToStruct([]string{"1"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res.ID != 1 {
		t.Error("wrong result")
	}

	sliceToStruct2 := New[T11](Params{
		NotCaseSensitive: true,
	})
	sliceToStruct2.SetFieldNames([]string{"iD_alias"})
	res2, err := sliceToStruct2.ToStruct([]string{"1"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res2.ID != 1 {
		t.Error("wrong result")
	}

	sliceToStruct3 := New[T12](Params{
		NotCaseSensitive: true,
	})
	sliceToStruct3.SetFieldNames([]string{"iD_alias"})
	res3, err := sliceToStruct3.ToStruct([]string{"1"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res3.ID != 1 {
		t.Error("wrong result")
	}

	sliceToStruct4 := New[T12](Params{
		NotCaseSensitive: true,
	})
	res4, err := sliceToStruct4.ToStruct([]string{"1"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res4.ID != 1 {
		t.Error("wrong result")
	}
}

func TestSqlNullTime(t *testing.T) {
	sliceToStruct := New[T14](Params{})
	res, err := sliceToStruct.ToStruct([]string{"01.02.2002"})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if res.Date.Time.Unix() != 1012521600 {
		t.Error("res.Date is wrong")
	}
}
