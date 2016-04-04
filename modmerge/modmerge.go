package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ScottBrooks/modmerge"
)

var baseKey = flag.String("base", "chitin.key", "base game key file")
var name = flag.String("name", "sod-dlc", "name of the mod to load")

func main() {
	flag.Parse()

	if !promptForPermission(*name) {
		log.Printf("Exiting without making any changes")
		return
	}

	// Step 1: Backup chitin.key
	err := backup(*baseKey)
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: Load in chitin.key
	baseKeyIn, err := os.Open(*baseKey)
	if err != nil {
		log.Fatal(err)
	}
	defer baseKeyIn.Close()

	bk, err := mergekeys.OpenKEY(baseKeyIn, "")
	if err != nil {
		log.Fatal(err)
	}
	baseKeyIn.Close()

	// Step 3: Load mod key, and attempt to extract the mod if possible
	modKeyIn, err := os.Open(*name + ".key")
	if os.IsNotExist(err) {
		log.Printf("No %s.key found, attempting to extract %s", *name, *name)
		err = attemptExtractMod(*name)
		if err != nil {
			log.Fatal(err)
		}
		modKeyIn, err = os.Open(*name + ".key")
	} else {
		log.Fatalf("%s.key found, have you already run this tool?", *name)

	}
	if err != nil {
		log.Fatal(err)
	}

	mk, err := mergekeys.OpenKEY(modKeyIn, "")
	if err != nil {
		log.Fatal(err)
	}
	modKeyIn.Close()

	// Step 4: Merge our mod key into our base key
	err = bk.MergeWith(mk, *name)
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: Re-create our base key
	outKey, err := os.Create(*baseKey)
	if err != nil {
		log.Fatal(err)
	}

	bk.Write(outKey)
	outKey.Close()

	// Step 6: Rename our mod to disable it
	err = os.Rename("dlc/"+*name+".zip", "dlc/"+*name+".disabled")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Conversion complete.\n")
}

func backup(name string) error {
	in, err := os.Open(name)
	if err != nil {
		return err
	}

	out, err := os.Create(name + ".bak")
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)

	return err
}

func attemptExtractMod(name string) error {
	r, err := zip.OpenReader(filepath.Join("dlc", name+".zip"))
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(name, os.ModeDir|0777)

	log.Printf("Extracting files")
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			os.MkdirAll(f.Name, f.Mode())
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}

		fname := f.Name
		chunks := strings.Split(f.Name, "/")
		if strings.ToLower(chunks[0]) == "data" {
			chunks[0] = name
		}
		if strings.ToLower(chunks[0]) == "mod.key" {
			chunks[0] = name + ".key"
		}
		fname = path.Join(chunks...)
		out, err := os.Create(fname)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		if err != nil {
			rc.Close()
			return err
		}

		out.Close()
		rc.Close()
		fmt.Printf(".")
	}
	fmt.Printf("\n")
	log.Printf("Finished extracting files")

	return nil
}

var msg = `
The Plan:
	1: Backup your chitin.key as chitin.key.bak
	2: Load our chitin.key
	3: Load our %s key file
	   If unsuccessful, attempt to extract dlc/mod.zip, and create %s.key
	   The data/ folder inside %s.zip will be placed inside the folder %s/
	4: Merge our mod %s into the chitin.key
	5: Overwrite our chitin.key
	6: Rename dlc/%s.zip to dlc/%s.disabled

`
var prompt = "Continue? [y/n]"

func promptForPermission(name string) bool {
	fmt.Printf(msg, name, name, name, name, name, name, name)
	reader := bufio.NewReader(os.Stdin)
	input := ""

	for !(input == "Y" || input == "y" || input == "N" || input == "n") {
		fmt.Printf(prompt)
		input, _ = reader.ReadString('\n')

		input = strings.TrimSpace(input)
	}

	return input == "Y" || input == "y"
}
