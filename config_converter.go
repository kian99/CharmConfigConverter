package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"text/template"

	charm "github.com/juju/charm/v9"
)

const (
	filename = "config.yaml"
	output   = "test.tf"
)

var global_series string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file, err := os.Open(filename)
	check(err)
	config, err := charm.ReadConfig(file)
	file.Close()
	check(err)
	f, err := os.Create(output)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = terraform_print_vars(config, w)
	w.Flush()
	check(err)
	fmt.Print("Done!\n")
	err = terraform_print_bundle()
	check(err)
}

func terraform_print_bundle() error {
	bundle, err := charm.ReadBundle("./bundle.yaml")
	if err != nil {
		return err
	}
	print("Printing bundle\n")
	print(bundle)
	return nil
}

func terraform_print_vars(config *charm.Config, w io.Writer) error {
	var type_map = map[string]string{"string": "string",
		"int":     "number",
		"float":   "number",
		"boolean": "bool"}
	for name, option := range config.Options {
		_, err := fmt.Fprintf(w, "variable \"%s\" {\n\ttype = %s", name, type_map[option.Type])
		if err != nil {
			return err
		}
		if option.Default != nil {
			if option.Type == "string" {
				//option.Default = "\"" + option.Default + "\""
				_, err = fmt.Fprintf(w, "\n\tdefault=%q", option.Default)
			} else {
				_, err = fmt.Fprintf(w, "\n\tdefault=%v", option.Default)
			}

			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(w, "\n}\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func templateTest(config *charm.Config) {
	//Using Go fmt
	new_f, err := os.Create("Temptest.tf")
	check(err)
	defer new_f.Close()
	new_w := bufio.NewWriter(new_f)
	doTemplate(config, new_w)
	new_w.Flush()
}

func doTemplate(config *charm.Config, w io.Writer) {
	t, err := template.New("template_name").Parse(variable_template)
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, config.Options)
	if err != nil {
		panic(err)
	}
}
