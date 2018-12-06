package qvalid

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"regexp"
	"strings"
)

// check bound
func (c *Constraint) checkBoundLimit(value float64, isLength bool) (bool, error) {
	upperLimitPass := false
	lowLimitPass := false

	if c.Gt != nil {
		if value > *c.Gt {
			lowLimitPass = true
		} else {
			if isLength {
				return false, fmt.Errorf("expect length > %v but get length: %v", *c.Gt, value)
			} else {
				return false, fmt.Errorf("expect value > %v but get value:%v", *c.Gt, value)
			}

		}
	}
	if c.Gte != nil {
		if value >= *c.Gte {
			lowLimitPass = true
		} else {
			if isLength {
				return false, fmt.Errorf("expect length >= %v but get length: %v", *c.Gte, value)
			} else {
				return false, fmt.Errorf("expect value >= %v but get value:%v", *c.Gte, value)
			}
		}
	}

	if c.Lt != nil {
		if value < *c.Lt {
			upperLimitPass = true
		} else {
			if isLength {
				return false, fmt.Errorf("expect length < %v but get length: %v", *c.Lt, value)
			} else {
				return false, fmt.Errorf("expect value < %v but get value:%v", *c.Lt, value)
			}
		}
	}

	if c.Lte != nil {
		if value <= *c.Lte {
			upperLimitPass = true
		} else {
			if isLength {
				return false, fmt.Errorf("expect length <= %v but get length: %v", *c.Lte, value)
			} else {
				return false, fmt.Errorf("expect value <= %v but get value:%v", *c.Lte, value)
			}
		}
	}

	switch {
	case upperLimitPass && lowLimitPass:
		return true, nil
	case !upperLimitPass && lowLimitPass:
		if c.hasUpperBoundLimit() {
			if isLength {
				return false, fmt.Errorf("length:%v fix low bound but miss upper bound", value)
			} else {
				return false, fmt.Errorf("value:%v fix low bound but miss upper bound", value)
			}
		} else {
			return true, nil
		}
	case upperLimitPass && !lowLimitPass:
		if c.hasLowBoundLimit() {
			if isLength {
				return false, fmt.Errorf("length:%v fix upper bound but miss low bound", value)
			} else {
				return false, fmt.Errorf("value:%v fix upper bound but miss low bound", value)
			}
		} else {
			return true, nil
		}
	case !upperLimitPass && !lowLimitPass:
		// check if has limit
		if c.hasBoundLimit() {
			return false, fmt.Errorf("value:%v logic not expected", value)
		} else {
			return true, nil
		}
	}
	return false, errors.New("not expected")

}

func (c *Constraint) hasBoundLimit() bool {
	if c.Gt != nil || c.Gte != nil || c.Lt != nil || c.Lte != nil {
		return true
	}
	return false
}

func (c *Constraint) hasLowBoundLimit() bool {
	if c.Gt != nil || c.Gte != nil {
		return true
	}
	return false
}

func (c *Constraint) hasUpperBoundLimit() bool {
	if c.Lt != nil || c.Lte != nil {
		return true
	}
	return false
}

func (c *Constraint) getLowBoundLimit() float64 {
	if c.Gt != nil {
		return *c.Gt
	}
	if c.Gte != nil {
		return *c.Gte
	}
	return 0
}

func (c *Constraint) getUpperBoundLimit() float64 {
	if c.Lt != nil {
		return *c.Lt
	}
	if c.Lte != nil {
		return *c.Lte
	}
	return 0
}

// for map string slice array, check length
// for number, check value
func (c *Constraint) checkValue(path string, v reflect.Value, t reflect.StructField) (bool, *ValidError) {
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		_, err := c.checkBoundLimit(float64(v.Len()), true)
		if err != nil {
			return false, &ValidError{
				Field: path + getTagName(t),
				Msg:   err.Error(),
			}
		}

		if v.Kind() == reflect.String {
			value := fmt.Sprintf("%v", v)
			// check attribute
			if c.Attr != nil {
				if regex, ok := stringRegexMap[*c.Attr]; ok {
					value := fmt.Sprintf("%v", v)
					isMatch := regex.Match([]byte(value))
					if !isMatch {
						return false, &ValidError{
							Field: path + getTagName(t),
							Msg:   fmt.Sprintf("value:%s not match attribute:%s", value, *c.Attr),
						}
					}
				}
			}

			if len(c.In) > 0 {
				if !isInStringSlice(value, c.In) {
					return false, &ValidError{
						Field: path + getTagName(t),
						Msg:   fmt.Sprintf("value:%s not in:%v", value, c.In),
					}
				}
			}
		}

	case reflect.Bool, reflect.Uintptr:
		return true, nil // ignore bool check
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := c.checkBoundLimit(float64(v.Int()), false)
		if err != nil {
			return false, &ValidError{
				Field: path + getTagName(t),
				Msg:   err.Error(),
			}
		}
		value := fmt.Sprintf("%v", v)
		if len(c.In) > 0 {
			if !isInStringSlice(value, c.In) {
				return false, &ValidError{
					Field: path + getTagName(t),
					Msg:   fmt.Sprintf("value:%s not in:%v", value, c.In),
				}
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err := c.checkBoundLimit(float64(v.Uint()), false)
		if err != nil {
			return false, &ValidError{
				Field: path + getTagName(t),
				Msg:   err.Error(),
			}
		}
	case reflect.Float32, reflect.Float64:
		_, err := c.checkBoundLimit(float64(v.Float()), false)
		if err != nil {
			return false, &ValidError{
				Field: path + getTagName(t),
				Msg:   err.Error(),
			}
		}
	case reflect.Interface, reflect.Ptr:
		return true, nil // ignore interface
	}

	// not required and empty is valid
	return true, nil
}

type Constraint struct {
	Lt     *float64 `yaml:"lt"`
	Lte    *float64 `yaml:"lte"`
	Gt     *float64 `yaml:"gt"`
	Gte    *float64 `yaml:"gte"`
	Equal  *float64 `yaml:"eq"`
	In     []string `yaml:"in"`
	Prefix *string  `yaml:"prefix"`
	Suffix *string  `yaml:"suffix"`
	Regex  *string  `yaml:"regex"`
	Attr   *string  `yaml:"attr"`
}

var bracketMarkRegex = `\[[a-zA-Z1-9-, ]+\]`
var bracketHolder = `bholder1xx`
var yamlTag = `: `

// get constraint from tag
func GetConstraintFromTag(tag string) (*Constraint, error) {
	c := Constraint{}

	// get square bracket content of 'in'
	reg := regexp.MustCompile(bracketMarkRegex)
	squareBracketContent := reg.FindString(tag)

	transformData := strings.Replace(tag, squareBracketContent, bracketHolder, 1)

	transformDatas := strings.Split(transformData, ",")

	trimTransformData := ""
	for _, v := range transformDatas {
		trimTransformData += strings.Trim(v, " ") + "\n"
	}

	// fill back
	trimTransformData = strings.Replace(trimTransformData, bracketHolder, squareBracketContent, 1)

	trimTransformData = strings.Replace(trimTransformData, "=", yamlTag, -1)

	err := yaml.Unmarshal([]byte(trimTransformData), &c)

	if c.Lte != nil && c.Lt != nil {
		return nil, errors.New("lt and lte can't both set")
	}
	if c.Gte != nil && c.Gt != nil {
		return nil, errors.New("gt and gt can't both set")
	}

	if c.hasLowBoundLimit() && c.hasUpperBoundLimit() {
		if c.getLowBoundLimit() >= c.getUpperBoundLimit() {
			return nil, errors.New("upper and lower bound limit illegal")
		}
	}

	if c.hasBoundLimit() && len(c.In) > 1 {
		return nil, errors.New("bound limit and 'in' can't both set")
	}

	return &c, err
}

func isInStringSlice(s string, data []string) bool {
	for _, v := range data {
		if v == s {
			return true
		}
	}
	return false
}
