package qvalid

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestGetConstraintFromTag(t *testing.T) {
	Convey("TestGetConstraintFromTag", t, func() {
		var originErrData = `lt=1, lte=1.0, gt=2, in=[a,b]`
		_, err := GetConstraintFromTag(originErrData)
		So(err, ShouldNotBeNil)

		var originData = `lt=5, gt=1, attr=email, in=[a,b]`
		c, err := GetConstraintFromTag(originData)
		So(err, ShouldEqual, nil)

		Convey("test constraint number", func() {

			var i int
			var value reflect.Value
			var isPass bool

			Convey("lt=5, gt=1 test value=2", func() {
				i = 2
				value = reflect.ValueOf(i)
				isPass, err := c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			Convey("lt=5, gt=1 test value=1", func() {
				i = 1
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			Convey("lt=5, gt=1 test value=5", func() {
				i = 5
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			originData = `lte=5, gte=1, attr=email, in=[a,b]`
			c, err = GetConstraintFromTag(originData)
			So(err, ShouldEqual, nil)

			Convey("lte=5, gte=1 test value=1", func() {
				i = 1
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			Convey("lte=5, gte=1 test value=5", func() {
				i = 5
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			originData = `lte=5, gt=1, attr=email, in=[a,b]`
			c, err = GetConstraintFromTag(originData)
			So(err, ShouldEqual, nil)

			Convey("lte=5, gt=1 test value=1", func() {
				i = 1
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			Convey("lte=5, gt=1 test value=5", func() {
				i = 5
				value = reflect.ValueOf(i)
				isPass, err = c.checkBoundLimit(float64(value.Int()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})
		})
		Convey("test constraint string", func() {

			var str string
			var value reflect.Value
			var isPass bool

			Convey("lt=5, gt=1 test value=ab", func() {
				str = "ab"
				value = reflect.ValueOf(str)
				isPass, err := c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			Convey("lt=5, gt=1 test value=a", func() {
				str = "a"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			Convey("lt=5, gt=1 test value=abcde", func() {
				str = "abcde"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			originData = `lte=5, gte=1, attr=email, in=[a,b]`
			c, err = GetConstraintFromTag(originData)
			So(err, ShouldEqual, nil)

			Convey("lte=5, gte=1 test value=a", func() {
				str = "a"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			Convey("lte=5, gte=1 test value=abcde", func() {
				str = "abcde"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})

			originData = `lte=5, gt=1, attr=email, in=[a,b]`
			c, err = GetConstraintFromTag(originData)
			So(err, ShouldEqual, nil)

			Convey("lte=5, gt=1 test value=a", func() {
				str = "a"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldNotBeNil)
				So(isPass, ShouldEqual, false)
			})

			Convey("lte=5, gt=1 test value=abcde", func() {
				str = "abcde"
				value = reflect.ValueOf(str)
				isPass, err = c.checkBoundLimit(float64(value.Len()))
				So(err, ShouldBeNil)
				So(isPass, ShouldEqual, true)
			})
		})
	})
}
