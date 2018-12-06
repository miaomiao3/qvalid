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

### rule 
1. As for bound limit, it means length of string/array/slice/map, and value of numbers(int/uint/float...)
2. If 'in' was set, do not set bound limit  
3. Comma `'` was reserved except of `in` property

|prpp|des|comment|
|---|---|---|
|lt|little than, upper bound limit | u can set lt **or** lte!  |
|lte|little than or equal, upper bound limit| u can set lt **or** lte!  |
|gt|greater than, lower bound limit| u can set gt **or** gte!  |
|gte|greater than or equal, lower bound limit| u can set gt **or** gte!  |
|in|must in one of the list item. item character must be numeric or alpha|If 'in' was set, do not set bound limit |
|attr|when the field is string, it works to some known attribute like email, ip .etc|read code for more|


### supported string attrs as follows, read code more:
```go
const (
	StringTypeEmail        = "email"
	StringTypeAlpha        = "alpha"
	StringTypeUpperAlpha   = "upper_alpha"
	StringTypeLowerAlpha   = "lower_alpha"
	StringTypeAlphaNumeric = "alpha_numeric"
	StringTypeNumeric      = "numeric"
	StringTypeInt          = "int"
	StringTypeFloat        = "float"
	StringTypeHex          = "hex"
	StringTypeAscii        = "ascii"
	StringTypeVisibleAscii = "visible_ascii"
	StringTypeBytes        = "bytes"
	StringTypeBase64       = "base64"
	StringTypeDNS          = "dns"
	StringTypeVersion      = "version"
	StringTypeIp           = "ip"
	StringTypePort         = "port"
	StringTypeURL          = "url"
)
```


## Examples
First, define some struct:
```go

type Dog struct {
	Name      string            `valid:"in=[rose,tulip]" json:"name"`
	Color     string            `valid:"lt=5, gte=3" json:"color"`
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


```

### validate struct with simple field
like Dog
```go

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

```
output:
```sh
validateSimpleField
    illegal input and result:
        isPass:false
        validErrors:
            err:0 --> &{Field:.name Msg:value: not in:[rose tulip]}
            err:1 --> &{Field:.color Msg:expect length >= 3 but get length: 0}
            err:2 --> &{Field:.weight Msg:expect value >= 10 but get value:0}
            err:3 --> &{Field:.clothes Msg:value:0 not in:[1 3 5]}
            err:4 --> &{Field:.NickNames Msg:expect length > 1 but get length: 0}
            err:5 --> &{Field:.Relations Msg:expect length > 1 but get length: 0}
            err:6 --> &{Field:.Email Msg:value: not match attribute:email}

    legal input and result:
        isPass:true

```


### validate embedded struct
like FakeFood
```go

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

```

output

```sh

validateEmbedStruct
    illegal input and result:
        isPass:false
        validErrors:
            err:0 --> &{Field:.Leaf.name Msg:value: not in:[rose tulip]}
            err:1 --> &{Field:.MainLeaf.name Msg:value: not in:[rose tulip]}

    legal input and result:
        isPass:true
        
```


### validate slice embedded struct
like Food
```go

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

```

output

```sh

validateSliceEmbedStruct
    illegal input and result:
        isPass:false
        validErrors:
            err:0 --> &{Field:.Leafs[0].name Msg:value: not in:[rose tulip]}

    legal input and result:
        isPass:true
        
```

### sample of bad tag
like Person
```go

func badTag() {
	bad := &BadTag{}

	isPass, validErrors := qvalid.ValidateStruct(bad)
	fmt.Println("badTag")
	checkAndDumpValidErrors(isPass, validErrors)
}

```

output
```sh

badTag
    illegal input and result:
        isPass:false
        validErrors:
            err:0 --> &{Field:[qvalid] GetConstraintFromTag Msg:lt and lte can't both set}
            err:1 --> &{Field:[qvalid] GetConstraintFromTag Msg:gt and gt can't both set}
            err:2 --> &{Field:[qvalid] GetConstraintFromTag Msg:bound limit and 'in' can't both set}
            err:3 --> &{Field:[qvalid] GetConstraintFromTag Msg:upper and lower bound limit illegal}
            
```

for more details, see example dir.

## TODO:
1. check pointer loop
2. customized field validator