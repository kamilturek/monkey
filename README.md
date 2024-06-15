# Monkey

<p align="center">
    <img src="img/monkey.webp" width="400">
</p>

A Monkey programming language interpreter from the book [Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball.

## Standard Types

Monkey supports several basic data types.

### Integer

```
let x = 12345;
```

### Boolean

```
let x = true;
let y = false;
```

## Operators

Monkey includes a variety of operators for performing arithmetic and comparisons.

### Basic Arithemtic Operators

You can use the following arithmetic operators: `+`, `-`, `*`, and `/`.

```
let x = 1 + 2 - 3 * 4;
```

### Comparison Operators

Monkey supports comparison operators such as `>`, `<`, `==`, and `!=`.

```
let x = 5 > 5;
let y = 5 < 5;
let z = 5 == 5;
let v = 5 != 5;
```

### Operator Precedence

The following table shows the operator precedence in Monkey, from lowest to highest:

| Precedence Level | Operators       | Description                |
|------------------|-----------------|----------------------------|
| 6 (Highest)      | Function calls  | Function calls             |
| 5                | Prefix `-`, `!` | Unary operations           |
| 4                | `*`, `/`        | Multiplication and Division|
| 3                | `+`, `-`        | Addition and Subtraction   |
| 2                | `<`, `>`        | Comparison                 |
| 1 (Lowest)       | `==`, `!=`      | Equality                   |

### Grouped Expressions

You can use parentheses to influence the order of executing arithmetic operations.

```
let x = (2 / (5 + 5));
```
