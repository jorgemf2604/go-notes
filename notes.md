# The GO Programming Language
# 1. Basics 
# 1.1 Names 
## 1.1.1 Names of variables, constants, functions, types, etc. 
- They must begin with a letter or an underscore.
- Case matters: if a name begins with an upper-case letter, it is exported, which means that it will be visible and accessible outside of its own package.
- Use of camelcase (or pascalcase) when exported.  
## 1.1.2 Names of packages
They are always in lowercase. 
## 1.1.3 Keywords
Go has 25 keywords that cannot be used as names: break, func, go, import, ...
## 1.1.4 Predeclared names 
There are more names for built-in types (like int, string, bool, error, ...), constants (like true, false, nil, ...) and functions (like len, cap, new, append, ...) that are not reserved, so they can be used in declarations. On occasitions redeclaring them might make sense, although normally is not what you want and you should be careful. 

# 1.2 Lifetime, Scope, packages, files, imports
## 1.2.1 Lifetime of Variables 
- The lifetime of a package-level variable is the entire execution of the program.
- Local variables: typically allocated on the stack. Their lifetime ends when the function returns and the garbage collector reclaim the variable's storage. The exception is when a local variable (e.g. returned from a function or captured by a closure) escapes to the heap (extending its lifetime). 

```
package main

import "fmt"

var myGlobalPointer *int

func main() {
	f()
	fmt.Println(*myGlobalPointer) // 1
  // X escapes from f to the heap. 
}

func f() {
	var x int = 1
	myGlobalPointer = &x
}
```

## 1.2.2 Lifetime vs Scope
Lifetime: how long the variable exist in memory. 
Scope: Where in the code the variable is visible. Golang uses lexical scope, the scope is determined at compile time based code structure. 

```
func demo() {
    x := 10 // scope: inside demo()
    if true {
        y := 20 // scope: inside this if-block only
        fmt.Println(x, y)
    }
    // y is out of scope here
}

```

Here though count is declared inside makeCounter (the scope is limited to that function), its lifetime continues because the closure keeps a reference to it. It scapes to the heap. 
```
func makeCounter() func() int {
    count := 0 // declared inside function
    return func() int {
        count++
        return count
    }
}

```

We can use the scope of a variable to determine its visibility (where it can be used) 99% of the time. But sometimes observing the scope is not enough, we also need to keep an eye on the lifetime of the variable (escaping to the heap, extends the lifetime of a variable). 

## 1.2.3 Package level declarations vs Local level declarations 
Variables, constants, types and functions can be declared at package level. In this case these entities will be visible throughout all the files of the package. If this entities are declared inside the function, we are talking about local declarations and they are visible only within the function in with they are declared. 

## 1.2.4 Declaring a variable outside a function
As a general rule, we should only declare variables in the package scope when they are inmmutable:
- variables declared outside functions have package scope (they can be accessed and modified from anywhere in the program). It makes it harder to track changes/read and debug.
- concurrency is more difficult: when using concurrency in go we can have multiple goroutines accessing the same variable. If that variable has multiple points of modification across the package, the process of synchronization gets more complex.

## 1.2.5 Shadowing a variable
Shadowing happens when you declare a new variable with the same name as an existing one, but in a narrower scope. The inner variable hides the outer one.
```
x := 10

if true {
    x := 20 // shadows the outer x
    fmt.Println(x) // prints 20
}

fmt.Println(x) // prints 10

```
### The most commong shadowing bug
```
var err error

if something {
    value, err := compute() // new err shadows outer err
    if err != nil {
        return err
    }
}

fmt.Println(err) // still nil, even if compute() failed

```
How to avoid accidental shadowing:
1. Use `=` instead of `:=`
2. Declare variables earlier if needed
```
var err error
var value int

if something {
    value, err = compute() // reassigns outer err
    if err != nil {
        return err
    }
}

fmt.Println(err) // now reflects the real error

```

## 1.2.6 Packages, folders, files and modules 
- A package is a collection of go files under a folder. All the `.go` files in the folder share the same package name. All .go files in the same folder must use the same package name. The folder name does not matter to the compiler. The folder name and package name can be different (although it is often the same -> strong developer's convention). There are cases when it is OK to have a different names (e.g. the `internal` folder, a special folder in GO for internal use only, EXTERNAL packages (outside the app/module) are not allowed to import it. Internal will have inside different packages with different names such as auth or bd) 

- two kinds of packages: executable and libraries. To build a runnable program (a binary executable), Go requires a package main and function named main() inside that package. Several go files in the same folder, importing the package main, will merge them into one executable (as long as there is only one main() function per package main). The other packages are not exucutable, there are libraries tto be imported.   
- why packages: modularity (break programs into logical pieces), reusability (write once, import it anywhere), namespace control (avoid name collisions), encapsulation (control what is private and what is public) 

- The folder name has relevance when importing a package. The import path is based on the folder path, not the package name.  
```
import "github.com/jorgemf2604/myapp/utils"

// Go does not look for a package named utils. It looks for a folder named utils inside your module. Inside the folder the .go files can declare any package name (e.g. `helpers` package). You will use it like this:

helpers.DoSomething()
```
This gives you flexibility. You can have a folder net/http but the package name is just http.

- A Module is a collections of packages. 
`go.mod` defines the module path.
```
module github.com/jorgemf2604/myapp
```

Other people will import your module (or a package) using that path.
```
import "github.com/jorgemf2604/myapp/utils"  // we are only importing the package utils in our module. 
```
## 1.2.7 Imports
Go looks for imports in three places: 
1. the standard library:
```
import "fmt"     // standard library
```
2. Your module: 
```
package main

import (
	"fmt"

	u "jorge.martin/hello/my-utils"
)

func main() {
	fmt.Println(u.Add(10, 5))
}

// the name of my module inside go.mod is `jorge.martin/hello`
// the name of my package is utils 
// the name of the folder that contains the files of the utils package is `my-utils`
// the alias used to refer to the package contained in the import path is `u`
// The import path (what we write in quotes) is the module path (jorge.martin/hello) + the folder path (/my-utils). Go uses the import path to locate the code in disk or online (a github account is common). It does not care what the package is named inside that folder.
// Because the import path can be a mouthful, you can use an alias to refer to it (it can be totally different from the package name inside the folder)
```
3. External packages (Go fetches these from the internet and stores them in your module cache.):
```
import "github.com/gin-gonic/gin"
```

### Type of imports 
1. Single:
```
import "fmt"
```
2. Multiple: 
```
import (
    "fmt"
    "math"
)
```
3. Renaming imports (alias):
```
import m "math"
```
4. Blank identifier (they run the package's init() function but does not bind any name, used in databases config):
```
import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver; imported for side‑effects only
)

```
Why is the black identifier useful ? 
The driver package runs its init() function. That init() registers the driver with database/sql under a name like "postgres". You never call anything from the driver package directly. If you imported it normally, Go would complain that the package is unused. So the blank identifier tells Go: “Import this package, run its initialization, but don’t bind it to a name (it is ok if I do not use anything from the package directly).”

5. Dot import (import all exported names in your file's namespace). Not recommended.
```
import . "fmt"

Println("hello") // no fmt.Println

```
## 1.2.8 Package initialization 
Go will load packages in dependency order (A package is initialized only after all the packages it depends on have been initialized): If main imports A, and A imports B, and B imports C, the initialization order is C -> B -> A -> main. The initialization of each package involves: 
- Evaluate all package-level variables 
- Run all init() functions in that package

So the full sequence is: 
1. Load packages in dependency order
2. Initialize package-level variables (they will do this in the order they appear. If there are multiple files in the same package, go will initialize variables in lexical file name order, that is alphabetically by file name)
3. Run all init() functions in that package
4. Move on to the next package (repeat if needed)
5. Finally run main.main()

## 1.2.9 init() function 
init() is a special function in Go that automatically runs (you never calls it yourself) after all package-level variables are initialized, before any other code in that package is used, and before main.main() runs. 

- You can have multiple init() functions in the same file or across multiple files in the same package (they run in lexical file name order)
- You cannot call init() yourself.
- You cannot pass arguments.
- You cannot return values from it. 
- It runs once per package even when multiple packages import it. 
- It runs after package-leve variables initialization.
- If runs in dependency order: if main imports A, and A imports B, the order will be B.init() -> A.init() -> main.init() -> main.main() 

# 1.3 Variable declarations 

1. Basic declaration 
```
var name string  // initialize to ""
```
2. Declaration with initialization 
```
var name string = "Jorge"
```
3. Type inference 
```
var name = "Hello World"
```
4. Short variable declaration (only inside functions)
```
message := "Hello"
```
5. Multiple variable declaration
```
var i, j, k int
var b, f, s = true, 2.3, "four"
var file, err = os.Open(fileName) // a set of variables being initialized by calling a functions that return multiple values
```
6. Grouped declaration
```
var (
  a int
  b int = 1
  c string = "hello"
)
```

## 1.3.1 Difference between short variable declaration (:=) and normal declaration (with var)
- The short variable declaration is only valid inside a function, where var declarations work everywhere (inside of funtions and at package level). 
- The short variable declaration requires an intial value (and the type will be inferred from that value), whereas var can be used with an intial value (explicit or inferred) or without an initial value (zero value). 

```
file1, err := os.Open(filename)
// ...
file2, err := os.Open(anotherFileName)    // file2 is declared and err is reassigned. 
```
- Do not confuse a muti-variable declaration with a tuple assignment
```
i, j := 0, 1   // variable declaration
i, j = j, i    // swap values of i and j
```


## 1.3.2 Short variable gotcha
It declares a new variable every time you use it. 
```
x := 10 
x := 20 // ❌ compile error: x already declared
```
BUT as long as there is at least one new variable on the lefthand side of the := , it will reassign existing variables.
```
	x := 10
	x, y := 30, "Hello"    (reassign of x, declaration of y)
```

## 1.3.3 Tuple Assigment 
With tuple assignment, we can assign several variables at one. 
```
package main

import "fmt"

func main() {
	x, y := 0, 1    // tuple 'declaration-assignment'
	fmt.Println(x, y)
}
```

```
package main

import "fmt"

func main() {
	x, y := 0, 1
	x, y = 10, 20     // tuple assignment
	fmt.Println(x, y) // 10, 20
}
```

It is very useful to swap values without creating an itermediate variable.
```
package main

import "fmt"

func main() {
	fmt.Println(fib(6)) // The 6th element of the fibonacci series is 8 [1, 1, 2, 3, 5, 8]
}

func fib(n int) int {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y // we swap the values here, without the need of creating an intermediate variable.
	}
	return x
}
```

It is useful when certain functions calls return mutiple values.
```
f, err := os.Open("text.txt)
f, err = os.Open("foo.txt")
_, err = io.Copy(destination, source)  // we discard the byte count
```

## 1.3.4 Unused variables 
In Go, every declared LOCAL variable must be read.
```
package main

var age int32 = 45 // no problem

func main() {
	name := "Jorge" //  error: name declared and not used.
}

```
This unused variable check works well but it is not perfect.
```
func main() {
	x := 10     // this assignment isn't read
	x = 20
	fmt.Println(x)
	x = 30     //  this assignment isn't read
}
```
This check is for LOCAL VARIABLES, Go will not stop you from creating unread package-level variables or unread const (constants in GO are calculated at compile time and cannot have side effects. If a const is not used, it is simply not included in the compiled binary). 

## 1.3.5 const 
Go has a way to declare a value as immutable with the const keyword.

We can declare a const at package level or at function level.

const in Go are very limited (they are basically a way to give names to literals). const can only hold values that the compiler can figure out at compile time. - It can holds : numeric literals, true and false, strings, runes, expressions using previous elements, ... - It cannot hold anything else: we cannot using const to declare immutable arrays, slices, map or structs.

Constants can be typed or untyped. An untyped constant works like a literal (it is untyped and has a default type when no other type can be inferred).
```
// untyped constant
const x = 10   // default type int
var z float64 = x   // untyped

// typed constant: it can be assigned directly only to an int 
const typedX int = 10
```

# 1.4 Zero values 
In Go there are not uninitialized variables. If the expression is omited (e.g. in basic declarations), the value will be the zero value of the type. `0` for numbers, `false` for booleans, `""` for strings, and nil for interfaces and reference types (slice, pointer, map, channel, function). The zero value of an aggregate type like an array or a struct has the zero value of all its elements or fields.

```
func main() {
	var isAdmin bool
	var age int
	var name string
	var percentage float64
	fmt.Printf("%#v, %#v, %#.f %#v\n", isAdmin, age, percentage, name) // false, 0, 0. ""
}
```

# 1.5 Literals
A literal in Go is simply a fixed value written directly in your source code — a value that is not computed, returned from a function, or stored in a variable first. 
Go has 4 type of literals:
- Integer literal: they are base 10 by default, but you can write them in different bases (Ob for binary, Oo for octal, or Ox for hexadecimal)
- Floating-point literal: again there are differnt ways of writting them.
- Rune literal: represents a single Unicode character. A rune type is an alias for the int32 type. It can be written in different ways:
	```
	unicode character   'a'          ---->  97 decimal
	8-bit   octal       '\141'       ---->  1*8^1 + 4*8^1 + 1*8^2 = 1 + 32 + 64 = 97 
	8-bit   hexadecimal '\x61'       ---->  1*16^0 + 6*16*1 = 1 + 96 = 97
	16-bit  hexadecimal '\x0061'     ---->  1*16^0 + 6*16*1 = 1 + 96 = 97
	32-bit  hexadecimal '\U00000061  ---->  1*16^0 + 6*16*1 = 1 + 96 = 97
	```
- String literals: there are 2 types (interpreted string literals and raw string literals).
	+ For interpreted string literals you should use double quotes. They interpret rune literals (both numerics and backslash escaped) into single charcters:
	```
   func main() {
	var msg1 string = "Greetings.\nThis is the character \x61"
	fmt.Println(msg1)
	// Greetings.
	// This is the character a
	}
	```
	+ If you need to inclide backlashes, double quotes or newlines in a string, use a raw string literal. They are delimited with backquotes. There are no escape characters in a raw string literal, all chatacters are included as is.
	```
	func main() {
	var msg1 string = `Greetins and 
	"Salutations"`
	fmt.Println(msg1)
	// Greetings and
	//      "Salutations"
	}
	```

--- 
### Literals are untyped 
Literals in Go are untyped. They can be used with variables whose type is compatible with the literal. Being untyped goes only so far, it must makes sense (e.g. you cannot assing a string literal to a int variable) and size limitations also exist (e.g. you cannot assign the literal 1000 to a variable of type byte).

```
func main() {
	var x float64 = 10
	var z int32 = 10
	var y float64 = 2.5 * 10
	fmt.Printf("Type of x: %T, type of y: %T, type of z: %T\n", x, y, z)
	// Type of x: float64, type of y: float64, type of z: int32
}
// In this example, both literals 10 are untyped and can be assgined to a float64, or a int32 variable, or used in a float64 expression alongside a float literal.
```
### There is a default type for literals
If there are is nothing in an declaration that makes clear what the type of the literal is, the literal defaults to a type.
```
func main() {
	a := 10                                         // int
	b := 1.45                                       // float64
	c := 'a'                                        // rune(int32)
	d := "Jorge"                                    // string
	e := `This is "great"`                          // string
	fmt.Printf("%T, %T, %T, %T, %T", a, b, c, d, e) // int, float64, int32, string, string
}
```

# 1.6 Pointers 
A variable is a named container for some value (it has a memory address and a name). A pointer variable is a variable that holds the memory address of another variable.
 
We use the `&` operator to get the address of a variable and the `*` operator to dereference (access the value stored in that address).

```
package main
import "fmt"

func main() {
    var x int = 10
    var p *int = &x   // p points to x
    fmt.Println("Address of x:", p)     //  Address of x: 0x10328000
    fmt.Println("Value of x via pointer:", *p)   // Value of x via pointer: 10
}

```
With a pointer we can read or update the value of the variable indirectly, whithout using or even knowing the name of the variable. 
Passing a pointer argument to a function makes it possible for the function to update the variable that was indirectly passed. Something similar when we passed reference types like slices, maps or channels, and even structs, arrays or interfaces that contains these types. 

```
package main

func main() {
	v := 1
	increment(&v)
	println(v) // 2

}

func increment(aValue *int) {
	*aValue++
}

```

The zero pointer for a pointer OF ANY TYPE is `nil`. Comparing against nil is common to check if a pointer is valid. If we dereference a nil pointer GO we wil panic.   
```
var p *int
if p == nil {
    fmt.Println("Pointer is nil")
}
```

```
var p *int
fmt.Println(*p)   
// In terminal =>  panic: runtime error: invalid memory address or nil pointer dereference
```

## 1.6.1 Escape analysis
In languages like C, returning the address of a local variable is unsafe because local variables are allocated on the stack, which is destroyed/reclaimed after the function returns. Dereferencing such pointer leads to underfined behaviour. 
In Go, it is safe for a function to return the address of a local variable. Go's compiler and runtime perform an escape analysis. If a local variable's address is returned, the compiler automatically allocates that variable on the heap instead of the stack. This ensures the variable remain valid after the functions exists. 
```
package main
import "fmt"

func newInt() *int {
    x := 10
    return &x // safe: x escapes to the heap
}

func main() {
    p := newInt()
    fmt.Println(*p) // Output: 10
}

```

# 1.7 The new function

The expression `new(T)` creates an unnamed variable of type T, initializes it to its zero value of T, and returns its address, which is a value of type `*T`. 

```
package main
import "fmt"

func main() {
    p := new(int)   // p is of type *int
    fmt.Println(*p) // Output: 0 (zero value of int)

    *p = 42
    fmt.Println(*p) // Output: 42
}

```

With the `new` function there is no need to invent and declare a dummy name. So instead of this

```
func newInt() *int {
  var dummy int
  return &dummy
}

```
we can do this:
```
func newInt() *int {
  return new(int)
}
```

The new function is not very used because the most common unnamed variables are of struct types, for wich using a composite literal with & is more flexible. 
```
package main

import "fmt"

func main() {
	type Person struct {
		Name string
		Age  int
	}

	p1 := &Person{Name: "Jorge", Age: 45}
	p2 := &Person{} // if you omit fields, they will take zero values

	fmt.Println(*p1) // {Jorge 45}
	fmt.Println(*p2) // {      0 }
}

```

It is also a predeclared name, not a keyword, so it can be redefined unintentionally.

```
func myFunction(old, new int) int {
  return new - old
}

// Inside myFunction the `new` function is unavailable.

```

# 1.8 Type declarations 
Apart from the types provided by the language, we can create new names types with the same underlying type as an existing type. 

```
package main

import "fmt"

// in these case the types are exported (accessible from other packages)
type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/6 + 32)
}

func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func main() {
	fmt.Printf("%.2f\n", CToF(100)) // 182.00
	fmt.Printf("%.2f\n", FToC(78))  // 25.56
}
```

Why don't we use float64 for both? Distinguishing the types makes it possible to avoid errors like inadvertenly combining temperatures in the two different scales. The functions accept and return different types, inside the body you must explicitily use type conversion. 
- By doing this, you give the compiler extra information and it will stop you from accidentally mixed them (type safety). 
- Apart from type safety, named types makes code more readable (self-documentation). 
- You can also define methods on named types, enriching their behaviour. 
```
package main

import "fmt"

type Celsius float64
type Fahrenheit float64

func (c Celsius) String() string {
	return fmt.Sprintf("%gC", c)
}

func main() {
	var myTemperature Celsius = 33
	fmt.Println(myTemperature.String()) // 33C
}
```

Named types are truly conveninet when dealing with complicated types, like structs. 

# 1.9 Type conversion
1. Go does not have automatic type promotion. Go forces you to convert them explicitly. 
```
var x int = 10
var y float64 = 3.5 
z := x + y // ❌ compile error
```
```
z := float64(x) + y // ✔️ explicit conversion
```
As a result of this, Go does not truthy/falsy values. Go never converts types implicitly, not even to booleans. 

---

2. A string can be converted back and forth to a slice of bytes or a slice of runes.
```
fmt.Println(string([]rune("Hello, 世界")[8])) // 界
```

---

3. A rune or a byte can be converted to a string
```
func main() {
	var a rune = 'x'
	fmt.Println(a) // 120
	var s string = string(a)
	fmt.Println(s) // x
	var b byte = 'y'
	fmt.Println(b) // 121
	var s2 string = string(b)
	fmt.Println(s2) // y
}
```

---
4. Common mistake: conversion from int to string yields a string of one rune, not a string of digits. 
```
func main() {
	var x int = 65
	var y = string(x)
	fmt.Println(y)   // A not "65"
}
```
Use strconv.Itoa() instead.

```
package main

import (
	"fmt"
	"strconv"
)

func main() {
	num := 123
	str := strconv.Itoa(num)
	fmt.Println(str) // Output: "123"
}
```
# 1.10 Categories of data types 
1. Basic data types (numbers, strings and booleans)
2. Aggregate data types (arrays and structs)
3. Reference types (pointers, slices, maps, functions and channels)
4. Interface types 

# 2. Basic data types
1. Numbers
2. Strings 
3. Booleans 

# 2.1 Numbers
We can divide them in 3 groups:
- Integer types
- Floating-point types
- Complex types

## 2.1.1 Integer types 
- int8
- int16
- int32
- int64
- uint8
- uint16
- uint32
- uint64
- byte = is an alias for uint8
- int = on a 32-bit CPU is a int32 and on most^1 64-bit CPU is a int64. Integer literals default to being of type int.
- uint = same rules as int but it is unsigned (0 or positive).
- rune = alias for an int32
- uintptr = it is used to hold pointer-sized integers. A 32-bit pointer on a 32-bit system and a 64-bit pointer on a 64-bit system.

NOTE: In the beginning, when GO was created, the language allowed more flexibility and some compilers could choose different int/uint sizes depending on implementation details. So you may read some outdated info in books like "Some uncommon 64-bit CPU architectures use a int32. Go supports some of them: (e.g. amd64p32)". That is not true anymore, the GO teams treats the size of int/uint as effectively fixed by architecture. All major compilers (gc, gccgo, tinygo) follow the same rule. Therefore, an int/uint is 32 bits in a 32-bit system and 64 bits in a 64-bit system. 
### 2.1.1.1 Zero value of an int
It is zero 

```
package main

import "fmt"

func main() {
	var number int
	fmt.Println(number) // 0
}

```

## 2.1.1.2 Default type of a integer literal
It is an int

```
package main

import "fmt"

func main() {
	number := 34
	fmt.Printf("The value %v is of type %[1]T\n", number) // The value 34 is of type int
}

```

## 2.1.1.3 Range of int
- Signed numbers: -2^(n-1) to 2^(n-1) - 1
- Unsigned numbers: 0 to 2^(n) - 1 

NOTE: GO provide unsigned numbers but it is a good idea to use singned numbers (like int) for quantities that cannot be negative, such as the length of an array. For example, the built-in `len` function returns a signed int. The alternative would be problematic:

```
medals := []string{"gold", "silver", "bronze"}
for i := len(medals) - 1; i >= 0; i-- {
	fmt.Println(medals[i])   // bronze, silver, gold
} 

// After the 3 interation, in which i == 0, the i-- statement would cause i not to become -1 but 2^64-1, and the program would panic trying to access an element outside the bounds of the slice.
```

## 2.1.1.4 Operators 
1. There are binary operators for:
- arithmetic
- logic
- comparison
- bitwise operations

2. There are unary operators: (+ and -). Both operators do not mutate the original value, they produce a new value. The unary plus (+x) returns the value unchanged (does nothing). The unary minus (-x) negates a number: if x is positive, -x is negative, if x is negative, -x is positive, if x is zero, -x is still zero. Both operators preverse the operands' type (e.g. if x is an int, -x will be a int).

There is a order of predecende in these operators (do not memorize it, use parenthesis to clarify)

## 2.1.1.5 Type conversion with int
Normally, explicit conversion is required to convert a value from one type to another. 
Converting a floating-point number to an int will discard any fractional part, truncating towards zero. 

## 2.1.1.6 int are comparable and can be ordered. 
Integers (and all of the rest basic types: booleans, numbers, strings) are comparable (we can use == or != on them). Furthermore, integers, floating-point numbers and strings are ordered (we can use >, <, etc.). No other data types are ordered. 

## 2.1.1.7 Module operator
The sign of the remainder is always the same as the sign of the dividend.
5%3 is 2
5%-3 is 2
-5%3 is -2
-5%-3 is -2

## 2.1.1.8 Division operator
- Integer division truncates the result :  5/4 is 1
- If we divide an integer by zero, GO will panic

# 2.1.2 Floating-Point Numbers
Go provides two sizes of floating-point numbers, `float32` and `float64`. A float32 provides approximately six decimal difits of precision, whereas a float64 provides about 15 digits. 

## 2.1.2.1 Zero value of floating-point numbers
It is zero

```
package main

import "fmt"

func main() {
	var number float64
	fmt.Println(number) // 0
}
```

## 2.1.2.2 Default type of floating type literals 
float64

```
package main

import "fmt"

func main() {
	number := 5.21
	fmt.Printf("The type of %v is %T\n", number, number) // The type of 5.21 is float64
}
```
## 2.1.2.3 Floating-point division
Dividing a non-zero floating-point variable by 0 returns +Inf or -Inf.
Dividing a a floating-point variable set to zero by 0 retuns NaN (not a number)

```
func main() {
	var num1 float64 = 5.5
	var num2 float64 = 0
	result := num1 / 0
	fmt.Println(result) // +Inf
	result2 := num2 / 0
	fmt.Println(result2) // NaN
}
```

## 2.1.2.4 Floating-point comparisons
Go let us use == and != to compare floats, but do not do it. Instead, define a maximum allowed variance (epsilon) and see if the difference between the floats is less that that.

package main

import (
    "fmt"
    "math"
)

func areEqual(a, b, epsilon float64) bool {
    return math.Abs(a-b) <= epsilon
}

func main() {
    num1 := 0.1 + 0.2
    num2 := 0.3
    epsilon := 1e-9 // Define your tolerance (e.g., 1e-9)

    if areEqual(num1, num2, epsilon) {
        fmt.Println("The numbers are equal within the specified tolerance.")
    } else {
        fmt.Println("The numbers are NOT equal.")
    }
}

I HAVE A DOUBT THOUGH: this shouldn't work. Why does it work ?

```
func main() {
	num3 := 0.1 + 0.2
	num4 := 0.3
	fmt.Printf("num3: %.70f\n", num3)
	fmt.Printf("num4: %.70f\n", num4)
	fmt.Println(num3 == num4)
}

// num3: 0.2999999999999999888977697537484345957636833190917968750000000000000000
// num4: 0.2999999999999999888977697537484345957636833190917968750000000000000000
// true
```

This happends only because Go treats numeric literals as high‑precision constants and reduces both sides to the same exact rational value before converting to float64. So both become the same binary64 number. This is a special case, don't rely on this, always use tolerance. 

# 2.2 Strings
There are two ways of thinking about strings:
- What they are under the hood: a slice of bytes (alias for uint8).
- What they look to us: a slice of runes (alias for int32 that represents UTF-8 encoded Unicode code point).

## 2.2.1 The length of a string
When we use the len() function on a string, we get the length in bytes. If what we want is to know the number of characters we can use `utf8.RuneCountInString` to know the number of characters in a string.

```
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	str := "Hello, 世界!"
	length := utf8.RuneCountInString(str)
	fmt.Printf("The number of characters in the string is: %d\n", length)
	fmt.Printf("The number of bytes in the string is: %d\n", len(str))
}

// The number of characters in the string is: 10
// The number of bytes in the string is: 14
```

## 2.2.2 Indixing and slicing with strings 
Indexing a tring gives you a byte, not the rune/character. Slicing counts positions in bytes. Use len, index or slicing with strings only when you know you are dealing with a string that does not contain characters that take more that one byte. Another sollution is to use a for-range loop to iterate over code points and to use functions in the `strings` and `unicode.utf8` packages. 

Some solutions
-------------

1. How to know the length of a string: Use utf8.RuneCountInString to know the number of characters in a string.
2. How to index:  Convert your string into an slice of runes, index the rune you want and covert it into a string.
```
package main

import "fmt"

func main() {
	fmt.Println(string("Hello"[1]))             //    e  ASCII only
	fmt.Println(string([]rune("Hello, 世界")[1])) //  e  UTF-8
	fmt.Println(string([]rune("Hello, 世界")[8])) //  界   UTF-8

	fmt.Println([]rune("Hello, 世界"))            // [72 101 108 108 111 44 32 19990 30028]
	fmt.Println([]rune("Hello, 世界")[8])         // 30028
	fmt.Println(string([]rune("Hello, 世界")[8])) // 界
}

```
3. How to slice: convert the string to a slice of runes, slice the rune, convert the rune slice into a string.
```
func main() {
	str := "Hello, 世界"
	runes := []rune(str)
	sliced := string(runes[7:9]) // Slicing the Unicode string
	fmt.Println(sliced)          // Output: 世界
}

```

## 2.2.3 Unicode and UTF-8
Unicode is a standard, UTF-8 is an encoding. 

- Unicode is a massive catalog of characters and each character has a unique number (code point).
- UTF-8 is one of several ways of encode Unicode code points into bytes. They are diffent encodings for Unicode (e.g. UTF-32 wich uses 4 bytes for every code point). UTF-8 is a variable encoding: from 1 to 4. It lets you use only 1 byte to represent the Unicode characters which values are below 128 (most letters, numbers and punctuation in the English language). It expands up to for bytes for more complex characters. UTF-8 was created by Ken Thompson and Robert Pike, 2 of the creators of Go.

A string is a sequence of bytes, but GO assumes those bytes represent UTF-8 encoded text. Go expects UTF-8:
- Source files must be UTF-8
- String literals are parsed as UTF-8
- for range loop assumes UTF-8
- Many standard library functions assume UTF-8
But Go does not enforce UTF-8:
- Invalid UTF-8 is allowed. 
- You can treat a string as a row binary data. 

Go is Unicode-friendly but it can be low-level and byte-oriented.

## 2.2.4 String literals 
There are two types: interpreted literals and raw string literals. 

## 2.2.5 Strings are immutable 
Strings in GO are immutable. 

```
s := "hello"
s[0] = 'H' // compile error: cannot assign to s[0]

```

```
s := "hello"
t := "H" + s[1:]
fmt.Println(t) // Hello
// this is allowed: you are creating a new string
```

Workaround: strings are immutable but you can convert the string to a slice of runes, change the rune un the slice, and convert the slice of runes into a new string.

```
r := []rune("你好")
r[0] = '您'
s := string(r)
fmt.Println(s) // 您好

```

## 2.2.6 Strings under the hood and complexity of slicing and the len function
Under the hood, a string is a two-field descriptor:
- with a pointer to immutable bytes 
- the length (number of bytes)

```
// pseudocode representation of a string under the hood
type string struct {
    ptr *byte  // pointer to the first byte
    len int    // number of bytes
}

```

- Strings store bytes, not characters.
- String does not store capacity (like slices). 


This is why:
1. slicing a string does not copy bytes (zero copy slicing)
```
s := "hello world"
t := s[:5] // "hello"

```
t becomes a new string header:
ptr → points to the same underlying byte array as s
len → 5

No bytes are copied.
Only the descriptor (pointer + length) is new. That is why slicing a string in GO has O(1) complexity.

***This is possible because strings are immutable, so sharing memory is always safe.***

2. len(s) is an O(1) operaration. Why are we storing the length of the array? We could use just a pointer and have the len() function iterate over every element of the array until reaching some kind of null terminator (like the C language does). Go prefers to store the pointer and the length in a heavier 'object' so retrieving the len is a O(1) operation instead of O(n). In C using the strlen function in a loop condition can be diabolical.

## 2.2.7 Complexity of string concatenation
In languages with immutable strings, each concatenation must allocate a new string and copy the contents of both operands. This makes a single concatenation O(M+N), and repeated concatenation in a loop O(N²) overall.

Languages with mutable strings or builder types avoid this cost by appending into a growable buffer. In these cases, repeated concatenation is O(N).

Go’s strings are immutable, so concatenation using + or fmt.Sprintf inside a loop is O(N²). To avoid this, Go provides strings.Builder, which uses a growable byte buffer. Appending to a strings.Builder is amortized O(1), and building a string of total length N is O(N).

NOTE: What does amortized O(1) means? => An operation is amortized O(1) amortized when most individual operations are constant time, occasionally an operation is more expensive (like O(n)), but spread across many operations, the average cost of the operations stays constant. Classic example: 	Dynamic array growth. Imagine a slice doubling its capacity when full: Most append calls just write a value → O(1). Occasionally, the slice must grow: allocate a new array, copy all existing elements → O(n). But this expensive step happens rarely (logarithmically), so the average cost per append stays constant.

# 2.3 Booleans 
- A boolean has only two possible values, true or false. 
- There is no implicit conversion from a boolean value to a numeric value or vice versa. There are no truthy/falsy values in GO.
- The zero value of a boolean is false 
```
package main

import "fmt"

func main() {
	var isAdmin bool
	fmt.Printf("The value %v is of type %[1]T\n", isAdmin) // The value false is of type bool
}
```

# 3. Aggregate types
# 3.1 Arrays 
1. An array in GO is a fixed-size contiguous block of memory. 
2. They are composed of one type only (we cannot mix different data types in the array).
3. The size of the array is baked into the type: [3]int and [4]int are different types.   
4. Arrays are value types (the actual data is stored in the variable, assigning it to another variable copies the entire value, passing it to a function passes a copy, mutating the copy does not affect the original).
5. It is a contiguous block of memory with no header, no metadata, no pointer indirection. This makes indexing extremely fast ( a[i] = \*(base + i\*sizeof(T)) )
```
var a [4]int

// In memory looks like this:
// | int | int | int | int |
```
## 3.1.1 Arrays uses
Some use cases:
- Hash keys. Cryptographic hashes return fixed-size outputs, and Go’s crypto packages use arrays for this reason.
- Avoid allocations by keeping the data on the stack (if they don't escape). Because arrays are value types with fixed size, the compiler can keep them on the stack, avoiding heap allocations.
- Interop with the C language where fixed-sized buffer matter. C APIs often expect fixed-size arrays, not pointers to dynamically sized memory.
- Embeding arrays inside structs for predictable memory layout. When you embed arrays inside structs, the struct becomes a pure value type with no internal pointers. The struct has a stable, predictable memory layout. No hidden allocations. Can be copied safely. Can be compared if all fields are comparable.

Even if you are not interested in these cases, arrays are important to understand because they provide backing store for slices, that are more used.  

## 3.1.2 Array declaration
1. Basic array declaration 
```
var a [5]int
// declares an array of exactly 5 ints
// all elements are zero-initialized:  [0 0 0 0 0]
// the type is [5]int
```

2. Declaration with initialization
```
var a = [3]int{1, 2, 3}
b := [3]int{4, 5, 6}    // short syntax
```

3. Inferred length 
```
a := [...]int{10, 20, 30, 40}
```

4. Partial initialization
```
a := [5]int{2: 100}   // [0 0 100 0 0]

```

5. Arrays of structs 
```
type Point struct {
    X, Y int
}

points := [2]Point{
    {1, 2},
    {3, 4},
}

```

5. Nested arrays
```
matrix := [3][3]int{}

```

==================== Arrays in other situations =======================

6. Arrays of structs 
```
type Point struct {
    X, Y int
}

points := [2]Point{
    {1, 2},
    {3, 4},
}

```

7. Array declaration inside a struct
```
type Packet struct {
    Header [4]byte
    Body   [128]byte
}
// This is common in protocols, crypto, and binary formats.
```

8. Array of pointers 
```
var arr [3]*int
// Still a fixed-size array — only the element type changes.
```

9. Array returned from a function
```
func makeID() [16]byte {
    var id [16]byte
    // fill id...
    return id
}
```

10. Array literal assigned to a variable
```
a := [3]int{1, 2, 3}
b := a // full copy

```

## 3.1.3 Comparability
Arrays are comparable. Arrays are value types, their size is part of the type, and the compiler knows exactly how to compare them. Two arrays can be compared if:
- They have the same length
- They have the same element type
- The element type itself is comparable 

```
[3]int == [3]int   // OK
[3]int == [4]int   // compile error (different types)

```
```
[2][]int           // not comparable (slice elements are not comparable)

```


As long as the struct is comparable, the array is comparable.
```
type Point struct {
    X, Y int
}

a := [2]Point{{1, 2}, {3, 4}}
b := [2]Point{{1, 2}, {3, 4}}

fmt.Println(a == b) // true

```

## 3.1.4 Working with arrays is hard
1. The size of the array is part of the type of the array and you cannot use type conversion.

```
var a [3]int
var b [4]int

a = b // compile error: cannot use b (type [4]int) as type [3]int

```

```
func main() {
	var x = [4]int{1, 2, 3, 4}
	var y = [3]int{1, 2, 3}
	x = [4]int(y) // error: IvalidConversion
}
```
2. You cannot use a variable to specify the size of an array.
```
n := 5
var a [n]int // compile error: n is not a constant

```

3. You cannot write a function that works with arrays of any size, you must specify the size. If you want to write a function that works for a collection of any size use slices or generics.

```
func makeArray(n int) [n]int { // compile error
    return [n]int{}
}

```

4. No dynamic resizing.

# 3.2 Structs 
## 3.2.1 Structs under the hood
1. A continous block of memory 
```
type User struct {
    Age   int32
    Score int64
    Flag  bool
}

// memory
//   | Age (4 bytes)  | padding (4 bytes) |
//   | Score (8 bytes)                   |
//   | Flag (1 byte)  | padding (7 bytes) |
// total 24 bytes

// Everything is laid out in order, with no pointers or metadata unless you explicitly add them. Why padding? Go align fields so the CPU access them efficiently. 
```

The order of the fields is important and affect size. 

Example of normal order:

```
type Bad2 struct {
    A int32
    B int64
    C int32
}

//  | A (4) | padding (4) |
//  | B (8) |
//  | C (4) | padding (4) |
// total 24 bytes
```

Reordered: 

```
type Good2 struct {
    B int64
    A int32
    C int32
}

//  | B (8) |
//  | A (4) | C (4) |
//  16 bytes (same fields, 33% smaller in size)

```

Rule of thumb when sorting fields:
- Write fields from largest to smallest
- Group fields of the same size together
- Put booleans and bytes at the end. 

Doing this is matters in performance-critical code or with large arrays of structs (if the array has millions of elements, saving 8 bytes per struct is huge)

2. Struct are value types: In GO, a struct is not a reference type. A struct is	a value, just like an int or a float64. That means:
- Assigning a struct copies all its fields: that is why we should use pointers for large structs. 
- Passing a struct to a function copies it
- Methods on value receivers operate on a copy (does not mutate the original)

3. No hidden pointers unless you add them: A struct in GO is a flat block of memory (continuous bytes), laid out in field order, padded for alignment, and coppied as a whole when assignment. There is no hidden header, no reference count, no runtime metadata stored inside the struct. If a struct contains pointers, it is because you put them there. 

```
type Node struct {
    Value int
    Next  *Node
}

```
The GO compiler tracks pointers, but even in these cases the struct does not store metadata. The metadata lives inside the type descriptor, not in the struct instances. Struct instances stay as a pure data blob.  

4. The compiler try to keep them in the stack. Go compiler prefers stack allocation whenever possible because:
- The stack is fast
- Allocation is just pointer bumping (*)
- Deallocation is free (pop the stack frame)


All of this, have implications:
- If a struct does no escape its scope, it stays on the stack.
```
func f() {
    p := Point{1, 2}
    fmt.Println(p.X)
}

// p is stack allocated because it never leaves the function and no pointer to it escapes. This is true even for large structs. 

```

- If a struct escapes, it moves to the heap.
```
func f() *Point {
    p := Point{1, 2}
    return &p
}

// The compiler rewrites this as a heap allocation.
```
- Passing a struct by value helps keep it on the stack
```
func process(p Point) { ... }

```
- Pointer receivers do NOT automtically cause heap allocation.
```
func (p *Point) Move() { p.X++ }
// Calling this does not force Point onto the heap. The compiler only allocates on the heap if the pointer itselft escapes.	
```

- The struct size does not matter. Even a huge struct stays on the stack if it doesn’t escape:
```
type Big struct {
    Data [4096]byte
}

func f() {
    b := Big{} // stays on stack if no escape
}

```

5. Struct comparision. A struct type is automatically comparable when every field is comparable. Comparable types include: numeric types, booleans, strings, pointers, channels, arrays and structs (recursevely, if all fields are comparable). Not comparable: slices, map and functions. 

6. Structs and the GO ABI (small structs may be passed in registers). Go’s modern ABI (post‑Go 1.17) is register‑based and much more efficient than the older stack‑only ABI. Under the current GO ABI, small structs with register-friendly fields (e.g. ints) may be passed entirely in registers. 

ABI =>  think of ABI (Application Binary Interface) as the rules of the road that all compiled Go code must follow. 

## 3.2.2 Struct declaration
1. Declaring a struct with a name typed:
```
type Person struct {
    Name string
    Age  int
}

p := Person{"Alice", 30}
p2 := Person{Age: 40} // keyed literal
```

2. Anonymous Structs (without a type name):
```
p := struct {
    Name string
    Age  int
}{
    Name: "Alice",
    Age:  30,
}

```

```
var person struct {
	name string
	age int
	pet string
}

person.name = "Jorge"
person.age = 45
person.pet = "dog"

// another way od declaring an anonymous struct. The firt one is more common
```

Anonymous struct are useful when the struct is used only once, : (e.g. in testing or marshalling).
```
tests := []struct {
    input string
    want  int
}{
    {"hello", 5},
    {"go", 2},
}

```

```
payload := struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}{
    Name: user.Name,
    Age:  user.Age,
}

json.Marshal(payload)

```

## 3.2.3 Zero value 
All fields are set to their own zero values
```
type User struct {
    Name string
    Age  int
    Active bool
}

var u User   

// the value of Name is "", the value of Age is 0, and the value of Active is false. 

```
## 3.2.4 Struct conversion
Struct conversion is when you explicitly convert one struct type into another:
```
var b B = B(a)
// This is only allowed when the two struct types have identical field sets and identical field order.
```

This is only allowed if:
- They have the same number of fields
- The fields appear in the same order
- The fields have identical types
- The fields have identical tags (tags must match exactly)

Assigment requires the same type:
```
package main

func main() {
	type Employee struct {
		name string
		age  int
	}

	type Person struct {
		name string
		age  int
	}

	a := Employee{"Jorge", 34}
	var person1 Person = a // cannot use a (variable of struct type Employee) as Person value in variable declaration.
}
```

Conversion creates a copy of the target type: 
```
package main

import "fmt"

func main() {
	type Employee struct {
		name string
		age  int
	}

	type Person struct {
		name string
		age  int
	}

	a := Employee{"Jorge", 34}
	var person1 Person = Person(a)
	fmt.Println(person1) // {Jorge 34}
}

```

## 3.2.5 Name Struct literals
There are two types:
1. Keyed form (the order does not matter, and you don't need to specify all fields):
```
package main

import "fmt"

func main() {
	type Employee struct {
		name    string
		age     int
		balance int
	}

	employee1 := Employee{age: 45, name: "Jorge"}
	fmt.Println(employee1) // {Jorge  45   0}
}

```
2. Positional form (order matters, and you need to specify all fields):
```
package main

import "fmt"

func main() {
	type Employee struct {
		name    string
		age     int
		balance int
	}

	employee1 := Employee{"Jorge", 45, 1000}
	fmt.Println(employee1) // {Jorge  45   1000}
}
```

## 3.2.6 Methods
A method are functions associated with a type. A method is a function with a special parameter called a receiver. In the case of a struct types it looks like this:

```
package main

import "fmt"

type Employee struct {
	name    string
	age     int
	balance int
}

func (e *Employee) deposit(amount int) {
	e.balance += amount
}

func main() {
	employee1 := Employee{"Jorge", 45, 1000}
	employee1.deposit(2000)
	fmt.Println(employee1.balance) // 3000
}
```
### Two kinds of receivers 
1. Value receiver: the value receives a copy of the struct, does not mutate the original struct.
```
func (p Point) Length() int {
    return p.X*p.X + p.Y*p.Y
}

```
2. Pointer receiver: the value receives a pointer to the struct (useful, avoid copying large structs), mutation affects the original struct (required if we need to modify the original struct).  
```
func (p *Point) Move(dx, dy int) {
    p.X += dx
    p.Y += dy
}

```
### Automatic address-taking 
1. You can call a pointer-receiver method on a value:
```
p := Point{1, 2}
p.Move(10, 20) // Go automatically does (&p).Move(...)

```

2. You can call a value-receiver method on a pointer:
```
pp := &Point{1, 2}
pp.Length() // Go automatically dereferences . (*ptr).Method()

```

### Method sets 
A method set is the collection of methods that a type exposes.
Let’s say you have:
```
type T struct{}
func (t T)  A() {}
func (t *T) B() {}

```

The method set of T (a value) is only A(). The method set of *T (a pointer) is A() and B(). A pointer can always be dereferenced to get a value, but a value cannot always be addressed to get a pointer. 

Automatic address‑taking does NOT change method sets. You can call a pointer‑receiver method on a value, but this does not mean T has method B. 
```
t := T{}
t.B() // Go rewrites to (&t).B()   // even if we are writing t.B(), t does not have a B() method, &t does.

```
It’s just syntactic sugar for method calls. When checking interface implementation, Go uses the method set, not the sugar.
 

## 3.2.7 Embedding structs 
```
type A struct {
    X int
}

type B struct {
    A   // embedded
    Y int
}

```
Under the hood, Go treats this exactly as:
```
type B struct {
    A A
    Y int
}

```
This is not inheritance, just composition.

1. Embedding gives you field promotion and method promotion:
- Field promotion: you can access A fields directly through B:
```
b := B{}
b.X = 10   // promoted. The compiler rewrites this to b.A.X = 10
b.Y = 20

```
- Method promotion: 
```
// If A has methods:
func (a A) Foo() {}
func (a *A) Bar() {}

```
then 
```
b.Foo() // becomes b.A.Foo()
b.Bar() // becomes (&b.A).Bar()

```

### Compile‑time rewriting, not dynamic dispatch.
In Go, method calls on embedded structs (and methods in general) use compile‑time rewriting, not dynamic dispatch. This is one of the most important distinctions between Go and classical OOP languages like Java.

```
type A struct{}
func (A) Foo() {}

type B struct {
    A
}

b.Foo()  // Foo actually belongs to an embedded field

```
Go does not perform any runtime lookup. Instead, the compiler rewrites the call to: b.A.Foo(). This happens before the program runs. There is no vtable, no dynamic dispatch, no runtime method resolution. It’s just syntactic sugar. 

What “dynamic dispatch” would look like (but Go does NOT do)? In classical OOP: The runtime looks up the method in a vtable. The actual method called depends on the runtime type. Overriding changes behavior dynamically. Go does none of this for struct methods.

### Embedding vs inheritance 
1. Core idea:

| Concept    | 	Go embedding                            | OOP Inheritance                          |
| ---------  | ---------------------------------------- | ---------------------------------------- |
| What is it | Composition: one struct contains another | Hierarchy: on class derives from another |
| Relationship | Has-a | Is-a |
| Mechanism | Compile-time/ method promotion  | Runtime polymorphism + class hierarchy |

2. How methods are resolved:

| Aspect     | 	Go embedding                            | OOP Inheritance                          |
| ---------  | ---------------------------------------- | ---------------------------------------- |
| Method resolution | Compile time rewriting (b.A.Foo()) | Dynamic dispatch via vtables |
| Overrriding | 	Not supported | 	Fully supported | 
| Runtime cost | Zero | some overhead |


3. Mental model:

- Embedding = composition + compile‑time promotion  
- Inheritance = hierarchy + runtime polymorphism

### Method sets with embedding 
When you embed a struct inside another struct, Go promotes the embedded type’s methods into the outer type’s method set. But the promotion follows the same method‑set rules as if the embedded field were a normal field.

```
type A struct{}
func (A) Foo() {}
func (*A) Bar() {}
type B struct { 
	A 
}

// Method set of B (value): Foo() NOT Bar()
// Method set of *B (pointer): FOO AND Bar()
```

Automatic address‑taking does NOT change method sets: 

```
b := B{}
b.Bar() // allowed. B has not a Bar() method. This is rewrites to (&b).Bar()
```
### Embedding chains
An embedding chain is when a struct embeds another struct, which itself embeds another struct, and so on.
```
type A struct {
    X int
}

type B struct {
    A
}

type C struct {
    B
}

var c C c.X = 10   // The compiler rewrites this to: c.B.A.X = 10

```
# 4. Maps (reference types)
A map is a pointer to to a runtime map header (a descriptor), which in turn points to a bucket that holds the actual key/value pair.

map variable
    ↓
[ fat pointer descriptor ]
    ↓
pointer to hmap
    ↓
hmap struct
    ↓
pointer to buckets
    ↓
bucket array
    ↓
key/value pairs

##  Is a map variable a pointer? 
It is what we call a FAT pointer (a small descriptor containing a pointer to the hmap). It is not a RAW pointer type. The descriptor contains a pointer to the map header plus type metadata.


## What is inside the map header?
The GO runtime defines an internal struct called hmap (not accessible directly). It contains a lof of info:
- count: number of key/value pairs
- B: log2 of number of buckets 
- buckets: pointer to the bucket array
- oldbuckets: pointer used during incremental growth
- hash seed
- flags

## Implications
- Passing a map to a function passes the pointer (mutations inside the function affects the caller's map)
- Copying a map variable copies only the pointer (both variables refer to the same underlying structure)
- A nil map is just a nil pointer to an hmap (you can read from it, but not write to it)

## Time complexity
Lookup, insert and delete are average O(1)

## Values and keys
- Keys must be comparable (no slices, maps or functions)
- Values can be anything

# 4.1 Creating maps
1. Map literal
```
m := map[string]int{
    "a": 1,
    "b": 2,
}
```
2. make()
```
m := make(map[string]int)
fmt.Println(m)      // map[]
fmt.Println(len(m)) // 0
// Used when you want an empty map that you’ll populate dynamically.


n := make(map[string]int, 100)
// you can hint an expected size, this is not a hard limit, it is just a hint for bucket allocation.

```

```
package main

import "fmt"

func main() {
	m := make(map[string]int)
	fmt.Println(m) // map[]

	b := map[string]int{"Spain": 32321, "UK": 12312}
	fmt.Println(b) // map[Spain:32321 UK:12312]
}

```
3. Nil Map

```
package main

import "fmt"

func main() {
	var m map[string]int     // m == nil
	fmt.Println(len(m))      // A nil map has a length of 0
	fmt.Println(m["potato"]) // 0 (zero value of int)
}

// safe to read from: attempting to read a nil map always returns the zero value for the map's value.
// Unsafe to write to (panics)
// A nil map has a length of 0.

```

## Why are nil maps useful?
A nil maps means: this map has not been initialized, treat it as empty, but immutable. Nil maps behave like empty maps for reads. Writing to a nil map, panics.
1. Allows lazy allocation (you can defer the allocation until the moment you actually needed)
```
func (c *Cache) Set(k string, v int) {
    if c.m == nil {
        c.m = make(map[string]int)
    }
    c.m[k] = v
}
```

2. Maps, slices, channels, functions, and interfaces all have a meaningful zero value: nil. For maps specifically: nil map → “optional field not initialized yet”, empty map → “initialized but empty”. And the key: A nil map costs nothing. No allocation. No GC pressure. No constructor. This is why Go leans so heavily on nil-able types.

	- You can create a struct with no constructor and still have a meaningful state.

```
type User struct {
    Metadata map[string]string
}

// Metadata is nil, which means “no metadata provided”.
// No extra flags. No pointers. No wrappers.

```

	- Reads are safe:

```
u := User{}
fmt.Println(u.Metadata["role"]) // prints "", no panic
// A nil map behaves like an empty map for reads.

```
	- Writes require explicit allocation:
```
u.Metadata["role"] = "admin" // panic
// This prevents accidental mutation of uninitialized state.
```

	- So you can write idiomatic patterns like:

```
func (u *User) SetMeta(k, v string) {
    if u.Metadata == nil {
        u.Metadata = make(map[string]string)
    }
    u.Metadata[k] = v
}

```
	- In many languages optional fields require pointers, in GO map headers contains a pointer already. Using *map is almost always a mistake. 

```
type User struct {
    Metadata *map[string]string // unnecessary in Go
}

```

A practical example (optional configurations fields):
```
type Config struct {
    Flags map[string]bool
}

cfg := Config{} // zero value

if cfg.Flags["debug"] {
    // safe read, even though Flags is nil
}

if cfg.Flags == nil {
    cfg.Flags = make(map[string]bool)
}
cfg.Flags["debug"] = true

```
You get optionality for free, with no runtime overhead until you actually need the map (make the zero value useful). 
- Nil maps require no allocation
- They require no constructor 
- They require no extra memory 
- No Garbage Collector work
- Allow safe reads
- Allow LAZY writes  

# 4.2 Operations with a map
## Read
```
v := m["x"] // returns zero value if key not present

```
When you try to read the value assigned to a map key that was never set, the map returns the zero value for the map's value type. 

1. The fact that returns zero value if the key is not present can be useful.
- Counting:
```
counts := map[string]int{}
counts["x"]++ // safe even if "x" wasn't present

```
- set iplemented as map[T]bool:
```
type StringSet map[string]bool

func NewStringSet() StringSet {
    return make(StringSet)
}

func (s StringSet) Add(v string) {
    s[v] = true
}

func (s StringSet) Remove(v string) {
    delete(s, v)
}

func (s StringSet) Has(v string) bool {
    return s[v] // returns false if missing
}

func (s StringSet) Size() int {
    return len(s)
}

// usage
s := NewStringSet()

s.Add("apple")
s.Add("banana")

fmt.Println(s.Has("apple"))   // true
fmt.Println(s.Has("orange"))  // false

s.Remove("banana")
fmt.Println(s.Has("banana"))  // false

```
This works because you cannot have duplicates keys in a map and because looking for a key that does not exist it returns the zero value of a boolean (false)
 
2. But how can we differenciete between a zero (no key present) and a zero (that is the value of the key) ? We use v, ok (check existance)

```
m := map[string]int{
    "a": 0,
}

v, ok := m["a"]
fmt.Println(v, ok) // 0 true

v, ok = m["b"]
fmt.Println(v, ok) // 0 false

```

## Write
```
m["x"] = 42

```

## Delete
```
delete(m, "x")

```

## Emptying a map
The clear function empties the map and set its length to zero.

```
package main

import "fmt"

func main() {
	m := map[string]int{"hello": 5, "world": 10}
	clear(m)
	fmt.Println(m, len(m)) // map[]   0
}

```

## Check existence 

```
v, ok := m["x"]
if ok {
    fmt.Println("exists:", v)
}

```
## Iteration
Maps a re unordered. Iteration order is not guaranteed. 

```
for k, v := range m {
    fmt.Println(k, v)
}

```

## Conparing maps
Maps are not comparable. You can check if they are equal to nil, but you cannot check if two maps have identical keys and values using == or !=. GO 1.21 added a package called `maps` that contains `maps.Equal` that are useful for comparing maps. 

package main

import (
	"fmt"
	"maps"
)

func main() {
	m := map[string]int{"hello": 5, "world": 10}
	l := map[string]int{"world": 10, "hello": 5}
	fmt.Println(maps.Equal(m, l)) // true
}

# 4.3 Maps and concurrency

The built‑in map type is not safe for concurrent writes, and even some read/write combinations can corrupt the internal hash table. This isn’t a “maybe” problem — it’s a guaranteed race condition that the runtime actively detects and panics on.

- Maps are not thread‑safe. Go’s map implementation is optimized for speed, not safety. There are no internal locks, no atomic operations.

So what should you do instead?
- Option 1: Protect the map with a mutex
- Option 2: Use sync.Map (specialized for concurrent patterns)
- Option 3: Use channels to serialize access

# 4.4. Difference between passing a map and passing a pointer to a map to a function.
Almost always, there is no benefit to passing a pointer to a map. Maps are already reference types: A map value in Go is a small header containing a pointer to the underlying hash table + metadata. When you pass a map to a function, you’re copying only the header, not the entire table. The underlying data is shared. Modifying the map inside the function modifies the caller’s map. 

YOU DON'T NEED A POINTER TO MUTATE A MAP

But, if you insist, what is the difference?
- Passing a map:
	- You can mutate the map (insert, delete, update)
	- You cannot reassign the map itself (i.e., replace the header)
- Passing a pointer to a map:
	- You can mutate the map (same as above)
	- You can also replace the map entirely

Passing a map:

```
func reset(m map[string]int) {
    m = make(map[string]int) // modifies only the local copy of the header
}

func main() {
    m := map[string]int{"a": 1}
    reset(m)
    fmt.Println(m) // still {"a": 1}
}

```

Passing a pointer to a map:

```
func reset(pm *map[string]int) {
    *pm = make(map[string]int)
}

func main() {
    m := map[string]int{"a": 1}
    reset(&m)
    fmt.Println(m) // now empty
}

```

When should you pass a pointer to a map?
- When you need to replace the map entirely

# 4.5 Maps of structs (map[k]T) vs maps of pointers to structs (map[k]*T)
A map of pointers:
```
map[string]*User
// allocates one object per entry on the heap
```

A map of structs:
```
map[string]User
// stores the struct inside the map buckets, not as separate heap options. This means: fewer allocations, less GC work.
```

So map of structs:
1. Avoid unnecessary heap allocations
2. The zero value works beautifully
```
type Stats struct {
    Count int
    Sum   float64
}

m := map[string]Stats{}
s := m["x"]
fmt.Println(s) // {0 0}
```
3. They avoid accidental shared mutation
```
map[string]*User
// Mutating it mutates the shared state.

map[string]User
// you get a copy. You can’t accidentally mutate the map entry unless you explicitly reassign.
```

When maps of pointers are better:
- the struct is large
- you mutate fields frequently
- you want shared references
- you want to avoid copying large values

# 5. Slices (reference types)
- A slice is a descriptor that contains 3 pieces of information: a pointer to a backing array, the length of the slice and the capacity of backing array. 
- What makes slices so useful is that you can grow slices as needed. 

# 5.1 Declaring slices
1. Declare a nil slice
```
var s []int

// the length is 0
// the capacity is 0
// s == nil is true
// No backing array yet
// Append works with nil slice

s = append(s, 10)  

```
2. Declare an empty but non-nil slice

```
s := []int{}
// the length is 0
// the cap is 0
// s != nil 
// Backing arrays exist (zero-length)
// append works 

s = append(s, 10)

```

3. Create a slice with make
```
s := make([]int, 5)        // len=5, cap=5
t := make([]int, 5, 10)    // len=5, cap=10 (capacity of the backing array)

```
# Zero value 
# 5.* Slicing slices
# 5.* Append
# 5.* Capacity
# 5.* Make
# Emptying a slice
# Copying a slice 
# Comparing slices
# A slice that is passed to a function can have its contents modified, but the size of the original slice cannot change.
# You can pass a slice of any size to a function but you cannot do the same with an array
# What is the difference bewtween passing a slice and passing a pointer to a slice?
# Converting arrays to slices 
# Converting slices to arrays 
# 5.* Differences between arrays and slices