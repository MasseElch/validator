package validator

import (
	sql "database/sql/driver"
	"testing"
	"time"
)

func BenchmarkFieldSuccess(b *testing.B) {

	validate := New()

	s := "1"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(&s, "len=1")
	}
}

func BenchmarkFieldFailure(b *testing.B) {

	validate := New()

	s := "12"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(&s, "len=1")
	}
}

func BenchmarkFieldDiveSuccess(b *testing.B) {

	validate := New()

	m := []string{"val1", "val2", "val3"}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		validate.Var(m, "required,dive,required")
	}
}

func BenchmarkFieldDiveFailure(b *testing.B) {

	validate := New()

	m := []string{"val1", "", "val3"}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(m, "required,dive,required")
	}
}

func BenchmarkFieldCustomTypeSuccess(b *testing.B) {

	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{
		Name: "1",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(val, "len=1")
	}
}

func BenchmarkFieldCustomTypeFailure(b *testing.B) {

	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(val, "len=1")
	}
}

func BenchmarkFieldOrTagSuccess(b *testing.B) {

	validate := New()

	s := "rgba(0,0,0,1)"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(s, "rgb|rgba")
	}
}

func BenchmarkFieldOrTagFailure(b *testing.B) {

	validate := New()

	s := "#000"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Var(s, "rgb|rgba")
	}
}

func BenchmarkStructLevelValidationSuccess(b *testing.B) {

	validate := New()
	validate.RegisterStructValidation(StructValidationTestStructSuccess, TestStruct{})

	tst := &TestStruct{
		String: "good value",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(tst)
	}
}

func BenchmarkStructLevelValidationFailure(b *testing.B) {

	validate := New()
	validate.RegisterStructValidation(StructValidationTestStruct, TestStruct{})

	tst := &TestStruct{
		String: "good value",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(tst)
	}
}

func BenchmarkStructSimpleCustomTypeSuccess(b *testing.B) {

	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{
		Name: "1",
	}

	type Foo struct {
		Valuer   valuer `validate:"len=1"`
		IntValue int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{Valuer: val, IntValue: 7}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructSimpleCustomTypeFailure(b *testing.B) {

	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{}

	type Foo struct {
		Valuer   valuer `validate:"len=1"`
		IntValue int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{Valuer: val, IntValue: 3}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructPartialSuccess(b *testing.B) {

	validate := New()

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Name")
	}
}

func BenchmarkStructPartialFailure(b *testing.B) {

	validate := New()

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "NickName")
	}
}

func BenchmarkStructExceptSuccess(b *testing.B) {

	validate := New()

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Nickname")
	}
}

func BenchmarkStructExceptFailure(b *testing.B) {

	validate := New()

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Name")
	}
}

func BenchmarkStructSimpleCrossFieldSuccess(b *testing.B) {

	validate := New()

	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	test := &Test{
		Start: now,
		End:   then,
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(test)
	}
}

func BenchmarkStructSimpleCrossFieldFailure(b *testing.B) {

	validate := New()

	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)

	test := &Test{
		Start: now,
		End:   then,
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(test)
	}
}

func BenchmarkStructSimpleCrossStructCrossFieldSuccess(b *testing.B) {

	validate := New()

	type Inner struct {
		Start time.Time
	}

	type Outer struct {
		Inner     *Inner
		CreatedAt time.Time `validate:"eqcsfield=Inner.Start"`
	}

	now := time.Now().UTC()

	inner := &Inner{
		Start: now,
	}

	outer := &Outer{
		Inner:     inner,
		CreatedAt: now,
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(outer)
	}
}

func BenchmarkStructSimpleCrossStructCrossFieldFailure(b *testing.B) {

	validate := New()

	type Inner struct {
		Start time.Time
	}

	type Outer struct {
		Inner     *Inner
		CreatedAt time.Time `validate:"eqcsfield=Inner.Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	inner := &Inner{
		Start: then,
	}

	outer := &Outer{
		Inner:     inner,
		CreatedAt: now,
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(outer)
	}
}

func BenchmarkStructSimpleSuccess(b *testing.B) {

	validate := New()

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructSimpleFailure(b *testing.B) {

	validate := New()

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(invalidFoo)
	}
}

func BenchmarkStructSimpleSuccessParallel(b *testing.B) {

	validate := New()

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(validFoo)
		}
	})
}

func BenchmarkStructSimpleFailureParallel(b *testing.B) {

	validate := New()

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(invalidFoo)
		}
	})
}

func BenchmarkStructComplexSuccess(b *testing.B) {

	validate := New()

	tSuccess := &TestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(tSuccess)
	}
}

func BenchmarkStructComplexFailure(b *testing.B) {

	validate := New()

	tFail := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate.Struct(tFail)
	}
}

func BenchmarkStructComplexSuccessParallel(b *testing.B) {

	validate := New()

	tSuccess := &TestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(tSuccess)
		}
	})
}

func BenchmarkStructComplexFailureParallel(b *testing.B) {

	validate := New()

	tFail := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(tFail)
		}
	})
}
