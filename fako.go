package fako

import (
	"reflect"

	"github.com/serenize/snaker"
)

//Fill fills all the fields that have a fako: tag
func Fill(elems ...interface{}) {
	for _, elem := range elems {
		FillElem(elem)
	}
}

//FillElem provides a way to fill a simple interface
func FillElem(strukt interface{}) {
	fillWithDetails(strukt, []string{}, []string{})
}

//FillOnly fills fields that have a fako: tag and its name is on the second argument array
func FillOnly(strukt interface{}, fields ...string) {
	fillWithDetails(strukt, fields, []string{})
}

//FillExcept fills fields that have a fako: tag and its name is not on the second argument array
func FillExcept(strukt interface{}, fields ...string) {
	fillWithDetails(strukt, []string{}, fields)
}

//FillByMap fills all the fields that by using maps
func FillByMap(strukt interface{},typemap map[string]string) {
	fillWithDetailsAndMap(strukt, []string{}, []string{}, typemap)
}

func fillWithDetails(strukt interface{}, only []string, except []string) {
	fillWithDetailsAndMap(strukt, only, except, make(map[string]string))
}

func fillWithDetailsAndMap(strukt interface{}, only []string, except []string, typemap map[string]string) {
	elem := reflect.ValueOf(strukt).Elem()
	elemT := reflect.TypeOf(strukt).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldt := elemT.Field(i)
		fakeType := fieldt.Tag.Get("fako")

		if fakeType == "" {
			if ftype, ok := typemap[fieldt.Name]; ok {
				fakeType = ftype
			}
		}

		if fakeType != "" {
			fakeType = snaker.SnakeToCamel(fakeType)
			function := findFakeFunctionFor(fakeType)

			inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
			notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

			if field.CanSet() && fakeType != "" && inOnly && notInExcept {
				field.SetString(function())
			}
		}

	}
}
