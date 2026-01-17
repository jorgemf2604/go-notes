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

## 1.2.5 Scope (book) => add it above
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
- how importing works 
- exported vs exported
## 1.2.8 Package initialization 






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

