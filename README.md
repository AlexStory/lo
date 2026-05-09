# lo

`lo` is a simple Lisp-like interpreter written in Go.

## Features

- **REPL**: Interactive Read-Eval-Print Loop.
- **File Execution**: Run `.lo` scripts from the command line.
- **Basic Types**: Integers, Floats, Strings, Booleans, Lists, and Nil.
- **Functions**: Define your own functions with `defn`.
- **Variables**: Assign values with `def`.
- **Macros & Standard Library**:
    - `defmacro` for defining macros.
    - Quasiquoting with `` ` ``, `~` (unquote), and `~@` (unquote-splicing).
    - `stdlib.lo` is automatically loaded if present.
- **Improved Functions**:
    - Variadic arguments with `&`, e.g., `(defn f [x & rest] ...)`.
- **Control Flow & Bindings**:
    - `let` for local bindings.
    - `if` and `when` for conditional execution.
    - `->` and `->>` threading macros for cleaner nested calls.
- **Built-in Functions**:
    - Math operations (`+`, `-`, `*`, `/`).
    - Printing (`print`, `println`).
    - String conversion (`str`).
    - List manipulation (`head`, `concat`).
    - Map manipulation (`get`, `assoc`, `dissoc`).
    - Sequence operations (`first`, `rest`, `cons`, `count`, `map`, `filter`, `reduce`).
    - Predicates (`empty?`, `nil?`).
    - Comparison & Logic (`=`, `<`, `>`, `<=`, `>=`, `not`, `and`, `or`).
    - Scripting: File I/O (`slurp`, `spit`) and OS interaction (`env`).

## Installation

Ensure you have [Go](https://golang.org/) installed.

```bash
go build -o lo main.go
```

## Usage

### REPL

Start the interactive shell:

```bash
./lo
```

### Running a File

Execute a script:

```bash
./lo examples/add.lo 1 2
```

## Syntax Overview

### Functions

Define your own functions with `defn` or anonymous functions with `\`.

```lisp
(defn add [x y]
    (+ x y))

(map (\ [x] (* x 2)) [1 2 3])
```

#### Lambda Literal

A concise syntax for anonymous functions, similar to Clojure.

- `%` or `%1`: First argument.
- `%n`: n-th argument.

```lisp
(map #(+ % 1) [1 2 3]) ; [2 3 4]
(#(+ %1 %2) 5 10)      ; 15
```

### Variable Definition

```lisp
(def n 10)
```

### Local Bindings

`let` allows defining local variables.

```lisp
(let [x 1 y 2]
    (+ x y))
```

### Threading Macros

Improve readability of nested function calls.

```lisp
;; Thread-first (->): inserts result as first argument
(-> 5 (+ 1) (* 2)) ; (* (+ 5 1) 2) => 12

;; Thread-last (->>): inserts result as last argument
(->> [1 2 3] (map #(+ % 1)) (filter #(> % 2))) ; [3 4]
```

### Projects

You can create a new project structure with `lo new`:

```bash
./lo new my-app
```

This creates a directory `my-app` with:
- `src/main.lo`: Entry point with a `main` function.
- `stdlib.lo`: A copy of the standard library.

To run the project, navigate into the directory and use `lo run`:

```bash
cd my-app
../lo run arg1 arg2
```

### Main Function

If a `main` function is defined in a file, it will be executed when the file is run, receiving command-line arguments.

```lisp
(defn main [x y]
    (def n (+ x y))
    (println n))
```

### Modules

`lo` has a simple module system. Filenames define namespaces.

- `(import "file.lo")`: Imports a file relative to the current file.
- `(import module)`: In a project, imports `src/module.lo` or `src/module/main.lo`.

Example:

`src/math.lo`:
```lisp
(defn add [x y] (+ x y))
```

`src/main.lo`:
```lisp
(import math)

(defn main []
  (println (math/add 1 2)))
```

Symbols from imported modules are prefixed with the module name and a slash, e.g., `math/add`.

### Arithmetic

```lisp
(+ 1 2)
(- 10 5)
(* 2 3 4)
(/ 100 2 2)
```

### Keywords and Maps

```lisp
(def person {:name "Alice" :age 30})
(println (get person :name))
;; or use keyword as getter
(println (:name person))
;; with default values
(println (get person :missing "not found"))
(println (:missing person "not found"))
```

### Map Manipulation

Maps are immutable. `assoc` and `dissoc` return new maps.

```lisp
(def m {:a 1})
(def m2 (assoc m :b 2 :c 3)) ; {:a 1 :b 2 :c 3}
(def m3 (dissoc m2 :a :b))    ; {:c 3}
```

### Scripting

```lisp
(spit "hello.txt" "Hello World")
(println (slurp "hello.txt"))
(println (env "PATH"))
```
