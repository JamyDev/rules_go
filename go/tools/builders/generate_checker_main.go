/* Copyright 2018 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Generates a checker binary for Bazel Go rules.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

var codeTpl = `
package main

import (
{{range $importPath := .}}
	_ "{{$importPath}}"
{{- end}}
)
`

func run(args []string) error {
	checkImportPaths := multiFlag{}
	flags := flag.NewFlagSet("generate_checker_main", flag.ExitOnError)
	out := flags.String("output", "", "output file to write (defaults to stdout)")
	flags.Var(&checkImportPaths, "check_importpath", "import path of a check library")
	// TODO(samueltan): use the config file to generate whitelist tables in the checker.
	flags.String("config", "", "checker config file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *out == "" {
		return errors.New("must provide output file")
	}

	outFile := os.Stdout
	var cErr error
	outFile, err := os.Create(*out)
	if err != nil {
		return fmt.Errorf("os.Create(%q): %v", *out, err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			cErr = fmt.Errorf("error closing %s: %v", outFile.Name(), err)
		}
	}()

	tpl := template.Must(template.New("source").Parse(codeTpl))
	if err := tpl.Execute(outFile, checkImportPaths); err != nil {
		return fmt.Errorf("template.Execute failed: %v", err)
	}
	return cErr
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
