// Package nanboxing implements methods to box and unbox stuff in NaN.
//
// Like it is done in Javascript languages to save variables with their type,
// it can be useful to save different things in the NaN scope of float64 type
package nanboxing

import (
	"fmt"
	"math"
	"unsafe"
)

// Box is the nanboxing top level type
type Box float64

// Object
type Object map[string]Box

// Array
type Array []Box

// tag defines the Box types
type Tag uint64

const (
	TagNumber Tag = 0x1
	tagString Tag = 0x2
	TagBool   Tag = 0x3
	TagArray  Tag = 0x4
	TagObject Tag = 0x5
	TagNull   Tag = 0xF
)

/*
 payload width is enough to store :
 - bool
 - int32
 - pointer (http://en.wikipedia.org/wiki/X86-64#Virtual_address_space_details)
   - to array
   - to object
*/
const (
	TagMask     uint64 = 0xF
	NaNMask     uint64 = 0x7FF8000000000000
	PayloadMask uint64 = 0x00007FFFFFFFFFFF
)

const (
	TagShift uint64 = 47
)

/*
 String prints a Box
*/
func (x *Box) String() string {
	return fmt.Sprintf("Box is %X, address is %v\n", *(*uint64)(unsafe.Pointer(x)), (*uint64)(unsafe.Pointer(x)))
}

/*
 Tag retrieve the tag of a box
*/
func (x *Box) Tag() Tag {
	return Tag(((*(*uint64)(unsafe.Pointer(x))) >> TagShift) & TagMask)
}

/*
 SetTag change the tag of a box
*/
func (x *Box) SetTag(t Tag) {
	ui64 := (NaNMask | (uint64(t) << TagShift) | (x.Payload() & PayloadMask))
	*x = *(*Box)(unsafe.Pointer(&ui64))
}

/*
 Payload retrieve the payload of a box
*/
func (x *Box) Payload() uint64 {
	return (*(*uint64)(unsafe.Pointer(x))) & PayloadMask
}

func (x *Box) SetPayload(u uint64) {
	ui64 := NaNMask | (uint64(x.Tag()) << TagShift) | (u & PayloadMask)
	*x = *(*Box)(unsafe.Pointer(&ui64))
}

/*
 NewFloat create a float box
*/
func NewFloat(f float64) Box {
	return Box(f)
}

/*
 ToFloat retrieve the float value of a box
*/
func (x *Box) ToFloat() float64 {
	return float64(*x)
}

/*
 IsFloat64 test if a box type is float
*/
func (x *Box) IsFloat64() bool {
	return !math.IsNaN(float64(*x))
}

/*
 NewNumber create a int32 box
*/
func NewNumber(n int32) Box {
	ui := *(*uint32)(unsafe.Pointer(&n))
	ui64 := NaNMask | (uint64(TagNumber) << TagShift) | (uint64(ui) & PayloadMask)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToNumber retrieve the int32 value of a box
*/
func (x *Box) ToNumber() int32 {
	return int32(x.Payload())
}

/*
 IsNumber test if a box type is number
*/
func (x *Box) IsNumber() bool {
	return !x.IsFloat64() && (x.Tag() == TagNumber)
}

/*
  NewString create a string box

 The string pointer is saved. String memory has to be keep by the box caller
*/
func NewString(s string) Box {
	p := uintptr(unsafe.Pointer(&s))
	ui64 := NaNMask | (uint64(tagString) << TagShift) | (uint64(p) & PayloadMask)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToString retrieve the string pointer of a box
*/
func (x *Box) ToString() string {
	return *(*string)(unsafe.Pointer(uintptr(x.Payload())))
}

/*
 IsString Test if a box type is String
*/
func (x *Box) IsString() bool {
	return !x.IsFloat64() && (x.Tag() == tagString)
}

/*
 NewBool create a bool box
*/
func NewBool(b bool) Box {
	// fmt.Printf("ui 1 is %X address is %v\n", b, &b)
	ui := *(*uint8)(unsafe.Pointer(&b))
	// fmt.Printf("ui 2 is %X address is %v\n", ui, &ui)
	ui64 := NaNMask | (uint64(TagBool) << TagShift) | (uint64(ui) & PayloadMask)
	// fmt.Printf("ui 3 is %X address is %v\n", ui64, &ui64)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToBool retrieve the bool value of a box
*/
func (x *Box) ToBool() bool {
	return !(0 == x.Payload())
}

/*
 IsBool test if a box type is bool
*/
func (x *Box) IsBool() bool {

	return !x.IsFloat64() && (x.Tag() == TagBool)
}

/*
 NewArray create an array box

 The array is saved as a pointer to a slice
*/
func NewArray(a Array) Box {
	p := uintptr(unsafe.Pointer(&a))
	ui64 := NaNMask | (uint64(TagArray) << TagShift) | (uint64(p) & PayloadMask)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToArray retrieve the array pointer of a box
*/
func (x *Box) ToArray() Array {
	return *(*Array)(unsafe.Pointer(uintptr(x.Payload())))
}

/*
 IsArray test if a box type is array
*/
func (x *Box) IsArray() bool {
	return !x.IsFloat64() && (x.Tag() == TagArray)
}

/*
 NewObject create an object box

 The object is saved as a pointer to a map
*/
func NewObject(o Object) Box {
	p := uintptr(unsafe.Pointer(&o))
	ui64 := NaNMask | (uint64(TagObject) << TagShift) | (uint64(p) & PayloadMask)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToObject retrieve the object pointer of a box
*/
func (x *Box) ToObject() Object {
	return *(*Object)(unsafe.Pointer(uintptr(x.Payload())))
}

/*
 IsObject tests if a box type is object
*/
func (x *Box) IsObject() bool {
	return !x.IsFloat64() && (x.Tag() == TagObject)
}

/*
 NewNull create a null box
*/
func NewNull() Box {
	ui64 := NaNMask | (uint64(TagNull) << TagShift) | (0 & PayloadMask)
	return *(*Box)(unsafe.Pointer(&ui64))
}

/*
 ToNull return nil, always nil
*/
func (x *Box) ToNull() unsafe.Pointer {
	return nil
}

/*
 IsNull tests if a box is null
*/
func (x *Box) IsNull() bool {
	return !x.IsFloat64() && (x.Tag() == TagNull)
}

/*
 IsPointer tests if a box is a pointer
*/
func (x *Box) IsPointer() bool {
	return !x.IsFloat64() && ((x.Tag() == TagArray) || (x.Tag() == TagObject))
}

func (x *Box) ToPointer() uintptr {
	return uintptr(x.Payload())
}
