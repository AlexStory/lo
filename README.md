# lo

`lo` is a simple Lisp-like interpreter written in Go.

## Features

- **REPL**: Interactive Read-Eval-Print Loop.
- **File Execution**: Run `.lo` scripts from the command line.
- **Basic Types**: Integers, Floats, Strings, Booleans, Lists, and Nil.
- **Functions**: Define your own functions with `defn`.
- **Variables**: Assign values with `def`.
- **Modern Types**: Keywords (e.g., `:key`) and Maps (e.g., `{:key value}`).
- **Built-in Functions**:
    - Math operations (`+`, `-`, `*`, `/`).
    - Printing (`print`, `println`).
    - String conversion (`str`).
    - List manipulation (`head`).
    - Map manipulation (`get`, `assoc`, `dissoc`).
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

### Function Definition

```lisp
(defn add [x y]
    (+ x y))
```

### Variable Definition

```lisp
(def n 10)
```

### Main Function

If a `main` function is defined in a file, it will be executed when the file is run, receiving command-line arguments.

```lisp
(defn main [x y]
    (def n (+ x y))
    (println n))
```

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
