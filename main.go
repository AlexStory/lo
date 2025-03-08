package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
			runFile(args[0], args[1:])
		}
	}
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	fmt.Println("lo v0.0.1")
	fmt.Println("type '.exit' to exit")
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

func runFile(path string, args []string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// read all contents of the file to a string
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %s\n", err)
		os.Exit(1)
	}
	l := lexer.New(string(contents), file.Name())
	p := parser.New(l)
	program := p.Parse()

	if len(p.Errors) != 0 {
		printParserErrors(p.Errors)
		os.Exit(1)
	}

	env := object.NewEnvironment()
	_ = eval.Eval(program, env)

	main, ok := env.Get("main")
	if ok && main.Type() == object.FUNCTION_OBJ {
		s := fmt.Sprintf("(main %s)", strings.Join(args, " "))
		l = lexer.New(s, "main")
		p = parser.New(l)
		program = p.Parse()

		eval.Eval(program, env)

	}

}

func printParserErrors(errors []parser.ParseError) {
	for _, msg := range errors {
		fmt.Println("\t" + msg.Msg)
	}
}
