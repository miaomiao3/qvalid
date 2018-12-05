## qvalid
mini and helpful tool to validate struct's exported fields


## Feature
1. validate field value of numbers(int/uint/float...)
2. validate field length of string/array/slice/map
3. support **in** check
4. when a field is slice and its element is struct/struct_pointer, qvalid auto validate this struct related element
5. when a field is string, support attribute check. e.g. email/ip/email... 
6. pretty field output msg, use json tag first as field name

## Install
`
go get -u github.com/miaomiao3/qvalid
`

## Syntax

As for bound limit, it means length of string/array/slice/map, and value of numbers(int/uint/float...)
If 'in' was set, do not set bound limit  
Comma `'` was reserved except of `in` property

|prpp|des|comment|
|---|---|---|
|lt|little than, upper bound limit | u can set lt **or** lte!  |
|lte|little than or equal, upper bound limit| u can set lt **or** lte!  |
|gt|greater than, lower bound limit| u can set gt **or** gte!  |
|gte|greater than or equal, lower bound limit| u can set gt **or** gte!  |
|in|must in one of list item. |If 'in' was set, do not set bound limit |
|attr|then the field is string, it works to some known attribute like email, ip .etc|read code for more|
|custom|customize callback, must be unique in one struct|read code for more|




## Examples
First, define some struct:
```go

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

```

### validate struct with one filed
like Leaf
```go

func validate1Field() {
	leaf := &Leaf{}
	isPass, validErrors := qvalid.ValidateStruct(leaf)
	fmt.Println("validate1Field isPass:%v", isPass)
	dumpValidErrors(validErrors)
}
```
output:
```sh
validate1Field isPass:%v false
validErrors:
err:0 --> &{Field:.color Msg:expect length > 2 but get length: 0}
***************

```


### validate embedded struct
like FakeFood
```go

func validateEmbedStruct() {
	food := &FakeFood{}
	isPass, validErrors := qvalid.ValidateStruct(food)
	fmt.Println("validateEmbedStruct isPass:%v", isPass)
	dumpValidErrors(validErrors)
}
```

output

```sh
validateEmbedStruct isPass:%v false
validErrors:
err:0 --> &{Field:.Protein Msg:expect length > 1 but get length: 0}
err:1 --> &{Field:.MainLeaf.color Msg:expect length > 2 but get length: 0}
***************
```


### validate slice embedded struct
like Food
```go
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
```

output

```sh
validateSliceEmbedStruct isPass:%v false
validErrors:
err:0 --> &{Field:.Protein Msg:expect length > 1 but get length: 0}
err:1 --> &{Field:.Leafs[0].color Msg:expect length > 2 but get length: 0}
***************
```

### mix validate
like Person
```go

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

```

output
```go
mixField isPass:%v false
validErrors:
err:0 --> &{Field:.name Msg:expect length > 1 but get length: 0}
err:1 --> &{Field:.age Msg:expect value > 20 but get value:0}
err:2 --> &{Field:[qvalid] GetConstraintFromTag Msg:bound limit and 'in' can't both set}
err:3 --> &{Field:.Addr Msg:value: not in:[aa bb]}
err:4 --> &{Field:.Email Msg:value: not match attr:email}
err:5 --> &{Field:.Weight Msg:expect value > 1 but get value:0}
err:6 --> &{Field:.Nicks Msg:expect length > 1 but get length: 0}
err:7 --> &{Field:.Food.Protein Msg:expect length > 1 but get length: 0}
err:8 --> &{Field:.Food.Leafs[0].color Msg:expect length > 2 but get length: 0}
err:9 --> &{Field:.foods Msg:expect length > 1 but get length: 1}
err:10 --> &{Field:.foods[0].Protein Msg:expect length > 1 but get length: 0}
err:11 --> &{Field:.foods[0].Leafs Msg:expect length >= 1 but get length: 0}
***************

```

for more details, see example dir.

## TODO:
1. check loop pointer check
2. customer callback(now developing)