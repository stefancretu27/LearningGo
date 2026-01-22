package go_didactical_apps

/*
*  Interface:
*	- an interface type defines a set of types. 
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
*	- the type set defined by such an interface is the set of types which implement all of those methods, and the corresponding method set consists exactly of 
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
*	- if  amethod is defined on pointer rceeiver, regardless it is called on pointer or value, the receiever will be treated as pointer and
* changes made to data memebrs persist
*	- if a method is defined on value receiver, changes made to data members will not persist, regardless it is called on pointer or value receiver
*	- if one struct embeds another, it can "inherit" its fields and methods, but it still remains a distinct type.
*	- a struct cannot directly qualify as another struct type because each struct type is distinct, even if they have identical fields. 
* Go enforces strict type safety, so two structs with the same fields are still considered different types and cannot be used interchangeably without explicit conversion.
*/

// declare a type constraint interface
type Number interface {
	int | float64
}

//declare another type constraint, with constraint to implement the listed methods
type MixedInterface interface {
	~int | ~int16 | ~int32 | ~int64 | string
	About(name string)
}

// use type constraint interface as generic parameter for func
func Add[T Number](a, b T) T {
	return a + b
}

// declare anpther type constraint interface
type Integer interface {
	Number | int16 | int32 | int64 | uint
}

// declare a basic interface
type IFace interface {
	PrintName(name string)
}

// attempt to implement interface using type constraint interface => cannot use type Integer outside a type constraint: interface contains type constraints
// func (i Integer) PrintName(name string) {
// 	fmt.Println("[PrintName] Integer", name)
// }

// attempt to implement interface using basic type => cannot define new methods on non-local type string
// func (s string) PrintName(name string) {
// 	fmt.Println("[PrintName] String", name)
// }

// as an alternative, basic types can implement interfaces (be extended with behaviors) if custom types are defined on them
type MyString string
func (s MyString) PrintName(name string) {
	fmt.Println("[PrintName] MyString", name)
}

type MySlice []int
func (s MySlice) PrintName(name string) {
	fmt.Println("[PrintName] Slice", name)
}

// custom function type whose signature matches the IFace only function's signature
type MyFunc func(name string)
func (mf MyFunc) PrintName(name string) {
	fmt.Println("[PrintName] MyFunc", name)
	mf(name)
}

type MyStruct struct {
	data map[string]int
	name string
}

func (ms *MyStruct) PrintName(name string) {
	fmt.Println("[PrintName] MyStruct", name)
}

//declare generic func
func GenericFunc[T Integer](data string, someInt T) {
	fmt.Println("[GenericFunc]", data, someInt)

	switch any(someInt).(type) {
	case uint:
		fmt.Println("[GenericFunc] unsigned")
	default:
		fmt.Println("[GenericFunc] signed")
	}
}

//declare a generic struct
type GenericStruct[T Integer] struct {
	data T
}

//generic struct implements interface
func (gs GenericStruct[T]) PrintName(name string) {
	gs.data = 27
	fmt.Println("[GenericStruct]", name, gs.data)
}

// generic interface with generic struct implementing the interface
type GenericInterface[N Number] interface {
	GenericMethod(value N)
}

type GenericIfaceStruct[N Number] struct {
	data N
}

func (gs GenericIfaceStruct[N]) GenericMethod(value N) {
	fmt.Println("[GenericMethod]:", value+gs.data)
}

type GenericIFace[T any] interface {
	PrintName()
	GetData() T
	SetData(input T)
}

type A[T any] struct {
	data T
	name string
}

func (a *A[T]) PrintName() {
	fmt.Println(a.name)
}

func (a *A[T]) GetData() T {
	return a.data
}

func (a *A[T]) SetData(input T) {
	a.data = input
}

//B will access directly A's methods. Those, being defined on A's fields, will access only A's fields. Inspite same naming of A and B fields, no name collision occurs
type B[T any] struct {
	A[T]
	data T
	name string
}

//C will access A's methods by dereferencing a field of type A[T]
type C[T any] struct {
	a    A[T]
	data T
	name string
}

func Interfaces_Structs() {
	//cannot use type Number outside a type constraint: interface contains type constraints
	//var numberInterface Number

	var ifaceVar IFace
	//cannot use map[string]float32{…} (value of type map[string]float32) as IFace value in assignment: map[string]float32 does not implement IFace (missing method PrintName)
	//ifaceVar = map[string]int{"pi": 3.14159}
	//cannot use []int{…} (value of type []int) as IFace value in assignment: []int does not implement IFace (missing method PrintName)
	//ifaceVar = []int{1, 2, 3, 4, 5}

	//custom types to basic types inetrface implementer
	var mySlice MySlice = []int{1, 2, 3, 4, 5}
	ifaceVar = mySlice
	ifaceVar.PrintName("interface variable used to call method after MySlice is assigned to it")
	
	//func interface implementer
	var myFunc MyFunc = func(name string) {
		fmt.Println("MyFunc", name)
	}
	ifaceVar = myFunc
	ifaceVar.PrintName("interface variable used to call method after MyFunc is assigned to it")

	//pointer to struct interface implementer
	var myStruct MyStruct = MyStruct{
		_data: map[string]float32{"pi": 3.14159},
		_name: "MyStruct",
	}
	ifaceVar = myStruct
	ifaceVar.PrintName("interface variable used to call method after MyStruct is assigned to it")

	//generics
	//structs always require type specification
	var gs GenericStruct[int16] = GenericStruct[int16]{30}
	gs.PrintName("generic struct instance implements interface")
	//type inferred from arguments
	GenericFunc("generic func arg", uint(14))

	var gi_pi GenericInterface[float32]
	gi_pi = GenericIfaceStruct[float32]{data: 3.14159}
	gi_pi.GenericMethod(1)

	var b B[int] = B[int]{A: A[int]{data: -3,
		name: "A int",
	},
		data: 27,
		name: "B int"}
	b.PrintName()
	fmt.Println(b.GetData())

	var c C[float32] = C[float32]{a: A[float32]{data: 3.14159,
		name: "A float32",
	},
		data: 2.7182,
		name: "C float32"}
	c.a.PrintName()
	fmt.Println(c.a.GetData())

}