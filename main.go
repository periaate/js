package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dop251/goja"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: js <arrow function>")
		return
	}
	var scanner *bufio.Scanner
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	}
	os.Args = os.Args[1:]
	vm := goja.New()
	expr := fmt.Sprintf("const fn = %s", os.Args[len(os.Args)-1])

	for _, arg := range os.Args[:len(os.Args)-1] {
		_, err := vm.RunString(arg)
		if err != nil {
			log.Fatalln("failed to parse expression", err)
		}
	}

	_, err := vm.RunString(expr)
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
			panic(err)
		}
		str := res.ToString().String()
		if len(str) != 0 {
			fmt.Println(res)
		}
	}
}
