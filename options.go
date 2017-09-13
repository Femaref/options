package options

import (
    "reflect"
    "github.com/hashicorp/go-multierror"
    "fmt"
    "github.com/femaref/helper/mapstruct"
)

type Options map[string]interface{}

type Expectation struct {
    Field string
    Type reflect.Type
    Def interface{}
    Optional bool
}

func New() *Expectations {
    return &Expectations{}
}

type Expectations struct {
    expectations []Expectation
}

func (this *Expectations) Expect(field string, t reflect.Type, def interface{}, optional bool) (*Expectations) {
    this.expectations = append(this.expectations, Expectation{Field:field, Type:t, Def: def, Optional: optional})
    return this
}

func (this Expectations) Parse(target interface{}, opt Options) error {
    var errs error


    for _, e := range this.expectations {
        val, fieldPresent := opt[e.Field]
        // field not optional, no default set, not present -> kill it
        if !e.Optional && !fieldPresent {
            errs = multierror.Append(errs, fmt.Errorf("Expected mandatory field %s", e.Field))
            continue
        }

        if e.Optional && !fieldPresent {
            opt[e.Field] = e.Def
        }

        if fieldPresent {
            reft := reflect.TypeOf(val)

            if reft != e.Type && !reft.ConvertibleTo(e.Type) {
                errs = multierror.Append(errs, fmt.Errorf("Expected field %s to contain %v, got %v", e.Field, e.Type, reft.Kind()))
                continue
            }
        }
    }

    if errs != nil {
        return errs
    }

    return mapstruct.MapToStructv2(target, opt)
}
