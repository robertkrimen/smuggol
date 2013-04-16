// This file was AUTOMATICALLY GENERATED by terst-import (smuggol) from github.com/robertkrimen/terst

/* 
Package terst is a terse (terst = test + terse), easy-to-use testing library for Go.

terst is compatible with (and works via) the standard testing package: http://golang.org/pkg/testing

	import (
		"testing"
		. "github.com/robertkrimen/terst"
	)

	func Test(t *testing.T) {
		Terst(t) // Associate terst methods with t (the current testing.T)

		Is(getApple(), "apple") // Pass
		Is(getOrange(), "orange") // Fail: emits nice-looking diagnostic 

		Compare(1, ">", 0) // Pass
		Compare(1, "==", 1.0) // Pass
	}

	func getApple() string {
		return "apple"
	}

	func getOrange() string {
		return "apple" // Intentional mistake
	}

At the top of your testing function, call Terst(), passing the testing.T you receive as the first argument:

	func TestExample(t *testing.T) {
		Terst(t)
		...
	}

After you initialize with the given *testing.T, you can use the following to test:

	Is
	IsNot
	Equal
	Unequal
	IsTrue
	IsFalse
	Like
	Unlike
	Compare

Each of the methods above can take an additional (optional) argument,
which is a string describing the test. If the test fails, this
description will be included with the test output For example:

	Is(2 + 2, float32(5), "This result is Doubleplusgood")

	--- FAIL: Test (0.00 seconds)
		test.go:17: This result is Doubleplusgood
			Failed test (Is)
			     got: 4 (int)
			expected: 5 (float32)

Future

	- Add Catch() for testing panic()
	- Add Same() for testing via .DeepEqual && == (without panicking?)
	- Add StrictCompare to use {}= scoping
	- Add BigCompare for easier math/big.Int testing?
	- Support the complex type in Compare()
	- Equality test for NaN?
	- Better syntax for At*
	- Need IsType/TypeIs

*/
package terst

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

func (self *Tester) hadResult(result bool, test *test, onFail func()) bool {
	if self.selfTesting {
		expect := true
		if self.failIsPassing {
			expect = false
		}
		if expect != result {
			self.Log(fmt.Sprintf("Expect %v but got %v (%v) (%v) (%v)\n", expect, result, test.kind, test.have, test.want))
			onFail()
			self._fail()
		}
		return result
	}
	if !result {
		onFail()
		self._fail()
	}
	return result
}

// IsTrue is DEPRECATED by:
//
//      Is(..., true)
//
func IsTrue(have bool, description ...interface{}) bool {
	return terstTester().IsTrue(have, description...)
}

// IsTrue is DEPRECATED by:
//
//      Is(..., true)
//
func (self *Tester) IsTrue(have bool, description ...interface{}) bool {
	return self.trueOrFalse(true, have, description...)
}

// IsFalse is DEPRECATED by:
//
//      Is(..., false)
//
func IsFalse(have bool, description ...interface{}) bool {
	return terstTester().IsFalse(have, description...)
}

// IsFalse is DEPRECATED by:
//
//      Is(..., false)
//
func (self *Tester) IsFalse(have bool, description ...interface{}) bool {
	return self.trueOrFalse(false, have, description...)
}

func (self *Tester) trueOrFalse(want bool, have bool, description ...interface{}) bool {
	kind := "IsTrue"
	if want == false {
		kind = "IsFalse"
	}
	test := newTest(kind, have, want, description)
	didPass := have == want
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForIsTrue(test))
	})
}

// Fail will fail immediately, reporting a test failure with the (optional) description
func Fail(description ...interface{}) bool {
	return terstTester().Fail(description...)
}

// Fail will fail immediately, reporting a test failure with the (optional) description
func (self *Tester) Fail(description ...interface{}) bool {
	return self.fail(description...)
}

func (self *Tester) fail(description ...interface{}) bool {
	kind := "Fail"
	test := newTest(kind, false, false, description)
	didPass := false
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForFail(test))
	})
}

// FailNow will fail immediately, triggering testing.FailNow() and optionally reporting a test failure with description
func FailNow(description ...interface{}) bool {
	return terstTester().FailNow(description...)
}

// FailNow will fail immediately, triggering testing.FailNow() and optionally reporting a test failure with description
func (self *Tester) FailNow(description ...interface{}) bool {
	return self.failNow(description...)
}

func (self *Tester) failNow(description ...interface{}) bool {
	if len(description) > 0 {
		kind := "FailNow"
		test := newTest(kind, false, false, description)
		didPass := false
		self.hadResult(didPass, test, func() {
			self.Log(self.failMessageForFail(test))
		})
	}
	self.TestingT.FailNow()
	return false
}

// Equal tests have against want via ==:
//
//		Equal(have, want) // Pass if have == want
//
// No special coercion or type inspection is done.
//
// If the type is incomparable (e.g. type mismatch) this will panic.
func Equal(have, want interface{}, description ...interface{}) bool {
	return terstTester().Equal(have, want, description...)
}

func (self *Tester) Equal(have, want interface{}, description ...interface{}) bool {
	return self.equal(have, want, description...)
}

func (self *Tester) equal(have, want interface{}, description ...interface{}) bool {
	test := newTest("==", have, want, description)
	didPass := have == want
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForEqual(test))
	})
}

// Unequal tests have against want via !=:
//
//		Unequal(have, want) // Pass if have != want
//
// No special coercion or type inspection is done.
//
// If the type is incomparable (e.g. type mismatch) this will panic.
func Unequal(have, want interface{}, description ...interface{}) bool {
	return terstTester().Unequal(have, want, description...)
}

func (self *Tester) Unequal(have, want interface{}, description ...interface{}) bool {
	return self.unequal(have, want, description...)
}

func (self *Tester) unequal(have, want interface{}, description ...interface{}) bool {
	test := newTest("!=", have, want, description)
	didPass := have != want
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForIs(test))
	})
}

// Is tests <have> against <want> in different ways, depending on the
// type of <want>.
//
// If <want> is a string, then it will first convert
// <have> to a string before doing the comparison:
//
//		Is(fmt.Sprintf("%v", have), want) // Pass if have == want
//
// Otherwise, Is is a shortcut for:
//
//		Compare(have, "==", want)
//
// If <want> is a slice, struct, or similar, Is will perform a reflect.DeepEqual() comparison.
func Is(have, want interface{}, description ...interface{}) bool {
	return terstTester().Is(have, want, description...)
}

// TODO "slice, struct, or similar" What is similar?

func (self *Tester) Is(have, want interface{}, description ...interface{}) bool {
	return self.isOrIsNot(true, have, want, description...)
}

// IsNot tests <have> against <want> in different ways, depending on the
// type of <want>.
//
// If <want> is a string, then it will first convert
// <have> to a string before doing the comparison:
//
//		IsNot(fmt.Sprintf("%v", have), want) // Pass if have != want
//
// Otherwise, Is is a shortcut for:
//
//		Compare(have, "!=", want)
//
// If <want> is a slice, struct, or similar, Is will perform a reflect.DeepEqual() comparison.
func IsNot(have, want interface{}, description ...interface{}) bool {
	return terstTester().IsNot(have, want, description...)
}

// TODO "slice, struct, or similar" What is similar?

func (self *Tester) IsNot(have, want interface{}, description ...interface{}) bool {
	return self.isOrIsNot(false, have, want, description...)
}

func (self *Tester) isOrIsNot(wantIs bool, have, want interface{}, description ...interface{}) bool {
	test := newTest("Is", have, want, description)
	if !wantIs {
		test.kind = "IsNot"
	}
	didPass := false
	switch want.(type) {
	case string:
		didPass = stringValue(have) == want
	default:
		didPass, _ = compare(have, "{}* ==", want)
	}
	if !wantIs {
		didPass = !didPass
	}
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForIs(test))
	})
}

// Like tests <have> against <want> in different ways, depending on the
// type of <want>.
//
// If <want> is a string, then it will first convert
// <have> to a string before doing a regular expression comparison:
//
//		Like(fmt.Sprintf("%v", have), want) // Pass if regexp.Match(want, have)
//
// Otherwise, Like is a shortcut for:
//
//		Compare(have, "{}~ ==", want)
//
// If <want> is a slice, struct, or similar, Like will perform a reflect.DeepEqual() comparison.
func Like(have, want interface{}, description ...interface{}) bool {
	return terstTester().Like(have, want, description...)
}

func (self *Tester) Like(have, want interface{}, description ...interface{}) bool {
	return self.likeOrUnlike(true, have, want, description...)
}

// Unlike tests <have> against <want> in different ways, depending on the
// type of <want>.
//
// If <want> is a string, then it will first convert
// <have> to a string before doing a regular expression comparison:
//
//		Unlike(fmt.Sprintf("%v", have), want) // Pass if !regexp.Match(want, have)
//
// Otherwise, Unlike is a shortcut for:
//
//		Compare(have, "{}~ !=", want)
//
// If <want> is a slice, struct, or similar, Unlike will perform a reflect.DeepEqual() comparison.
func Unlike(have, want interface{}, description ...interface{}) bool {
	return terstTester().Unlike(have, want, description...)
}

func (self *Tester) Unlike(have, want interface{}, description ...interface{}) bool {
	return self.likeOrUnlike(false, have, want, description...)
}

func (self *Tester) likeOrUnlike(wantLike bool, have, want interface{}, description ...interface{}) bool {
	test := newTest("Like", have, want, description)
	if !wantLike {
		test.kind = "Unlike"
	}
	didPass := false
	switch want0 := want.(type) {
	case string:
		haveString := stringValue(have)
		didPass, error := regexp.Match(want0, []byte(haveString))
		if !wantLike {
			didPass = !didPass
		}
		if error != nil {
			panic("regexp.Match(" + want0 + ", ...): " + error.Error())
		}
		want = fmt.Sprintf("(?:%v)", want) // Make it look like a regular expression
		return self.hadResult(didPass, test, func() {
			self.Log(self.failMessageForMatch(test, stringValue(have), stringValue(want), wantLike))
		})
	}
	didPass, operator := compare(have, "{}~ ==", want)
	if !wantLike {
		didPass = !didPass
	}
	test.operator = operator
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForLike(test, stringValue(have), stringValue(want), wantLike))
	})
}

// Compare will compare <have> to <want> with the given operator. The operator can be one of the following:
//
//		  ==
//		  !=
//		  <
//		  <=
//		  >
//		  >=
//
// Compare is not strict when comparing numeric types,
// and will make a best effort to promote <have> and <want> to the
// same type.
//
// Compare will promote int and uint to big.Int for testing
// against each other.
//
// Compare will promote int, uint, and float to float64 for
// float testing.
//
// For example:
//	
//		Compare(float32(1.0), "<", int8(2)) // A valid test
//
//		result := float32(1.0) < int8(2) // Will not compile because of the type mismatch
//
func Compare(have interface{}, operator string, want interface{}, description ...interface{}) bool {
	return terstTester().Compare(have, operator, want, description...)
}

func (self *Tester) Compare(have interface{}, operator string, want interface{}, description ...interface{}) bool {
	return self.compare(have, operator, want, description...)
}

func (self *Tester) compare(left interface{}, operatorString string, right interface{}, description ...interface{}) bool {
	operatorString = strings.TrimSpace(operatorString)
	test := newTest("Compare "+operatorString, left, right, description)
	didPass, operator := compare(left, operatorString, right)
	test.operator = operator
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForCompare(test))
	})
}

type (
	compareScope int
)

const (
	compareScopeEqual compareScope = iota
	compareScopeTilde
	compareScopeAsterisk
)

type compareOperator struct {
	scope      compareScope
	comparison string
}

var newCompareOperatorRE *regexp.Regexp = regexp.MustCompile(`^\s*(?:((?:{}|#)[*~=])\s+)?(==|!=|<|<=|>|>=)\s*$`)

func newCompareOperator(operatorString string) compareOperator {

	if operatorString == "" {
		return compareOperator{compareScopeEqual, ""}
	}

	result := newCompareOperatorRE.FindStringSubmatch(operatorString)
	if result == nil {
		panic(fmt.Errorf("Unable to parse %v into a compareOperator", operatorString))
	}

	scope := compareScopeAsterisk
	switch result[1] {
	case "#*", "{}*":
		scope = compareScopeAsterisk
	case "#~", "{}~":
		scope = compareScopeTilde
	case "#=", "{}=":
		scope = compareScopeEqual
	}

	comparison := result[2]

	return compareOperator{scope, comparison}
}

func compare(left interface{}, operatorString string, right interface{}) (bool, compareOperator) {
	pass := true
	operator := newCompareOperator(operatorString)
	comparator := newComparator(left, operator, right)
	// FIXME Confusing
	switch operator.comparison {
	case "==":
		pass = comparator.IsEqual()
	case "!=":
		pass = !comparator.IsEqual()
	default:
		if comparator.HasOrder() {
			switch operator.comparison {
			case "<":
				pass = comparator.Compare() == -1
			case "<=":
				pass = comparator.Compare() <= 0
			case ">":
				pass = comparator.Compare() == 1
			case ">=":
				pass = comparator.Compare() >= 0
			default:
				panic(fmt.Errorf("Compare operator (%v) is invalid", operator.comparison))
			}
		} else {
			pass = false
		}
	}
	return pass, operator
}

// Compare / Comparator

type compareKind int

const (
	kindInterface compareKind = iota
	kindInteger
	kindUnsignedInteger
	kindFloat
	kindString
	kindBoolean
)

func comparatorValue(value interface{}) (reflect.Value, compareKind) {
	reflectValue := reflect.ValueOf(value)
	kind := kindInterface
	switch value.(type) {
	case int, int8, int16, int32, int64:
		kind = kindInteger
	case uint, uint8, uint16, uint32, uint64:
		kind = kindUnsignedInteger
	case float32, float64:
		kind = kindFloat
	case string:
		kind = kindString
	case bool:
		kind = kindBoolean
	}
	return reflectValue, kind
}

func toFloat(value reflect.Value) float64 {
	switch result := value.Interface().(type) {
	case int, int8, int16, int32, int64:
		return float64(value.Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(value.Uint())
	case float32, float64:
		return float64(value.Float())
	default:
		panic(fmt.Errorf("toFloat( %v )", result))
	}
	panic(0)
}

func toInteger(value reflect.Value) *big.Int {
	switch result := value.Interface().(type) {
	case int, int8, int16, int32, int64:
		return big.NewInt(value.Int())
	case uint, uint8, uint16, uint32, uint64:
		yield := big.NewInt(0)
		yield.SetString(fmt.Sprintf("%v", value.Uint()), 10)
		return yield
	default:
		panic(fmt.Errorf("toInteger( %v )", result))
	}
	panic(0)
}

func toString(value reflect.Value) string {
	switch result := value.Interface().(type) {
	case string:
		return result
	default:
		panic(fmt.Errorf("toString( %v )", result))
	}
	panic(0)
}

func toBoolean(value reflect.Value) bool {
	switch result := value.Interface().(type) {
	case bool:
		return result
	default:
		panic(fmt.Errorf("toBoolean( %v )", result))
	}
	panic(0)
}

type aComparator interface {
	Compare() int
	HasOrder() bool
	IsEqual() bool
	CompareScope() compareScope
}

type baseComparator struct {
	hasOrder bool
	operator compareOperator
}

func (self *baseComparator) Compare() int {
	panic(fmt.Errorf("Invalid .Compare()"))
}
func (self *baseComparator) HasOrder() bool {
	return self.hasOrder
}
func (self *baseComparator) CompareScope() compareScope {
	return self.operator.scope
}
func comparatorWithOrder(operator compareOperator) *baseComparator {
	return &baseComparator{true, operator}
}
func comparatorWithoutOrder(operator compareOperator) *baseComparator {
	return &baseComparator{false, operator}
}

type interfaceComparator struct {
	*baseComparator
	left  interface{}
	right interface{}
}

func (self *interfaceComparator) IsEqual() bool {
	if self.CompareScope() != compareScopeEqual {
		return reflect.DeepEqual(self.left, self.right)
	}
	return self.left == self.right
}

type floatComparator struct {
	*baseComparator
	left  float64
	right float64
}

func (self *floatComparator) Compare() int {
	if self.left == self.right {
		return 0
	} else if self.left < self.right {
		return -1
	}
	return 1
}
func (self *floatComparator) IsEqual() bool {
	return self.left == self.right
}

type integerComparator struct {
	*baseComparator
	left  *big.Int
	right *big.Int
}

func (self *integerComparator) Compare() int {
	return self.left.Cmp(self.right)
}
func (self *integerComparator) IsEqual() bool {
	return 0 == self.left.Cmp(self.right)
}

type stringComparator struct {
	*baseComparator
	left  string
	right string
}

func (self *stringComparator) Compare() int {
	if self.left == self.right {
		return 0
	} else if self.left < self.right {
		return -1
	}
	return 1
}
func (self *stringComparator) IsEqual() bool {
	return self.left == self.right
}

type booleanComparator struct {
	*baseComparator
	left  bool
	right bool
}

func (self *booleanComparator) IsEqual() bool {
	return self.left == self.right
}

func newComparator(left interface{}, operator compareOperator, right interface{}) aComparator {
	leftValue, _ := comparatorValue(left)
	rightValue, rightKind := comparatorValue(right)

	// The simplest comparator is comparing interface{} =? interface{}
	targetKind := kindInterface
	// Are left and right of the same kind?
	// (reflect.Value.Kind() is different from compareKind)
	scopeEqual := leftValue.Kind() == rightValue.Kind()
	scopeTilde := false
	scopeAsterisk := false
	if scopeEqual {
		targetKind = rightKind // Since left and right are the same, the targetKind is Integer/Float/String/Boolean
	} else {
		// Examine the prefix of reflect.Value.Kind().String() to see if there is a similarity of 
		// the left value to right value
		lk := leftValue.Kind().String()
		hasPrefix := func(prefix string) bool {
			return strings.HasPrefix(lk, prefix)
		}

		switch right.(type) {
		case float32, float64:
			// Right is float*
			if hasPrefix("float") {
				// Left is also float*
				targetKind = kindFloat
				scopeTilde = true
			} else if hasPrefix("int") || hasPrefix("uint") {
				// Left is a kind of numeric (int* or uint*)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		case uint, uint8, uint16, uint32, uint64:
			// Right is uint*
			if hasPrefix("uint") {
				// Left is also uint*
				targetKind = kindInteger
				scopeTilde = true
			} else if hasPrefix("int") {
				// Left is an int* (a numeric)
				targetKind = kindInteger
				scopeAsterisk = true
			} else if hasPrefix("float") {
				// Left is an float* (a numeric)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		case int, int8, int16, int32, int64:
			// Right is int*
			if hasPrefix("int") {
				// Left is also int*
				targetKind = kindInteger
				scopeTilde = true
			} else if hasPrefix("uint") {
				// Left is a uint* (a numeric)
				targetKind = kindInteger
				scopeAsterisk = true
			} else if hasPrefix("float") {
				// Left is an float* (a numeric)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		default:
			// Right is a non-numeric
			// Can only really compare string to string or boolean to boolean, so
			// we will either have a string/boolean/interfaceComparator
		}
	}

	/*fmt.Println("%v %v %v %v %s %s", operator.scope, same, sibling, family, leftValue, rightValue)*/
	{
		mismatch := false
		switch operator.scope {
		case compareScopeEqual:
			mismatch = !scopeEqual
		case compareScopeTilde:
			mismatch = !scopeEqual && !scopeTilde
		case compareScopeAsterisk:
			mismatch = !scopeEqual && !scopeTilde && !scopeAsterisk
		}
		if mismatch {
			targetKind = kindInterface
		}
	}

	switch targetKind {
	case kindFloat:
		return &floatComparator{
			comparatorWithOrder(operator),
			toFloat(leftValue),
			toFloat(rightValue),
		}
	case kindInteger:
		return &integerComparator{
			comparatorWithOrder(operator),
			toInteger(leftValue),
			toInteger(rightValue),
		}
	case kindString:
		return &stringComparator{
			comparatorWithOrder(operator),
			toString(leftValue),
			toString(rightValue),
		}
	case kindBoolean:
		return &booleanComparator{
			comparatorWithoutOrder(operator),
			toBoolean(leftValue),
			toBoolean(rightValue),
		}
	}

	// As a last resort, we can always compare left (interface{}) to right (interface{})
	return &interfaceComparator{
		comparatorWithoutOrder(operator),
		left,
		right,
	}
}

// failMessage*

func (self *Tester) failMessageForIsTrue(test *test) string {
	test.findFileLineFunction(self)
	return formatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, stringValue(test.have), stringValue(test.want))
}

func (self *Tester) failMessageForFail(test *test) string {
	test.findFileLineFunction(self)
	return formatMessage(`
        %s:%d: %s 
           Failed test (%s)
    `, test.file, test.line, test.Description(), test.kind)
}

func typeKindString(value interface{}) string {
	reflectValue := reflect.ValueOf(value)
	kind := reflectValue.Kind().String()
	result := fmt.Sprintf("%T", value)
	if kind == result {
		if kind == "string" {
			return ""
		}
		return fmt.Sprintf(" (%T)", value)
	}
	return fmt.Sprintf(" (%T=%s)", value, kind)
}

func (self *Tester) failMessageForCompare(test *test) string {
	test.findFileLineFunction(self)
	return formatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  %v%s
                       %s
                  %v%s
    `, test.file, test.line, test.Description(), test.kind, test.have, typeKindString(test.have), test.operator.comparison, test.want, typeKindString(test.want))
}

func (self *Tester) failMessageForEqual(test *test) string {
	return self.failMessageForIs(test)
}

func (self *Tester) failMessageForIs(test *test) string {
	test.findFileLineFunction(self)
	return formatMessage(`
        %s:%d: %v
           Failed test (%s)
                  got: %v%s
             expected: %v%s
    `, test.file, test.line, test.Description(), test.kind, test.have, typeKindString(test.have), test.want, typeKindString(test.want))
}

func (self *Tester) failMessageForMatch(test *test, have, want string, wantMatch bool) string {
	test.findFileLineFunction(self)
	expect := "  like"
	if !wantMatch {
		expect = "unlike"
	}
	return formatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %v%s
               %s: %s
    `, test.file, test.line, test.Description(), test.kind, have, typeKindString(have), expect, want)
}

func (self *Tester) failMessageForLike(test *test, have, want string, wantLike bool) string {
	test.findFileLineFunction(self)
	if !wantLike {
		want = "Anything else"
	}
	return formatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %v%s
             expected: %v%s
    `, test.file, test.line, test.Description(), test.kind, have, typeKindString(have), want, typeKindString(want))
}

// ...

type Tester struct {
	TestingT *testing.T

	sanityChecking bool
	selfTesting    bool
	failIsPassing  bool

	testEntry  uintptr
	focusEntry uintptr
}

var _terstTester *Tester = nil

func findTestEntry() uintptr {
	height := 2
	for {
		functionPC, _, _, ok := runtime.Caller(height)
		function := runtime.FuncForPC(functionPC)
		functionName := function.Name()
		if !ok {
			return 0
		}
		if index := strings.LastIndex(functionName, ".Test"); index >= 0 {
			// Assume we have an instance of TestXyzzy in a _test file
			return function.Entry()
		}
		height += 1
	}
	return 0
}

// Focus will focus the entry point of the test to the current method.
//
// This is important for test failures in getting feedback on which line was at fault.
//
// Consider the following scenario:
//
//		func testingMethod( ... ) {
//			Is( ..., ... )
//		}
//
//		func TestExample(t *testing.T) {
//			Terst(t)
//
//			testingMethod( ... )
//			testingMethod( ... ) // If something in testingMethod fails, this line number will come up
//			testingMethod( ... )
//		}
//	
// By default, when a test fails, terst will report the outermost line that led to the failure.
// Usually this is what you want, but if you need to drill down, you can by inserting a special
// call at the top of your testing method:
//
//		func testingMethod( ... ) {
//			Terst().Focus() // Grab the global Tester and tell it to focus on this method
//			Is( ..., ... ) // Now if this test fails, this line number will come up
//		}
//
func (self *Tester) Focus() {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		function := runtime.FuncForPC(pc)
		self.focusEntry = function.Entry()
	}
}

//		Terst(*testing.T)
//
// Create a new terst Tester and return it.  Associate calls to Is, Compare, Like, etc. with the newly created terst.
//
//		Terst()
//
// Return the current Tester (if any).
//
//		Terst(nil)
//
// Clear out the current Tester (if any).
func Terst(terst ...interface{}) *Tester {
	if len(terst) == 0 {
		return terstTester()
	} else {
		if terst[0] == nil {
			_terstTester = nil
			return nil
		}
		_terstTester = newTester(terst[0].(*testing.T))
		_terstTester.enableSanityChecking()
		_terstTester.testEntry = findTestEntry()
		_terstTester.focusEntry = _terstTester.testEntry
	}
	return _terstTester
}

func terstTester() *Tester {
	if _terstTester == nil {
		panic("_terstTester == nil")
	}
	return _terstTester.checkSanity()
}

func newTester(t *testing.T) *Tester {
	return &Tester{
		TestingT: t,
	}
}

func formatMessage(message string, argumentList ...interface{}) string {
	message = fmt.Sprintf(message, argumentList...)
	message = strings.TrimLeft(message, "\n")
	message = strings.TrimRight(message, " \n")
	return message + "\n\n"
}

// Log is a utility method that will append the given output to the normal output stream.
func (self *Tester) Log(output string) {
	outputValue := reflect.ValueOf(self.TestingT).Elem().FieldByName("output")
	output_ := outputValue.Bytes()
	output_ = append(output_, output...)
	*(*[]byte)(unsafe.Pointer(outputValue.UnsafeAddr())) = output_
}

func (self *Tester) _fail() {
	self.TestingT.Fail()
}

func (self *Tester) enableSanityChecking() *Tester {
	self.sanityChecking = true
	return self
}

func (self *Tester) disableSanityChecking() *Tester {
	self.sanityChecking = false
	return self
}

func (self *Tester) enableSelfTesting() *Tester {
	self.selfTesting = true
	return self
}

func (self *Tester) disableSelfTesting() *Tester {
	self.selfTesting = false
	return self
}

func (self *Tester) failIsPass() *Tester {
	self.failIsPassing = true
	return self
}

func (self *Tester) passIsPass() *Tester {
	self.failIsPassing = false
	return self
}

func (self *Tester) checkSanity() *Tester {
	if self.sanityChecking && self.testEntry != 0 {
		foundEntryPoint := findTestEntry()
		if self.testEntry != foundEntryPoint {
			panic(fmt.Errorf("TestEntry(%v) does not match foundEntry(%v): Did you call Terst when entering a new Test* function?", self.testEntry, foundEntryPoint))
		}
	}
	return self
}

func (self *Tester) findDepth() int {
	height := 1 // Skip us
	for {
		pc, _, _, ok := runtime.Caller(height)
		function := runtime.FuncForPC(pc)
		if !ok {
			// Got too close to the sun
			if false {
				for ; height > 0; height-- {
					pc, _, _, ok := runtime.Caller(height)
					fmt.Printf("[%d %v %v]", height, pc, ok)
					if ok {
						function := runtime.FuncForPC(pc)
						fmt.Printf(" => [%s]", function.Name())
					}
					fmt.Printf("\n")
				}
			}
			return 1
		}
		functionEntry := function.Entry()
		if functionEntry == self.focusEntry || functionEntry == self.testEntry {
			return height - 1 // Not the surrounding test function, but within it
		}
		height += 1
	}
	return 1
}

// test

type test struct {
	kind        string
	have        interface{}
	want        interface{}
	description []interface{}
	operator    compareOperator

	file       string
	line       int
	functionPC uintptr
	function   string
}

func newTest(kind string, have, want interface{}, description []interface{}) *test {
	operator := newCompareOperator("")
	return &test{
		kind:        kind,
		have:        have,
		want:        want,
		description: description,
		operator:    operator,
	}
}

func (self *test) findFileLineFunction(tester *Tester) {
	self.file, self.line, self.functionPC, self.function, _ = atFileLineFunction(tester.findDepth())
}

func (self *test) Description() string {
	description := ""
	if len(self.description) > 0 {
		description = fmt.Sprintf("%v", self.description[0])
	}
	return description
}

func findPathForFile(file string) string {
	terstBase := os.ExpandEnv("$TERST_BASE")
	if len(terstBase) > 0 && strings.HasPrefix(file, terstBase) {
		file = file[len(terstBase):]
		if file[0] == '/' || file[0] == '\\' {
			file = file[1:]
		}
		return file
	}

	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[index+1:]
	} else if index = strings.LastIndex(file, "\\"); index >= 0 {
		file = file[index+1:]
	}

	return file
}

func atFileLineFunction(callDepth int) (string, int, uintptr, string, bool) {
	pc, file, line, ok := runtime.Caller(callDepth + 1)
	function := runtime.FuncForPC(pc).Name()
	if ok {
		file = findPathForFile(file)
		if index := strings.LastIndex(function, ".Test"); index >= 0 {
			function = function[index+1:]
		}
	} else {
		pc = 0
		file = "?"
		line = 1
	}
	return file, line, pc, function, ok
}

// Conversion

func integerValue(value interface{}) int64 {
	return reflect.ValueOf(value).Int()
}

func unsignedIntegerValue(value interface{}) uint64 {
	return reflect.ValueOf(value).Uint()
}

func floatValue(value interface{}) float64 {
	return reflect.ValueOf(value).Float()
}

func stringValue(value interface{}) string {
	return fmt.Sprintf("%v", value)
}