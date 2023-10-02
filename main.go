package main

import (
	"log"
	"os"

	"github.com/ejuju/jus/pkg/jul"
)

func main() {
	// Start in REPL or file mode
	from := os.Stdin
	if len(os.Args) > 1 {
		var err error
		from, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()
	}

	// Execute code from stdin or file
	vm := jul.NewVM()
	err := vm.Execute(from)
	if err != nil {
		log.Println(err)
	}
}
