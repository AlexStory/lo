package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"lo/eval"
	"lo/lexer"
	"lo/object"
	"lo/parser"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		runRepl()
	} else {
		switch args[0] {
		default:
			fmt.Printf("Unknown command: %s\n", args[0])
		}
	}
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	fmt.Println("lo v0.0.1")
	fmt.Println()

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		switch line {
		case ".exit":
			return
		}

		l := lexer.New(line, "repl")
		p := parser.New(l)
		program := p.Parse()

		if len(p.Errors) != 0 {
			printParserErrors(p.Errors)
			continue
		}

		result := eval.Eval(program, env)
		if result != nil {
			fmt.Println(result.Inspect())
		}
	}
}

func printParserErrors(errors []parser.ParseError) {
	for _, msg := range errors {
		fmt.Println("\t" + msg.Msg)
	}
}
