package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"lo/eval"
	"lo/expand"
	"lo/lexer"
	"lo/object"
	"lo/parser"
)

var globalEvaluator = &eval.Evaluator{}
var expander = expand.New(globalEvaluator)

func init() {
	globalEvaluator.Expander = expander
	eval.SetEvaluator(globalEvaluator)
}

//go:embed stdlib.lo
var stdlibEmbed embed.FS

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		runRepl()
		return
	}

	switch args[0] {
	case "new":
		if len(args) < 2 {
			fmt.Println("Usage: lo new <project-name>")
			os.Exit(1)
		}
		createNewProject(args[1])
	case "run":
		runProject(args[1:])
	default:
		runFile(args[0], args[1:])
	}
}

func createNewProject(name string) {
	err := os.Mkdir(name, 0755)
	if err != nil {
		fmt.Printf("Error creating project directory: %s\n", err)
		os.Exit(1)
	}

	err = os.Mkdir(name+"/src", 0755)
	if err != nil {
		fmt.Printf("Error creating src directory: %s\n", err)
		os.Exit(1)
	}

	mainContent := `(defn main [& args]
  (println "Hello from " "` + name + `!")
  (println "Arguments: " args))
`
	err = os.WriteFile(name+"/src/main.lo", []byte(mainContent), 0644)
	if err != nil {
		fmt.Printf("Error creating main.lo: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created new project: %s\n", name)
}

func runProject(args []string) {
	mainPath := "src/main.lo"
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		fmt.Println("Error: src/main.lo not found")
		os.Exit(1)
	}
	runFile(mainPath, args)
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()
	loadStdlib(env)

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

		expanded, err := expander.ExpandMacros(program, env)
		if err != nil {
			fmt.Printf("expansion error: %s\n", err)
			continue
		}

		result := globalEvaluator.Eval(expanded, env)
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
	contents, err := os.ReadFile(path)
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
	loadStdlib(env)

	expanded, err := expander.ExpandMacros(program, env)
	if err == nil {
		_ = globalEvaluator.Eval(expanded, env)
	}

	main, ok := env.Get("main")
	if ok && main.Type() == object.FUNCTION_OBJ {
		var s string
		if len(args) > 0 {
			quotedArgs := make([]string, len(args))
			for i, arg := range args {
				quotedArgs[i] = fmt.Sprintf("%q", arg)
			}
			s = fmt.Sprintf("(main %s)", strings.Join(quotedArgs, " "))
		} else {
			s = "(main)"
		}
		l = lexer.New(s, "main")
		p = parser.New(l)
		program = p.Parse()

		expanded, err = expander.ExpandMacros(program, env)
		if err == nil {
			globalEvaluator.Eval(expanded, env)
		}

	}

}

func printParserErrors(errors []parser.ParseError) {
	for _, msg := range errors {
		fmt.Println("\t" + msg.Msg)
	}
}

func loadStdlib(env *object.Environment) {
	contents, err := stdlibEmbed.ReadFile("stdlib.lo")
	if err != nil {
		return
	}

	l := lexer.New(string(contents), "stdlib.lo")
	p := parser.New(l)
	program := p.Parse()

	if len(p.Errors) == 0 {
		expanded, err := expander.ExpandMacros(program, env)
		if err == nil {
			globalEvaluator.Eval(expanded, env)
		}
	}
}
