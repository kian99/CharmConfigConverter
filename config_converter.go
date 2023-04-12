package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	charm "github.com/juju/charm/v9"
)

var type_map = map[string]string{
	"string":  "string",
	"int":     "number",
	"float":   "number",
	"boolean": "bool",
}

var sensitive_keywords = []string{
	"password",
	"cert",
	"certificate",
	"key",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Please provide a input file and output file.")
		fmt.Printf("E.g. %s ./config.yaml ./output.tf\n", os.Args[0])
		os.Exit(1)
	}
	input := args[0]
	output := args[1]
	file, err := os.Open(input)
	check(err)
	defer file.Close()
	config, err := charm.ReadConfig(file)
	check(err)
	f, err := os.Create(output)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = terraform_print_vars(config, w)
	w.Flush()
	check(err)
	fmt.Print("Done!\n")
}

func terraform_print_vars(config *charm.Config, w io.Writer) error {

	for name, option := range config.Options {
		_, err := fmt.Fprintf(w, "variable \"%s\" {\n\ttype = %s", name, type_map[option.Type])
		if err != nil {
			return err
		}
		if option.Default != nil {
			if option.Type == "string" {
				_, err = fmt.Fprintf(w, "\n\tdefault=%q", option.Default)
			} else {
				_, err = fmt.Fprintf(w, "\n\tdefault=%v", option.Default)
			}

			if err != nil {
				return err
			}
		}
		if option.Description != "" {
			option.Description = strings.Replace(option.Description, "\n", " ", -1)
			option.Description = strings.Replace(option.Description, "\"", "", -1)
			_, err = fmt.Fprintf(w, "\n\tdescription=%q", option.Description)
			if err != nil {
				return err
			}
		}
		if containsAny(name, sensitive_keywords) {
			_, err = fmt.Fprint(w, "\n\tsensitive=true")
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

func containsAny(name string, checkList []string) bool {
	for _, item := range checkList {
		if strings.Contains(name, item) {
			return true
		}
	}
	return false
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
