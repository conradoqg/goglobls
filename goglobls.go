package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/kr/fs"
	glob "github.com/obeattie/ohmyglob"
	yaml "gopkg.in/yaml.v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: goglobls [options] [path...]\n")
	flag.PrintDefaults()
}

// Config Configuration structure
type Config struct {
	Types []struct {
		Name  string   `yaml:"name"`
		Paths []string `yaml:"paths"`
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	flag.Usage = usage
	configFile := flag.String("config", "", "Config file (required)")
	isSense := flag.Bool("i", false, "Is Insensitive (Unix Only)")

	var types arrayFlags
	flag.Var(&types, "type", "Type to filter (add more than one by repeating the option)")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		fmt.Printf("Path is missing.")
		os.Exit(2)
	}

	if *configFile == "" {
		usage()
		fmt.Printf("Config file is missing.")
		os.Exit(2)
	}

	data, err := ioutil.ReadFile(*configFile)

	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(2)
	}

	config := Config{}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	patterns := make([]string, 0)

	for _, _type := range config.Types {
		if len(types) > 0 {
			foundType, _ := inArray(_type.Name, types)
			if !foundType {
				continue
			}
		}
		for _, pattern := range _type.Paths {
			if runtime.GOOS == "windows" {
				patterns = append(patterns, strings.ToLower(pattern))
			} else {
				if *isSense == true {
					patterns = append(patterns, strings.ToLower(pattern))
				} else {
					patterns = append(patterns, pattern)
				}
			}
		}
	}

	set, _ := glob.CompileGlobSet(patterns, glob.DefaultOptions)

	walker := fs.Walk(args[0])

	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		var withFixedSeparator string

		if runtime.GOOS == "windows" {
			withFixedSeparator = strings.ToLower(strings.Replace(walker.Path(), "\\", "/", -1))
		} else {
			if *isSense == true {
				withFixedSeparator = strings.ToLower(walker.Path())
			} else {
				withFixedSeparator = walker.Path()
			}
		}

		if walker.Stat().Mode().IsRegular() && set.MatchString(withFixedSeparator) {
			fmt.Printf("%v\n", walker.Path())
		}
	}
}

func inArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}
