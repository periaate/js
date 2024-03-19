package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/dop251/goja"
	gf "github.com/jessevdk/go-flags"
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"Print verbose output."`
	Errors  bool `short:"e" long:"errors" description:"Panic on errors."`
}

func main() {
	opts := &Options{}
	parser := gf.NewParser(opts, gf.Default)
	parser.Name = "js"
	parser.Usage = "[options] <arrow function>"

	rest, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	if len(rest) < 2 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if !opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelError + 1)
	}

	var scanner *bufio.Scanner
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	}

	vm := goja.New()

	var t string
	if len(rest) == 1 {
		t = rest[0]
	} else {
		t = strings.Join(rest, " ")
	}

	expr := fmt.Sprintf("const fn = %s", t)

	_, err = vm.RunString(expr)
	if err != nil {
		log.Fatalln("failed to parse expression", err)
	}

	fn, ok := goja.AssertFunction(vm.Get("fn"))
	if !ok {
		panic("Not a function")
	}

	for scanner.Scan() {
		text := scanner.Text()

		res, err := fn(goja.Undefined(), vm.ToValue(text))
		if err != nil {
			slog.Error("error evaluating expression", "input", text, "error", err)
			if opts.Errors {
				panic(err)
			}
		}
		str := res.ToString().String()
		if len(str) != 0 {
			fmt.Println(res)
		}
	}
}
