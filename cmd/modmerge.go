package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ScottBrooks/mergekeys"
)

var baseKey = flag.String("base", "chitin.key", "base game key file")
var modKey = flag.String("mod", "mod.key", "mod key file")
var outputKey = flag.String("output", "output.key", "modified output key file")

func main() {
	flag.Parse()

	baseKeyIn, err := os.Open(*baseKey)
	if err != nil {
		log.Fatal(err)
	}
	defer baseKeyIn.Close()

	modKeyIn, err := os.Open(*modKey)
	if err != nil {
		log.Fatal(err)
	}
	defer modKeyIn.Close()

	bk, err := mergekeys.OpenKEY(baseKeyIn, "")
	if err != nil {
		log.Fatal(err)
	}

	mk, err := mergekeys.OpenKEY(modKeyIn, "")
	if err != nil {
		log.Fatal(err)
	}

	err = bk.MergeWith(mk)
	if err != nil {
		log.Fatal(err)
	}

	outKey, err := os.Create(*outputKey)
	if err != nil {
		log.Fatal(err)
	}

	bk.Write(outKey)

	outKey.Close()

	fmt.Printf("Conversion complete.  Rename %s to chitin.key\n", *outputKey)
}
