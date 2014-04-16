package nanboxing

import (
	"math"
	"testing"
)

func TestBoolBox(t *testing.T) {
	trueBox := NewBool(true)
	falseBox := NewBool(false)
	if false == math.IsNaN(float64(trueBox)) {
		t.Error("The box should be NaN")
	}
	if false == math.IsNaN(float64(falseBox)) {
		t.Error("The box should be NaN")
	}
	if true != trueBox.ToBool() {
		t.Error("trueBox value should be true")
	}
	if false != falseBox.ToBool() {
		t.Error("falseBox value should be false")
	}
}

func TestNumberBox(t *testing.T) {
	{
		zeroBox := NewNumber(0)
		if false == math.IsNaN(float64(zeroBox)) {
			t.Error("The box should be NaN")
		}
		if 0 != zeroBox.ToNumber() {
			t.Error("zeroBox value should be 0")
		}
	}
	{
		negativeBox := NewNumber(-1)
		if false == math.IsNaN(float64(negativeBox)) {
			t.Error("The box should be NaN")
		}
		if -1 != negativeBox.ToNumber() {
			t.Error("negativeBox value should be -1")
		}
	}
	{
		positiveBox := NewNumber(1)
		if false == math.IsNaN(float64(positiveBox)) {
			t.Error("The box should be NaN")
		}
		if 1 != positiveBox.ToNumber() {
			t.Error("positiveBox value should be 1")
		}
	}
	{
		maxNegativeBox := NewNumber(math.MinInt32)
		if false == math.IsNaN(float64(maxNegativeBox)) {
			t.Error("The box should be NaN")
		}
		if math.MinInt32 != maxNegativeBox.ToNumber() {
			t.Errorf("maxNegativeBox value should be %d", math.MinInt32)
		}
	}
	{
		maxPositiveBox := NewNumber(math.MaxInt32)
		if false == math.IsNaN(float64(maxPositiveBox)) {
			t.Error("The box should be NaN")
		}
		if math.MaxInt32 != maxPositiveBox.ToNumber() {
			t.Errorf("maxPositiveBox value should be %d", math.MaxInt32)
		}
	}
}

func TestObject(t *testing.T) {
	objectSource := make(Object)
	b1 := NewNumber(1)
	b2 := NewNumber(2)
	objectSource["key1"] = b1
	objectSource["key2"] = b2

	objectBox := NewObject(objectSource)
	if false == math.IsNaN(float64(objectBox)) {
		t.Error("The box should be NaN")
	}
	objectTarget := objectBox.ToObject()
	if len(objectSource) != len(objectTarget) {
		t.Error("objectTarget and objectSource should have the same length")
	}

	// To address elements in the map we have to go through a temporary variable
	// This is why the test set back to original variables the entries from the
	// objectTarget
	// objectTarget["key1"].ToNumber() or objectTarget["key2"].ToNumber() are forbidden
	b1 = objectTarget["key1"]
	b2 = objectTarget["key2"]
	if 1 != b1.ToNumber() {
		t.Error("objectTarget[\"key1\"] should be equal to objectSource[\"key1\"")
	}
	if 2 != b2.ToNumber() {
		t.Error("objectTarget[\"key2\"] should be equal to objectSource[\"key2\"")
	}
}

func TestArray(t *testing.T) {
	arraySource := Array{NewNumber(1), NewNumber(2)}
	arrayBox := NewArray(arraySource)
	if false == math.IsNaN(float64(arrayBox)) {
		t.Error("The box should be NaN")
	}
	arrayTarget := arrayBox.ToArray()
	numberBox1 := arrayTarget[0]
	numberBox2 := arrayTarget[1]
	if 1 != numberBox1.ToNumber() {
		t.Error("arrayTarget[0] should be equal to arraySource[0]")
	}
	if 2 != numberBox2.ToNumber() {
		t.Error("arrayTarget[1] should be equal to arraySource[1]")
	}
}
func TestTag(t *testing.T) {
	box := NewBool(true)
	tag := box.Tag()
	if TagBool != tag {
		t.Errorf("The tag should be Bool, box value %x, NaN value %x, test %x", box, math.NaN(), uint64(0x7FF8000000000001))
	}
	maxPositiveBox := NewNumber(math.MaxInt32)
	tag = maxPositiveBox.Tag()
	if TagNumber != tag {
		t.Errorf("The tag should be Number, box value %x.")
	}
	(&maxPositiveBox).SetTag(tagString)
	tag = maxPositiveBox.Tag()
	if tagString != tag {
		t.Errorf("The tag should be String, box value %x.")
	}
	(&maxPositiveBox).SetTag(TagNumber)
	tag = maxPositiveBox.Tag()
	if math.MaxInt32 != maxPositiveBox.ToNumber() {
		t.Errorf("maxPositiveBox value should be %d", math.MaxInt32)
	}
}

func TestPayload(t *testing.T) {
	maxPositiveBox := NewNumber(math.MinInt32)
	if math.MinInt32 != maxPositiveBox.ToNumber() {
		t.Errorf("maxPositiveBox value should be %d", math.MinInt32)
	}
	maxPositiveBox.SetPayload(uint64(math.MaxInt32))
	if TagNumber != maxPositiveBox.Tag() {
		t.Errorf("The tag should be Number, box value %x.")
	}
	if math.MaxInt32 != maxPositiveBox.ToNumber() {
		t.Errorf("maxPositiveBox value should be %d", math.MaxInt32)
	}
}

func TestIs(t *testing.T) {
	{
		trueBox := NewBool(true)
		if true == trueBox.IsFloat64() {
			t.Error("trueBox should not be a valid float64")
		}
		if false == trueBox.IsBool() {
			t.Error("trueBox should be a valid bool box")
		}
	}
	{
		maxNegativeBox := NewNumber(math.MinInt32)
		if false == maxNegativeBox.IsNumber() {
			t.Error("maxNegativeBox box type should be number")
		}
		randomPositiveBox := NewNumber(12) // Or almost random
		if false == randomPositiveBox.IsNumber() {
			t.Error("randomPositiveBox box type should be number")
		}
	}
	{
		laLaLaLaObject := make(Object)
		laLaLaLaObject["three"] = NewString("Once")
		laLaLaLaObject["little"] = NewString("upon")
		laLaLaLaObject["pigs"] = NewString("a time...")
		laLaLaLaBox := NewObject(laLaLaLaObject)
		if false == laLaLaLaBox.IsObject() {
			t.Error("laLaLaLaObject box type should be object")
		}
	}
	{
		leetArray := Array{NewNumber(1), NewNumber(3), NewNumber(3), NewNumber(7)}
		leetBox := NewArray(leetArray)
		if false == leetBox.IsArray() {
			t.Error("leetArray box type should be array")
		}
	}
}
