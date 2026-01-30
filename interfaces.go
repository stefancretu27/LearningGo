package go_didactical_apps

/*
*  Interface:
*	- In Go, an interface is a type that specifies a set of methods. 
*	- all types that implement an interface, must implement all methods encapsulated by it. These types are referred to as the set of types defined by the interface
*	or the set of types which sastisfy the interface
*	- a variable of interface type can store a value of any type that is in the type set of the interface. Such a type is said to implement the interface. 
*	- the value of an uninitialized variable of interface type is nil.
*	- an interface type is specified by a list of interface elements. An interface element is either a method or a type element, where a type element is a 
*	union of one or more type terms.
*	- all types implement the empty interface which stands for the set of all (non-interface) types: interface{}
*
* Types That Can Implement Interfaces
*	- Structs (most common)
*	- Pointers to structs
*	- Custom types based on basic types: int, string, slices, maps etc
*	- functions: if the interface has exactly one method, and the function's signature matches that method, then a function type can implement that interface.
*	- basic types (string, int) cannot implement interfaces => "cannot define new methods on non-local type". Thus, they require custom type definition
*		type MyInt int
*		func (i MyInt) InterfaceMethod() {
*		}
*	- a given type can implement multiple interfaces => it implements all methods of those interfaces
*
* Basic interfaces:
*	- in its most basic form an interface specifies a (possibly empty) list of methods. 
*	- can create variables of such types. They are nil. Used to hold objects that implement interface
*		eg: var ifaceVar MyInterface
*	- the type set defined by such an interface is the set of types which implement all of those methods, and the corresponding methods set consists exactly of 
*	the methods specified by the interface.
*	- a struct or other type can have more methods than those required by an interface, and it still implements the interface as long as it satisfies all the 
*	methods declared in the interface.
*	- the names (implicitly the signatures, as overloading is not supported) of the methods in an interface must be unique => no methods with same name
*	Also, the names cannot be blank (_)
*
* General interfaces:
*	- in their most general form, an interface element may also be an arbitrary type term T, or a term of the form ~T specifying the underlying type T,
*	 or a union of terms t1|t2|…|tn (eg: int | int32 | int64 etc)
*	- mostly used as type constraints from generics
*	- cannot create variable whose types are such interfaces
*	- they can contain method's declarations. Thus, still are to be used as type constraints for the specified types and which implement the listed methods
*	- in a term of the form ~T, the underlying type of T must be itself, and T cannot be an interface or type parameter
*Eg:	type MyInt int
*
*		interface {
*			~[]byte  // the underlying type of []byte is itself
*			~MyInt   // illegal: the underlying type of MyInt is not MyInt
*			~error   // illegal: error is an interface
*		}
*
*		The Float interface represents all floating-point types (including any named types whose underlying types are either float32 or float64).
*		type Float interface {
*			~float32 | ~float64
*		}
*/

/*
* Structs:
*	- Go does not allow method overloading — meaning you cannot define two methods with the same name on the same type, 
*	even if one uses a value receiver and the other a pointer receiver.
*	- if a method is defined on pointer receiver, regardless it is called on pointer or value, the receiever will be treated as pointer and
* changes made to data memebrs persist
*	- if a method is defined on value receiver, changes made to data members will not persist, regardless it is called on pointer or value receiver
*	- Embeding: -> if one struct embeds another, it can access its fields and methods, regardless they are public or private, but it still remains a distinct type. 
* 				-> an embedded struct is included in another struct without explicitly naming it as a field. This promotes the fields and methods of the embedded struct 
* to the outer struct, allowing direct access without referencing the embedded struct's name. 
*				-> the outer struct can have own implementation of (some/all) the embedded struct's methods (sort of overriding). The methods of the embedded struct 
* can be accessed using its name: outterInst.EmbeddedStruct.Method(). If the embeded struct implements the inetrface on ptr, the outer struct must do so, for interface
* compatibility
*	type EmbedStruct struct{
*		Base
*		<other field> <field type>}
*	- Nesting: one struct contains an instance of another struct, it has access to its methods and fields via that instance. Hence, if the nested struct implements an interface,
*	the nesting struct doesn't, so it must explicitly implement it
*	type NestInst struct{
*		inst Base
*		<other field> <field type>}
*	- a struct cannot directly qualify as another struct type because each struct type is distinct, even if they have identical fields. 
* Go enforces strict type safety, so two structs with the same fields are still considered different types and cannot be used interchangeably without explicit conversion.
*/

// declare a type constraint interface
type Number interface {
	int | float64
}

//declare another type constraint, with constraint to implement the listed method
type MixedInterface interface {
	~int | ~int16 | ~int32 | ~int64 | string
	About(name string)
}

// attempt to implement interface using type constraint interface => cannot use type Number outside a type constraint: interface contains type constraints
// func (i Number) PrintName(name string) {
// 	fmt.Println("[PrintName] NUuber", name)
// }

// attempt to implement interface using basic type => cannot define new methods on non-local type string
// func (s string) PrintName(name string) {
// 	fmt.Println("[PrintName] String", name)
// }

// as an alternative, basic types can implement interfaces (be extended with behaviors) if custom types are defined on them

// -----------------------------General interfaces used as type constraints---------------
// General interface, encapsulating a union of types, used as type constraint in generics
type Integer interface {
	~int | ~int16 | ~int32 | ~int64 | ~uint | ~uint16 | ~uint32 | ~uint64
}

// General interface used for type constraints in generics (otehr interfaces, structs, custom types)
type Numeric interface {
	Integer | ~float32 | ~float64
}

// -----------------------------Function type implementing an interface------------------------
// Define an interface holding exactly one method so it can be implemented by a function, whose signature matches it
// This is a common pattern in Go’s standard library as it allows to pass plain functions as interface implementations => avoid boilerplate structs
// Boilerplate refers to sections of code that are repeated in multiple places with little or no modification.
type FuncIface[N Numeric] interface {
	Add(lhs, rhs N) N
}

type MySum[N Numeric] func(a, b N) N

// Implement the above interface, using custom function type as implementation entity, with the receiver being called within the implementation, thus wrapping it
func (ms MySum[N]) Add(lhs, rhs N) N {
	return ms(lhs, rhs)
}

func add(a, b int) int {
	return a + b
}

func FuncIfaceImpl() {
	fmt.Println("-----------------------------Function type implementing an interface------------------------")
	// Convert existing function to custom function type, necessary to call the interface's method
	ms := MySum[int](add)
	fmt.Println("[Add][int]:", ms.Add(43, 44))

	var ms2 MySum[float64] = MySum[float64](func(a, b float64) float64 {
		return a + b
	})
	fmt.Println("[Add][float64]:", ms2.Add(3.14159, 2.7182))
}

// Define a generic interface
type SimpleIface[T Numeric] interface {
	PrintName()
	GetValue(index interface{}) T
	SetValue(value T, index interface{})
	GetValues() T
	SetValues(value T)
}

// -----------------------------Pointer to custom data type implementing an interface------------------------
// Basic data types cannot implement interfaces, but aliases to them can (aka custom data types)
type MyUint16 uint16

// Use ptr to receiver such that the setter has effect. The value itself is gotten or set
func (mu MyUint16) PrintName() {
	fmt.Println("[PrintName][MyUint16]: Custom type (alias) to uint16")
}

func (mu *MyUint16) GetValue(index interface{}) MyUint16 {
	if mu != nil {
		return *mu
	}

	var zeroValue MyUint16
	return zeroValue
}

func (mu *MyUint16) SetValue(value MyUint16, index interface{}) {
	if mu != nil {
		*mu = value
	}
}

func (mu *MyUint16) GetValues() MyUint16 {
	if mu != nil {
		return *mu
	}

	var zeroValue MyUint16
	return zeroValue
}

func (mu *MyUint16) SetValues(value MyUint16) {
	if mu != nil {
		*mu = value
	}
}

type MyInt32 int32

func (mi MyInt32) PrintName() {
	fmt.Println("[PrintName][MyInt32]: Custom type (alias) to int32")
}

func (mi MyInt32) GetValue(index interface{}) MyInt32 {
	return mi
}

func (mi MyInt32) SetValue(value MyInt32, index interface{}) {
	mi = value
}

func (mi MyInt32) GetValues() MyInt32 {
	return mi
}

func (mi MyInt32) SetValues(value MyInt32) {
	mi = value
}

func CustomDataTypeIfaceImpl() {
	fmt.Println("-----------------------------Pointer to custom data type implementing an interface------------------------")
	var sIface SimpleIface[MyUint16]
	var mu16Ptr *MyUint16 = new(MyUint16)
	*mu16Ptr = 3216
	sIface = mu16Ptr // &MyUint16(3216) => error: cannot take address of constant 3216 of uint16 type MyUint16
	sIface.PrintName()
	fmt.Println("[GetValue][MyUint16]:", sIface.GetValue(nil))
	sIface.SetValue(64, nil) //directly use uint16 input
	fmt.Println("[GetValue][MyUint16]:", sIface.GetValue(nil))
	var mu MyUint16 = 27
	mu.SetValue(13, nil)     //persistent change due to ptr receiver, although caller is instance
	sIface.SetValue(mu, nil) //use MyUint16 input
	fmt.Println("[GetValue][][MyUint16]:", sIface.GetValue(nil), mu)

	fmt.Println("-----------------------------Custom data type implementing an interface------------------------")
	var mi1 MyInt32 = 4
	mi1.PrintName()
	mi1.SetValue(55, nil)
	fmt.Println("[GetValue][][MyInt32]:", mi1.GetValue(nil), mi1)
}

// -----------------------------Custom type, with map underlying type, implementing an interface------------------------
type MyMap[N Numeric] map[string]N

func (mm MyMap[N]) PrintName() {
	fmt.Println("[PrintName][MyMap[N]]: generic custom data type whose underlying type is a map")
}

func (mm MyMap[N]) GetValue(index interface{}) N {
	if index != nil {
		//any(value).(T) — Type Assertion, not a type switch => can be used outside switch type. It extracts the concrete value from an interface if it matches type T.
		//any is just a type alias for interface{} in Go 1.18+. Writing any(value) is simply a type conversion of value to interface{}.
		//value.(T) is also a type assertion and works if value is an interface. any(value).(T) just converts value to interface{}
		if localIdx, assertionOk := any(index).(string); assertionOk {
			return mm[localIdx]
		}
	}
	//Cannot return nil or empty object for a type N or for any/interface{}, as not all types can be nil (nil can be pointers, slices, map, chan, func, interface)
	//Thus, it is returned the zero value for that type, which is held by an uninitialized instance of it.
	//It works for all types — for pointers it will be nil, for numbers it will be 0, for strings it will be ""
	var zeroValue N
	return zeroValue
}

func (mm MyMap[N]) SetValue(value N, index interface{}) {
	if index != nil {
		if localIdx, assertionSuccess := any(index).(string); assertionSuccess {
			mm[localIdx] = value
		}
	}
}

func (mm MyMap[N]) GetValues() MyMap[N] {
	return mm
}

func (mm MyMap[N]) SetValues(inputMap MyMap[N]) {
	mm = inputMap
}

func CustomMapIfaceImpl() {
	fmt.Println("-----------------------------Ccustom type, with slice underlying type, implementing a map------------------------")
	var mm MyMap[float64] = MyMap[float64]{
		"pi":             3.14159,
		"euler's number": 2.7182,
	}
	mm.PrintName()
	mm.SetValue(-3.309, "inserted value")
	fmt.Println("[GetValues][MyMap[N]]:", mm.GetValues())

	doubleMap := map[string]float64{
		"double pi":             2 * 3.14159,
		"double euler's number": 2 * 2.7182,
	}
	mm = MyMap[float64](doubleMap)
	fmt.Println("[][MyMap[N]]:", mm)

	minusMap := map[string]float64{
		"-pi":             -3.14159,
		"-euler's number": -2.7182,
	}
	// Not persistent due to instance receiver, not pointer receiver
	mm.SetValues(minusMap)
	fmt.Println("[GetValues][MyMap[N]]:", mm.GetValues())
	mm.SetValue(-273.01, "zero kelvin")
	fmt.Println("[GetValue][MyMap[N]]:", mm["-pi"], mm["double pi"], mm.GetValue("zero kelvin"))
	fmt.Println("[][MyMap[N]]:", mm)
	//The underlying type can be assigned to the custom type
	mm = map[string]float64{
		"x": 5,
		"y": 6}
	fmt.Println("[GetValues][][MyMap[N]]:", mm.GetValues(), mm)
}

// -----------------------------Pointer to custom type, with slice underlying type, implementing an interface------------------------
// Define a generic custom data type whose underlying impl is a slice
type MySlice[T Integer] []T

func (ms *MySlice[T]) PrintName() {
	fmt.Println("[PrintName][MySlice[T]]: generic custom data type whose underlying type is a slice")
}

func (ms *MySlice[T]) GetValue(index interface{}) T {
	if ms != nil && index != nil {
		if localIdx, assertionOk := any(index).(int); assertionOk {
			return (*ms)[localIdx]
		}
	}

	var zeroValue T
	return zeroValue
}

func (ms *MySlice[T]) SetValue(in T, index interface{}) {
	if ms != nil && index != nil {
		if localIdx, assertionOk := any(index).(int); assertionOk {
			(*ms)[localIdx] = in
		}
	}
}

func (ms *MySlice[T]) GetValues() MySlice[T] {
	if ms != nil {
		return *ms
	}

	var zeroValue MySlice[T]
	return zeroValue
}

func (ms *MySlice[T]) SetValues(inSlice MySlice[T]) {
	if ms != nil {
		*ms = inSlice
	}
}

func CustomSliceIFaceImpl() {
	fmt.Println("-----------------------------Pointer to custom type, with slice underlying type, implementing an interface------------------------")
	var ms *MySlice[int16]                   //= new(MySlice[int16])
	ms = &MySlice[int16]{11, 22, 33, 44, 55} //or *ms = MySlice[int16]{11, 22, 33, 44, 55} if memory is allocated, so it can eb dereferenced
	ms.PrintName()
	fmt.Println("[GetValue][MySlice[T]]", ms.GetValue(0))
	fmt.Println("[GetValues][MySlice[T]]", ms.GetValues())
	ms.SetValues(MySlice[int16]{-12, -13, -14, -15, -16, -17, -18})
	ms.SetValue(273, 3)
	fmt.Println("[GetValues][MySlice[T]]", ms.GetValues())
	*ms = []int16{0, 1, 2}
	fmt.Println("[GetValues][MySlice[T]]", ms.GetValues())
	ms.SetValues([]int16{-1, -3, -4, -5, -8})
	ms.SetValue(int16(2222), 2)
	fmt.Println("[GetValues][MySlice[T]]", ms.GetValues())
	someSlice := []int16{888, 999, 000, 555}
	*ms = someSlice
	fmt.Println("[GetValues][MySlice[T]]", ms.GetValues())
}

// -----------------------------Pointer to struct implementing an interface------------------------
// Basic interface with generic parameter
type Iface[T Integer] interface {
	PrintType()
	SetValue(value T)
	GetValue() T
	SetName(name string)
	GetName() string
	SetFields(value T, name string)
	PrintFields()
}

// Generic struct implementing the above interface
type Impl[T Integer] struct {
	data T
	name string
}

// As the interface entails setter methods, the interface is implemented on pointer receiver
func (i *Impl[T]) SetValue(value T) {
	i.data = value
}

func (i *Impl[T]) GetValue() T {
	return i.data
}

func (i *Impl[T]) SetName(name string) {
	i.name = name
}

func (i *Impl[T]) GetName() string {
	return i.name
}

func (i *Impl[T]) SetFields(value T, name string) {
	i.data = value
	i.name = name
}

func (i *Impl[T]) PrintFields() {
	fmt.Println("[PrintFields][Impl[T]] Data:", i.data, "name:", i.name)
}

func (i *Impl[T]) PrintType() {

	// Unlike value.(T), where T can be any specific type (int, string, slice, custom type), value.(type) is used only within switch case statements.
	// Here, as in the case of any(value).(T), any(value) is a simple type conversion that converts value's type to interface{} type
	switch any(i).(type) {
	case uint, uint16, uint32, uint64:
		fmt.Println("Unsigned")
	case int, int16, int32, int64:
		fmt.Println("Integer")
	case float32, float64:
		fmt.Println("Floating point")
	default:
		fmt.Println("Unknown type")
	}

}

func GenericStructIfaceImpl() {
	fmt.Println("-----------------------------Pointer to struct implementing an interface------------------------")
	var ifInt Iface[int]
	ifInt = &Impl[int]{
		data: 1,
		name: "<iface var 1>",
	}

	ifInt.PrintFields()
	ifInt.PrintType()

	implInt := Impl[int]{
		data: 2,
		name: "<impl var 2>",
	}
	implInt.PrintFields()
	implInt.SetFields(3, "<impl var changed to 3>")
	implInt.PrintFields()
	implInt.SetValue(4)
	implInt.PrintFields()
	implInt.PrintType()

	// As the interface is implemented on pointer receiver, any variable of its type can hold pointers to entities implementing it
	ifInt = &implInt
	ifInt.SetName("ifInt pointing to implInt ptr, not instance")
	ifInt.SetValue(5)
	fmt.Println("[Print members' values][Iface[T]]", "name:", ifInt.GetName(), "| value:", ifInt.GetValue())
}

// -----------------------------Nested structs-----------------------
type GenericInterface[N Numeric] interface {
	PrintName()
	GetData() N
	SetData(value N)
}

type Base[N Numeric] struct {
	data N
}

func (b *Base[N]) PrintName() {
	fmt.Println("[PrintName] Base struct", b.data)
}

func (b *Base[N]) GetData() N {
	return b.data
}

func (b *Base[N]) SetData(value N) {
	b.data = value
}

func (b *Base[N]) DoubleData() {
	b.data = b.data * 2
}

type EmbedStruct[N Numeric] struct {
	Base[N]
	name string
}

func (es *EmbedStruct[N]) PrintName() {
	fmt.Println("[PrintName][EmbedStruct] ", es.data)
}

type NestInstance[N Numeric] struct {
	inst Base[N]
	name string
}

func (ni *NestInstance[N]) PrintName() {
	fmt.Println("[PrintName]:", ni.name, ni.inst.data)
}

func (ni *NestInstance[N]) GetData() N {
	return ni.inst.data
}

func (ni *NestInstance[N]) SetData(value N) {
	ni.inst.data = value
}

func TestPolymorphism[N Numeric](input GenericInterface[N], value N) {
	input.PrintName()
	input.SetData(value)
	fmt.Println("[TestPolymorphism][GetData]:", input.GetData())
}

func NestedStructs() {
	fmt.Println("// -----------------------------Nested structs-----------------------")
	b := &Base[int]{data: 27}
	b.DoubleData()
	ns := &EmbedStruct[int]{Base: Base[int]{
		data: 28},
		name: "embed base struct",
	}
	ns.DoubleData()

	ni := &NestInstance[int]{inst: Base[int]{
		data: 13},
		name: "nested instance"}
	TestPolymorphism(b, -1)
	TestPolymorphism(ns, -2)
	TestPolymorphism(ni, -3)
}

func main() {
	NestedStructs()

	FuncIfaceImpl()
	CustomDataTypeIfaceImpl()
	CustomMapIfaceImpl()
	CustomSliceIFaceImpl()
	GenericStructIfaceImpl()
}
