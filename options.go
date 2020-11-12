package options

import (
	"fmt"
	"reflect"

	"github.com/femaref/helper/mapstruct"
	"github.com/hashicorp/go-multierror"
)

type Options map[string]interface{}

type Expectation struct {
	Field    string
	Type     reflect.Type
	Def      interface{}
	Optional bool
}

func New() *Expectations {
	return &Expectations{}
}

type Expectations struct {
	expectations []Expectation
}

func (this *Expectations) Expect(field string, t reflect.Type, def interface{}, optional bool) *Expectations {
	this.expectations = append(this.expectations, Expectation{Field: field, Type: t, Def: def, Optional: optional})
	return this
}

var interface_type = reflect.TypeOf((*interface{})(nil)).Elem()

func (this Expectations) Parse(target interface{}, opt Options) error {
	var errs error

	target_t := reflect.TypeOf(target)
	target_v := reflect.ValueOf(target)

	if target_t.Kind() != reflect.Ptr {
		return fmt.Errorf("target needs to be a pointer to something, is %T", target)
	}

	pted_to_value := reflect.ValueOf(target).Elem()
	pted_to_type := pted_to_value.Type()

	inline_set := pted_to_type.Kind() == reflect.Map

	if inline_set {
		if pted_to_type.Key().Kind() != reflect.String {
			return fmt.Errorf("when targetting map, key needs to be string, is %v", pted_to_type.Key().Kind())
		}

		if !pted_to_type.Elem().ConvertibleTo(interface_type) {
			return fmt.Errorf("when targetting map, value needs to be interface{}, is %v", pted_to_type.Elem())
		}

		if pted_to_value.IsZero() {
			pted_to_value = reflect.MakeMap(pted_to_type)
			target_v.Elem().Set(pted_to_value)
		}
	}

	for _, e := range this.expectations {
		val, fieldPresent := opt[e.Field]
		// field not optional, no default set, not present -> kill it
		if !e.Optional && !fieldPresent {
			errs = multierror.Append(errs, fmt.Errorf("Expected mandatory field %s", e.Field))
			continue
		}

		if e.Optional && !fieldPresent {
			opt[e.Field] = e.Def

			if inline_set {
				pted_to_value.SetMapIndex(reflect.ValueOf(e.Field), reflect.ValueOf(e.Def))
			}
		}

		if fieldPresent {
			reft := reflect.TypeOf(val)

			if reft != e.Type && !reft.ConvertibleTo(e.Type) {
				errs = multierror.Append(errs, fmt.Errorf("Expected field %s to contain %v, got %v", e.Field, e.Type, reft.Kind()))
				continue
			}

			if inline_set {
				pted_to_value.SetMapIndex(reflect.ValueOf(e.Field), reflect.ValueOf(val).Convert(e.Type))
			}
		}
	}

	if errs != nil {
		return errs
	}

	if inline_set {
		return nil
	}

	switch pted_to_type.Kind() {
	case reflect.Struct:
		return mapstruct.MapToStructv2(target, opt)
	}
	return fmt.Errorf("can't parse into %T", target)

}
