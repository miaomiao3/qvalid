package qvalid

import (
	"fmt"
	"reflect"
	"strings"
)

// result will be equal to `false` if there are any errors.
func ValidateStruct(s interface{}) (bool, []*ValidError) {
	return validateStruct("", s)
}

const systemTips = "[qvalid]"

func validateStruct(path string, s interface{}) (bool, []*ValidError) {
	if s == nil {
		return true, nil
	}
	result := true
	newPath := path + "."
	validErrors := make([]*ValidError, 0)

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		validErrors = append(validErrors, &ValidError{
			Msg: fmt.Sprintf("input must be structs, but get %s", val.Kind()),
		})

		return false, validErrors
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if typeField.PkgPath != "" {
			continue // Private field
		}
		if valueField.Kind() == reflect.Interface {
			valueField = valueField.Elem()
		}
		if (valueField.Kind() == reflect.Struct || (valueField.Kind() == reflect.Ptr && valueField.Elem().Kind() == reflect.Struct)) &&
			typeField.Tag.Get(validTag) != "-" {
			isTypeValid, validErrs := validateStruct(newPath+getTagName(typeField), valueField.Interface())
			if len(validErrs) > 0 {
				validErrors = append(validErrors, validErrs...)
			}
			result = result && isTypeValid
			continue
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()
		}

		isTypeValid, validErrs := typeCheck(newPath, valueField, typeField)
		if len(validErrs) > 0 {
			validErrors = append(validErrors, validErrs...)
		}
		result = result && isTypeValid
	}
	return result, validErrors
}

// don't check invalid value
func typeCheck(path string, v reflect.Value, t reflect.StructField) (isValid bool, validErrors []*ValidError) {
	if !v.IsValid() {
		return false, nil
	}

	validErrors = make([]*ValidError, 0)
	tag := t.Tag.Get(validTag)

	//  if '-',  ignored
	switch tag {
	case "-":
		return true, nil
	}

	constraint, err := GetConstraintFromTag(tag)
	if err != nil {
		validErrors = append(validErrors, &ValidError{
			Field: systemTips + " GetConstraintFromTag",
			Msg:   err.Error(),
		})
		return
	}

	//TODO check loop

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		isPass, validErr := constraint.checkValue(path, v, t)
		if validErr != nil {
			validErrors = append(validErrors, validErr)
		}
		isValid = isPass
		return

	case reflect.Map:
		// map只检查元素数量，因为key的类型不确定，value的元素也不确定
		isPass, validErr := constraint.checkValue(path, v, t)
		if validErr != nil {
			validErrors = append(validErrors, validErr)
		}
		isValid = isPass
		return

	case reflect.Slice, reflect.Array:
		// only trace when slice element is struct
		result := true

		isPass, validErr := constraint.checkValue(path, v, t)
		if validErr != nil {
			validErrors = append(validErrors, validErr)
		}

		result = result && isPass

		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Kind() == reflect.Struct || (v.Index(i).Kind() == reflect.Ptr && v.Index(i).Elem().Kind() == reflect.Struct) {
				isPass, validErrs := validateStruct(path+fmt.Sprintf("%s[%d]", getTagName(t), i), v.Index(i).Interface())
				if len(validErrs) > 0 {
					validErrors = append(validErrors, validErrs...)
				}
				result = result && isPass
			}
		}
		isValid = result
		return
	case reflect.Interface:
		// If the value is an interface then encode its element
		if v.IsNil() {
			return true, nil
		}
		return ValidateStruct(v.Interface())
	case reflect.Ptr:
		// If the value is a pointer then check its element
		if v.IsNil() {
			return true, nil
		}
		return typeCheck(path, v.Elem(), t)
	case reflect.Struct:
		return validateStruct("", v.Interface())
	default:
		validErrors = append(validErrors, &ValidError{
			Msg: "unsupported type",
		})
		return
	}
}

// json tag first
func getTagName(t reflect.StructField) string {
	jsonTagStr := t.Tag.Get("json")
	if len(jsonTagStr) > 0 {
		jsonTagStrs := strings.Split(jsonTagStr, ",")
		if len(jsonTagStrs) > 0 && len(jsonTagStrs[0]) > 0 {
			return jsonTagStrs[0]
		}
	}
	return t.Name
}

type ValidError struct {
	Field string
	Msg   string
}
