# Monkey

<p align="center">
    <img src="img/monkey.webp" width="400" alt="logo">
</p>

A Monkey programming language interpreter from the book
[Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball.

## Table of Contents

- [Usage](#usage)
- [Standard Types](#standard-types)
  - [Integer](#integer)
  - [Boolean](#boolean)
  - [String](#string)
  - [Array](#array)
- [Operator Precedence](#operator-precedence)
  - [Operator Precedence](#operator-precedence)
  - [Grouped Expressions](#grouped-expressions)
- [Flow Control](#flow-control)
  - [If Expressions](#if-expressions)
- [Functions](#functions)
  - [Defining Functions](#defining-functions)
  - [Anonymous Functions](#anonymous-functions)
  - [First-Class Functions](#first-class-functions)

## Usage

Monkey can be used in the REPL mode (Read-Eval-Print Loop).
Simply run the `monkey` command without any arguments:

```sh
$ monkey
Hello username! This is the Monkey programming language!
Feel free to type in commands
>> 
```

## Standard Types

Monkey supports several basic data types.

### Integer

```javascript
let x = 12345;
```

You can use the following operators to perform arithmetic operations on integers
in Monkey: `+`, `-`, `*`, and `/`.

```javascript
let x = 1 + 2 - 3 * 4;
```

To compare integers, comparison operators such as `>`, `>=`, `<`, `<=`, `==`,
and `!=` can be used.

```javascript
let a = 5 > 5;   // false
let b = 5 >= 5;  // true
let c = 5 < 5;   // false
let d = 5 <= 5;  // true
let e = 5 == 5;  // true
let f = 5 != 5;  // false
```

### Boolean

```javascript
let x = true;
let y = false;
```

Booleans can be compared using the `==` and `!=` operators.

```javascript
let x = true;

if (x == true) {
    puts("x is true");
}
```

### String

```javascript
let x = "Hello World";
```

Strings, similarly to other types, can be compared using commparison operators
(`==` and `!=`).

The `+` operator concatenates two or more strings.

```javascript
let x = "Hello" + " " + "World!"; // "Hello World!"
```

### Array

Arrays can be used to store a sequence of objects that can be accessed
by an index.

```javascript
let names = ["John", "Bob", "Alice"];
names[0]; // "John"
```

Different object types can be mixed within a single array.

```javascript
let myArray = ["John", 1, fn(x) { return x; }, false];
```

## Operator Precedence

The following table shows the operator precedence in Monkey, from lowest to highest:

| Precedence Level | Operators            | Description                |
|------------------|----------------------|----------------------------|
| 6 (Highest)      | Function calls       | Function calls             |
| 5                | Prefix `-`, `!`      | Unary operations           |
| 4                | `*`, `/`             | Multiplication and Division|
| 3                | `+`, `-`             | Addition and Subtraction   |
| 2                | `<`, `>`, `<=`, `>=` | Comparison                 |
| 1 (Lowest)       | `==`, `!=`           | Equality                   |

### Grouped Expressions

You can use parentheses to influence the order of executing arithmetic operations.

```javascript
let x = (2 / (5 + 5));
```

## Flow Control

### If Expressions

Monkey supports `if` expressions for flow control. An `if` expression evaluates
a condition and executes the corresponding block of code.

The syntax for an `if` expression is as follows:

```javascript
if (condition) { 
    // block of code 
} else { 
    // optional else block 
}
```

- The `else` block is optional.
- Each block can contain multiple expressions or statements.
- The value of the `if` expression is the value of the last expression in the
executed block.

#### Example

```javascript
let x = 10;
let y = 20;

let max = if (x > y) {
    x
} else {
    y
};
```

In this example, max will be set to `20` because `y` is greater than `x`.

## Functions

Monkey supports functions, which can be defined, assigned to variables, called,
and passed as arguments to other functions.

### Defining Functions

Functions can be defined using the `fn` keyword:

```javascript
let add = fn(x, y) {
    return x + y;
};
```

The `return` statement is optional. If omitted, the last expression in the
function block will be returned:

```javascript
let add = fn(x, y) {
    x + y;
};
```

### Anonymous Functions

Functions can be called immediately without being assigned to a variable:

```javascript
let three = fn(x, y) { 
    x + y; 
}(1, 2);
```

### First-Class Functions

Functions in Monkey are first-class objects, meaning they can be passed as
arguments to other functions:

```javascript
let sayHello = fn() {
    print("hello");
};

let callTwice = fn(f) {
    f();
    f();
};

callTwice(sayHello);
```
