package main

import (
	"fmt"
	"github.com/miaomiao3/qvalid"
)

func main() {
	validateSimpleField()
	validateEmbedStruct()
	validateSliceEmbedStruct()
	badTag()
}

type Dog struct {
	Name      string            `valid:"in=[rose,tulip]" json:"name"`
	Color     string            `valid:"lt=5, gte=4" json:"color"`
	Weight    float64           `valid:"lt=100, gte=10" json:"weight"`
	Clothes   int               `valid:"in=[1,3,5]" json:"clothes"`
	NickNames []string          `valid:"lt=5, gt=1"`
	Relations map[string]string `valid:"lt=5, gt=1"`
	Email     string            `valid:"attr=email"`
	from      string            `json:"from" valid:"lt=10, gt=1"` // unexported, will be ignored by qvalid
}

type BadTag struct {
	Err1 string `valid:"lt=10, lte=1"`            // this will cause [qvalid] error msg
	Err2 string `valid:"gt=10, gte=1"`            // this will cause [qvalid] error msg
	Err3 string `valid:"lt=10, gt=1, in=[aa,bb]"` // this will cause [qvalid] error msg
	Err4 string `valid:"lt=1, gte=1"`             // this will cause [qvalid] error msg
}

type FakeFood struct {
	Leaf     Leaf
	MainLeaf *Leaf
}

type Food struct {
	Leafs []Leaf `valid:"gte=1"`
}

type Leaf struct {
	Name string `valid:"in=[rose,tulip]" json:"name"`
}

func validateSimpleField() {
	dog := &Dog{}
	isPass, validErrors := qvalid.ValidateStruct(dog)
	fmt.Println("validateSimpleField")
	checkAndDumpValidErrors(isPass, validErrors)

	newFlower := &Dog{
		Name:      "rose",
		Color:     "gray",
		Weight:    30.0,
		Clothes:   3,
		NickNames: []string{"wangcai", "dawang"},
		Relations: map[string]string{
			"owner": "cy",
			"birth": "2018",
		},
		Email: "google@gmail.com",
	}
	isPass, validErrors = qvalid.ValidateStruct(newFlower)
	checkAndDumpValidErrors(isPass, validErrors)
}

func validateEmbedStruct() {
	food := &FakeFood{
		MainLeaf: &Leaf{},
	}
	isPass, validErrors := qvalid.ValidateStruct(food)
	fmt.Println("validateEmbedStruct")
	checkAndDumpValidErrors(isPass, validErrors)

	newFakeFood := FakeFood{
		Leaf: Leaf{
			Name: "rose",
		},
		MainLeaf: &Leaf{
			Name: "rose",
		},
	}
	isPass, validErrors = qvalid.ValidateStruct(newFakeFood)
	checkAndDumpValidErrors(isPass, validErrors)
}

func validateSliceEmbedStruct() {
	food := &Food{
		Leafs: []Leaf{ // if Leafs is empty, qvalid do not check empty slice field, so set 1 element to test
			Leaf{},
		},
	}
	isPass, validErrors := qvalid.ValidateStruct(food)
	fmt.Println("validateSliceEmbedStruct")
	checkAndDumpValidErrors(isPass, validErrors)

	newFood := Food{
		Leafs: []Leaf{
			Leaf{
				Name: "rose",
			},
		},
	}
	isPass, validErrors = qvalid.ValidateStruct(newFood)
	checkAndDumpValidErrors(isPass, validErrors)
}

func badTag() {
	bad := &BadTag{}

	isPass, validErrors := qvalid.ValidateStruct(bad)
	fmt.Println("badTag")
	checkAndDumpValidErrors(isPass, validErrors)
}

func checkAndDumpValidErrors(isPass bool, validErrors []*qvalid.ValidError) {
	if !isPass {
		fmt.Println("    illegal input and result:")
	} else {
		fmt.Println("    legal input and result:")
	}
	fmt.Printf("        isPass:%v\n", isPass)
	if len(validErrors) > 0 {
		fmt.Println("        validErrors:")
		for k, v := range validErrors {
			fmt.Printf("            err:%d --> %+v\n", k, v)
		}
		fmt.Println("")
	}

}
