package main

import (
	"fmt"
	"qvalid"
)

func main() {
	validate1Field()
	validateEmbedStruct()
	validateSliceEmbedStruct()
	mixField()
}

type Person struct {
	Name          string   `json:"name" valid:"lt=10, gt=1"`
	from          string   `json:"from" valid:"lt=10, gt=1"` // unexported, will be ignore
	Age           int      `json:"age" valid:"lt=30, gt=20"`
	AddrSyntaxErr string   `valid:"lt=10, gt=1, in=[aa,bb]"` // this will cause [qvalid] error msg
	Addr          string   `valid:"in=[aa,bb]"`
	Email         string   `valid:"attr=email"`
	Weight        int      `valid:"lt=10, gt=1"`
	Nicks         []string `valid:"lt=5, gt=1"`
	Food          Food
	PFood         *Food
	Foods         []Food `json:"foods"  valid:"gt=1"`
}

type Food struct {
	Protein string `valid:"lt=10, gt=1"`
	Leafs   []Leaf `valid:"gte=1"`
}
type FakeFood struct {
	Protein  string `valid:"lt=10, gt=1"`
	MainLeaf Leaf   `valid:"gte=1"`
}

type Leaf struct {
	Color string `valid:"lt=5, gt=2" json:"color"`
}

func validate1Field() {
	leaf := &Leaf{}
	isPass, validErrors := qvalid.ValidateStruct(leaf)
	fmt.Println("validate1Field isPass:%v", isPass)
	dumpValidErrors(validErrors)
}

func validateEmbedStruct() {
	food := &FakeFood{}
	isPass, validErrors := qvalid.ValidateStruct(food)
	fmt.Println("validateEmbedStruct isPass:%v", isPass)
	dumpValidErrors(validErrors)
}

func validateSliceEmbedStruct() {
	food := &Food{
		Leafs: []Leaf{
			Leaf{},
		},
	}
	isPass, validErrors := qvalid.ValidateStruct(food)
	fmt.Println("validateSliceEmbedStruct isPass:%v", isPass)
	dumpValidErrors(validErrors)
}

func mixField() {
	p := &Person{
		Name: "",
		Age:  0,
		Food: Food{
			Leafs: []Leaf{
				Leaf{},
			},
		},
		Foods: []Food{{Protein: ""}},
	}

	isPass, validErrors := qvalid.ValidateStruct(p)
	fmt.Println("mixField isPass:%v", isPass)
	dumpValidErrors(validErrors)
}

func dumpValidErrors(validErrors []*qvalid.ValidError) {
	fmt.Println("validErrors:")
	for k, v := range validErrors {
		fmt.Printf("err:%d --> %+v\n", k, v)
	}
	fmt.Println("***************")
}
